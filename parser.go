package js

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

type Parser struct {
	s       *Scanner
	current *Token
	next    *Token

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
func (p *Parser) peekPrecedence() int {
	if v, ok := precedences[p.next.Type]; ok {
		return v
	}
	return Lowest
}
func NewParser(s *Scanner) *Parser {
	p := &Parser{
		s: s,
	}

	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(Ident, p.parseName)
	p.registerPrefix(String, p.parseStringLiteral)
	p.registerPrefix(Number, p.parseNumberLiteral)
	p.registerPrefix(OpenCurly, p.parseHashLiteral)
	p.registerPrefix(OpenBracket, p.parseArrayLiteral)

	p.registerPrefix(Dot, p.parseDotExpression)

	p.registerInfix(Add, p.parseInfixExpression)
	p.registerInfix(Assign, p.parseInfixExpression)
	p.registerInfix(OpenBracket, p.parseIndexExpression)
	p.nextToken()
	p.nextToken()
	return p
}
func (p *Parser) parseDotExpression() Expression {

	log.Printf("parseDotExpression %v %v", p.current, p.next)
	return nil
}
func (p *Parser) parseIndexExpression(left Expression) Expression {
	exp := &IndexExpression{Token: p.current, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(Lowest)
	if !p.expectPeek(CloseBracket) {
		return nil
	}

	return exp
}
func (p *Parser) parseNumberLiteral() Expression {
	n, err := strconv.Atoi(p.current.Value)
	if err != nil {
		log.Println("error number conversion", err)
	}
	return &NumberLiteral{Token: p.current, Value: int64(n)}
}
func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.current, Value: p.current.Value}
}
func (p *Parser) parseName() Expression {
	return &Identifier{Token: p.current, Value: p.current.Value}
}

func (p *Parser) nextToken() {
	p.current = p.next
	p.next = p.s.NextToken()
}

func (p *Parser) Parse() *Program {

	program := &Program{}
	program.Statements = []Statement{}

	for p.current.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.current.Type {
	case Var:
		return p.parseVariable()
	case CommentLine:
		return nil
	case CommentBlock:
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {

	// expression statement
	// foo[0] = 1300
	stmt := &ExpressionStatement{Token: p.current}
	stmt.Expression = p.parseExpression(Lowest)
	// dot expression
	// foo.bar
	if p.peekTokenIs(Dot) {
		p.nextToken()
		right := p.parseInfixExpression(stmt.Expression)
		stmt.Expression = right
	}

	if p.peekTokenIs(Assign) {

		p.nextToken()

		right := p.parseInfixExpression(stmt.Expression)
		stmt.Expression = right
	}

	if p.peekTokenIs(Semi) {
		p.nextToken()
	}
	return stmt
}
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.current.Type]; ok {
		return p
	}

	return Lowest
}
func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.current,
		Operator: p.current.Value,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseVariable() *VariableStatement {
	stmt := &VariableStatement{Token: p.current}
	if !p.expectPeek(Ident) {
		return nil
	}
	stmt.Name = &Identifier{Token: p.current, Value: p.current.Value}
	if !p.expectPeek(Assign) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(Lowest)
	if p.peekTokenIs(Semi) {
		p.nextToken()
	}
	return stmt
}
func (p *Parser) currentTokenIs(tokenType TokenType) bool {
	return p.current.Type == tokenType
}
func (p *Parser) peekTokenIs(tokenType TokenType) bool {
	return p.next.Type == tokenType
}

type AtomicCounter struct {
	sync.Mutex
	count int64
}

func (ac *AtomicCounter) Add() {
	ac.Lock()
	defer ac.Unlock()
	ac.count++
}

func (ac *AtomicCounter) Get() int64 {
	ac.Lock()
	defer ac.Unlock()
	return ac.count
}

var counter = AtomicCounter{count: 0}

func (p *Parser) parseExpression(priority int) Expression {

	if p.currentTokenIs(CommentLine) || p.currentTokenIs(CommentBlock) {
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
			fmt.Fprintf(bugFile, "%v; error at (%d,%d)", &p.s.buf, p.s.line, p.s.column)
			os.Exit(1)
		}
		log.Printf("no prefix for %v (%d,%d)", p.current, p.s.line, p.s.column)
		counter.Add()
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(Semi) && priority < p.peekPrecedence() {
		infix := p.infixParseFns[p.next.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()

		leftExp = infix(leftExp)
	}
	return leftExp
}
func (p *Parser) expectPeek(tokenType TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	}
	return false
}
func (p *Parser) parseArrayLiteral() Expression {
	array := &ArrayLiteral{Token: p.current}
	array.Elements = p.parseExpressionList(CloseBracket)
	return array
}

func (p *Parser) parseExpressionList(end TokenType) []Expression {
	list := []Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(Lowest))
	for p.peekTokenIs(Comma) {
		p.nextToken() // ,
		p.nextToken() // Expression
		list = append(list, p.parseExpression(Lowest))
	}

	if !p.expectPeek(end) {
		return nil
	}
	return list
}
func (p *Parser) parseHashLiteral() Expression {

	hash := &HashLiteral{Token: p.current}
	hash.Pairs = make(map[Expression]Expression)

	for !p.peekTokenIs(CloseCurly) {
		p.nextToken() // eat open curly

		key := p.parseExpression(Lowest)

		if !p.expectPeek(Colon) {
			return nil
		}

		p.nextToken() // EAT Colon

		value := p.parseExpression(Lowest)

		hash.Pairs[key] = value
		if !p.peekTokenIs(CloseCurly) && !p.expectPeek(Comma) {
			return nil
		}

	}

	if !p.expectPeek(CloseCurly) {
		return nil
	}

	return hash
}
