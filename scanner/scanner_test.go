package scanner

import (
	token2 "github.com/bundgaard/js/token"
	"strings"
	"testing"
)

func TestScannerWithEscaped(t *testing.T) {
	s := New(strings.NewReader(`

	var t = "Hello, \"World!\" String";
	`))

	tk := s.NextToken()
	isToken(t, tk, token2.Var)

	tk = s.NextToken()
	isToken(t, tk, token2.Ident)

	tk = s.NextToken()
	isToken(t, tk, token2.Assign)

	tk = s.NextToken()
	isToken(t, tk, token2.String)

	if tk.Value != "Hello, \\\"World!\\\" String" {
		t.Errorf("expected %v (%d). got %v (%d)", "Hello, \\\"World!\\\"", len("Hello, \\\"World!\\\""), []byte(tk.Value), len(tk.Value))
	}
}
func isToken(t *testing.T, tk *token2.Token, expected token2.Type) {
	t.Helper()
	if tk.Type != expected {
		t.Errorf("expected %q. got %q", expected, tk)
	}
}
func TestScanner(t *testing.T) {
	s := New(strings.NewReader(`
	/* 
	Hello World
	*/

	// Line comment

	var s = 12345678;
	var t = "Hello, World!";
	`))

	token := s.NextToken()
	if token.Type != token2.Var {
		t.Errorf("expected variable %s", token.Type.String())
	}

	token = s.NextToken()
	if token.Type != token2.Ident {
		t.Errorf("expected identifier %s", token.Type.String())
	}

	token = s.NextToken()
	if token.Type != token2.Assign {
		t.Errorf("expected Assign %s", token.Type.String())
	}

	token = s.NextToken()
	if token.Type != token2.Number {
		t.Errorf("expected number %s", token.Type.String())
	}

	token = s.NextToken()
	if token.Type != token2.Semi {
		t.Errorf("expected Semi. Got %q", token.Type)
	}

	token = s.NextToken()
	if token.Type != token2.Var {
		t.Errorf("expected Var got %q", token.Type)
	}

	token = s.NextToken()
	if token.Type != token2.Ident {
		t.Errorf("expected Identifier. Got %q", token.Type)
	}

	token = s.NextToken()
	if token.Type != token2.Assign {
		t.Errorf("expected Assign. Got %q", token.Type)
	}

	token = s.NextToken()
	if token.Type != token2.String {
		t.Errorf("expected String. Got %q", token.Type)
	}

	token = s.NextToken()
	if token.Type != token2.Semi {
		t.Errorf("expected Semi. Got %q", token.Type)
	}
}
