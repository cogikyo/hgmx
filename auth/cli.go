package auth

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/internal/webhost"
	"github.com/cogikyo/hgmx/pg"
)

const usage = "usage: auth link [--origin URL] <email|username> | auth reset [--origin URL] <email|username>"

// Command takes link|reset [--origin URL] <who>, without the "auth" prefix.
// --origin is required only with more than one configured host.
// Who is an email or an existing username; both print a secret, one-use, 30-minute sign-in link.
// An unknown email creates a user; an unknown username is an error.
// Links identify users by ID, so an email is not required.
// Reset unlinks Google and revokes passkeys, sessions, codes and links on all hosts.
// Host.Layout is not required.
func Command(ctx context.Context, db *gorm.DB, c Config, args []string, out io.Writer) error {
	if len(args) == 0 || args[0] != "link" && args[0] != "reset" {
		return errors.New(usage)
	}
	flags := flag.NewFlagSet("auth "+args[0], flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	raw := flags.String("origin", "", "host origin the link signs in to")
	if err := flags.Parse(args[1:]); err != nil || flags.NArg() != 1 {
		return errors.New(usage)
	}
	hosts, _, err := c.check()
	if err != nil {
		return err
	}
	h, err := target(hosts, *raw)
	if err != nil {
		return err
	}
	who := flags.Arg(0)
	email, ok := parseEmail(who)
	var want *pg.ID
	if !ok {
		name := strings.TrimSpace(who)
		var u User
		res := db.WithContext(ctx).Raw("SELECT "+userColumns+` FROM auth_users WHERE lower(username COLLATE "C") = ?`, foldASCII(name)).Scan(&u)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("auth: no user has the email or username %q", who)
		}
		email, want = u.Email, &u.ID
	}
	t, u, err := issueLink(ctx, db, c.Key, h, email, want, args[0] == "reset")
	if errors.Is(err, errDenied) {
		return fmt.Errorf("auth: %s may not sign in to %s", who, h.origin)
	}
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "%s/auth/link#%s\nThe link signs %s in once and expires in 30 minutes.\n", h.origin, t, u.name())
	return err
}

func target(hosts map[string]*host, raw string) (*host, error) {
	if raw == "" {
		if len(hosts) == 1 {
			for _, h := range hosts {
				return h, nil
			}
		}
		return nil, fmt.Errorf("auth: --origin is required; one of %s", strings.Join(slices.Sorted(maps.Keys(hosts)), ", "))
	}
	origin, _, err := webhost.Parse(raw)
	if err != nil {
		return nil, err
	}
	if h, ok := hosts[origin]; ok {
		return h, nil
	}
	return nil, fmt.Errorf("auth: %s is not a configured host; one of %s", raw, strings.Join(slices.Sorted(maps.Keys(hosts)), ", "))
}
