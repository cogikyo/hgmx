package auth

import (
	"context"
	"errors"
	"fmt"
	"io"

	"gorm.io/gorm"
)

const usage = "usage: auth link <email> | auth reset <email>"

// Command takes ["link", email] or ["reset", email], without the "auth" prefix.
// Both create a missing user and print a secret, one-use, 30-minute sign-in link.
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
	if !ok {
		return fmt.Errorf("auth: invalid email %q", args[1])
	}
	t, err := issueLink(ctx, db, c.Key, email, args[0] == "reset")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "%s/auth/link#%s\nThe link signs %s in once and expires in 30 minutes.\n", origin, t, email)
	return err
}
