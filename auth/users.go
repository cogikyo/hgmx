package auth

import (
	"crypto/rand"
	stdmail "net/mail"
	"regexp"
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

var nameChars = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

func foldASCII(s string) string {
	return strings.Map(func(r rune) rune {
		if 'A' <= r && r <= 'Z' {
			return r + 'a' - 'A'
		}
		return r
	}, s)
}

func ownEmail(name, email string) bool { return email != "" && foldASCII(name) == email }

func validName(name, email string) bool {
	if strings.Contains(name, "@") {
		return ownEmail(name, email)
	}
	return nameChars.MatchString(name)
}

func (u User) name() string {
	if u.Username != "" {
		return u.Username
	}
	return u.Email
}

// lock must precede user row locks to avoid cycles with reset and email changes.
func lock(tx *gorm.DB, email string) error {
	return tx.Exec("SELECT pg_advisory_xact_lock(hashtextextended(?, 0))", "auth:"+email).Error
}

const userColumns = "id, coalesce(email, '') AS email, coalesce(username, '') AS username, created"

func findUser(tx *gorm.DB, email string) (User, error) {
	var u User
	err := tx.Raw("SELECT "+userColumns+" FROM auth_users WHERE email = ?", email).Take(&u).Error
	return u, err
}

func hold(tx *gorm.DB, id pg.ID) (User, error) {
	var u User
	err := tx.Raw("SELECT "+userColumns+" FROM auth_users WHERE id = ? FOR NO KEY UPDATE", id).Take(&u).Error
	return u, err
}

func createUser(tx *gorm.DB, email string) (User, error) {
	handle := make([]byte, 16)
	rand.Read(handle)
	u := User{ID: pg.NewID(), Email: email, Username: email, Created: time.Now()}
	err := tx.Exec("INSERT INTO auth_users (id, email, username, handle, created) VALUES (?, ?, ?, ?, ?)", u.ID, u.Email, u.Username, handle, u.Created).Error
	return u, err
}

func claim(tx *gorm.DB, column string, v any) error {
	return tx.Exec("UPDATE auth_users SET username = email WHERE "+column+" = ? AND username IS NULL AND email IS NOT NULL", v).Error
}

type holder struct {
	User
	handle []byte
	creds  []webauthn.Credential
}

func (h *holder) WebAuthnID() []byte                         { return h.handle }
func (h *holder) WebAuthnName() string                       { return h.name() }
func (h *holder) WebAuthnDisplayName() string                { return h.name() }
func (h *holder) WebAuthnCredentials() []webauthn.Credential { return h.creds }

func loadHolder(db *gorm.DB, column string, v any) (*holder, error) {
	var row struct {
		User
		Handle []byte
	}
	if err := db.Raw("SELECT "+userColumns+", handle FROM auth_users WHERE "+column+" = ?", v).Take(&row).Error; err != nil {
		return nil, err
	}
	var keys []passkey
	if err := db.Where("user_id = ?", row.ID).Find(&keys).Error; err != nil {
		return nil, err
	}
	h := &holder{User: row.User, handle: row.Handle}
	for _, k := range keys {
		h.creds = append(h.creds, k.Credential.V)
	}
	return h, nil
}
