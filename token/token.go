package token

type TokenType string

const (
	LET       TokenType = "LET"
	IDENT     TokenType = "IDENT"
	ASSIGN    TokenType = "="
	INT       TokenType = "INT"
	SEMICOLON TokenType = ";"
	COMMA     TokenType = ","
	PLUS      TokenType = "+"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	FUNCTION  TokenType = "FUNCTION"

	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
)

type Token struct {
	Type    TokenType
	Literal string
}
