package parser

import (
	"fmt"
	"github.com/jatin-malik/yal/ast"
	"github.com/jatin-malik/yal/lexer"
	"github.com/jatin-malik/yal/token"
	"strconv"
)

type prefixParsingFunction func() ast.Expression
type infixParsingFunction func(ast.Expression) ast.Expression

type Parser struct {
	lexer         *lexer.Lexer
	curToken      token.Token
	peekToken     token.Token
	Errors        []string
	prefixParsers map[token.TokenType]prefixParsingFunction
	infixParsers  map[token.TokenType]infixParsingFunction
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer}
	parser.curToken = lexer.NextToken()
	parser.peekToken = lexer.NextToken()
	parser.Errors = []string{}
	parser.prefixParsers = make(map[token.TokenType]prefixParsingFunction)
	parser.infixParsers = make(map[token.TokenType]infixParsingFunction)

	// Register parsing functions for tokens here
	// Prefix parsers
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.TRUE, parser.parseBooleanLiteral)
	parser.registerPrefix(token.FALSE, parser.parseBooleanLiteral)
	parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiteral)
	parser.registerPrefix(token.MACRO, parser.parseMacroLiteral)
	parser.registerPrefix(token.IF, parser.parseIfElseConditional)
	parser.registerPrefix(token.STRING, parser.parseStringLiteral)
	parser.registerPrefix(token.LBRACKET, parser.parseArrayLiteral)
	parser.registerPrefix(token.LBRACE, parser.parseHashLiteral)

	// Infix parsers
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.EQ, parser.parseInfixExpression)
	parser.registerInfix(token.NEQ, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)
	parser.registerInfix(token.LPAREN, parser.parseCallExpression)
	parser.registerInfix(token.LBRACKET, parser.parseIndexExpression)

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
		stmt = p.parseExpressionStatement()
	}
	return &stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.curToken,
	}

	p.Next()
	name := p.parseExpression(LowestPrecedence)
	if ident, ok := name.(*ast.Identifier); ok {
		stmt.Name = ident
	} else {
		return nil
	}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.Next()
	stmt.Right = p.parseExpression(LowestPrecedence)
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
	stmt.Value = p.parseExpression(LowestPrecedence)
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.curToken,
	}

	stmt.Expr = p.parseExpression(LowestPrecedence)
	// Semicolon is optional for an expression statement for convenience in REPL
	if p.peekToken.Type == token.SEMICOLON {
		p.Next()
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
		p.Errors = append(p.Errors, fmt.Sprintf("cannot parse %q as integer", p.curToken.Literal))
		return nil
	}
	exp.Value = value
	return exp
}

func (p *Parser) parseStringLiteral() ast.Expression {
	exp := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	return exp
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	exp := &ast.BooleanLiteral{Token: p.curToken}
	exp.Value = p.curToken.Literal == "true"
	return exp
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	exp := &ast.FunctionLiteral{
		Token: p.curToken,
	}

	p.Next()
	exp.Parameters = p.parseFunctionParams()
	p.Next()
	exp.Body = p.parseBlockStatement()
	return exp
}

func (p *Parser) parseMacroLiteral() ast.Expression {
	exp := &ast.MacroLiteral{
		Token: p.curToken,
	}

	p.Next()
	exp.Parameters = p.parseFunctionParams()
	p.Next()
	exp.Body = p.parseBlockStatement()
	return exp
}

func (p *Parser) parseIfElseConditional() ast.Expression {
	exp := &ast.IfElseConditional{
		Token: p.curToken,
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.Next()
	exp.Condition = p.parseExpression(LowestPrecedence)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	p.Next()
	exp.Consequence = p.parseBlockStatement()
	if p.peekToken.Type == token.ELSE {
		p.Next()
		p.Next()
		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseFunctionParams() []*ast.Identifier {
	if p.curToken.Type != token.LPAREN {
		return nil
	}

	p.Next()
	var identifiers []*ast.Identifier
	for p.curToken.Type != token.RPAREN {
		ident := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		identifiers = append(identifiers, ident)
		p.Next()
		if p.curToken.Type == token.COMMA {
			p.Next()
		}
	}
	return identifiers
}

// parseCommaSeparatedExpressions parses comma separated expressions and stops parsing when encounters endToken.
func (p *Parser) parseCommaSeparatedExpressions(endToken token.TokenType) []ast.Expression {
	p.Next()
	var arguments []ast.Expression
	for p.curToken.Type != endToken {
		arg := p.parseExpression(LowestPrecedence)
		arguments = append(arguments, arg)
		p.Next()
		if p.curToken.Type == token.COMMA {
			p.Next()
		}
	}
	return arguments
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	if p.curToken.Type != token.LBRACE {
		return nil
	}
	p.Next()
	program := &ast.BlockStatement{
		Token: p.curToken,
	}
	var statements []ast.Statement

	for p.curToken.Type != token.EOF && p.curToken.Type != token.RBRACE {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, *stmt)
		}
		p.Next()
	}

	if p.curToken.Type != token.RBRACE {
		return nil
	}

	program.Statements = statements
	return program
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	var leftExp ast.Expression
	if prefixParser, ok := p.prefixParsers[p.curToken.Type]; ok {
		leftExp = prefixParser()
	} else {
		p.Errors = append(p.Errors, fmt.Sprintf("no prefix parsing function registered for %s", p.curToken))
		return leftExp
	}

	for p.peekToken.Type != token.SEMICOLON && getTokenPrecedence(p.peekToken.Type) > precedence {
		if infixParser, exists := p.infixParsers[p.peekToken.Type]; !exists {
			return leftExp
		} else {
			p.Next()
			leftExp = infixParser(leftExp)
		}

	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	pe := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	curPrecedence := getTokenPrecedence(p.curToken.Type)
	p.Next()
	pe.Right = p.parseExpression(curPrecedence)
	return pe
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.Next()
	if p.curToken.Type == token.RPAREN {
		return nil
	}
	exp := p.parseExpression(LowestPrecedence)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	ce := &ast.CallExpression{
		Token:    p.curToken,
		Function: left,
	}
	ce.Arguments = p.parseCommaSeparatedExpressions(token.RPAREN)
	return ce
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	ie := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}
	p.Next()
	ie.Index = p.parseExpression(LowestPrecedence)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return ie
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	al := &ast.ArrayLiteral{
		Token: p.curToken,
	}
	al.Elements = p.parseCommaSeparatedExpressions(token.RBRACKET)
	return al
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hl := &ast.HashLiteral{
		Token: p.curToken,
	}

	p.Next()
	pairs := make(map[ast.Expression]ast.Expression)
	for p.curToken.Type != token.EOF && p.curToken.Type != token.RBRACE {
		key := p.parseExpression(LowestPrecedence)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.Next()
		val := p.parseExpression(LowestPrecedence)
		pairs[key] = val
		p.Next()
		if p.curToken.Type == token.COMMA {
			p.Next()
		}
	}

	hl.Pairs = pairs
	return hl
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	ie := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	curPrecedence := getTokenPrecedence(p.curToken.Type)
	p.Next()
	ie.Right = p.parseExpression(curPrecedence)
	return ie
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekToken.Type == tokenType {
		p.Next()
		return true
	} else {
		errMsg := fmt.Sprintf("expected token %s, got %s", tokenType, p.peekToken)
		p.Errors = append(p.Errors, errMsg)
		return false
	}
}
