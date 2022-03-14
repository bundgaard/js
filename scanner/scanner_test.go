package scanner

import (
	token2 "github.com/bundgaard/js/token"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	s := NewScanner(strings.NewReader(`
	/* 
	Hello World
	*/

	// Line comment

	var s = 12345678;
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
}
