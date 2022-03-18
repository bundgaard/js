package object

import (
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/token"
	"strings"
)

type Function struct {
	Token       *token.Token
	Name        string
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment *Environment
}

func (fl *Function) Type() Type { return FunctionType }
func (fl *Function) Inspect() string {
	var out strings.Builder
	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.Token.Value)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}
