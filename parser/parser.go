package parser

import (
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/scanner"
	"github.com/bundgaard/js/token"
	"io"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	s       *scanner.Scanner
	current *token.Token
	next    *token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) peekPrecedence() int {
	if v, ok := ast.Precedences[p.next.Type]; ok {
		return v
	}
	return ast.Lowest
}
func NewParser(rd io.RuneReader) *Parser {
	p := &Parser{
		s: scanner.New(rd),
	}

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.Ident, p.parseName)
	p.registerPrefix(token.String, p.parseStringLiteral)
	p.registerPrefix(token.Number, p.parseNumberLiteral)
	p.registerPrefix(token.OpenCurly, p.parseHashLiteral)
	p.registerPrefix(token.OpenBracket, p.parseArrayLiteral)

	p.registerPrefix(token.Dot, p.parseDotExpression)

	p.registerInfix(token.Add, p.parseInfixExpression)
	p.registerInfix(token.Mul, p.parseInfixExpression)
	p.registerInfix(token.Div, p.parseInfixExpression)
	p.registerInfix(token.Sub, p.parseInfixExpression)

	p.registerInfix(token.Assign, p.parseInfixExpression)
	p.registerInfix(token.OpenBracket, p.parseIndexExpression)
	p.registerInfix(token.OpenParen, p.parseCallExpression)
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.current = p.next
	p.next = p.s.NextToken()
}

func (p *Parser) Parse() *ast.Program {

	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.current.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}
	return program
}

func (p *Parser) curPrecedence() int {
	if p, ok := ast.Precedences[p.current.Type]; ok {
		return p
	}

	return ast.Lowest
}

func (p *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return p.current.Type == tokenType
}
func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.next.Type == tokenType
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	}
	return false
}
