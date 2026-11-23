package auth

import (
	"encoding/json"
	"errors"
	"go/ast"
	"go/parser"
	gotoken "go/token"
	"go/types"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connectproto"
	"google.golang.org/protobuf/proto"

	authv1 "github.com/cogikyo/hgmx/proto/hgmx/auth/v1"
)

func TestEverySentinelHasReason(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	cliOnly := map[string]bool{"errMoved": true}
	sentinels := map[string]bool{}
	keys := map[string]bool{}
	fset := gotoken.NewFileSet()
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, name, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, decl := range f.Decls {
			g, ok := decl.(*ast.GenDecl)
			if !ok || g.Tok != gotoken.VAR {
				continue
			}
			for _, spec := range g.Specs {
				v := spec.(*ast.ValueSpec)
				for i, value := range v.Values {
					if call, ok := value.(*ast.CallExpr); ok && types.ExprString(call.Fun) == "errors.New" {
						sentinels[v.Names[i].Name] = true
					}
					if lit, ok := value.(*ast.CompositeLit); ok && v.Names[i].Name == "reasons" {
						for _, elt := range lit.Elts {
							keys[types.ExprString(elt.(*ast.KeyValueExpr).Key)] = true
						}
					}
				}
			}
		}
	}
	if len(sentinels) == 0 || len(keys) == 0 {
		t.Fatalf("found %d sentinels and %d reasons; the scan is broken", len(sentinels), len(keys))
	}
	for name := range sentinels {
		if !keys[name] && !cliOnly[name] {
			t.Errorf("%s has no entry in reasons, so neither transport can word it", name)
		}
	}
}

func TestReasonsRenderInBothTransports(t *testing.T) {
	for sentinel, want := range reasons {
		t.Run(sentinel.Error(), func(t *testing.T) {
			if want.text == "" || want.status < 400 || want.code == 0 {
				t.Fatalf("incomplete reason %+v", want)
			}

			w := httptest.NewRecorder()
			problem(w, sentinel)
			var body struct{ Problem string }
			if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if w.Code != want.status || body.Problem == "" {
				t.Errorf("browser rendered %d %q", w.Code, body.Problem)
			}

			var ce *connect.Error
			if !errors.As(fault(t.Context(), sentinel), &ce) {
				t.Fatal("fault did not return a connect error")
			}
			if ce.Code() != want.code || ce.Message() != want.text {
				t.Errorf("app rendered %v %q", ce.Code(), ce.Message())
			}
			if want.detail != nil && len(ce.Details()) == 0 {
				t.Error("app response lost its detail")
			}
		})
	}
}

func TestTransportDifferences(t *testing.T) {
	for err, want := range map[error]string{
		errExhausted:   "Too many wrong codes for this address. Use the link in the email, or try again tomorrow.",
		errUnconfirmed: "Tick the box to confirm deleting your account.",
	} {
		if _, got := shown(err); got != want {
			t.Errorf("browser %v: %q, want %q", err, got, want)
		}
		var ce *connect.Error
		if !errors.As(fault(t.Context(), err), &ce) || ce.Message() == want {
			t.Errorf("app %v: %v, want the app wording", err, ce)
		}
	}

	for err, want := range map[error]proto.Message{
		errRecent: &authv1.SignInAgain{Message: "Sign in again to make this change."},
		errLast:   &authv1.LastSignInMethod{Message: "This is your only way to sign in. Add another one first."},
	} {
		var ce *connect.Error
		if !errors.As(fault(t.Context(), err), &ce) {
			t.Fatalf("%v: fault did not return a connect error", err)
		}
		details := ce.Details()
		if len(details) != 1 {
			t.Fatalf("%v: %d details, want 1", err, len(details))
		}
		got, derr := connectproto.UnmarshalErrorDetail(details[0])
		if derr != nil {
			t.Fatal(derr)
		}
		if !proto.Equal(got, want) {
			t.Errorf("%v detail = %v, want %v", err, got, want)
		}
	}
}
