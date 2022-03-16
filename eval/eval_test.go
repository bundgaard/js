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

var usch = "test1 "  + "test 2";
println(z);
println(zz);
println(zzz);
println(zzzz);
println(len("hello, World"));

var u = {"hej": 1, "med": "two"};

var o = ["hej", "med", "dig"];

`))
	environment := object.NewEnvironment()
	evaluatedProgram := Eval(p.Parse(), environment)
	t.Log(evaluatedProgram)
	environment.ForEach(func(key string, val object.Object) {
		t.Logf("%v -> %v", key, val.Inspect())
	})
}

func TestEvalBuiltin(t *testing.T) {
	code := `println("Hello, World!");`
	e := object.NewEnvironment()
	newParser := parser.NewParser(strings.NewReader(code))
	eProgram := Eval(newParser.Parse(), e)

	t.Logf("%v", eProgram)

}
