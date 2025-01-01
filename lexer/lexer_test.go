package lexer_test

import (
	"testing"

	"github.com/jatin-malik/make-thy-interpreter/lexer"
	"github.com/jatin-malik/make-thy-interpreter/token"
)

func TestLexer(t *testing.T) {

	t.Run("single line input", func(t *testing.T) {

		input := "+=(){};"
		l := lexer.New(input)

		tests := []struct {
			expectedTokenType token.TokenType
			expectedLiteral   string
		}{
			{token.PLUS, "+"},
			{token.ASSIGN, "="},
			{token.LPAREN, "("},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RBRACE, "}"},
			{token.SEMICOLON, ";"},
			{token.EOF, string(byte(0))},
		}

		for _, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedTokenType {
				t.Errorf("expected %q, got %q", tt.expectedTokenType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("expected %q, got %q", tt.expectedLiteral, tok.Literal)
			}
		}
	})

	t.Run("multi line input", func(t *testing.T) {

		input := `let five=5;
			let ten = 10;
			let add = fn(x,y){
				x+y;
			}
			let result = add(five,ten);`

		l := lexer.New(input)

		tests := []struct {
			expectedTokenType token.TokenType
			expectedLiteral   string
		}{
			{token.LET, "let"},
			{token.IDENT, "five"},
			{token.ASSIGN, "="},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "ten"},
			{token.ASSIGN, "="},
			{token.INT, "10"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "add"},
			{token.ASSIGN, "="},
			{token.FUNCTION, "fn"},
			{token.LPAREN, "("},
			{token.IDENT, "x"},
			{token.COMMA, ","},
			{token.IDENT, "y"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.IDENT, "x"},
			{token.PLUS, "+"},
			{token.IDENT, "y"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.LET, "let"},
			{token.IDENT, "result"},
			{token.ASSIGN, "="},
			{token.IDENT, "add"},
			{token.LPAREN, "("},
			{token.IDENT, "five"},
			{token.COMMA, ","},
			{token.IDENT, "ten"},
			{token.RPAREN, ")"},
			{token.SEMICOLON, ";"},
			{token.EOF, string(byte(0))},
		}

		for _, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedTokenType {
				t.Errorf("expected %q, got %q", tt.expectedTokenType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("expected %q, got %q", tt.expectedLiteral, tok.Literal)
			}
		}
	})

}
