package parser

import (
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/token"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {

	tests := []struct {
		Input    string
		Expected *ast.Program
	}{
		{Input: `var s = 1234;`, Expected: &ast.Program{
			Statements: []ast.Statement{
				&ast.VariableStatement{
					Token: &token.Token{Type: token.Var, Value: "var"},
					Name: &ast.Identifier{
						Token: token.New(token.Ident, "s"),
						Value: "s"},
					Value: &ast.NumberLiteral{
						Token: token.New(token.Number, "1234"),
						Value: 1234,
					},
				},
			},
		}},
	}

	for idx, test := range tests {
		p := New(strings.NewReader(test.Input))

		program := p.Parse()
		if !helperDeepEqual(t, test.Expected, program) {
			t.Errorf("test[%04d] expected %v. got %v", idx, test.Expected, program)
		}

	}

}

func helperDeepEqual(t *testing.T, expected, got *ast.Program) bool {
	t.Helper()
	result := false
	if len(got.Statements) != len(expected.Statements) {
		t.Log("expected and received not the same length")
		return false
	}
	for idx := range expected.Statements {
		if expected.Statements[idx].String() == got.Statements[idx].String() {
			result = true
		}
	}
	return result
}
