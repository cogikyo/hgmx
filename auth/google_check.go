package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cogikyo/hgmx/credential"
)

const unknownRefreshToken = "hgmx-credential-check-unknown-refresh-token"

// CheckGoogle probes token-endpoint client authentication against a decoy secret.
// Success does not verify browser sign-in or redirect registration.
func CheckGoogle(ctx context.Context, clientID, secret string) error {
	return checkGoogle(ctx, googleIssuer.TokenURL, clientID, secret)
}

func checkGoogle(ctx context.Context, tokenURL, clientID, secret string) error {
	if clientID == "" {
		return &credential.Failure{Reason: "Google client ID is not configured"}
	}
	control, ok := decoy(secret)
	if !ok {
		return &credential.Failure{Reason: "unverified: Google client secret has no characters to vary for a decoy"}
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}

	rejected, code, err := googleToken(ctx, client, tokenURL, clientID, control)
	if err != nil {
		return err
	}
	if rejected != "invalid_client" {
		return &credential.Failure{Reason: "unverified: Google did not reject a decoy secret as invalid_client", Code: code}
	}
	reply, code, err := googleToken(ctx, client, tokenURL, clientID, secret)
	switch {
	case err != nil:
		return err
	case reply == "invalid_grant":
		return nil
	case reply == "invalid_client" || reply == "unauthorized_client" || reply == "deleted_client":
		return &credential.Failure{Reason: "Google rejected the client ID or secret: " + reply, Code: code}
	}
	return &credential.Failure{Reason: "unverified: unexpected Google token response", Code: code}
}

func googleToken(ctx context.Context, client *http.Client, tokenURL, clientID, secret string) (string, int, error) {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {unknownRefreshToken},
		"client_id":     {clientID},
		"client_secret": {secret},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", 0, &credential.Failure{Reason: "Google token request could not be built"}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		var timeout interface{ Timeout() bool }
		if errors.As(err, &timeout) && timeout.Timeout() {
			return "", 0, &credential.Failure{Reason: "Google token endpoint timed out"}
		}
		return "", 0, &credential.Failure{Reason: "Google token endpoint unreachable"}
	}
	defer resp.Body.Close()
	if resp.StatusCode < 400 {
		return "", resp.StatusCode, &credential.Failure{Reason: "unverified: Google token endpoint did not refuse the probe", Code: resp.StatusCode}
	}
	var body struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&body); err != nil {
		return "", resp.StatusCode, &credential.Failure{Reason: "unverified: Google token error was not JSON", Code: resp.StatusCode}
	}
	switch body.Error {
	case "invalid_client", "unauthorized_client", "deleted_client", "invalid_grant":
		return body.Error, resp.StatusCode, nil
	}
	return "other", resp.StatusCode, nil
}

func decoy(secret string) (string, bool) {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"
	var vary []int
	for i := strings.IndexByte(secret, '-') + 1; i < len(secret); i++ {
		if strings.IndexByte(alphabet, secret[i]) >= 0 {
			vary = append(vary, i)
		}
	}
	if len(vary) == 0 {
		return "", false
	}
	b := []byte(secret)
	for string(b) == secret {
		for _, i := range vary {
			b[i] = alphabet[rand.IntN(len(alphabet))]
		}
	}
	return string(b), true
}
