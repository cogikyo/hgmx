package credential

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/cogikyo/hgmx/cli"
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
		err := r.command(context.Background(), flag.NewFlagSet("site", flag.ContinueOnError), []string{"check"}, checks, &out)
		if strings.Contains(out.String()+fmt.Sprint(err), secret) {
			t.Errorf("%s: output contains the credential value:\n%s", name, out.String())
		}
		if want := map[bool]int{true: 0, false: 1}[fail == nil]; cli.Exit(err) != want {
			t.Errorf("%s: exit %d, want %d: %v\n%s", name, cli.Exit(err), want, err, out.String())
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
		err := r.command(context.Background(), fs, append([]string{"check", "--"}, tc.flags...), checks, &out)
		if cli.Exit(err) != tc.code || out.String() != tc.out {
			t.Errorf("flags %q: exit %d, output %q; want %d, %q", tc.flags, cli.Exit(err), out.String(), tc.code, tc.out)
		}
	}
}

func TestCommandExit(t *testing.T) {
	t.Setenv("CREDENTIALS_DIRECTORY", "")
	t.Setenv("TEST_SMTP_PASSWORD", "x")

	var r Registry
	smtp := r.Declare(Definition{Name: "smtp-password", Env: "TEST_SMTP_PASSWORD", Purpose: "p", Source: Issued{}, Check: "smtp"}, Always, "site")
	for _, tc := range []struct {
		args  []string
		check error
		code  int
	}{
		{nil, nil, 2},
		{[]string{"rotate"}, nil, 2},
		{[]string{"manifest", "extra"}, nil, 2},
		{[]string{"check", "unknown"}, nil, 2},
		{[]string{"check", "--", "-nope"}, nil, 2},
		{[]string{"check"}, errors.New("rejected"), 1},
		{[]string{"check"}, nil, 0},
		{[]string{"manifest"}, nil, 0},
	} {
		checks := Checks{smtp.Name(): func(context.Context, string) error { return tc.check }}
		err := Run(&r, flag.NewFlagSet("site", flag.ContinueOnError), checks)(context.Background(), tc.args)
		if got := cli.Exit(err); got != tc.code {
			t.Errorf("%q, check error %v: exit %d (%v), want %d", tc.args, tc.check, got, err, tc.code)
		}
	}
}
