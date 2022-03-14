package ast

import "github.com/bundgaard/js/token"

const (
	_ int = iota
	Lowest
	Equals
	LessGreater
	Sum
	Product
	Prefix
	Call
	Index
)

var Precedences = map[token.TokenType]int{
	token.Add:         Sum,
	token.Sub:         Sum,
	token.Mul:         Product,
	token.Div:         Product,
	token.OpenBracket: Index,
	token.OpenParen:   Call,
}
