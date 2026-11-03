package mail

import "github.com/cogikyo/hgmx/credential"

var PasswordCredential = credential.Definition{
	Name:    "smtp-password",
	Env:     "HGMX_SMTP_PASSWORD",
	Purpose: "Authenticate outbound sign-in mail to the SMTP relay.",
	Source:  credential.Issued{Format: "Opaque SMTP authentication token"},
	Check:   "smtp",
}
