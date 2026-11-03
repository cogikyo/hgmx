package credential

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"strings"
	"testing"
)

func TestCheckNeverPrintsValue(t *testing.T) {
	const secret = "s3cret-value-xyz"
	t.Setenv("CREDENTIALS_DIRECTORY", "")
	t.Setenv("TEST_SMTP_PASSWORD", secret)
	t.Setenv("TEST_GOOGLE_SECRET", secret)

	var r Registry
	smtp := r.Declare(Definition{Name: "smtp-password", Env: "TEST_SMTP_PASSWORD", Purpose: "p", Source: Issued{}, Check: "smtp"}, Always, "site")
	google := r.Declare(Definition{Name: "google-client-secret", Env: "TEST_GOOGLE_SECRET", Purpose: "p", Source: Issued{}, Check: "google"}, Always, "site")

	for name, fail := range map[string]error{
		"plain error": errors.New("auth failed for " + secret),
		"failure":     &Failure{Reason: "SMTP authentication rejected", Code: 535},
		"pass":        nil,
	} {
		var out bytes.Buffer
		checks := Checks{
			smtp.Name():   func(_ context.Context, v string) error { return fail },
			google.Name(): func(_ context.Context, v string) error { return fail },
		}
		code := Command(context.Background(), &r, flag.NewFlagSet("site", flag.ContinueOnError), []string{"check"}, checks, &out)
		if strings.Contains(out.String(), secret) {
			t.Errorf("%s: output contains the credential value:\n%s", name, out.String())
		}
		if want := map[bool]int{true: 0, false: 1}[fail == nil]; code != want {
			t.Errorf("%s: exit %d, want %d:\n%s", name, code, want, out.String())
		}
	}
}
