package parser

import (
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/token"
)

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.current, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(ast.Lowest)
	if !p.expectPeek(token.CloseBracket) {
		return nil
	}

	return exp
}
