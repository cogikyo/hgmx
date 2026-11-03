package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/pg"
)

const usage = "usage: auth link <email|username> | auth reset <email|username>"

// Command takes ["link", who] or ["reset", who], without the "auth" prefix.
// Who is an email or an existing username; both print a secret, one-use, 30-minute sign-in link.
// An unknown email creates a user; an unknown username is an error.
// Reset also removes passkeys and sessions and invalidates codes and links.
// Config.Layout is not required.
func Command(ctx context.Context, db *gorm.DB, c Config, args []string, out io.Writer) error {
	if len(args) != 2 || args[0] != "link" && args[0] != "reset" {
		return errors.New(usage)
	}
	origin, err := c.check()
	if err != nil {
		return err
	}
	email, ok := parseEmail(args[1])
	var want *pg.ID
	if !ok {
		name := strings.TrimSpace(args[1])
		var u User
		res := db.WithContext(ctx).Raw("SELECT "+userColumns+` FROM auth_users WHERE lower(username COLLATE "C") = ?`, foldASCII(name)).Scan(&u)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("auth: no user has the email or username %q", args[1])
		}
		if u.Email == "" {
			return fmt.Errorf("auth: %s has no email, which sign-in links need", name)
		}
		email, want = u.Email, &u.ID
	}
	t, u, err := issueLink(ctx, db, c.Key, email, want, args[0] == "reset")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "%s/auth/link#%s\nThe link signs %s in once and expires in 30 minutes.\n", origin, t, u.name())
	return err
}
