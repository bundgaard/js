package ast

import (
	"bytes"
	"github.com/bundgaard/js/token"
	"strings"
)

type ArrayLiteral struct {
	Token    *token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Value
}
func (al *ArrayLiteral) String() string {
	var (
		out      bytes.Buffer
		elements []string
	)

	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
