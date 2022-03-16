package parser

import (
	"fmt"
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/token"
	"log"
	"os"
	"time"
)

func (p *Parser) parseExpression(priority int) ast.Expression {

	if p.currentTokenIs(token.CommentLine) || p.currentTokenIs(token.CommentBlock) {
		p.nextToken()
	}
	prefix := p.prefixParseFns[p.current.Type]
	if prefix == nil {
		if counter.Get() == 10 {

			filename := time.Now().Format("20060102T15") + ".debug"
			log.Println("something weird happened, saved to file", filename)
			bugFile, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer bugFile.Close()
			fmt.Fprintf(bugFile, "%v; error at (%d,%d)", &p.s.Buf, p.s.Line, p.s.Column)
			os.Exit(1)
		}
		log.Printf("no prefix for %v (%d,%d)", p.current, p.s.Line, p.s.Column)
		counter.Add()
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.Semi) && priority < p.peekPrecedence() {
		infix := p.infixParseFns[p.next.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()

		leftExp = infix(leftExp)
	}
	return leftExp
}
