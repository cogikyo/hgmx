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

func TestCheckSkipsOffFeatures(t *testing.T) {
	t.Setenv("CREDENTIALS_DIRECTORY", "")
	t.Setenv("TEST_GOOGLE_SECRET", "")

	var r Registry
	on := r.Feature("google", Flag("-google-client-id").Set())
	google := r.Declare(Definition{Name: "google-client-secret", Env: "TEST_GOOGLE_SECRET", Purpose: "p", Source: Issued{}, Check: "google"}, on, "site")
	checks := Checks{google.Name(): func(context.Context, string) error { return nil }}

	for _, tc := range []struct {
		flags []string
		code  int
		out   string
	}{
		{nil, 0, "google-client-secret: skipped; feature google is off\n"},
		{[]string{"-google-client-id", "x"}, 1, "google-client-secret: missing from $CREDENTIALS_DIRECTORY and env TEST_GOOGLE_SECRET\n"},
	} {
		fs := flag.NewFlagSet("site", flag.ContinueOnError)
		fs.String("google-client-id", "", "")
		var out bytes.Buffer
		code := Command(context.Background(), &r, fs, append([]string{"check", "--"}, tc.flags...), checks, &out)
		if code != tc.code || out.String() != tc.out {
			t.Errorf("flags %q: exit %d, output %q; want %d, %q", tc.flags, code, out.String(), tc.code, tc.out)
		}
	}
}
