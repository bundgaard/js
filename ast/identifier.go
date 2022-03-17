package ast

import "github.com/bundgaard/js/token"

type Identifier struct {
	Token *token.Token
	Value string
}

func (id *Identifier) expressionNode()      {}
func (id *Identifier) TokenLiteral() string { return id.Token.Value }
func (id *Identifier) String() string {
	return id.Value
}
