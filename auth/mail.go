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
	if a.mail != nil {
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

func (a *Auth) codeMail(ctx context.Context, to, code string) mail.Message {
	return message(ctx, to,
		fmt.Sprintf("%s is your %s sign-in code", code, a.c.Name),
		fmt.Sprintf("Your %s sign-in code is %s.\n\nIt expires in 10 minutes. If you didn't ask for it, you can ignore this email.\n", a.c.Name, code),
		codeMail(a.c.Name, code))
}

func (a *Auth) signedInMail(ctx context.Context, to string) mail.Message {
	when := time.Now().UTC().Format("2 Jan 2006 15:04 MST")
	url := a.origin + accountPath
	return message(ctx, to,
		fmt.Sprintf("New sign-in to %s with an email code", a.c.Name),
		fmt.Sprintf("Someone signed in to your %s account with an email code on %s.\n\nIf it wasn't you, review your passkeys and sign out other devices at %s\n", a.c.Name, when, url),
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
