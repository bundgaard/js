package js

import (
	"github.com/bundgaard/js/eval"
	"github.com/bundgaard/js/object"
	"github.com/bundgaard/js/parser"
	"strings"
)

func New(data string) (object.Object, *object.Environment) {
	p := parser.New(strings.NewReader(data))
	e := object.NewEnvironment()
	return eval.Eval(p.Parse(), e), e
}
