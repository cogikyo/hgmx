package auth

import "github.com/cogikyo/hgmx/credential"

var KeyCredential = credential.Definition{
	Name:    "auth-key",
	Env:     "HGMX_AUTH_KEY",
	Purpose: "Sign account verification codes and sign-in links.",
	Source:  credential.Generate{Bytes: 32},
}

var GoogleCredential = credential.Definition{
	Name:    "google-client-secret",
	Env:     "HGMX_GOOGLE_CLIENT_SECRET",
	Purpose: "Authenticate the Google sign-in OAuth client.",
	Source:  credential.Issued{Format: "Google web OAuth client secret; commonly GOCSPX-"},
	Check:   "google",
}
