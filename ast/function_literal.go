package ast

import (
	"github.com/bundgaard/js/token"
	"strings"
)

type FunctionLiteral struct {
	Token      *token.Token
	Name       string
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Value }
func (fl *FunctionLiteral) String() string {
	var (
		out    strings.Builder
		params []string
	)

	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	if fl.Name != "" {
		out.WriteString(" " + fl.Name + " ")
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}
