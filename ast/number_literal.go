package ast

import (
	"fmt"
	"github.com/bundgaard/js/token"
)

type NumberLiteral struct {
	Token *token.Token
	Value int64
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Value }
func (nl *NumberLiteral) String() string {
	return fmt.Sprint(nl.Value)
}
