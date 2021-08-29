package js

import (
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
	if token.Type != Var {
		t.Errorf("expected variable %s", token.Type.String())
	}

	token = s.NextToken()
	if token.Type != Ident {
		t.Errorf("expected identifier %s", token.Type.String())
	}

	token = s.NextToken()
	if token.Type != Assign {
		t.Errorf("expected Assign %s", token.Type.String())
	}

	token = s.NextToken()
	if token.Type != Number {
		t.Errorf("expected number %s", token.Type.String())
	}
}
