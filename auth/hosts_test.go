package auth

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cogikyo/hgmx/pg"
	"github.com/go-webauthn/webauthn/webauthn"
)

func signedIn(t *testing.T, a *Auth, origin string, user pg.ID) *http.Cookie {
	token, err := startSession(a.db, httptest.NewRequest("POST", "/", nil), at(a, origin), user)
	if err != nil {
		t.Fatal(err)
	}
	return &http.Cookie{Name: sessionCookie, Value: token}
}

func serve(a *Auth, method, host, path, body string, cookie *http.Cookie, header ...string) int {
	mux := http.NewServeMux()
	a.Mount(mux)
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Host = host
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for i := 0; i+1 < len(header); i += 2 {
		r.Header.Set(header[i], header[i+1])
	}
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w.Code
}

func TestSessionStaysOnHost(t *testing.T) {
	a := setup(t)
	u, err := createUser(a.db, "op@example.com")
	if err != nil {
		t.Fatal(err)
	}
	cookie := signedIn(t, a, admin, u.ID)
	for host, want := range map[string]int{
		"admin.example.com":     http.StatusOK,
		"ADMIN.example.com:443": http.StatusOK,
		"example.com":           http.StatusSeeOther,
		"example.net":           http.StatusMisdirectedRequest,
	} {
		if got := serve(a, "GET", host, accountPath, "", cookie); got != want {
			t.Errorf("admin session on %s: status %d, want %d", host, got, want)
		}
	}
	spoofed := serve(a, "GET", "example.com", accountPath, "", cookie,
		"X-Forwarded-Host", "admin.example.com", "Forwarded", "host=admin.example.com", "X-Forwarded-Proto", "https")
	if spoofed != http.StatusSeeOther {
		t.Errorf("admin session on the apex with forwarding headers naming admin: status %d, want 303", spoofed)
	}
}

func TestProofStaysOnHost(t *testing.T) {
	const email = "op@example.com"
	a := setup(t)
	at(a, apex).mode, at(a, admin).mode = Invite, Invite
	if _, err := createUser(a.db, email); err != nil {
		t.Fatal(err)
	}
	link := addCode(t, a, email, "123456")
	v, _, err := a.redeem(t.Context(), httptest.NewRequest("POST", "/auth/link", nil), at(a, admin), link)
	if err != nil || v.verdict != stale {
		t.Fatalf("apex link on admin: verdict %d, err %v; want stale", v.verdict, err)
	}
	g, err := a.checkCode(t.Context(), httptest.NewRequest("POST", "/", nil), at(a, admin), email, "123456")
	if err != nil || g.verdict != stale {
		t.Fatalf("apex code on admin: verdict %d, err %v; want stale", g.verdict, err)
	}
	if got := follow(t, a, link); got != accepted {
		t.Fatalf("apex link after the admin attempts: verdict %d, want accepted", got)
	}
}

func TestOlderPasskeysDontCount(t *testing.T) {
	a := setup(t)
	u, err := createUser(a.db, "op@example.com")
	if err != nil {
		t.Fatal(err)
	}
	keys := map[string]webauthn.Credential{
		"current": attested(t, a.rp, "current"),
		"admin":   attested(t, "admin.example.com", "admin"),
		"bare":    {ID: []byte("bare")},
	}
	for id, c := range keys {
		k := passkey{ID: []byte(id), UserID: u.ID, Name: id, Credential: pg.JSON[webauthn.Credential]{V: c}, Created: u.Created}
		if err := a.db.Create(&k).Error; err != nil {
			t.Fatal(err)
		}
	}
	if n, err := a.methods(a.db, at(a, apex), u); err != nil || n != 1 {
		t.Fatalf("usable methods %d, err %v; want 1", n, err)
	}
	cookie := signedIn(t, a, apex, u.ID)
	for _, c := range []struct {
		id   string
		want int
	}{{"current", http.StatusConflict}, {"admin", http.StatusSeeOther}, {"bare", http.StatusSeeOther}, {"current", http.StatusConflict}} {
		if got := serve(a, "POST", "example.com", "/auth/account/passkeys/"+encodeID([]byte(c.id))+"/remove", "", cookie); got != c.want {
			t.Errorf("remove %s: status %d, want %d", c.id, got, c.want)
		}
	}
}

func TestSendBudgetPerHost(t *testing.T) {
	const email = "op@example.com"
	a := setup(t)
	at(a, apex).mode, at(a, admin).mode = Invite, Invite
	if _, err := createUser(a.db, email); err != nil {
		t.Fatal(err)
	}
	send := func(origin string) bool {
		s, err := a.issueCode(t.Context(), at(a, origin), email)
		if err != nil {
			t.Fatal(err)
		}
		return s.limited
	}
	if send(apex) || !send(apex) {
		t.Fatal("apex: want one code, then a limited send")
	}
	if send(admin) {
		t.Fatal("admin send limited by the apex budget")
	}

	for range triesPerDay / triesPerCode {
		addCode(t, a, email, "111111")
		for range triesPerCode {
			try(t, a, email, "000000")
		}
	}
	v, err := a.checkCode(t.Context(), httptest.NewRequest("POST", "/", nil), at(a, admin), email, "000000")
	if err != nil || v.verdict != exhausted {
		t.Fatalf("admin guess after the apex spent the wrong-code budget: verdict %d, err %v; want exhausted", v.verdict, err)
	}
}

func TestSignoutOthersKeepsSiblingHost(t *testing.T) {
	a := setup(t)
	u, err := createUser(a.db, "op@example.com")
	if err != nil {
		t.Fatal(err)
	}
	cookie := signedIn(t, a, admin, u.ID)
	signedIn(t, a, admin, u.ID)
	signedIn(t, a, apex, u.ID)
	if got := serve(a, "POST", "admin.example.com", "/auth/account/signout-others", "", cookie); got != http.StatusSeeOther {
		t.Fatalf("sign out others: status %d", got)
	}
	var n struct{ Admin, Apex int }
	a.db.Raw("SELECT count(*) FILTER (WHERE origin = ?) AS admin, count(*) FILTER (WHERE origin = ?) AS apex FROM auth_sessions", admin, apex).Scan(&n)
	if n.Admin != 1 || n.Apex != 1 {
		t.Fatalf("after sign-out-others on admin: %d admin and %d apex sessions, want 1 and 1", n.Admin, n.Apex)
	}
}

func TestAllowRefusesAccount(t *testing.T) {
	a := setup(t)
	at(a, admin).allow = func(u User) bool { return u.Email == "op@example.com" }
	u, err := createUser(a.db, "buyer@example.com")
	if err != nil {
		t.Fatal(err)
	}
	at(a, admin).mode = Open
	cookie := signedIn(t, a, admin, u.ID)
	key := "/auth/account/passkeys/" + encodeID([]byte("k"))
	for _, path := range []string{
		"GET " + accountPath,
		"POST /auth/account/username",
		"POST /auth/account/signout-others",
		"POST /auth/account/delete",
		"POST " + key + "/rename",
		"POST " + key + "/remove",
		"POST /auth/account/passkeys/begin",
		"POST /auth/account/passkeys/finish",
	} {
		method, path, _ := strings.Cut(path, " ")
		if got := serve(a, method, "admin.example.com", path, "username=taken&name=x&confirm=yes", cookie); got != http.StatusForbidden {
			t.Errorf("%s %s for a refused user: status %d, want 403", method, path, got)
		}
	}
	if got, _ := findUser(a.db, u.Email); got.Username != u.Email {
		t.Fatalf("refused user's username changed to %q", got.Username)
	}

	c := config()
	c.Hosts[1].Allow = at(a, admin).allow
	err = Command(t.Context(), a.db, c, []string{"link", "--origin", admin, "new@example.com"}, io.Discard)
	if err == nil {
		t.Fatal("printed a link the admin host refuses")
	}
	if _, err := findUser(a.db, "new@example.com"); !errors.Is(err, pg.ErrNotFound) {
		t.Fatalf("refused link left a user behind: %v", err)
	}
}

func TestSiteRP(t *testing.T) {
	for _, c := range []struct {
		origins []string
		rp      string
	}{
		{[]string{"https://admin.example.com", "https://example.com"}, "example.com"},
		{[]string{"http://localhost:8080", "http://studio.localhost:8080"}, "localhost"},
		{[]string{"http://studio.localhost:8080"}, "studio.localhost"},
		{[]string{"http://127.0.0.1:8080"}, "127.0.0.1"},
		{[]string{"http://localhost:8080", "https://example.com"}, ""},
		{[]string{"http://127.0.0.1:8080", "http://localhost:8080"}, ""},
		{[]string{"https://192.0.2.1"}, ""},
		{[]string{"https://co.uk"}, ""},
	} {
		var hs []Host
		for _, o := range c.origins {
			hs = append(hs, Host{Origin: o})
		}
		_, rp, err := parseHosts(hs)
		if rp != c.rp || (err == nil) != (c.rp != "") {
			t.Errorf("%v: rp %q, err %v; want rp %q", c.origins, rp, err, c.rp)
		}
	}
}
