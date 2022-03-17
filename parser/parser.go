package parser

import (
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/scanner"
	"github.com/bundgaard/js/token"
	"io"
	"strings"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	s       *scanner.Scanner
	current *token.Token
	next    *token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func New(rd io.RuneReader) *Parser {
	p := &Parser{
		s: scanner.New(rd),
	}

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.Ident, p.parseName)
	p.registerPrefix(token.String, p.parseStringLiteral)
	p.registerPrefix(token.Number, p.parseNumberLiteral)
	p.registerPrefix(token.OpenCurly, p.parseHashLiteral)
	p.registerPrefix(token.OpenBracket, p.parseArrayLiteral)
	p.registerPrefix(token.Function, p.parseFunctionLiteral)

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

func NewString(data string) *Parser {
	return New(strings.NewReader(data))
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

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.current}
	// name
	p.nextToken()
	fn.Name = p.current.Value
	if !p.expectPeek(token.OpenParen) {
		return nil
	}
	fn.Parameters = p.parseFunctionArguments()
	if !p.expectPeek(token.OpenCurly) {
		return nil
	}
	fn.Body = p.parseBlockStatement()
	return fn
}

func (p *Parser) parseFunctionArguments() []*ast.Identifier {
	var identifiers []*ast.Identifier

	if p.peekTokenIs(token.CloseParen) {
		p.nextToken()

		return identifiers
	}
	p.nextToken() // eat OpenParen

	ident := &ast.Identifier{Token: p.current, Value: p.current.Value}
	identifiers = append(identifiers, ident)
	for p.peekTokenIs(token.Comma) {
		p.nextToken() // Eat Comma

		p.nextToken() // become thing after comma

		ident := &ast.Identifier{Token: p.current, Value: p.current.Value}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.CloseParen) {
		return nil
	}
	return identifiers
}
