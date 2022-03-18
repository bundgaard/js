package eval

import (
	"github.com/bundgaard/js/object"
	"github.com/bundgaard/js/parser"
	"log"
	"strings"
	"testing"
)

func TestEvalMapObject(t *testing.T) {
	p := parser.New(strings.NewReader(`
var u = {"hej": 1, "med": "two"};

`))
	environment := object.NewEnvironment()
	evaluatedProgram := Eval(p.Parse(), environment)
	t.Log(evaluatedProgram)
	environment.ForEach(func(key string, val object.Object) {

		t.Logf("%v -> %v", key, val.Inspect())
	})
}
func TestEvalArrayObject(t *testing.T) {
	p := parser.New(strings.NewReader(`
var o = ["hej", "med", "dig"];
var u = [1,2,3,4];
`))
	environment := object.NewEnvironment()
	evaluatedProgram := Eval(p.Parse(), environment)
	t.Log(evaluatedProgram)
	environment.ForEach(func(key string, val object.Object) {
		t.Logf("%v -> %v", key, val.Inspect())
	})
	log.Printf("")
}
func TestEval(t *testing.T) {
	p := parser.New(strings.NewReader(`var x = 100;
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
	newParser := parser.New(strings.NewReader(code))
	eProgram := Eval(newParser.Parse(), e)

	t.Logf("%v", eProgram)

}

func TestEvalFunction1(t *testing.T) {

	p := parser.New(strings.NewReader(`fn hej() { println("hello, world!"); }`))
	env := object.NewEnvironment()
	result := Eval(p.Parse(), env)

	t.Logf("Function result %v", result)
	t.Logf("Environment %v", env)
}

func TestEvalFunction(t *testing.T) {

	p := parser.NewString(`
fn hej() {
println("hello, world!");
}

hej();`)

	output, env := WithEnvironment(p.Parse())
	t.Logf("%v %v", output, env)

	t.Logf("%c", 127)
	t.Logf("%c", 219)
}

func TestBoolean(t *testing.T) {
	p := parser.NewString(`var sand = true;
println(sand);`)

	output, env := WithEnvironment(p.Parse())
	got, err := env.GetBool("sand")
	if err != nil {
		t.Errorf("%v", err)
	}
	if !got {
		t.Errorf("expected true. got %t", got)
	}
	t.Logf("%v %v", output, env)
}

func TestEvalWithIndexOperation(t *testing.T) {
	p := parser.NewString(`var x = ["Yahoo"];
println(x[0]);`)

	output, environ := WithEnvironment(p.Parse())
	t.Logf("%v %v", output, environ)

}

func TestEvalWithHashIndexOperation(t *testing.T) {
	p := parser.NewString(`var x = {"hej": "dig"};
println(x["hej"]);`)

	output, env := WithEnvironment(p.Parse())
	t.Logf("%v %v", output, env)

}
