package parser

import (
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/token"
)

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {

	exp := &ast.CallExpression{
		Token:    p.current,
		Function: function,
	}
	exp.Arguments = p.parseCallArguments()
	return exp
}
func (p *Parser) parseCallArguments() []ast.Expression {
	// log.Printf("parseCallArguments  %v %s %s", p.current, p.current.Value, p.current.Type)
	var args []ast.Expression

	if p.peekTokenIs(token.CloseParen) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(ast.Lowest))
	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(ast.Lowest))
	}

	if !p.expectPeek(token.CloseParen) {
		return nil
	}
	// log.Printf("call arguments %v", args)
	return args
}
