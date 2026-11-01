package auth

import (
	"crypto/rand"
	stdmail "net/mail"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/pg"
)

func parseEmail(raw string) (string, bool) {
	email := strings.ToLower(strings.TrimSpace(raw))
	if len(email) > 254 || strings.ContainsAny(email, "\x00\r\n") {
		return "", false
	}
	addr, err := stdmail.ParseAddress(email)
	if err != nil || addr.Name != "" || addr.Address != email {
		return "", false
	}
	return email, true
}

// lock coordinates code flows and final account checks with reset by address.
func lock(tx *gorm.DB, email string) error {
	return tx.Exec("SELECT pg_advisory_xact_lock(hashtextextended(?, 0))", "auth:"+email).Error
}

func findUser(tx *gorm.DB, email string) (User, error) {
	var u User
	err := tx.Raw("SELECT id, email, created FROM auth_users WHERE email = ?", email).Take(&u).Error
	return u, err
}

func createUser(tx *gorm.DB, email string) (User, error) {
	handle := make([]byte, 16)
	rand.Read(handle)
	u := User{ID: pg.NewID(), Email: email, Created: time.Now()}
	err := tx.Exec("INSERT INTO auth_users (id, email, handle, created) VALUES (?, ?, ?, ?)", u.ID, u.Email, handle, u.Created).Error
	return u, err
}

type holder struct {
	User
	handle []byte
	creds  []webauthn.Credential
}

func (h *holder) WebAuthnID() []byte                         { return h.handle }
func (h *holder) WebAuthnName() string                       { return h.Email }
func (h *holder) WebAuthnDisplayName() string                { return h.Email }
func (h *holder) WebAuthnCredentials() []webauthn.Credential { return h.creds }

func loadHolder(db *gorm.DB, column string, v any) (*holder, error) {
	var row struct {
		ID      pg.ID
		Email   string
		Created time.Time
		Handle  []byte
	}
	if err := db.Raw("SELECT id, email, created, handle FROM auth_users WHERE "+column+" = ?", v).Take(&row).Error; err != nil {
		return nil, err
	}
	var keys []passkey
	if err := db.Where("user_id = ?", row.ID).Find(&keys).Error; err != nil {
		return nil, err
	}
	h := &holder{User: User{ID: row.ID, Email: row.Email, Created: row.Created}, handle: row.Handle}
	for _, k := range keys {
		h.creds = append(h.creds, k.Credential.V)
	}
	return h, nil
}
