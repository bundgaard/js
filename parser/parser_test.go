package parser

import (
	"encoding/json"
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
func TestParserNumberExpression(t *testing.T) {
	p := New(strings.NewReader("(1 + 2) * 100"))
	program := p.Parse()
	content, err := json.Marshal(program)
	if err != nil {
		t.Error(err)
	}
	expect := `{"Statements":[{"Token":{"Type":17,"Value":"("},"Expression":{"Token":{"Type":15,"Value":"*"},"Left":{"Token":{"Type":13,"Value":"+"},"Left":{"Token":{"Type":26,"Value":"1"},"Value":1},"Operator":"+","Right":{"Token":{"Type":26,"Value":"2"},"Value":2}},"Operator":"*","Right":{"Token":{"Type":26,"Value":"100"},"Value":100}}}]}`
	if string(content) != expect {
		t.Fail()
	}
	t.Logf("content %s", content)
	t.Logf("%#v\n\n%s", program, program)
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
