package ast

import (
	"github.com/bundgaard/js/token"
	"strings"
)

type BlockStatement struct {
	Token      *token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Value
}

func (bs *BlockStatement) String() string {
	var out strings.Builder

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
