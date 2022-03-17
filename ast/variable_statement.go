package ast

import (
	"bytes"
	"github.com/bundgaard/js/token"
)

type VariableStatement struct {
	Token *token.Token
	Name  *Identifier
	Value Expression
}

func (vs *VariableStatement) statementNode()       {}
func (vs *VariableStatement) TokenLiteral() string { return vs.Token.Value }
func (vs *VariableStatement) String() string {
	out := new(bytes.Buffer)
	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")
	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}
	return out.String()

}
