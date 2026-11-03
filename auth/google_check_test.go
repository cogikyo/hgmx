package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cogikyo/hgmx/credential"
)

func TestCheckGoogle(t *testing.T) {
	const id, secret = "client.apps.googleusercontent.com", "GOCSPX-real-secret"
	const hyphenated = "GOCSPX-real-secret-"
	strict := func(want string) func(*http.Request) string {
		return func(r *http.Request) string {
			if r.PostFormValue("client_id") == id && r.PostFormValue("client_secret") == want {
				return "invalid_grant"
			}
			return "invalid_client"
		}
	}
	format := func(want string) func(*http.Request) string {
		return func(r *http.Request) string {
			got := r.PostFormValue("client_secret")
			if len(got) != len(want) || !strings.HasPrefix(got, "GOCSPX-") {
				return "invalid_client"
			}
			return "invalid_grant"
		}
	}
	cases := []struct {
		name   string
		reply  func(*http.Request) string
		secret string
		pass   bool
	}{
		{"right secret", strict(secret), secret, true},
		{"right secret with a trailing hyphen", strict(hyphenated), hyphenated, true},
		{"wrong secret", strict(secret), "GOCSPX-wrong-secret", false},
		{"endpoint ignores the secret", func(*http.Request) string { return "invalid_grant" }, secret, false},
		{"endpoint checks only prefix and length", format(secret), secret, false},
		{"endpoint checks only prefix and length, trailing hyphen", format(hyphenated), hyphenated, false},
		{"secret with nothing to vary", strict("GOCSPX--"), "GOCSPX--", false},
	}
	var reply func(*http.Request) string
	url := fakeTokenEndpoint(t, func(r *http.Request) string { return reply(r) })
	for _, c := range cases {
		reply = c.reply
		err := checkGoogle(context.Background(), url, id, c.secret)
		var f *credential.Failure
		switch {
		case c.pass && err != nil:
			t.Errorf("%s: %v", c.name, err)
		case !c.pass && !errors.As(err, &f):
			t.Errorf("%s: got %v, want a Failure", c.name, err)
		case err != nil && strings.Contains(err.Error(), c.secret):
			t.Errorf("%s: failure echoes the secret: %q", c.name, err.Error())
		}
	}
}

func fakeTokenEndpoint(t *testing.T, reply func(*http.Request) string) string {
	ln, err := net.Listen("tcp", "127.0.0.1:18522")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": reply(r)})
	}))
	srv.Listener.Close()
	srv.Listener = ln
	srv.Start()
	t.Cleanup(srv.Close)
	return srv.URL
}
