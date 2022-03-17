package ast

import (
	"bytes"
	"github.com/bundgaard/js/token"
	"strings"
)

type HashLiteral struct {
	Token *token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Value }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	var pairs []string
	for k, v := range hl.Pairs {
		pairs = append(pairs, k.String()+":"+v.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
