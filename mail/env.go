package mail

import (
	"errors"
	"fmt"
	stdmail "net/mail"
	"os"
	"strings"

	"github.com/cogikyo/hgmx/credential"
)

// ErrNoSMTP means no SMTP setting is present, tested with errors.Is.
var ErrNoSMTP = errors.New("mail: HGMX_SMTP_* settings unset")

// FromEnv reads HGMX_SMTP_HOST, HGMX_SMTP_USER and HGMX_SMTP_FROM.
// It needs all three and password, or returns [ErrNoSMTP] if all four are empty.
// Partial settings or an invalid From address return an error.
func FromEnv(password string) (SMTP, error) {
	s := SMTP{Host: os.Getenv("HGMX_SMTP_HOST"), User: os.Getenv("HGMX_SMTP_USER"), Password: password, From: os.Getenv("HGMX_SMTP_FROM")}
	var missing []string
	for _, f := range []struct{ name, value string }{
		{"HGMX_SMTP_HOST", s.Host},
		{"HGMX_SMTP_USER", s.User},
		{"HGMX_SMTP_FROM", s.From},
		{fmt.Sprintf("credential %s (or env %s)", PasswordCredential.Name, PasswordCredential.Env), s.Password},
	} {
		if f.value == "" {
			missing = append(missing, f.name)
		}
	}
	switch len(missing) {
	case 0:
	case 4:
		return SMTP{}, ErrNoSMTP
	default:
		return SMTP{}, &credential.Failure{Reason: "mail: SMTP needs every setting; missing " + strings.Join(missing, ", ")}
	}
	if _, err := stdmail.ParseAddress(s.From); err != nil {
		return SMTP{}, &credential.Failure{Reason: fmt.Sprintf("mail: HGMX_SMTP_FROM %q is not an email address", s.From)}
	}
	return s, nil
}
