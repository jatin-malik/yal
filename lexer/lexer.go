package lexer

import "github.com/jatin-malik/yal/token"

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
	l.eatWhiteSpace() // whitespaces are just token separators for us

	if l.pos >= len(l.input) {
		return newToken(token.EOF, 0)
	}

	ch := l.input[l.pos]

	var tok token.Token

	switch ch {
	case '"':
		tok.Literal = l.readString()
		tok.Type = token.STRING
	case '+':
		tok = newToken(token.PLUS, ch)
	case ':':
		tok = newToken(token.COLON, ch)
	case '-':
		tok = newToken(token.MINUS, ch)
	case '/':
		tok = newToken(token.SLASH, ch)
	case '*':
		tok = newToken(token.ASTERISK, ch)
	case '!':
		nextCh := l.peekNextChar()
		if nextCh == '=' {
			l.pos++
			tok.Type = token.NEQ
			tok.Literal = "!="
		} else {
			tok = newToken(token.BANG, ch)
		}
	case '<':
		tok = newToken(token.LT, ch)
	case '>':
		tok = newToken(token.GT, ch)
	case '=':
		nextCh := l.peekNextChar()
		if nextCh == '=' {
			l.pos++
			tok.Type = token.EQ
			tok.Literal = "=="
		} else {
			tok = newToken(token.ASSIGN, ch)
		}
	case '(':
		tok = newToken(token.LPAREN, ch)
	case ')':
		tok = newToken(token.RPAREN, ch)
	case '[':
		tok = newToken(token.LBRACKET, ch)
	case ']':
		tok = newToken(token.RBRACKET, ch)
	case '{':
		tok = newToken(token.LBRACE, ch)
	case '}':
		tok = newToken(token.RBRACE, ch)
	case ';':
		tok = newToken(token.SEMICOLON, ch)
	case ',':
		tok = newToken(token.COMMA, ch)
	default:
		if isLetter(ch) {
			tok.Literal = l.readIdent()
			tok.Type = token.GetTokenFromName(tok.Literal)
			return tok
		} else if isDigit(ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, ch)
		}
	}

	l.pos++
	return tok

}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// isLetter defines the allowed characters in the language identifiers.
func isLetter(ch byte) bool {
	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch == '_') {
		return true
	}
	return false
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) readIdent() string {
	startingPos := l.pos
	for l.pos < len(l.input) && isLetter(l.input[l.pos]) {
		l.pos++
	}
	return l.input[startingPos:l.pos]
}

func (l *Lexer) readNumber() string {
	startingPos := l.pos
	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
	}
	return l.input[startingPos:l.pos]
}

func (l *Lexer) readString() string {
	l.pos++ // move on from starting quote literal
	startingPos := l.pos
	// TODO: throw error if string is unbounded and EOF comes before closing quote?
	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		l.pos++
	}
	return l.input[startingPos:l.pos]
}

func (l *Lexer) eatWhiteSpace() {
	for l.pos < len(l.input) && isWhiteSpace(l.input[l.pos]) {
		l.pos++
	}
}

func isWhiteSpace(ch byte) bool {
	if (ch == ' ') || (ch == '\t') || (ch == '\n') || (ch == '\r') {
		return true
	}
	return false
}

func (l *Lexer) peekNextChar() byte {
	lookupIdx := l.pos + 1
	if lookupIdx >= len(l.input) {
		return 0
	}
	return l.input[lookupIdx]
}
