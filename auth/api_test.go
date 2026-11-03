package auth

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"connectrpc.com/connect/v2"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/coreos/go-oidc/v3/oidc/oidctest"
	"github.com/go-webauthn/webauthn/protocol/webauthncbor"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/cogikyo/hgmx/mail"
	"github.com/cogikyo/hgmx/migrate"
	"github.com/cogikyo/hgmx/pg/pgtest"
	authv1 "github.com/cogikyo/hgmx/proto/hgmx/auth/v1"
	"github.com/cogikyo/hgmx/server"
	"github.com/cogikyo/hgmx/telemetry"
)

const (
	iosID         = "ios"
	androidID     = "android"
	adminID       = "admin-ios"
	webClient     = "web.apps.googleusercontent.com"
	iosClient     = "ios.apps.googleusercontent.com"
	androidClient = "android.apps.googleusercontent.com"
	apexHost      = "example.com"
	adminHost     = "admin.example.com"
	svc           = "/hgmx.auth.v1.AuthService/"
)

var (
	release   = sha256.Sum256([]byte("release signing cert"))
	debugCert = sha256.Sum256([]byte("debugCert signing cert"))
	apkOrigin = "android:apk-key-hash:" + base64.RawURLEncoding.EncodeToString(release[:])
)

type outbox struct {
	mu   sync.Mutex
	msgs []mail.Message
}

func (o *outbox) Send(_ context.Context, m mail.Message) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.msgs = append(o.msgs, m)
	return nil
}

type rig struct {
	t      *testing.T
	a      *Auth
	srv    *httptest.Server
	google *fakeGoogle
}

func appConfig() Config {
	c := config()
	c.Hosts[0].Mode = Open
	c.Hosts[0].Apps = []App{
		{ID: iosID, IOS: &IOS{Prefix: "ABCDE12345", Bundle: "com.example.app"}, Build: 10, Store: "https://apps.apple.com/app/id123", Google: []GoogleToken{{webClient, iosClient}}},
		{ID: androidID, Android: &Android{Package: "com.example.app", Signers: [][]byte{release[:]}}, Build: 10, Store: "https://play.google.com/store/apps/details?id=com.example.app", Google: []GoogleToken{{webClient, androidClient}}},
	}
	c.Hosts[1].Apps = []App{{ID: adminID, IOS: &IOS{Prefix: "ABCDE12345", Bundle: "com.example.admin"}, Build: 1, Store: "https://apps.apple.com/app/id124"}}
	return c
}

func newRig(t *testing.T, interceptors ...connect.ServerInterceptor) *rig {
	t.Setenv("PGDATABASE", "pgtest")
	db := pgtest.Open(t)
	if err := migrate.Command(t.Context(), db, []migrate.Track{Track}, []string{"up"}, io.Discard); err != nil {
		t.Fatal(err)
	}
	q := mail.NewQueue(&outbox{}, 64, slog.New(slog.DiscardHandler))
	t.Cleanup(func() { q.Close(context.Background()) })
	a, err := New(db, q, appConfig())
	if err != nil {
		t.Fatal(err)
	}
	fake := newFakeGoogle()
	gs := httptest.NewServer(fake)
	t.Cleanup(gs.Close)
	fake.start(gs.URL)
	a.tokens = fake.provider().keys

	api := a.API(interceptors...)
	mux := http.NewServeMux()
	a.Mount(mux)
	if err := api.Mount(mux); err != nil {
		t.Fatal(err)
	}
	mux.HandleFunc("POST "+svc+"Boom", func(http.ResponseWriter, *http.Request) { panic("boom before dispatch") })
	mux.HandleFunc("GET /page", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "<!doctype html><p>page") })
	mux.HandleFunc("GET /panic", func(http.ResponseWriter, *http.Request) { panic("page panic") })
	site := server.Site{Telemetry: telemetry.Disabled(slog.DiscardHandler), RPC: api.Server()}
	srv := httptest.NewServer(site.Handler(mux))
	t.Cleanup(srv.Close)
	return &rig{t: t, a: a, srv: srv, google: fake}
}

type reply struct {
	status  int
	header  http.Header
	body    []byte
	code    string
	details []string
}

func (r reply) detail(t *testing.T, kind string, into proto.Message) {
	t.Helper()
	var e struct {
		Details []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"details"`
	}
	json.Unmarshal(r.body, &e)
	for _, d := range e.Details {
		if d.Type == kind {
			raw, err := base64.RawStdEncoding.DecodeString(strings.TrimRight(d.Value, "="))
			if err != nil {
				t.Fatal(err)
			}
			if err := proto.Unmarshal(raw, into); err != nil {
				t.Fatal(err)
			}
			return
		}
	}
	t.Fatalf("no %s detail in %s", kind, r.body)
}

type hdr map[string]string

func meta(id string, build int) hdr {
	return hdr{appHeader: id, buildHeader: itoa(build)}
}

func itoa(n int) string { return strconv.Itoa(n) }

func (h hdr) with(k, v string) hdr {
	out := hdr{}
	for a, b := range h {
		out[a] = b
	}
	out[k] = v
	return out
}

func (h hdr) bearer(t string) hdr { return h.with("Authorization", "Bearer "+t) }

func (g *rig) raw(host, method, path, contentType string, body []byte, h hdr) reply {
	g.t.Helper()
	req, err := http.NewRequest(method, g.srv.URL+path, bytes.NewReader(body))
	if err != nil {
		g.t.Fatal(err)
	}
	req.Host = host
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, v := range h {
		if k == "Cookie" {
			req.Header.Add("Cookie", v)
			continue
		}
		req.Header.Set(k, v)
	}
	res, err := (&http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}).Do(req)
	if err != nil {
		g.t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	out := reply{status: res.StatusCode, header: res.Header.Clone(), body: b}
	for k, v := range res.Trailer {
		out.header[k] = v
	}
	if res.StatusCode != 200 {
		var e struct {
			Code string `json:"code"`
		}
		json.Unmarshal(b, &e)
		out.code = e.Code
	}
	return out
}

func (g *rig) rpc(host, method string, h hdr, req, out proto.Message) reply {
	g.t.Helper()
	body, err := protojson.Marshal(req)
	if err != nil {
		g.t.Fatal(err)
	}
	h = h.with("Connect-Protocol-Version", "1")
	r := g.raw(host, "POST", svcPath(method), "application/json", body, h)
	if ct := r.header.Get("Content-Type"); strings.Contains(ct, "html") {
		g.t.Fatalf("%s returned HTML: %s", method, r.body)
	}
	if r.status == 200 && out != nil {
		if err := protojson.Unmarshal(r.body, out); err != nil {
			g.t.Fatalf("%s: %v: %s", method, err, r.body)
		}
	}
	return r
}

func svcPath(method string) string {
	if strings.HasPrefix(method, "/") {
		return method
	}
	return svc + method
}

func (g *rig) want(r reply, code string) {
	g.t.Helper()
	if code == "" && r.status != 200 || code != "" && r.code != code {
		g.t.Fatalf("want %q, got %d %q: %s", code, r.status, r.code, r.body)
	}
	if code != "" && !strings.HasPrefix(r.header.Get("Content-Type"), "application/json") {
		g.t.Fatalf("fault content type %q", r.header.Get("Content-Type"))
	}
}

func (g *rig) signIn(host string, h hdr, email, label string) *authv1.SessionGrant {
	g.t.Helper()
	addCode(g.t, g.a, email, "123456")
	var res authv1.VerifyCodeResponse
	g.want(g.rpc(host, "VerifyCode", h, &authv1.VerifyCodeRequest{Email: email, Code: "123456", Label: label}, &res), "")
	return res.Grant
}

func (g *rig) browser(email string) string {
	g.t.Helper()
	u, err := findUser(g.a.db, email)
	if err != nil {
		g.t.Fatal(err)
	}
	tok, err := startSession(g.a.db, caller{host: at(g.a, apex)}, u.ID)
	if err != nil {
		g.t.Fatal(err)
	}
	return tok
}

func (g *rig) page(host, path, cookie string) reply {
	h := hdr{}
	if cookie != "" {
		h["Cookie"] = sessionCookie + "=" + cookie
	}
	return g.raw(host, "GET", path, "", nil, h)
}

func loggedIn(r reply) bool { return r.status == 200 }

type softKey struct {
	key    *ecdsa.PrivateKey
	id     []byte
	handle []byte
	rp     string
	count  uint32
}

func newKey(rp string) *softKey {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	id := make([]byte, 16)
	rand.Read(id)
	return &softKey{key: key, id: id, rp: rp}
}

func b64(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

func challenge(t *testing.T, options []byte) (string, string) {
	var o struct {
		PublicKey *struct {
			Challenge string `json:"challenge"`
			User      struct {
				ID string `json:"id"`
			} `json:"user"`
		} `json:"publicKey"`
		Challenge string `json:"challenge"`
		User      struct {
			ID string `json:"id"`
		} `json:"user"`
	}
	if err := json.Unmarshal(options, &o); err != nil {
		t.Fatal(err)
	}
	if o.PublicKey != nil {
		return o.PublicKey.Challenge, o.PublicKey.User.ID
	}
	return o.Challenge, o.User.ID
}

func (k *softKey) authData(flags byte, attested bool) []byte {
	rp := sha256.Sum256([]byte(k.rp))
	b := append(rp[:], flags)
	b = binary.BigEndian.AppendUint32(b, k.count)
	if attested {
		pub, _ := k.key.PublicKey.ECDH()
		raw := pub.Bytes()
		cose, _ := webauthncbor.Marshal(map[int]any{1: 2, 3: -7, -1: 1, -2: raw[1:33], -3: raw[33:]})
		b = append(b, make([]byte, 16)...)
		b = binary.BigEndian.AppendUint16(b, uint16(len(k.id)))
		b = append(b, k.id...)
		b = append(b, cose...)
	}
	return b
}

func clientData(typ, chal, origin string) []byte {
	b, _ := json.Marshal(map[string]any{"type": typ, "challenge": chal, "origin": origin, "crossOrigin": false})
	return b
}

func (k *softKey) create(t *testing.T, options []byte, origin string) []byte {
	chal, user := challenge(t, options)
	k.handle, _ = base64.RawURLEncoding.DecodeString(user)
	att, _ := webauthncbor.Marshal(map[string]any{"fmt": "none", "attStmt": map[string]any{}, "authData": k.authData(0x45, true)})
	b, _ := json.Marshal(map[string]any{"id": b64(k.id), "rawId": b64(k.id), "type": "public-key", "response": map[string]any{
		"clientDataJSON": b64(clientData("webauthn.create", chal, origin)), "attestationObject": b64(att)}})
	return b
}

func (k *softKey) get(t *testing.T, options []byte, origin string) []byte {
	chal, _ := challenge(t, options)
	k.count++
	cd := clientData("webauthn.get", chal, origin)
	ad := k.authData(0x05, false)
	sum := sha256.Sum256(cd)
	digest := sha256.Sum256(append(append([]byte{}, ad...), sum[:]...))
	sig, _ := ecdsa.SignASN1(rand.Reader, k.key, digest[:])
	b, _ := json.Marshal(map[string]any{"id": b64(k.id), "rawId": b64(k.id), "type": "public-key", "response": map[string]any{
		"clientDataJSON": b64(cd), "authenticatorData": b64(ad), "signature": b64(sig), "userHandle": b64(k.handle)}})
	return b
}

func (g *rig) nativeAdd(h hdr, tok string, k *softKey, origin string) reply {
	var begin authv1.BeginPasskeyAddResponse
	g.want(g.rpc(apexHost, "BeginPasskeyAdd", h.bearer(tok), &authv1.BeginPasskeyAddRequest{}, &begin), "")
	return g.rpc(apexHost, "FinishPasskeyAdd", h.bearer(tok), &authv1.FinishPasskeyAddRequest{Ceremony: begin.Options.Ceremony, CredentialJson: k.create(g.t, begin.Options.PublicKeyOptionsJson, origin)}, nil)
}

func (g *rig) nativeLogin(begin, finish hdr, host string, k *softKey, origin string) reply {
	var b authv1.BeginPasskeySigninResponse
	g.want(g.rpc(host, "BeginPasskeySignin", begin, &authv1.BeginPasskeySigninRequest{}, &b), "")
	return g.rpc(host, "FinishPasskeySignin", finish, &authv1.FinishPasskeySigninRequest{Ceremony: b.Options.Ceremony, CredentialJson: k.get(g.t, b.Options.PublicKeyOptionsJson, origin)}, nil)
}

func (g *rig) browserLogin(k *softKey, origin string) reply {
	r := g.raw(apexHost, "POST", "/auth/passkey/begin", "", nil, nil)
	if r.status != 200 {
		g.t.Fatalf("browser begin %d", r.status)
	}
	var ceremony string
	for _, c := range (&http.Response{Header: r.header}).Cookies() {
		if c.Name == ceremonyCookie {
			ceremony = c.Value
		}
	}
	return g.raw(apexHost, "POST", "/auth/passkey/finish", "application/json", k.get(g.t, r.body, origin), hdr{"Cookie": ceremonyCookie + "=" + ceremony, "Origin": apex})
}

func (g *rig) browserAdd(cookie string, k *softKey, origin string) reply {
	h := hdr{"Cookie": sessionCookie + "=" + cookie, "Origin": apex}
	r := g.raw(apexHost, "POST", "/auth/account/passkeys/begin", "", nil, h)
	if r.status != 200 {
		g.t.Fatalf("browser add begin %d %s", r.status, r.body)
	}
	var ceremony string
	for _, c := range (&http.Response{Header: r.header}).Cookies() {
		if c.Name == ceremonyCookie {
			ceremony = c.Value
		}
	}
	h["Cookie"] += "; " + ceremonyCookie + "=" + ceremony
	return g.raw(apexHost, "POST", "/auth/account/passkeys/finish", "application/json", k.create(g.t, r.body, origin), h)
}

func (g *rig) idToken(aud, azp, nonce string, exp time.Time, iss string) string {
	if iss == "" {
		iss = g.google.issuer
	}
	claims := map[string]any{"iss": iss, "aud": aud, "azp": azp, "sub": "google-sub", "email": "g@example.com", "email_verified": true, "iat": time.Now().Unix(), "exp": exp.Unix()}
	if nonce != "" {
		claims["nonce"] = nonce
	}
	b, _ := json.Marshal(claims)
	return oidctest.SignIDToken(g.google.key, "k", oidc.RS256, string(b))
}

func (g *rig) googleSignin(h hdr, token func(nonce string) string) reply {
	var b authv1.BeginGoogleSigninResponse
	g.want(g.rpc(apexHost, "BeginGoogleSignin", h, &authv1.BeginGoogleSigninRequest{}, &b), "")
	return g.rpc(apexHost, "FinishGoogleSignin", h, &authv1.FinishGoogleSigninRequest{Flow: b.Flow.Flow, IdToken: token(b.Flow.Nonce)}, nil)
}

func TestNativeSessionClass(t *testing.T) {
	g := newRig(t)
	ios, android := meta(iosID, 10), meta(androidID, 10)
	email := "class@example.com"
	tok := g.signIn(apexHost, ios, email, "").Token
	cookie := g.browser(email)
	g.want(g.rpc(apexHost, "GetAccount", ios.bearer(tok), &authv1.GetAccountRequest{}, nil), "")
	g.want(g.rpc(apexHost, "GetAccount", android.bearer(tok), &authv1.GetAccountRequest{}, nil), "unauthenticated")
	g.want(g.rpc(apexHost, "GetAccount", ios.bearer(cookie), &authv1.GetAccountRequest{}, nil), "unauthenticated")
	g.want(g.rpc(apexHost, "GetAccount", ios.with("Cookie", sessionCookie+"="+cookie), &authv1.GetAccountRequest{}, nil), "unauthenticated")
	if loggedIn(g.page(apexHost, "/auth/account", tok)) {
		t.Fatal("native bearer worked as a browser cookie")
	}
	if !loggedIn(g.page(apexHost, "/auth/account", cookie)) {
		t.Fatal("browser cookie stopped working")
	}
}

func TestNativePasskeyOrigin(t *testing.T) {
	g := newRig(t)
	ios, android := meta(iosID, 10), meta(androidID, 10)
	email := "origin@example.com"
	it := g.signIn(apexHost, ios, email, "").Token
	atok := g.signIn(apexHost, android, email, "").Token
	g.want(g.nativeAdd(android, atok, newKey("example.com"), apex), "invalid_argument")
	g.want(g.nativeAdd(ios, it, newKey("example.com"), apkOrigin), "invalid_argument")
	if r := g.browserAdd(g.browser(email), newKey("example.com"), apkOrigin); r.status != 400 {
		t.Fatalf("browser registration accepted an Android origin: %d", r.status)
	}
	k := newKey("example.com")
	g.want(g.nativeAdd(android, atok, k, apkOrigin), "")
	g.want(g.nativeLogin(android, android, apexHost, k, apkOrigin), "")
	g.want(g.nativeLogin(android, android, apexHost, k, apex), "unauthenticated")
	if r := g.browserLogin(k, apkOrigin); r.status == 204 {
		t.Fatal("browser login accepted an Android origin")
	}
}

func TestNativeGoogleAudience(t *testing.T) {
	g := newRig(t)
	ios := meta(iosID, 10)
	hour := time.Now().Add(time.Hour)
	g.want(g.googleSignin(ios, func(n string) string { return g.idToken(webClient, androidClient, n, hour, "") }), "unauthenticated")
	g.want(g.googleSignin(ios, func(n string) string { return g.idToken(webClient, iosClient, "other", hour, "") }), "unauthenticated")
	g.want(g.googleSignin(ios, func(n string) string { return g.idToken(webClient, iosClient, n, hour, "") }), "")
}

func TestRPCFaultsStayRPC(t *testing.T) {
	g := newRig(t)
	g.want(g.rpc(apexHost, "Boom", meta(iosID, 10), &authv1.GetAccountRequest{}, nil), "internal")
	g.want(g.rpc(apexHost, "/x.v1.Nope/Do", meta(iosID, 10), &authv1.GetAccountRequest{}, nil), "unimplemented")
	g.want(g.rpc("evil.test", "GetAccount", meta(iosID, 10), &authv1.GetAccountRequest{}, nil), "not_found")
}
