package ast

import "github.com/bundgaard/js/token"

type Boolean struct {
	Token *token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) String() string {
	return b.Token.Value
}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Value
}

type Null struct {
	Token *token.Token
	Value string
}

func (n *Null) expressionNode() {}
func (n *Null) String() string {
	return n.Value
}
func (n *Null) TokenLiteral() string {
	return n.Token.Value
}
