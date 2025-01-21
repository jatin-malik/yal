package parser

import (
	"fmt"
	"github.com/jatin-malik/make-thy-interpreter/ast"
	"github.com/jatin-malik/make-thy-interpreter/lexer"
	"github.com/jatin-malik/make-thy-interpreter/token"
	"strconv"
)

type prefixParsingFunction func() ast.Expression
type infixParsingFunction func(ast.Expression) ast.Expression

type Parser struct {
	lexer         *lexer.Lexer
	curToken      token.Token
	peekToken     token.Token
	errors        []string
	prefixParsers map[token.TokenType]prefixParsingFunction
	infixParsers  map[token.TokenType]infixParsingFunction
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer}
	parser.curToken = lexer.NextToken()
	parser.peekToken = lexer.NextToken()
	parser.errors = []string{}
	parser.prefixParsers = make(map[token.TokenType]prefixParsingFunction)
	parser.infixParsers = make(map[token.TokenType]infixParsingFunction)

	// Register parsing functions for tokens here
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)

	return parser
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParsingFunction) {
	p.prefixParsers[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParsingFunction) {
	p.infixParsers[tokenType] = fn
}

func (p *Parser) Next() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// ParseProgram is the top-level function to parse a program.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	var statements []ast.Statement

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, *stmt)
		}
		p.Next()
	}

	program.Statements = statements
	return program
}

func (p *Parser) parseStatement() *ast.Statement {
	var stmt ast.Statement
	switch p.curToken.Type {
	case token.LET:
		stmt = p.parseLetStatement()
	case token.RETURN:
		stmt = p.parseReturnStatement()
	default:
		stmt = nil
	}
	return &stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.curToken,
	}

	p.Next()
	name := p.parseExpression()
	if ident, ok := name.(*ast.Identifier); ok {
		stmt.Name = ident
	} else {
		return nil
	}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.Next()
	stmt.Right = p.parseExpression()
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}

	p.Next()
	stmt.Value = p.parseExpression()
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	return stmt
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return ident
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	exp := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("cannot parse %q as integer", p.curToken.Literal))
		return nil
	}
	exp.Value = value
	return exp
}

func (p *Parser) parseExpression() ast.Expression {
	var exp ast.Expression
	if prefixParser, ok := p.prefixParsers[p.curToken.Type]; ok {
		exp = prefixParser()
	}
	return exp
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekToken.Type == tokenType {
		p.Next()
		return true
	} else {
		errMsg := fmt.Sprintf("expected token %s, got %s", tokenType, p.peekToken.Type)
		p.errors = append(p.errors, errMsg)
		return false
	}
}
