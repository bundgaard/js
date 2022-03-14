package eval

import (
	"github.com/bundgaard/js/object"
	"github.com/bundgaard/js/parser"
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	p := parser.NewParser(strings.NewReader(`var x = 100;
var y = 100

var z = x * y;
var zz = x + y;
var zzz = x - y;
var zzzz = z / x;
// println(z);
`))
	environment := make(map[string]object.Object)
	evaluatedProgram := Eval(p.Parse(), environment)
	t.Log(evaluatedProgram)
	for k, v := range environment {
		t.Logf("%v -> %v", k, v.Inspect())
	}
}
