package auth

import (
	"crypto/sha256"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/go-webauthn/webauthn/protocol/webauthncbor"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/migrate"
	"github.com/cogikyo/hgmx/pg"
	"github.com/cogikyo/hgmx/pg/pgtest"
)

const (
	apex  = "https://example.com"
	admin = "https://admin.example.com"
)

func config() Config {
	layout := func(_ *http.Request, _ string, body templ.Component) templ.Component { return body }
	return Config{
		Name:  "Test",
		Key:   []byte(strings.Repeat("k", 32)),
		Hosts: []Host{{Origin: apex, Layout: layout}, {Origin: admin, Layout: layout}},
	}
}

func open(t *testing.T, db *gorm.DB) *Auth {
	a, err := New(db, nil, config())
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func setup(t *testing.T) *Auth {
	t.Setenv("PGDATABASE", "pgtest")
	db := pgtest.Open(t)
	if err := migrate.Command(t.Context(), db, []migrate.Track{Track}, []string{"up"}, io.Discard); err != nil {
		t.Fatal(err)
	}
	return open(t, db)
}

func at(a *Auth, origin string) *host { return a.hosts[origin] }

func addCode(t *testing.T, a *Auth, email, secret string) string {
	var c code
	err := a.db.Transaction(func(tx *gorm.DB) error {
		var err error
		c, err = newCode(tx, a.key, code{Origin: apex, Email: email, Purpose: signin}, secret, codeLife)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return c.ID.String() + "." + mailSecret(a.key, c)
}

func try(t *testing.T, a *Auth, email, guess string) verdict {
	v, err := a.checkCode(t.Context(), httptest.NewRequest("POST", "/auth/code/verify", nil), at(a, apex), email, guess)
	if err != nil {
		t.Fatal(err)
	}
	return v.verdict
}

func follow(t *testing.T, a *Auth, token string) verdict {
	v, _, err := a.redeem(t.Context(), httptest.NewRequest("POST", "/auth/link", nil), at(a, apex), token)
	if err != nil {
		t.Fatal(err)
	}
	return v.verdict
}

func invite(t *testing.T, email string) *Auth {
	a := setup(t)
	at(a, apex).mode = Invite
	if _, err := createUser(a.db, email); err != nil {
		t.Fatal(err)
	}
	return a
}

func TestWrongCodeBudget(t *testing.T) {
	const email = "a@example.com"
	a := invite(t, email)

	first := addCode(t, a, email, "111111")
	for i := range triesPerCode {
		if got := try(t, a, email, "000000"); got != wrong {
			t.Fatalf("first code, wrong try %d: verdict %d, want wrong", i+1, got)
		}
	}
	if got := try(t, a, email, "111111"); got != stale {
		t.Fatalf("right guess after %d wrong tries: verdict %d, want stale", triesPerCode, got)
	}

	addCode(t, a, email, "222222")
	for range triesPerDay - triesPerCode {
		try(t, a, email, "000000")
	}
	latest := addCode(t, a, email, "333333")
	if got := try(t, a, email, "333333"); got != exhausted {
		t.Fatalf("fresh code after %d wrong tries across codes: verdict %d, want exhausted", triesPerDay, got)
	}
	if got := follow(t, a, latest+"x"); got != stale {
		t.Fatalf("wrong link secret: verdict %d, want stale", got)
	}
	for _, link := range []string{first, latest} {
		if got := follow(t, a, link); got != accepted {
			t.Fatalf("link after the code budget is spent: verdict %d, want accepted", got)
		}
	}
}

func TestCodeAndLinkPaired(t *testing.T) {
	const email = "b@example.com"
	a := invite(t, email)

	link := addCode(t, a, email, "111111")
	if got := try(t, a, email, "111111"); got != accepted {
		t.Fatalf("code: verdict %d, want accepted", got)
	}
	if got := follow(t, a, link); got != stale {
		t.Fatalf("link after its code: verdict %d, want stale", got)
	}

	link = addCode(t, a, email, "222222")
	if got := follow(t, a, link); got != accepted {
		t.Fatalf("link: verdict %d, want accepted", got)
	}
	if got := try(t, a, email, "222222"); got != stale {
		t.Fatalf("code after its link: verdict %d, want stale", got)
	}

	link = addCode(t, a, email, "333333")
	r := httptest.NewRequest("POST", "/auth/code/verify", nil)
	type result struct {
		v   verified
		err error
	}
	got := make(chan result, 2)
	go func() {
		v, err := a.checkCode(t.Context(), r, at(a, apex), email, "333333")
		got <- result{v, err}
	}()
	go func() {
		v, _, err := a.redeem(t.Context(), r, at(a, apex), link)
		got <- result{v, err}
	}()
	x, y := <-got, <-got
	if x.err != nil || y.err != nil || (x.v.verdict == accepted) == (y.v.verdict == accepted) {
		t.Fatalf("concurrent code and link: verdicts %d and %d, errors %v and %v; want exactly one accepted", x.v.verdict, y.v.verdict, x.err, y.err)
	}
}

func TestNextStaysOnSite(t *testing.T) {
	for raw, want := range map[string]string{
		`/a/../\evil.example/`: accountPath,
		`/a/..//evil.example/`: accountPath,
		`/members/?u=//x`:      `/members/?u=//x`,
	} {
		w := httptest.NewRecorder()
		http.Redirect(w, httptest.NewRequest("POST", "/auth/link", nil), next(raw), http.StatusSeeOther)
		if got := w.Header().Get("Location"); got != want {
			t.Errorf("next(%q) redirects to %q, want %q", raw, got, want)
		}
	}
}

func TestCeremonySingleUse(t *testing.T) {
	a := setup(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/auth/passkey/begin", nil)
	if err := a.saveCeremony(w, r, login, nil, &webauthn.SessionData{Challenge: "c", Origin: apex}); err != nil {
		t.Fatal(err)
	}
	cookie := w.Result().Cookies()[0].Value

	if _, ok, err := takeCeremony(a.db, at(a, apex), cookie, login, a.library); err != nil || !ok {
		t.Fatalf("first take: ok %v, err %v", ok, err)
	}
	if _, ok, err := takeCeremony(a.db, at(a, apex), cookie, login, a.library); err != nil || ok {
		t.Fatalf("second take: ok %v, err %v; want a consumed ceremony", ok, err)
	}
}

func TestFinishLoginRequiresOwner(t *testing.T) {
	a := setup(t)
	owner, err := createUser(a.db, "owner@example.com")
	if err != nil {
		t.Fatal(err)
	}
	other, err := createUser(a.db, "other@example.com")
	if err != nil {
		t.Fatal(err)
	}
	cred := webauthn.Credential{ID: []byte("cred"), PublicKey: []byte("pk")}
	if err := a.db.Create(&passkey{ID: cred.ID, UserID: owner.ID, Name: "k", Credential: pg.JSON[webauthn.Credential]{V: cred}, Created: owner.Created}).Error; err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("POST", "/auth/passkey/finish", nil)

	cred.Authenticator.SignCount = 7
	finish := func(user pg.ID) error {
		return a.db.Transaction(func(tx *gorm.DB) error {
			_, err := finishLogin(tx, r, at(a, apex), user, &cred)
			return err
		})
	}
	if err := finish(other.ID); !errors.Is(err, pg.ErrNotFound) {
		t.Fatalf("finish for a user without the credential: err %v, want ErrNotFound", err)
	}
	var sessions int64
	a.db.Table("auth_sessions").Count(&sessions)
	if sessions != 0 {
		t.Fatalf("%d sessions after a refused finish, want 0", sessions)
	}

	if err := finish(owner.ID); err != nil {
		t.Fatal(err)
	}
	var k passkey
	a.db.Take(&k, "id = ?", cred.ID)
	if k.Credential.V.Authenticator.SignCount != 7 || k.Used == nil {
		t.Fatalf("owner finish stored sign count %d, used %v; want 7 and a time", k.Credential.V.Authenticator.SignCount, k.Used)
	}
}

func TestResetBeatsRedeem(t *testing.T) {
	a := setup(t)
	const email = "op@example.com"
	raw, _, err := issueLink(t.Context(), a.db, a.key, at(a, apex), email, nil, false)
	if err != nil {
		t.Fatal(err)
	}

	reset := a.db.Begin()
	defer reset.Rollback()
	if err := lock(reset, email); err != nil {
		t.Fatal(err)
	}
	if err := reset.Exec("DELETE FROM auth_sessions").Error; err != nil {
		t.Fatal(err)
	}
	type result struct {
		ok  bool
		err error
	}
	done := make(chan result, 1)
	go func() {
		v, _, err := a.redeem(t.Context(), httptest.NewRequest("POST", "/auth/link", nil), at(a, apex), raw)
		done <- result{v.verdict == accepted, err}
	}()
	for deadline := time.Now().Add(5 * time.Second); ; {
		var waiting int64
		if err := a.db.Raw("SELECT count(*) FROM pg_locks WHERE locktype = 'advisory' AND NOT granted").Scan(&waiting).Error; err != nil {
			t.Fatal(err)
		}
		if waiting > 0 {
			break
		}
		if len(done) > 0 || time.Now().After(deadline) {
			t.Fatal("redeem did not wait on the address lock while a reset held it")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := reset.Exec("UPDATE auth_codes SET consumed = now() WHERE email = ?", email).Error; err != nil {
		t.Fatal(err)
	}
	if err := reset.Commit().Error; err != nil {
		t.Fatal(err)
	}

	got := <-done
	if got.err != nil || got.ok {
		t.Fatalf("redeem during reset: ok %v, err %v; want refused", got.ok, got.err)
	}
	var sessions int64
	a.db.Table("auth_sessions").Count(&sessions)
	if sessions != 0 {
		t.Fatalf("%d sessions survive the reset, want 0", sessions)
	}
}

func TestOldBinaryRows(t *testing.T) {
	t.Setenv("PGDATABASE", "pgtest")
	db := pgtest.Open(t)
	up := func(steps []migrate.Step) {
		if err := migrate.Command(t.Context(), db, []migrate.Track{{Name: Track.Name, Steps: steps}}, []string{"up"}, io.Discard); err != nil {
			t.Fatal(err)
		}
	}
	insert := func(email string) {
		err := db.Exec("INSERT INTO auth_users (id, email, handle, created) VALUES (?, ?, ?, now())", pg.NewID(), email, []byte(email)).Error
		if err != nil {
			t.Fatalf("v1-shaped insert of %s: %v", email, err)
		}
	}
	up(Track.Steps[:1])
	insert("before@example.com")
	up(Track.Steps)
	insert("after@example.com")

	a := open(t, db)
	at(a, apex).mode = Invite
	for _, email := range []string{"before@example.com", "after@example.com"} {
		addCode(t, a, email, "111111")
		if got := try(t, a, email, "111111"); got != accepted {
			t.Fatalf("%s: verdict %d, want accepted", email, got)
		}
		u, err := findUser(a.db, email)
		if err != nil || u.Username != email {
			t.Fatalf("%s after sign-in: username %q, err %v; want the email", email, u.Username, err)
		}
	}
}

func TestUsernameAt(t *testing.T) {
	for _, c := range []struct {
		name, email string
		ok          bool
	}{
		{"Me@Example.com", "me@example.com", true},
		{"you@example.com", "me@example.com", false},
		{"you@example.com", "", false},
		{"me\uff20example.com", "me@example.com", false},
		{"s@example.com", "\u017f@example.com", false},
		{"cullyn.k-2_", "me@example.com", true},
	} {
		if got := validName(c.name, c.email); got != c.ok {
			t.Errorf("validName(%q, %q) = %v, want %v", c.name, c.email, got, c.ok)
		}
	}

	a := setup(t)
	at(a, apex).mode = Invite
	for _, email := range []string{"\u017f@example.com", "s@example.com", "S@example.org"} {
		if _, err := createUser(a.db, email); err != nil {
			t.Fatalf("default username for %q: %v", email, err)
		}
	}
	if _, err := createUser(a.db, "s@example.org"); !errors.Is(err, pg.ErrDuplicate) {
		t.Fatalf("username differing only in ASCII case: err %v, want ErrDuplicate", err)
	}

	u, err := createUser(a.db, "old@example.com")
	if err != nil {
		t.Fatal(err)
	}
	ch, err := a.confirm(t.Context(), at(a, apex), pending(t, a, u.ID, "new@example.com"))
	if err != nil || ch.verdict != accepted {
		t.Fatalf("confirm: verdict %d, err %v; want accepted", ch.verdict, err)
	}
	if _, err := createUser(a.db, "old@example.com"); err != nil {
		t.Fatalf("a new user for the released address: %v", err)
	}
	got, err := findUser(a.db, "new@example.com")
	if err != nil || got.Username != "new@example.com" {
		t.Fatalf("after the email change: username %q, err %v; want the new email", got.Username, err)
	}
}

func TestLastMethodRace(t *testing.T) {
	a := setup(t)
	u, err := createUser(a.db, "op@example.com")
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"k1", "k2"} {
		k := passkey{ID: []byte(id), UserID: u.ID, Name: id, Credential: pg.JSON[webauthn.Credential]{V: attested(t, a.rp, id)}, Created: u.Created}
		if err := a.db.Create(&k).Error; err != nil {
			t.Fatal(err)
		}
	}
	token, err := startSession(a.db, httptest.NewRequest("POST", "/", nil), at(a, apex), u.ID)
	if err != nil {
		t.Fatal(err)
	}

	block := a.db.Begin()
	defer block.Rollback()
	if _, err := hold(block, u.ID); err != nil {
		t.Fatal(err)
	}
	codes := make(chan int, 2)
	for _, id := range []string{"k1", "k2"} {
		go func() {
			r := httptest.NewRequest("POST", "/auth/account/passkeys/"+encodeID([]byte(id))+"/remove", nil)
			r.SetPathValue("id", encodeID([]byte(id)))
			r.AddCookie(&http.Cookie{Name: sessionCookie, Value: token})
			w := httptest.NewRecorder()
			a.remove(w, r, at(a, apex))
			codes <- w.Code
		}()
	}
	for deadline := time.Now().Add(5 * time.Second); ; {
		var waiting int64
		if err := a.db.Raw("SELECT count(*) FROM pg_locks WHERE NOT granted").Scan(&waiting).Error; err != nil {
			t.Fatal(err)
		}
		if waiting >= 2 {
			break
		}
		if len(codes) > 0 || time.Now().After(deadline) {
			t.Fatal("removals did not wait on the user row")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := block.Commit().Error; err != nil {
		t.Fatal(err)
	}

	x, y := <-codes, <-codes
	if x+y != http.StatusSeeOther+http.StatusConflict {
		t.Fatalf("concurrent removals of the last two passkeys: statuses %d and %d, want one 303 and one 409", x, y)
	}
	var left int64
	a.db.Model(&passkey{}).Where("user_id = ?", u.ID).Count(&left)
	if left != 1 {
		t.Fatalf("%d passkeys left, want 1", left)
	}
}

func pending(t *testing.T, a *Auth, user pg.ID, to string) string {
	token, err := startSession(a.db, httptest.NewRequest("POST", "/", nil), at(a, apex), user)
	if err != nil {
		t.Fatal(err)
	}
	s, _, err := lookup(a.db, apex, token)
	if err != nil {
		t.Fatal(err)
	}
	var c code
	err = a.db.Transaction(func(tx *gorm.DB) error {
		var err error
		c, err = newCode(tx, a.key, code{Origin: apex, Email: to, Purpose: verify, SessionID: &s.id}, "s", linkLife)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return c.ID.String() + ".s"
}

func TestOldResetCancelsEmailChange(t *testing.T) {
	a := setup(t)
	u, err := createUser(a.db, "victim@example.com")
	if err != nil {
		t.Fatal(err)
	}
	raw := pending(t, a, u.ID, "attacker@example.com")

	for _, q := range []string{
		"DELETE FROM auth_passkeys WHERE user_id = ?",
		"DELETE FROM auth_sessions WHERE user_id = ?",
		"DELETE FROM auth_ceremonies WHERE user_id = ?",
	} {
		if err := a.db.Exec(q, u.ID).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := a.db.Exec("UPDATE auth_codes SET consumed = now() WHERE email = ? AND consumed IS NULL", u.Email).Error; err != nil {
		t.Fatal(err)
	}

	ch, err := a.confirm(t.Context(), at(a, apex), raw)
	if err != nil || ch.verdict == accepted {
		t.Fatalf("confirm after a v1 reset: verdict %d, err %v; want refused", ch.verdict, err)
	}
	if got, err := findUser(a.db, "victim@example.com"); err != nil || got.ID != u.ID {
		t.Fatalf("victim's email after the refused change: %v, %v", got, err)
	}
}

func attested(t *testing.T, rp, id string) webauthn.Credential {
	sum := sha256.Sum256([]byte(rp))
	object, err := webauthncbor.Marshal(map[string]any{"fmt": "none", "attStmt": map[string]any{}, "authData": append(sum[:], 1, 0, 0, 0, 0)})
	if err != nil {
		t.Fatal(err)
	}
	return webauthn.Credential{ID: []byte(id), Attestation: webauthn.CredentialAttestation{Object: object}}
}
