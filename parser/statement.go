package parser

import (
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.current.Type {
	case token.Var:
		return p.parseVariable()
	case token.CommentLine:
		return nil
	case token.CommentBlock:
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {

	// expression statement
	// foo[0] = 1300
	stmt := &ast.ExpressionStatement{Token: p.current}
	stmt.Expression = p.parseExpression(ast.Lowest)
	// dot expression
	// foo.bar
	if p.peekTokenIs(token.Dot) {
		p.nextToken()
		right := p.parseInfixExpression(stmt.Expression)
		stmt.Expression = right
	}

	if p.peekTokenIs(token.Assign) {

		p.nextToken()

		right := p.parseInfixExpression(stmt.Expression)
		stmt.Expression = right
	}

	if p.peekTokenIs(token.Semi) {
		p.nextToken()
	}
	return stmt
}
