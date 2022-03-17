package ast

import "github.com/bundgaard/js/token"

type StringLiteral struct {
	Token *token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Value }
func (sl *StringLiteral) String() string {
	return sl.Value
}
