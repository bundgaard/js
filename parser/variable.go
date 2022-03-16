package parser

import (
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/token"
)

func (p *Parser) parseVariable() *ast.VariableStatement {
	stmt := &ast.VariableStatement{Token: p.current}
	if !p.expectPeek(token.Ident) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.current, Value: p.current.Value}
	if !p.expectPeek(token.Assign) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(ast.Lowest)
	if p.peekTokenIs(token.Semi) {
		p.nextToken()
	}
	return stmt
}
