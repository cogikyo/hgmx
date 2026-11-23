package auth

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"time"
	"unicode"

	"connectrpc.com/connect/v2"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-webauthn/webauthn/protocol"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/cogikyo/hgmx/mail"
	"github.com/cogikyo/hgmx/pg"
	appv1 "github.com/cogikyo/hgmx/proto/hgmx/app/v1"
	authv1 "github.com/cogikyo/hgmx/proto/hgmx/auth/v1"
	"github.com/cogikyo/hgmx/proto/hgmx/auth/v1/authv1connect"
)

type appService struct{ api *API }

func (s *appService) GetConfiguration(ctx context.Context, _ *appv1.GetConfigurationRequest) (*appv1.GetConfigurationResponse, error) {
	c := callOf(ctx)
	res := &appv1.GetConfigurationResponse{
		ApiMajors:      []int32{1},
		MinimumBuild:   int64(c.app.build),
		StoreUrl:       c.app.store,
		UpdateRequired: c.build < c.app.build,
	}
	for spec := range s.api.server.Specs() {
		if s.api.auth.enabled(c.caller, spec.Procedure) {
			res.Methods = append(res.Methods, spec.Procedure)
		}
	}
	slices.Sort(res.Methods)
	return res, nil
}

func (a *Auth) enabled(c caller, procedure string) bool {
	switch procedure {
	case authv1connect.AuthServiceSendCodeProcedure, authv1connect.AuthServiceVerifyCodeProcedure:
		return c.host.mode.codes()
	case authv1connect.AuthServiceChangeEmailProcedure:
		return a.mail != nil
	case authv1connect.AuthServiceBeginGoogleSigninProcedure, authv1connect.AuthServiceFinishGoogleSigninProcedure,
		authv1connect.AuthServiceBeginGoogleLinkProcedure, authv1connect.AuthServiceFinishGoogleLinkProcedure:
		return a.tokens != nil && len(c.app.google) > 0
	case authv1connect.AuthServiceUnlinkGoogleProcedure:
		return a.googleOn(c.host)
	case authv1connect.AuthServiceDeleteAccountProcedure:
		return c.host.mode == Open
	}
	return true
}

func (a *Auth) googleOn(h *host) bool {
	if a.google != nil {
		return true
	}
	for _, p := range a.apps {
		if p.host == h && len(p.google) > 0 && a.tokens != nil {
			return true
		}
	}
	return false
}

type authService struct{ a *Auth }

func labeled(c caller, raw string) (caller, error) {
	l := strings.TrimSpace(raw)
	if len([]rune(l)) > 64 || strings.ContainsFunc(l, func(r rune) bool { return !unicode.IsPrint(r) }) {
		return c, connect.NewError(connect.CodeInvalidArgument, "Device labels need at most 64 printable characters.")
	}
	if l != "" {
		c.label = &l
	}
	return c, nil
}

func recently(s session) error {
	if !s.recent() {
		return errRecent
	}
	return nil
}

func (s *authService) grant(ctx context.Context, c caller, token string) (*authv1.SessionGrant, error) {
	c.token = token
	cur, ok, err := s.a.find(ctx, c)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errStale
	}
	acc, err := s.a.protoAccount(ctx, c, cur)
	if err != nil {
		return nil, err
	}
	return &authv1.SessionGrant{Token: token, ExpiresAt: timestamppb.New(cur.expires), Account: acc}, nil
}

func (a *Auth) protoAccount(ctx context.Context, c caller, s session) (*authv1.Account, error) {
	p, err := a.loadProfile(a.db.WithContext(ctx), c.host, s.user)
	if err != nil {
		return nil, err
	}
	u := p.user
	acc := &authv1.Account{
		Id:          u.ID.String(),
		Username:    u.Username,
		DisplayName: u.name(),
		CreatedAt:   timestamppb.New(u.Created),
		RecentUntil: timestamppb.New(s.created.Add(recent)),
		Capabilities: &authv1.Capabilities{
			EmailCodes:    c.host.mode.codes(),
			ChangeEmail:   a.mail != nil,
			Google:        a.enabled(c, authv1connect.AuthServiceFinishGoogleSigninProcedure),
			DeleteAccount: c.host.mode == Open,
		},
	}
	if u.Email != "" {
		acc.Email = &u.Email
	}
	for _, k := range p.passkeys {
		usable := compatible(k.Credential.V, a.rp)
		key := &authv1.Passkey{Id: encodeID(k.ID), Name: k.Name, CreatedAt: timestamppb.New(k.Created), Usable: usable, Removable: !usable || p.methods > 1}
		if k.Used != nil {
			key.UsedAt = timestamppb.New(*k.Used)
		}
		acc.Passkeys = append(acc.Passkeys, key)
	}
	if p.linked {
		acc.Identities = append(acc.Identities, &authv1.Identity{Provider: authv1.Provider_PROVIDER_GOOGLE, Removable: p.methods > 1})
	}
	return acc, nil
}

func (s *authService) current(ctx context.Context, c call) (*authv1.Account, error) {
	cur, err := live(s.a.db.WithContext(ctx), c.caller, c.session.user.ID)
	if err != nil {
		return nil, fault(ctx, err)
	}
	acc, err := s.a.protoAccount(ctx, c.caller, cur)
	return acc, fault(ctx, err)
}

func (s *authService) BeginPasskeySignin(ctx context.Context, _ *authv1.BeginPasskeySigninRequest) (*authv1.BeginPasskeySigninResponse, error) {
	c := callOf(ctx)
	if !c.host.begins.allow(c.ip) {
		return nil, fault(ctx, errThrottled)
	}
	options, t, err := s.a.beginLogin(ctx, c.caller)
	if err != nil {
		return nil, fault(ctx, err)
	}
	raw, err := json.Marshal(options.Response)
	if err != nil {
		return nil, fault(ctx, err)
	}
	return &authv1.BeginPasskeySigninResponse{Options: &authv1.PasskeyOptions{Ceremony: t, PublicKeyOptionsJson: raw}}, nil
}

func (s *authService) FinishPasskeySignin(ctx context.Context, req *authv1.FinishPasskeySigninRequest) (*authv1.FinishPasskeySigninResponse, error) {
	c := callOf(ctx)
	k, err := labeled(c.caller, req.GetLabel())
	if err != nil {
		return nil, err
	}
	cer, ok, err := takeCeremony(s.a.db.WithContext(ctx), k, req.GetCeremony(), login, s.a.library)
	if err != nil {
		return nil, fault(ctx, err)
	}
	if !ok {
		return nil, fault(ctx, errExpired)
	}
	parsed, err := protocol.ParseCredentialRequestResponseBytes(req.GetCredentialJson())
	if err != nil {
		return nil, fault(ctx, errUnreadable)
	}
	t, err := s.a.passkeyLogin(ctx, k, cer, parsed)
	if err != nil {
		return nil, fault(ctx, err)
	}
	g, err := s.grant(ctx, k, t)
	if err != nil {
		return nil, fault(ctx, err)
	}
	return &authv1.FinishPasskeySigninResponse{Grant: g}, nil
}

func (s *authService) BeginPasskeyAdd(ctx context.Context, _ *authv1.BeginPasskeyAddRequest) (*authv1.BeginPasskeyAddResponse, error) {
	c := callOf(ctx)
	if err := recently(c.session); err != nil {
		return nil, fault(ctx, err)
	}
	if !c.host.begins.allow(c.ip) {
		return nil, fault(ctx, errThrottled)
	}
	options, t, err := s.a.beginAdd(ctx, c.caller, c.session)
	if err != nil {
		return nil, fault(ctx, err)
	}
	raw, err := json.Marshal(options.Response)
	if err != nil {
		return nil, fault(ctx, err)
	}
	return &authv1.BeginPasskeyAddResponse{Options: &authv1.PasskeyOptions{Ceremony: t, PublicKeyOptionsJson: raw}}, nil
}

func (s *authService) FinishPasskeyAdd(ctx context.Context, req *authv1.FinishPasskeyAddRequest) (*authv1.FinishPasskeyAddResponse, error) {
	c := callOf(ctx)
	cer, ok, err := takeCeremony(s.a.db.WithContext(ctx), c.caller, req.GetCeremony(), register, s.a.library)
	if err != nil {
		return nil, fault(ctx, err)
	}
	if !ok || cer.UserID == nil || *cer.UserID != c.session.user.ID || cer.SessionID == nil || *cer.SessionID != c.session.id {
		return nil, fault(ctx, errExpired)
	}
	if err := recently(c.session); err != nil {
		return nil, fault(ctx, err)
	}
	parsed, err := protocol.ParseCredentialCreationResponseBytes(req.GetCredentialJson())
	if err != nil {
		return nil, fault(ctx, errUnverified)
	}
	if err := s.a.finishAdd(ctx, c.caller, c.session, cer, parsed); err != nil {
		return nil, fault(ctx, err)
	}
	acc, err := s.current(ctx, c)
	return &authv1.FinishPasskeyAddResponse{Account: acc}, err
}

func (s *authService) RenamePasskey(ctx context.Context, req *authv1.RenamePasskeyRequest) (*authv1.RenamePasskeyResponse, error) {
	c := callOf(ctx)
	id, ok := parsePasskeyID(req.GetPasskeyId())
	if !ok {
		return nil, fault(ctx, errMissing)
	}
	if err := s.a.renamePasskey(ctx, c.caller, c.session.user.ID, id, req.GetName()); err != nil {
		return nil, fault(ctx, err)
	}
	acc, err := s.current(ctx, c)
	return &authv1.RenamePasskeyResponse{Account: acc}, err
}

func (s *authService) RemovePasskey(ctx context.Context, req *authv1.RemovePasskeyRequest) (*authv1.RemovePasskeyResponse, error) {
	c := callOf(ctx)
	if err := recently(c.session); err != nil {
		return nil, fault(ctx, err)
	}
	id, ok := parsePasskeyID(req.GetPasskeyId())
	if !ok {
		return nil, fault(ctx, errMissing)
	}
	if err := s.a.removePasskey(ctx, c.caller, c.session.user.ID, id); err != nil {
		return nil, fault(ctx, err)
	}
	acc, err := s.current(ctx, c)
	return &authv1.RemovePasskeyResponse{Account: acc}, err
}

const (
	signinIntent = "signin"
	linkIntent   = "link"
)

type nativeFlow struct {
	ID        []byte
	Provider  string    `gorm:"not null"`
	Intent    string    `gorm:"not null;check:auth_native_flows_intent_check,intent IN ('signin', 'link')"`
	Nonce     string    `gorm:"not null"`
	ClientID  string    `gorm:"not null"`
	Origin    string    `gorm:"not null"`
	SessionID *pg.ID    `gorm:"check:auth_native_flows_check,(intent = 'signin' AND session_id IS NULL) OR (intent = 'link' AND session_id IS NOT NULL)"`
	Expires   time.Time `gorm:"not null"`
}

func (nativeFlow) TableName() string { return "auth_native_flows" }

func (a *Auth) beginFlow(ctx context.Context, c caller, intent string, from *pg.ID) (*authv1.GoogleFlow, error) {
	if !c.host.begins.allow(c.ip) {
		return nil, errThrottled
	}
	db := a.db.WithContext(ctx)
	if err := sweep(db, a.library); err != nil {
		return nil, err
	}
	t, nonce := token(), token()
	f := nativeFlow{ID: hash(t), Provider: googleKind, Intent: intent, Nonce: nonce, ClientID: c.app.id, Origin: c.host.origin, SessionID: from, Expires: time.Now().Add(flowLife)}
	if err := db.Create(&f).Error; err != nil {
		return nil, err
	}
	return &authv1.GoogleFlow{Flow: t, Nonce: nonce}, nil
}

func (a *Auth) takeNative(ctx context.Context, c caller, t, intent string) (nativeFlow, bool, error) {
	var f nativeFlow
	err := a.db.WithContext(ctx).Raw("DELETE FROM auth_native_flows WHERE id = ? AND provider = ? AND intent = ? AND client_id = ? AND origin = ? RETURNING *",
		hash(t), googleKind, intent, c.app.id, c.host.origin).Scan(&f).Error
	return f, err == nil && f.ID != nil && time.Now().Before(f.Expires), err
}

func (a *Auth) googleToken(ctx context.Context, p *app, raw, nonce string) (identity, error) {
	tok, err := a.tokens.Verifier(&oidc.Config{SkipClientIDCheck: true}).Verify(ctx, raw)
	if err != nil {
		return identity{}, err
	}
	var claims struct {
		Presenter string `json:"azp"`
		Email     string `json:"email"`
		Verified  bool   `json:"email_verified"`
	}
	if err := tok.Claims(&claims); err != nil {
		return identity{}, err
	}
	switch {
	case tok.Subject == "":
		return identity{}, errors.New("auth: id_token has no subject")
	case len(tok.Audience) != 1 || !slices.Contains(p.google, GoogleToken{Audience: tok.Audience[0], Presenter: claims.Presenter}):
		return identity{}, errors.New("auth: id_token audience is not this app's")
	case subtle.ConstantTimeCompare([]byte(tok.Nonce), []byte(nonce)) != 1:
		return identity{}, errors.New("auth: id_token nonce mismatch")
	}
	return identity{subject: tok.Subject, email: claims.Email, verified: claims.Verified}, nil
}

func (a *Auth) proveGoogle(ctx context.Context, c caller, flow, raw, intent string) (nativeFlow, identity, error) {
	f, ok, err := a.takeNative(ctx, c, flow, intent)
	if err != nil {
		return f, identity{}, err
	}
	if !ok {
		return f, identity{}, errExpired
	}
	id, err := a.googleToken(ctx, c.app, raw, f.Nonce)
	if err != nil {
		return f, id, connect.NewError(connect.CodeUnauthenticated, "That Google sign-in didn't check out. Try again.").WithCause(err)
	}
	if !time.Now().Before(f.Expires) {
		return f, id, errExpired
	}
	return f, id, nil
}

func (s *authService) BeginGoogleSignin(ctx context.Context, _ *authv1.BeginGoogleSigninRequest) (*authv1.BeginGoogleSigninResponse, error) {
	f, err := s.a.beginFlow(ctx, callOf(ctx).caller, signinIntent, nil)
	if err != nil {
		return nil, fault(ctx, err)
	}
	return &authv1.BeginGoogleSigninResponse{Flow: f}, nil
}

func (s *authService) FinishGoogleSignin(ctx context.Context, req *authv1.FinishGoogleSigninRequest) (*authv1.FinishGoogleSigninResponse, error) {
	k, err := labeled(callOf(ctx).caller, req.GetLabel())
	if err != nil {
		return nil, err
	}
	_, id, err := s.a.proveGoogle(ctx, k, req.GetFlow(), req.GetIdToken(), signinIntent)
	if err != nil {
		return nil, fault(ctx, err)
	}
	v, err := s.a.enterGoogle(ctx, k, id)
	if err != nil {
		return nil, fault(ctx, err)
	}
	switch v.verdict {
	case accepted:
	case denied:
		return nil, fault(ctx, errDenied)
	case unknown:
		return nil, connect.NewError(connect.CodeNotFound, "No account here uses that Google account. Sign in another way, then link Google.")
	case taken:
		return nil, connect.NewError(connect.CodeAlreadyExists, "An account here already uses that email. Sign in the way you usually do, then link Google.")
	default:
		return nil, connect.NewError(connect.CodeAborted, "Google sign-in changed while we checked it. Try again.")
	}
	g, err := s.grant(ctx, k, v.token)
	if err != nil {
		return nil, fault(ctx, err)
	}
	return &authv1.FinishGoogleSigninResponse{Grant: g}, nil
}

func (s *authService) BeginGoogleLink(ctx context.Context, _ *authv1.BeginGoogleLinkRequest) (*authv1.BeginGoogleLinkResponse, error) {
	c := callOf(ctx)
	if err := recently(c.session); err != nil {
		return nil, fault(ctx, err)
	}
	f, err := s.a.beginFlow(ctx, c.caller, linkIntent, &c.session.id)
	if err != nil {
		return nil, fault(ctx, err)
	}
	return &authv1.BeginGoogleLinkResponse{Flow: f}, nil
}

func (s *authService) FinishGoogleLink(ctx context.Context, req *authv1.FinishGoogleLinkRequest) (*authv1.FinishGoogleLinkResponse, error) {
	c := callOf(ctx)
	f, id, err := s.a.proveGoogle(ctx, c.caller, req.GetFlow(), req.GetIdToken(), linkIntent)
	if err != nil {
		return nil, fault(ctx, err)
	}
	v, _, err := s.a.linkGoogle(ctx, c.caller, c.session, f.SessionID, id)
	if err != nil {
		return nil, fault(ctx, err)
	}
	switch v {
	case accepted:
	case denied:
		return nil, fault(ctx, errDenied)
	case taken:
		return nil, connect.NewError(connect.CodeAlreadyExists, "That Google account is linked to another account here.")
	case linked:
		return nil, connect.NewError(connect.CodeAlreadyExists, googleLinked)
	default:
		return nil, fault(ctx, errRecent)
	}
	acc, err := s.current(ctx, c)
	return &authv1.FinishGoogleLinkResponse{Account: acc}, err
}

func (s *authService) UnlinkGoogle(ctx context.Context, _ *authv1.UnlinkGoogleRequest) (*authv1.UnlinkGoogleResponse, error) {
	c := callOf(ctx)
	if err := recently(c.session); err != nil {
		return nil, fault(ctx, err)
	}
	if err := s.a.unlinkGoogle(ctx, c.caller, c.session.user.ID); err != nil {
		return nil, fault(ctx, err)
	}
	acc, err := s.current(ctx, c)
	return &authv1.UnlinkGoogleResponse{Account: acc}, err
}

func (s *authService) SendCode(ctx context.Context, req *authv1.SendCodeRequest) (*authv1.SendCodeResponse, error) {
	c := callOf(ctx)
	email, ok := parseEmail(req.GetEmail())
	if !ok {
		return nil, fault(ctx, errBadEmail)
	}
	if !c.host.sends.allow(c.ip) {
		return nil, fault(ctx, errThrottled)
	}
	err := s.a.issueCode(ctx, c.host, email, func(pin, _ string) mail.Message {
		return s.a.appCodeMail(ctx, email, pin)
	})
	if err != nil {
		return nil, fault(ctx, err)
	}
	return &authv1.SendCodeResponse{}, nil
}

func (s *authService) VerifyCode(ctx context.Context, req *authv1.VerifyCodeRequest) (*authv1.VerifyCodeResponse, error) {
	k, err := labeled(callOf(ctx).caller, req.GetLabel())
	if err != nil {
		return nil, err
	}
	email, ok := parseEmail(req.GetEmail())
	if !ok {
		return nil, fault(ctx, errBadEmail)
	}
	guess, ok := parseCode(req.GetCode())
	if !ok {
		return nil, fault(ctx, errBadCode)
	}
	if !k.host.verifies.allow(k.ip) {
		return nil, fault(ctx, errThrottled)
	}
	v, err := s.a.checkCode(ctx, k, email, guess)
	if err != nil {
		return nil, fault(ctx, err)
	}
	if err := v.codeErr(); err != nil {
		return nil, fault(ctx, err)
	}
	g, err := s.grant(ctx, k, v.token)
	if err != nil {
		return nil, fault(ctx, err)
	}
	return &authv1.VerifyCodeResponse{Grant: g}, nil
}

func (s *authService) GetAccount(ctx context.Context, _ *authv1.GetAccountRequest) (*authv1.GetAccountResponse, error) {
	acc, err := s.current(ctx, callOf(ctx))
	return &authv1.GetAccountResponse{Account: acc}, err
}

func (s *authService) SetUsername(ctx context.Context, req *authv1.SetUsernameRequest) (*authv1.SetUsernameResponse, error) {
	c := callOf(ctx)
	if err := s.a.saveUsername(ctx, c.caller, c.session, req.GetUsername()); err != nil {
		return nil, fault(ctx, err)
	}
	acc, err := s.current(ctx, c)
	return &authv1.SetUsernameResponse{Account: acc}, err
}

func (s *authService) ChangeEmail(ctx context.Context, req *authv1.ChangeEmailRequest) (*authv1.ChangeEmailResponse, error) {
	c := callOf(ctx)
	if err := recently(c.session); err != nil {
		return nil, fault(ctx, err)
	}
	to, err := newEmail(req.GetEmail(), c.session.user.Email)
	if err != nil {
		return nil, fault(ctx, err)
	}
	if !c.host.sends.allow(c.ip) {
		return nil, fault(ctx, errThrottled)
	}
	if err := s.a.issueVerify(ctx, c.caller, c.session.user.ID, to); err != nil {
		return nil, fault(ctx, err)
	}
	return &authv1.ChangeEmailResponse{}, nil
}

func (s *authService) ListSessions(ctx context.Context, _ *authv1.ListSessionsRequest) (*authv1.ListSessionsResponse, error) {
	c := callOf(ctx)
	rows, err := s.a.sessions(s.a.db.WithContext(ctx), c.host, c.session.user.ID)
	if err != nil {
		return nil, fault(ctx, err)
	}
	res := &authv1.ListSessionsResponse{}
	for _, r := range rows {
		item := &authv1.Session{
			Id:        r.ID.String(),
			Kind:      authv1.ClientKind_CLIENT_KIND_BROWSER,
			CreatedAt: timestamppb.New(r.Created),
			SeenAt:    timestamppb.New(r.Seen),
			ExpiresAt: timestamppb.New(r.Expires),
			Current:   r.ID == c.session.id,
		}
		if r.Label != nil {
			item.Label = *r.Label
		}
		if r.ClientID != nil {
			item.Kind = s.a.apps[*r.ClientID].kind
		}
		res.Sessions = append(res.Sessions, item)
	}
	return res, nil
}

func (s *authService) RevokeSession(ctx context.Context, req *authv1.RevokeSessionRequest) (*authv1.RevokeSessionResponse, error) {
	c := callOf(ctx)
	id, err := pg.ParseID(req.GetSessionId())
	if err != nil {
		return nil, fault(ctx, errMissing)
	}
	if err := s.a.revoke(ctx, c.caller, c.session.user.ID, &id); err != nil {
		return nil, fault(ctx, err)
	}
	return &authv1.RevokeSessionResponse{}, nil
}

func (s *authService) RevokeOtherSessions(ctx context.Context, _ *authv1.RevokeOtherSessionsRequest) (*authv1.RevokeOtherSessionsResponse, error) {
	c := callOf(ctx)
	if err := s.a.revoke(ctx, c.caller, c.session.user.ID, nil); err != nil {
		return nil, fault(ctx, err)
	}
	return &authv1.RevokeOtherSessionsResponse{}, nil
}

func (s *authService) SignOut(ctx context.Context, _ *authv1.SignOutRequest) (*authv1.SignOutResponse, error) {
	if err := end(s.a.db.WithContext(ctx), callOf(ctx).caller); err != nil {
		return nil, fault(ctx, err)
	}
	return &authv1.SignOutResponse{}, nil
}

func (s *authService) DeleteAccount(ctx context.Context, req *authv1.DeleteAccountRequest) (*authv1.DeleteAccountResponse, error) {
	c := callOf(ctx)
	if !req.GetConfirm() {
		return nil, fault(ctx, errUnconfirmed)
	}
	if err := recently(c.session); err != nil {
		return nil, fault(ctx, err)
	}
	if err := s.a.erase(ctx, c.caller, c.session); err != nil {
		return nil, fault(ctx, err)
	}
	return &authv1.DeleteAccountResponse{}, nil
}
