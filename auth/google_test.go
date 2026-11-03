package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/coreos/go-oidc/v3/oidc/oidctest"

	"github.com/cogikyo/hgmx/pg"
)

type fakeClaims struct {
	Sub      string `json:"sub"`
	Email    string `json:"email,omitempty"`
	Verified bool   `json:"email_verified"`
}

type grant struct {
	claims    fakeClaims
	nonce     string
	challenge string
}

type fakeGoogle struct {
	mu     sync.Mutex
	key    *rsa.PrivateKey
	keys   oidctest.Server
	issuer string
	next   fakeClaims
	deny   bool
	codes  map[string]grant
}

func newFakeGoogle() *fakeGoogle {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	f := &fakeGoogle{key: key, codes: map[string]grant{}}
	f.keys.PublicKeys = []oidctest.PublicKey{{PublicKey: key.Public(), KeyID: "k", Algorithm: oidc.RS256}}
	return f
}

func (f *fakeGoogle) start(issuer string) {
	f.issuer = issuer
	f.keys.SetIssuer(issuer)
}

func (f *fakeGoogle) as(c fakeClaims) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.next, f.deny = c, false
}

func (f *fakeGoogle) provider() *provider {
	return newProvider(googleKind, oidc.ProviderConfig{
		IssuerURL:  f.issuer,
		AuthURL:    f.issuer + "/auth",
		TokenURL:   f.issuer + "/token",
		JWKSURL:    f.issuer + "/keys",
		Algorithms: []string{oidc.RS256},
	}, "client", "secret")
}

func (f *fakeGoogle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	switch r.URL.Path {
	case "/auth":
		q := r.URL.Query()
		back, _ := url.Parse(q.Get("redirect_uri"))
		v := url.Values{"state": {q.Get("state")}}
		if f.deny {
			v.Set("error", "access_denied")
		} else {
			code := rand.Text()
			f.codes[code] = grant{claims: f.next, nonce: q.Get("nonce"), challenge: q.Get("code_challenge")}
			v.Set("code", code)
		}
		back.RawQuery = v.Encode()
		http.Redirect(w, r, back.String(), http.StatusFound)
	case "/token":
		id, secret, ok := r.BasicAuth()
		if !ok {
			id, secret = r.FormValue("client_id"), r.FormValue("client_secret")
		}
		g, found := f.codes[r.FormValue("code")]
		delete(f.codes, r.FormValue("code"))
		sum := sha256.Sum256([]byte(r.FormValue("code_verifier")))
		if id != "client" || secret != "secret" || !found || base64.RawURLEncoding.EncodeToString(sum[:]) != g.challenge {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid_grant"}`))
			return
		}
		claims, _ := json.Marshal(map[string]any{
			"iss": f.issuer, "aud": "client", "sub": g.claims.Sub, "nonce": g.nonce,
			"email": g.claims.Email, "email_verified": g.claims.Verified,
			"iat": time.Now().Unix(), "exp": time.Now().Add(time.Hour).Unix(),
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token": "at", "token_type": "Bearer", "expires_in": 3600,
			"id_token": oidctest.SignIDToken(f.key, "k", oidc.RS256, string(claims)),
		})
	default:
		f.keys.ServeHTTP(w, r)
	}
}

type googleRig struct {
	a    *Auth
	mux  *http.ServeMux
	fake *fakeGoogle
}

func withGoogle(t *testing.T, mode Mode) googleRig {
	a := setup(t)
	at(a, apex).mode = mode
	fake := newFakeGoogle()
	srv := httptest.NewServer(fake)
	t.Cleanup(srv.Close)
	fake.start(srv.URL)
	a.google = fake.provider()
	mux := http.NewServeMux()
	a.Mount(mux)
	return googleRig{a: a, mux: mux, fake: fake}
}

var refresh = regexp.MustCompile(`href="([^"]+)"`)

func (g googleRig) begin(t *testing.T, path string, cookies ...*http.Cookie) (string, *http.Cookie) {
	r := httptest.NewRequest("POST", path, strings.NewReader("next=%2Fmembers"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, c := range cookies {
		r.AddCookie(c)
	}
	w := httptest.NewRecorder()
	g.mux.ServeHTTP(w, r)
	m := refresh.FindStringSubmatch(w.Body.String())
	if w.Code != http.StatusOK || m == nil {
		t.Fatalf("begin %s: %d %s", path, w.Code, w.Body)
	}
	for _, c := range w.Result().Cookies() {
		if c.Name == flowCookie {
			return html.UnescapeString(m[1]), c
		}
	}
	t.Fatal("no flow cookie")
	return "", nil
}

func (g googleRig) finish(t *testing.T, to string, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	res, err := client.Get(to)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	r := httptest.NewRequest("GET", res.Header.Get("Location"), nil)
	for _, c := range cookies {
		r.AddCookie(c)
	}
	w := httptest.NewRecorder()
	g.mux.ServeHTTP(w, r)
	return w
}

func (g googleRig) signin(t *testing.T, c fakeClaims) *httptest.ResponseRecorder {
	g.fake.as(c)
	to, flow := g.begin(t, "/auth/google")
	return g.finish(t, to, flow)
}

func (g googleRig) link(t *testing.T, session *http.Cookie, c fakeClaims) *httptest.ResponseRecorder {
	g.fake.as(c)
	to, flow := g.begin(t, "/auth/account/google", session)
	return g.finish(t, to, flow, session)
}

func (g googleRig) session(t *testing.T, user pg.ID) *http.Cookie {
	token, err := startSession(g.a.db, httptest.NewRequest("GET", "/", nil), at(g.a, apex), user)
	if err != nil {
		t.Fatal(err)
	}
	return &http.Cookie{Name: sessionCookie, Value: token}
}

func (g googleRig) owner(t *testing.T, subject string) pg.ID {
	id, _, err := owner(g.a.db, googleKind, subject)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func TestGoogleOpen(t *testing.T) {
	g := withGoogle(t, Open)
	if _, err := createUser(g.a.db, "taken@example.com"); err != nil {
		t.Fatal(err)
	}

	w := g.signin(t, fakeClaims{Sub: "s1", Email: "taken@example.com", Verified: true})
	if w.Code != http.StatusConflict || g.owner(t, "s1") != (pg.ID{}) {
		t.Fatalf("email match: %d, owner %v", w.Code, g.owner(t, "s1"))
	}

	for _, sub := range []string{"s2", "s3"} {
		w := g.signin(t, fakeClaims{Sub: sub, Email: "new.person@gmail.com", Verified: true})
		if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/members" {
			t.Fatalf("%s: %d %s", sub, w.Code, w.Body)
		}
	}
	var users []User
	if err := g.a.db.Raw("SELECT u.id, coalesce(u.email, '') AS email, u.username FROM auth_users u JOIN auth_identities i ON i.user_id = u.id ORDER BY i.subject").Scan(&users).Error; err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 || users[0].Email != "" || users[1].Email != "" || users[0].Username != "new.person" || !strings.HasPrefix(users[1].Username, "new.person-") {
		t.Fatalf("google users: %+v", users)
	}
}

func TestGoogleLinkConflict(t *testing.T) {
	g := withGoogle(t, Invite)
	a, err := createUser(g.a.db, "a@example.com")
	if err != nil {
		t.Fatal(err)
	}
	b, err := createUser(g.a.db, "b@example.com")
	if err != nil {
		t.Fatal(err)
	}

	if w := g.signin(t, fakeClaims{Sub: "s1"}); w.Code != http.StatusForbidden {
		t.Fatalf("unknown identity: %d", w.Code)
	}
	if w := g.link(t, g.session(t, a.ID), fakeClaims{Sub: "s1"}); w.Code != http.StatusSeeOther {
		t.Fatalf("link a: %d %s", w.Code, w.Body)
	}
	if w := g.link(t, g.session(t, b.ID), fakeClaims{Sub: "s1"}); w.Code != http.StatusConflict {
		t.Fatalf("link b: %d", w.Code)
	}
	if g.owner(t, "s1") != a.ID {
		t.Fatal("identity moved")
	}
	if w := g.signin(t, fakeClaims{Sub: "s1"}); w.Code != http.StatusSeeOther {
		t.Fatalf("linked sign-in: %d", w.Code)
	}
}

func TestResetKillsLink(t *testing.T) {
	const email = "a@example.com"
	g := withGoogle(t, Invite)
	u, err := createUser(g.a.db, email)
	if err != nil {
		t.Fatal(err)
	}
	if w := g.link(t, g.session(t, u.ID), fakeClaims{Sub: "s1"}); w.Code != http.StatusSeeOther {
		t.Fatalf("link: %d", w.Code)
	}

	session := g.session(t, u.ID)
	g.fake.as(fakeClaims{Sub: "s2"})
	to, flow := g.begin(t, "/auth/account/google", session)
	if _, _, err := issueLink(t.Context(), g.a.db, g.a.key, at(g.a, apex), email, nil, true); err != nil {
		t.Fatal(err)
	}
	g.finish(t, to, flow, session)

	var n int64
	if err := g.a.db.Table("auth_identities").Count(&n).Error; err != nil || n != 0 {
		t.Fatalf("identities after reset: %d %v", n, err)
	}
}

func TestResetGoogleOnly(t *testing.T) {
	g := withGoogle(t, Open)
	for _, sub := range []string{"s1", "s2"} {
		if w := g.signin(t, fakeClaims{Sub: sub, Email: sub + "@gmail.com", Verified: true}); w.Code != http.StatusSeeOther {
			t.Fatalf("%s: %d", sub, w.Code)
		}
	}
	command := func(verb, who string) string {
		var out strings.Builder
		if err := Command(t.Context(), g.a.db, config(), []string{verb, "--origin", apex, who}, &out); err != nil {
			t.Fatal(err)
		}
		token, _, _ := strings.Cut(strings.SplitN(out.String(), "#", 2)[1], "\n")
		return token
	}
	other := command("link", "s2")
	token := command("reset", "s1")

	var n int64
	if err := g.a.db.Raw("SELECT (SELECT count(*) FROM auth_identities WHERE subject = 's1') + (SELECT count(*) FROM auth_sessions WHERE user_id = ?)", g.owner(t, "s1")).Scan(&n).Error; err != nil || n != 0 {
		t.Fatalf("after reset: %d %v", n, err)
	}
	v, _, err := g.a.redeem(t.Context(), httptest.NewRequest("POST", "/auth/link", nil), at(g.a, apex), token)
	if err != nil || v.verdict != accepted || v.user.Username != "s1" || v.user.Email != "" {
		t.Fatalf("reset link: %v %+v", err, v)
	}
	if follow(t, g.a, token) != stale {
		t.Fatal("reset link reused")
	}
	if follow(t, g.a, other) != accepted {
		t.Fatal("another email-less user's link was revoked")
	}
}
