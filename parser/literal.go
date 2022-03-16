package parser

import (
	"github.com/bundgaard/js/ast"
	"log"
	"strconv"
)

func (p *Parser) parseNumberLiteral() ast.Expression {
	n, err := strconv.Atoi(p.current.Value)
	if err != nil {
		log.Println("error number conversion", err)
	}
	return &ast.NumberLiteral{Token: p.current, Value: int64(n)}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.current, Value: p.current.Value}
}
func (p *Parser) parseName() ast.Expression {
	return &ast.Identifier{Token: p.current, Value: p.current.Value}
}
