package parser

import (
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/token"
	"log"
)

func (p *Parser) parseHashLiteral() ast.Expression {
	log.Printf("parseHashLiteral %q", p.current)
	hash := &ast.HashLiteral{Token: p.current}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.CloseCurly) {
		p.nextToken() // eat open curly
		log.Printf("\t%q", p.current)
		key := p.parseExpression(ast.Lowest)
		log.Printf("\t%q", key)
		log.Printf("\t%q", p.current)
		if !p.expectPeek(token.Colon) {
			return nil
		}

		p.nextToken() // EAT Colon

		value := p.parseExpression(ast.Lowest)

		hash.Pairs[key] = value
		if !p.peekTokenIs(token.CloseCurly) && !p.expectPeek(token.Comma) {
			return nil
		}

	}

	if !p.expectPeek(token.CloseCurly) {
		return nil
	}

	return hash
}
