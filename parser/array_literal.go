package parser

import (
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/token"
)

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.current}
	array.Elements = p.parseExpressionList(token.CloseBracket)
	return array
}
