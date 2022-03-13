package js

import (
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	p := NewParser(strings.NewReader(`var x = 100;
var y = 100

var z = x + y;
println(z);
`))
	environment := make(map[string]Object)
	evaluatedProgram := Eval(p.Parse(), environment)
	t.Log(evaluatedProgram)
	for k, v := range environment {
		t.Logf("%v -> %v", k, v.Inspect())
	}
}
