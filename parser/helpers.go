package parser

import (
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/token"
)

func (p *Parser) curPrecedence() int {
	if p, ok := ast.Precedences[p.current.Type]; ok {
		return p
	}

	return ast.Lowest
}

func (p *Parser) currentTokenIs(tokenType token.Type) bool {
	return p.current.Type == tokenType
}
func (p *Parser) peekTokenIs(tokenType token.Type) bool {
	return p.next.Type == tokenType
}

func (p *Parser) expectPeek(tokenType token.Type) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) nextToken() {
	p.current = p.next
	p.next = p.s.NextToken()
}

func (p *Parser) peekPrecedence() int {
	if v, ok := ast.Precedences[p.next.Type]; ok {
		return v
	}
	return ast.Lowest
}
