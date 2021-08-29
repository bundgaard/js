package js

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {

	tests := []struct {
		Input    string
		Expected *Program
	}{
		{Input: `var s = 1234;`, Expected: &Program{
			Statements: []Statement{
				&VariableStatement{
					Token: &Token{Type: Var, Value: "var"},
					Name: &Identifier{
						Token: newToken(Ident, "s"),
						Value: "s"},
					Value: &NumberLiteral{
						Token: newToken(Number, "1234"),
						Value: 1234,
					},
				},
			},
		}},
	}

	for idx, test := range tests {
		p := NewParser(strings.NewReader(test.Input))

		program := p.Parse()
		t.Logf("program %T %+v %T %+v", program, program, test.Expected, test.Expected)
		if !helperDeepEqual(t, test.Expected, program) {
			t.Errorf("test[%04d] expected %v. got %v", idx, test.Expected, program)
		}

	}

}

func helperDeepEqual(t *testing.T, expected, got *Program) bool {
	t.Helper()
	result := false
	if len(got.Statements) != len(expected.Statements) {
		t.Log("expected and received not the same length")
		return false
	}
	for idx := range expected.Statements {
		if expected.Statements[idx] == got.Statements[idx] {
			result = true
		}
	}
	return result
}
