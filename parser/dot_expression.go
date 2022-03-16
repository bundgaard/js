package parser

import (
	"github.com/bundgaard/js/ast"
	"log"
)

func (p *Parser) parseDotExpression() ast.Expression {

	log.Printf("parseDotExpression %v %v", p.current, p.next)
	return nil
}
