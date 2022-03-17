package parser

import (
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/token"
)

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	var list []ast.Expression

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(ast.Lowest))

	for p.peekTokenIs(token.Comma) {
		p.nextToken() // ,
		p.nextToken() // Expression
		list = append(list, p.parseExpression(ast.Lowest))
	}

	if !p.expectPeek(end) {
		return nil
	}
	return list
}
