package parser

import (
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/scanner"
	"github.com/bundgaard/js/token"
	"io"
	"log"
	"strconv"
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
func (p *Parser) curPrecedence() int {
	if p, ok := ast.Precedences[p.current.Type]; ok {
		return p
	}

	return ast.Lowest
}

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
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.current}
	array.Elements = p.parseExpressionList(token.CloseBracket)
	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	var list []ast.Expression

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(ast.Lowest))
	for p.peekTokenIs(token.Comma) {
		p.nextToken() // ,
		p.nextToken() // Expression
		list = append(list, p.parseExpression(ast.Lowest))
	}

	if !p.expectPeek(end) {
		return nil
	}
	return list
}
func (p *Parser) parseHashLiteral() ast.Expression {

	hash := &ast.HashLiteral{Token: p.current}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.CloseCurly) {
		p.nextToken() // eat open curly

		key := p.parseExpression(ast.Lowest)

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
