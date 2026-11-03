package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"connectrpc.com/connect/v2/connectproto"
	"google.golang.org/protobuf/proto"

	"github.com/cogikyo/hgmx/clientip"
	"github.com/cogikyo/hgmx/pg"
	appv1 "github.com/cogikyo/hgmx/proto/hgmx/app/v1"
	"github.com/cogikyo/hgmx/proto/hgmx/app/v1/appv1connect"
	authv1 "github.com/cogikyo/hgmx/proto/hgmx/auth/v1"
	"github.com/cogikyo/hgmx/proto/hgmx/auth/v1/authv1connect"
)

type Access int

const (
	Public        Access = iota + 1 // No sign-in or minimum-build gate.
	Anonymous                       // Minimum build, without resolving a user.
	Authenticated                   // Minimum build, live bearer and Host.Allow.
	exit
)

const (
	appHeader   = "HGMX-App-ID"    // Required public selector, not attestation.
	buildHeader = "HGMX-App-Build" // Required integer build, not marketing version.
	readLimit   = 1 << 20
)

// API applies host, app and access checks when mounted with [API.Mount].
type API struct {
	auth   *Auth
	server *connect.Server
	access map[string]Access
	errors *connecthttp.ErrorWriter
}

func (a *Auth) API(interceptors ...connect.ServerInterceptor) *API {
	p := &API{
		auth:   a,
		server: connect.NewServer(append([]connect.ServerInterceptor{internal}, interceptors...)...),
		access: map[string]Access{},
		errors: connecthttp.NewErrorWriter(connecthttp.WithRequireConnectProtocolHeader()),
	}
	authv1connect.RegisterAuthServiceHandler(p.server, &authService{a})
	appv1connect.RegisterAppServiceHandler(p.server, &appService{p})
	p.Allow(Public, appv1connect.AppServiceGetConfigurationProcedure)
	p.Allow(Anonymous,
		authv1connect.AuthServiceBeginPasskeySigninProcedure,
		authv1connect.AuthServiceFinishPasskeySigninProcedure,
		authv1connect.AuthServiceBeginGoogleSigninProcedure,
		authv1connect.AuthServiceFinishGoogleSigninProcedure,
		authv1connect.AuthServiceSendCodeProcedure,
		authv1connect.AuthServiceVerifyCodeProcedure,
	)
	p.Allow(Authenticated,
		authv1connect.AuthServiceBeginPasskeyAddProcedure,
		authv1connect.AuthServiceFinishPasskeyAddProcedure,
		authv1connect.AuthServiceRenamePasskeyProcedure,
		authv1connect.AuthServiceRemovePasskeyProcedure,
		authv1connect.AuthServiceBeginGoogleLinkProcedure,
		authv1connect.AuthServiceFinishGoogleLinkProcedure,
		authv1connect.AuthServiceUnlinkGoogleProcedure,
		authv1connect.AuthServiceGetAccountProcedure,
		authv1connect.AuthServiceSetUsernameProcedure,
		authv1connect.AuthServiceChangeEmailProcedure,
		authv1connect.AuthServiceListSessionsProcedure,
		authv1connect.AuthServiceRevokeSessionProcedure,
		authv1connect.AuthServiceRevokeOtherSessionsProcedure,
		authv1connect.AuthServiceDeleteAccountProcedure,
	)
	p.access[authv1connect.AuthServiceSignOutProcedure] = exit
	return p
}

func (p *API) Server() *connect.Server { return p.server }

func (p *API) Allow(access Access, procedures ...string) {
	if access < Public || access > Authenticated {
		panic("auth: Allow needs Public, Anonymous or Authenticated")
	}
	for _, proc := range procedures {
		if _, ok := p.access[proc]; ok {
			panic("auth: " + proc + " already has an Access")
		}
		p.access[proc] = access
	}
}

// Mount rejects missing or unmatched access policies before adding routes.
func (p *API) Mount(mux *http.ServeMux) error {
	registered := map[string]bool{}
	for spec := range p.server.Specs() {
		if _, ok := p.access[spec.Procedure]; !ok {
			return fmt.Errorf("auth: RPC %s has no Access; call API.Allow", spec.Procedure)
		}
		registered[spec.Procedure] = true
	}
	for proc := range p.access {
		if !registered[proc] {
			return fmt.Errorf("auth: RPC %s has an Access but no handler", proc)
		}
	}
	connecthttp.Mount(guarded{mux, p}, p.server, connecthttp.WithRequireConnectProtocolHeader(), connecthttp.WithReadMaxBytes(readLimit))
	return nil
}

type guarded struct {
	mux *http.ServeMux
	api *API
}

func (g guarded) Handle(pattern string, h http.Handler) {
	g.mux.Handle(pattern, g.api.edge(pattern, h))
}

type callKey struct{}

type call struct {
	caller
	session session
	build   int
}

func (p *API) edge(procedure string, next http.Handler) http.Handler {
	access := p.access[procedure]
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		h, ok := p.auth.host(r)
		if !ok {
			p.errors.Write(w, r, connect.NewError(connect.CodeNotFound, "unknown host"))
			return
		}
		if access == 0 {
			next.ServeHTTP(w, r)
			return
		}
		c, err := p.auth.call(r, h, procedure, access)
		if err != nil {
			p.errors.Write(w, r, fault(r.Context(), err))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), callKey{}, c)))
	})
}

func (a *Auth) call(r *http.Request, h *host, procedure string, access Access) (call, error) {
	id, ok := single(r.Header, appHeader)
	if !ok {
		return call{}, connect.NewError(connect.CodeInvalidArgument, "send one "+appHeader+" header")
	}
	app := a.apps[id]
	if app == nil || app.host != h {
		return call{}, connect.NewError(connect.CodeInvalidArgument, "this app is not registered here")
	}
	if !a.enabled(caller{host: h, app: app}, procedure) {
		return call{}, connect.NewError(connect.CodeUnimplemented, "this site doesn't offer that")
	}
	raw, ok := single(r.Header, buildHeader)
	build, err := strconv.Atoi(raw)
	if !ok || err != nil || build < 1 || strings.TrimLeft(raw, "0123456789") != "" {
		return call{}, connect.NewError(connect.CodeInvalidArgument, "send one numeric "+buildHeader+" header")
	}
	if access != Public && access != exit && build < app.build {
		return call{}, outdated(app)
	}
	token, err := bearer(r.Header)
	if err != nil {
		return call{}, err
	}
	if _, err := r.Cookie(sessionCookie); err == nil && token != "" {
		return call{}, connect.NewError(connect.CodeInvalidArgument, "send a bearer token or a session cookie, not both")
	}
	c := call{caller: caller{host: h, app: app, token: token, ip: clientip.From(r)}, build: build}
	if access == Public || access == Anonymous {
		return c, nil
	}
	if token == "" {
		return call{}, connect.NewError(connect.CodeUnauthenticated, "sign in to continue")
	}
	s, ok, err := a.find(r.Context(), c.caller)
	switch {
	case err != nil:
		return call{}, err
	case !ok:
		return call{}, errStale
	case access == Authenticated && !h.admits(s.user):
		return call{}, errDenied
	}
	c.session = s
	return c, nil
}

func single(h http.Header, key string) (string, bool) {
	v := h.Values(key)
	return strings.Join(v, ""), len(v) == 1 && v[0] != ""
}

func bearer(h http.Header) (string, error) {
	v := h.Values("Authorization")
	if len(v) == 0 {
		return "", nil
	}
	scheme, token, _ := strings.Cut(v[0], " ")
	if len(v) > 1 || !strings.EqualFold(scheme, "Bearer") || !wellFormed(token) {
		return "", connect.NewError(connect.CodeUnauthenticated, "send one Authorization: Bearer token")
	}
	return token, nil
}

func wellFormed(token string) bool {
	return len(token) == 43 && strings.Trim(token, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_") == ""
}

func outdated(p *app) error {
	return detailed(connect.NewError(connect.CodeFailedPrecondition, "Update the app to keep going."),
		&appv1.UpdateRequired{MinimumBuild: int64(p.build), StoreUrl: p.store, Message: "Update the app to keep going."})
}

func detailed(err *connect.Error, detail proto.Message) *connect.Error {
	d, derr := connectproto.NewErrorDetail(detail)
	if derr != nil {
		return err
	}
	return err.WithDetail(d)
}

func callOf(ctx context.Context) call {
	return ctx.Value(callKey{}).(call)
}

// Caller returns the user resolved for an authenticated RPC, if any.
// Public and Anonymous calls do not resolve a user, even with a bearer.
func Caller(ctx context.Context) (User, bool) {
	c, ok := ctx.Value(callKey{}).(call)
	return c.session.user, ok && c.session.id != pg.ID{}
}

func fault(ctx context.Context, err error) error {
	var ce *connect.Error
	switch {
	case err == nil:
		return nil
	case errors.As(err, &ce):
		return ce
	case errors.Is(err, errStale), errors.Is(err, pg.ErrNotFound):
		return connect.NewError(connect.CodeUnauthenticated, "Your session ended. Sign in again.")
	case errors.Is(err, errRecent):
		return detailed(connect.NewError(connect.CodeFailedPrecondition, "Sign in again to make this change."),
			&authv1.SignInAgain{Message: "Sign in again to make this change."})
	case errors.Is(err, errLast):
		return detailed(connect.NewError(connect.CodeFailedPrecondition, "This is your only way to sign in. Add another one first."),
			&authv1.LastSignInMethod{Message: "This is your only way to sign in. Add another one first."})
	case errors.Is(err, errDenied):
		return connect.NewError(connect.CodePermissionDenied, "This account can't sign in here.")
	case errors.Is(err, errMissing):
		return connect.NewError(connect.CodeNotFound, "That item no longer exists.")
	case errors.Is(err, pg.ErrDuplicate):
		return connect.NewError(connect.CodeAlreadyExists, "That already exists.")
	}
	slog.ErrorContext(ctx, "auth: rpc failed", "err", err)
	return connect.NewError(connect.CodeInternal, "Something went wrong. Try again.").WithCause(err)
}

func internal(next connect.ServerFunc) connect.ServerFunc {
	return func(ctx context.Context, spec connect.Spec, stream connect.ServerStream) error {
		err := next(ctx, spec, stream)
		var ce *connect.Error
		if err == nil || errors.As(err, &ce) {
			return err
		}
		slog.ErrorContext(ctx, "rpc failed", "procedure", spec.Procedure, "err", err)
		return connect.NewError(connect.CodeInternal, "Something went wrong. Try again.").WithCause(err)
	}
}
