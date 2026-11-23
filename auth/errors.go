package auth

import (
	"errors"
	"net/http"

	"connectrpc.com/connect/v2"
	"google.golang.org/protobuf/proto"

	"github.com/cogikyo/hgmx/pg"
	authv1 "github.com/cogikyo/hgmx/proto/hgmx/auth/v1"
)

var (
	errStale       = errors.New("auth: session ended or was revoked")
	errRecent      = errors.New("auth: the session is too old for this change")
	errDenied      = errors.New("auth: the host's Allow refused the user")
	errMissing     = errors.New("auth: row not found")
	errLast        = errors.New("auth: the last sign-in method stays")
	errChanged     = errors.New("auth: account email changed mid-request")
	errOwned       = errors.New("auth: site data still references the account")
	errNameTaken   = errors.New("auth: username taken")
	errKeyExists   = errors.New("auth: passkey already saved")
	errExpired     = errors.New("auth: ceremony or flow expired")
	errThrottled   = errors.New("auth: too many attempts from this network")
	errTooSoon     = errors.New("auth: address mailed too recently")
	errUnknown     = errors.New("auth: passkey not registered")
	errOrigin      = errors.New("auth: passkey response from another origin")
	errUnreadable  = errors.New("auth: unreadable passkey response")
	errUnverified  = errors.New("auth: passkey could not be verified")
	errWrongCode   = errors.New("auth: wrong code")
	errOldCode     = errors.New("auth: code expired, replaced, or unknown")
	errExhausted   = errors.New("auth: daily wrong-code budget spent")
	errBadEmail    = errors.New("auth: invalid email")
	errSameEmail   = errors.New("auth: email unchanged")
	errBadCode     = errors.New("auth: malformed code")
	errKeyName     = errors.New("auth: invalid passkey name")
	errName        = errors.New("auth: invalid username")
	errNameAt      = errors.New("auth: username has @ but is not the account email")
	errUnconfirmed = errors.New("auth: deletion not confirmed")
)

type reason struct {
	status      int
	code        connect.Code
	text        string
	browserText string
	detail      func(string) proto.Message
}

var ended = reason{status: http.StatusUnauthorized, code: connect.CodeUnauthenticated, text: "Your session ended. Sign in again."}

var reasons = map[error]reason{
	errStale:       ended,
	pg.ErrNotFound: ended, // User deleted mid-request; other missing rows use errMissing.
	errRecent: {
		status: http.StatusForbidden, code: connect.CodeFailedPrecondition,
		text:   "Sign in again to make this change.",
		detail: func(m string) proto.Message { return &authv1.SignInAgain{Message: m} },
	},
	errDenied:  {status: http.StatusForbidden, code: connect.CodePermissionDenied, text: "This account can't sign in here."},
	errMissing: {status: http.StatusNotFound, code: connect.CodeNotFound, text: "That item no longer exists."},
	errLast: {
		status: http.StatusConflict, code: connect.CodeFailedPrecondition,
		text:   "This is your only way to sign in. Add another one first.",
		detail: func(m string) proto.Message { return &authv1.LastSignInMethod{Message: m} },
	},
	errChanged:    {status: http.StatusConflict, code: connect.CodeAborted, text: "Your email changed while saving. Try again."},
	errOwned:      {status: http.StatusConflict, code: connect.CodeFailedPrecondition, text: "This account still owns data on this site, so it can't be deleted here."},
	errNameTaken:  {status: http.StatusConflict, code: connect.CodeAlreadyExists, text: "That username is taken. Choose another."},
	errKeyExists:  {status: http.StatusConflict, code: connect.CodeAlreadyExists, text: "That passkey is already saved."},
	errExpired:    {status: http.StatusBadRequest, code: connect.CodeInvalidArgument, text: "That took too long. Try again."},
	errThrottled:  {status: http.StatusTooManyRequests, code: connect.CodeResourceExhausted, text: "Too many attempts from this network. Try again later."},
	errTooSoon:    {status: http.StatusTooManyRequests, code: connect.CodeResourceExhausted, text: "We sent an email recently. Use the latest one, or try again in a little while."},
	errUnknown:    {status: http.StatusUnauthorized, code: connect.CodeUnauthenticated, text: "That passkey is not registered here."},
	errOrigin:     {status: http.StatusUnauthorized, code: connect.CodeUnauthenticated, text: "That passkey response doesn't belong to this app."},
	errUnreadable: {status: http.StatusBadRequest, code: connect.CodeInvalidArgument, text: "That passkey response was not readable."},
	errUnverified: {status: http.StatusBadRequest, code: connect.CodeInvalidArgument, text: "The passkey couldn't be verified. Try adding it again."},
	errWrongCode:  {status: http.StatusUnprocessableEntity, code: connect.CodeInvalidArgument, text: "That code isn't right. Check the latest email and try again."},
	errOldCode:    {status: http.StatusUnprocessableEntity, code: connect.CodeInvalidArgument, text: "That code has expired or was replaced. Ask for a new one."},
	errExhausted: {
		status: http.StatusTooManyRequests, code: connect.CodeResourceExhausted,
		text:        "Too many wrong codes for this address. Try again tomorrow.",
		browserText: "Too many wrong codes for this address. Use the link in the email, or try again tomorrow.",
	},
	errBadEmail:  {status: http.StatusUnprocessableEntity, code: connect.CodeInvalidArgument, text: "Enter a valid email address."},
	errSameEmail: {status: http.StatusUnprocessableEntity, code: connect.CodeInvalidArgument, text: "That's already your email."},
	errBadCode:   {status: http.StatusUnprocessableEntity, code: connect.CodeInvalidArgument, text: "Enter the 6-digit code from the email."},
	errKeyName:   {status: http.StatusUnprocessableEntity, code: connect.CodeInvalidArgument, text: "Passkey names need 1 to 64 characters."},
	errName:      {status: http.StatusUnprocessableEntity, code: connect.CodeInvalidArgument, text: "Use 1–64 letters (A–Z), digits, dots, dashes or underscores."},
	errNameAt:    {status: http.StatusUnprocessableEntity, code: connect.CodeInvalidArgument, text: "Use your account email as-is, or choose a username without @."},
	errUnconfirmed: {
		status: http.StatusUnprocessableEntity, code: connect.CodeInvalidArgument,
		text:        "Confirm deleting your account.",
		browserText: "Tick the box to confirm deleting your account.",
	},
}

func reasonOf(err error) (reason, bool) {
	for e, r := range reasons {
		if errors.Is(err, e) {
			return r, true
		}
	}
	return reason{}, false
}

func shown(err error) (int, string) {
	r, ok := reasonOf(err)
	if !ok {
		panic(err)
	}
	if r.browserText != "" {
		return r.status, r.browserText
	}
	return r.status, r.text
}

func (r reason) rpc() *connect.Error {
	e := connect.NewError(r.code, r.text)
	if r.detail != nil {
		return detailed(e, r.detail(r.text))
	}
	return e
}
