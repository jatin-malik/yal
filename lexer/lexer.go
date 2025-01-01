package lexer

import "github.com/jatin-malik/make-thy-interpreter/token"

type Lexer struct {
	input string // the input to the lexer i.e the source code
	pos   int    // the current position to read from
}

func New(input string) *Lexer {
	return &Lexer{
		input: input,
	}
}

func (l *Lexer) NextToken() token.Token {
	if l.pos >= len(l.input) {
		return newToken(token.EOF, 0)
	}

	ch := l.input[l.pos]

	var tok token.Token

	switch ch {
	case '+':
		tok = newToken(token.PLUS, ch)
	case '=':
		tok = newToken(token.ASSIGN, ch)
	case '(':
		tok = newToken(token.LPAREN, ch)
	case ')':
		tok = newToken(token.RPAREN, ch)
	case '{':
		tok = newToken(token.LBRACE, ch)
	case '}':
		tok = newToken(token.RBRACE, ch)
	case ';':
		tok = newToken(token.SEMICOLON, ch)
	default:
		tok = newToken(token.ILLEGAL, ch)
	}

	l.pos++
	return tok

}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
