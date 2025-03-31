package token

type TokenType string

const (
	// Single char tokens
	ASSIGN    TokenType = "="
	SEMICOLON TokenType = ";"
	COMMA     TokenType = ","
	PLUS      TokenType = "+"
	MINUS     TokenType = "-"
	SLASH     TokenType = "/"
	ASTERISK  TokenType = "*"
	BANG      TokenType = "!"
	LT        TokenType = "<"
	GT        TokenType = ">"
	LPAREN    TokenType = "("
	LBRACKET  TokenType = "["
	RBRACKET  TokenType = "]"
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	COLON     TokenType = ":"

	// Double char tokens
	EQ  TokenType = "=="
	NEQ TokenType = "!="

	// Keywords
	LET      TokenType = "LET"
	INT      TokenType = "INT"
	STRING   TokenType = "STRING"
	FUNCTION TokenType = "FUNCTION"
	MACRO    TokenType = "MACRO"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	RETURN   TokenType = "RETURN"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	LOOP     TokenType = "LOOP"

	// Others
	IDENT   TokenType = "IDENT"
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"macro":  MACRO,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"loop":   LOOP,
}

func GetTokenFromName(name string) TokenType {
	if token, ok := keywords[name]; ok {
		return token
	}
	return IDENT
}
