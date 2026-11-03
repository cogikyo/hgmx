package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/a-h/templ"

	"github.com/cogikyo/hgmx/mail"
)

func (a *Auth) notify(m mail.Message) {
	if a.mail != nil && m.To != "" {
		a.mail.Enqueue(m)
	}
}

func message(ctx context.Context, to, subject, text string, html templ.Component) mail.Message {
	var b strings.Builder
	if err := html.Render(ctx, &b); err != nil {
		panic(err)
	}
	return mail.Message{To: to, Subject: subject, Text: text, HTML: b.String()}
}

func (a *Auth) codeMail(ctx context.Context, to, code, link string) mail.Message {
	return message(ctx, to,
		fmt.Sprintf("%s is your %s sign-in code", code, a.c.Name),
		fmt.Sprintf("Your %s sign-in code is %s.\n\nOr sign in with this link: %s\n\nBoth expire in 10 minutes and work once. If you didn't ask for them, you can ignore this email.\n", a.c.Name, code, link),
		codeMail(a.c.Name, code, link))
}

func (a *Auth) signedInMail(ctx context.Context, to string) mail.Message {
	when := time.Now().UTC().Format("2 Jan 2006 15:04 MST")
	url := a.origin + accountPath
	return message(ctx, to,
		fmt.Sprintf("New sign-in to %s by email", a.c.Name),
		fmt.Sprintf("Someone signed in to your %s account by email on %s.\n\nIf it wasn't you, review your passkeys and sign out other devices at %s\n", a.c.Name, when, url),
		signedInMail(a.c.Name, when, url))
}

func (a *Auth) passkeyMail(ctx context.Context, to, name string, added bool) mail.Message {
	url := a.origin + accountPath
	change, subject := "removed from", fmt.Sprintf("Passkey removed from your %s account", a.c.Name)
	if added {
		change, subject = "added to", fmt.Sprintf("Passkey added to your %s account", a.c.Name)
	}
	return message(ctx, to, subject,
		fmt.Sprintf("The passkey %q was %s your %s account.\n\nIf you didn't do this, review your passkeys at %s\n", name, change, a.c.Name, url),
		passkeyMail(a.c.Name, name, change, url))
}

func (a *Auth) verifyMail(ctx context.Context, to, link string) mail.Message {
	return message(ctx, to,
		fmt.Sprintf("Confirm your new %s email", a.c.Name),
		fmt.Sprintf("Confirm %s as the email for your %s account with this link: %s\n\nThe link expires in 30 minutes and works once. If you didn't ask for this, you can ignore this email.\n", to, a.c.Name, link),
		verifyMail(a.c.Name, to, link))
}

func (a *Auth) changedMail(ctx context.Context, from, to string) mail.Message {
	url := a.origin + accountPath
	return message(ctx, from,
		fmt.Sprintf("Your %s email was changed", a.c.Name),
		fmt.Sprintf("The email for your %s account was changed to %s.\n\nIf you didn't do this, review your passkeys and sign out other devices at %s\n", a.c.Name, to, url),
		changedMail(a.c.Name, to, url))
}

func (a *Auth) googleMail(ctx context.Context, to string, linked bool) mail.Message {
	url := a.origin + accountPath
	change, subject := "removed from", fmt.Sprintf("Google account removed from your %s account", a.c.Name)
	if linked {
		change, subject = "linked to", fmt.Sprintf("Google account linked to your %s account", a.c.Name)
	}
	return message(ctx, to, subject,
		fmt.Sprintf("A Google account was %s your %s account.\n\nIf you didn't do this, review your account at %s\n", change, a.c.Name, url),
		googleMail(a.c.Name, change, url))
}
