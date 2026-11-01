package auth

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/migrate"
	"github.com/cogikyo/hgmx/pg"
	"github.com/cogikyo/hgmx/pg/pgtest"
)

func setup(t *testing.T) *Auth {
	t.Setenv("PGDATABASE", "pgtest")
	db := pgtest.Open(t)
	if err := migrate.Command(t.Context(), db, []migrate.Track{Track}, []string{"up"}, io.Discard); err != nil {
		t.Fatal(err)
	}
	a, err := New(db, nil, Config{
		Name:   "Test",
		Origin: "https://example.com",
		Key:    []byte(strings.Repeat("k", 32)),
		Layout: func(_ *http.Request, _ string, body templ.Component) templ.Component { return body },
	})
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func addCode(t *testing.T, a *Auth, email, secret string) {
	err := a.db.Transaction(func(tx *gorm.DB) error {
		_, err := newCode(tx, a.c.Key, email, signin, secret, codeLife)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
}

func try(t *testing.T, a *Auth, email, guess string) verdict {
	v, err := a.checkCode(t.Context(), httptest.NewRequest("POST", "/auth/code/verify", nil), email, guess)
	if err != nil {
		t.Fatal(err)
	}
	return v.verdict
}

func TestWrongCodeBudget(t *testing.T) {
	a := setup(t)
	const email = "a@example.com"
	if _, err := createUser(a.db, email); err != nil {
		t.Fatal(err)
	}

	addCode(t, a, email, "111111")
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
	addCode(t, a, email, "333333")
	if got := try(t, a, email, "333333"); got != exhausted {
		t.Fatalf("fresh code after %d wrong tries across codes: verdict %d, want exhausted", triesPerDay, got)
	}
}

func TestCeremonySingleUse(t *testing.T) {
	a := setup(t)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/auth/passkey/begin", nil)
	if err := a.saveCeremony(w, r, login, nil, &webauthn.SessionData{Challenge: "c"}); err != nil {
		t.Fatal(err)
	}
	cookie := w.Result().Cookies()[0].Value

	if _, ok, err := takeCeremony(a.db, cookie, login, a.library); err != nil || !ok {
		t.Fatalf("first take: ok %v, err %v", ok, err)
	}
	if _, ok, err := takeCeremony(a.db, cookie, login, a.library); err != nil || ok {
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
			_, err := finishLogin(tx, r, user, &cred)
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
	raw, err := issueLink(t.Context(), a.db, a.c.Key, email, false)
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
		_, ok, err := a.redeem(t.Context(), httptest.NewRequest("POST", "/auth/link", nil), raw)
		done <- result{ok, err}
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
