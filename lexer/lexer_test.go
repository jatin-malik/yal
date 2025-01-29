package lexer_test

import (
	"testing"

	"github.com/jatin-malik/yal/lexer"
	"github.com/jatin-malik/yal/token"
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
			let result = add(five,ten);
			-/*!<>50
			
			if (5<10){
				return true
			}else{
				return false
			}
			5==5
			5!=10
			"hello"
			"hello world"
			[100,"hello",true]`

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
			{token.MINUS, "-"},
			{token.SLASH, "/"},
			{token.ASTERISK, "*"},
			{token.BANG, "!"},
			{token.LT, "<"},
			{token.GT, ">"},
			{token.INT, "50"},

			{token.IF, "if"},
			{token.LPAREN, "("},
			{token.INT, "5"},
			{token.LT, "<"},
			{token.INT, "10"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.TRUE, "true"},
			{token.RBRACE, "}"},
			{token.ELSE, "else"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.FALSE, "false"},
			{token.RBRACE, "}"},
			{token.INT, "5"},
			{token.EQ, "=="},
			{token.INT, "5"},
			{token.INT, "5"},
			{token.NEQ, "!="},
			{token.INT, "10"},
			{token.STRING, "hello"},
			{token.STRING, "hello world"},
			{token.LBRACKET, "["},
			{token.INT, "100"},
			{token.COMMA, ","},
			{token.STRING, "hello"},
			{token.COMMA, ","},
			{token.TRUE, "true"},
			{token.RBRACKET, "]"},
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
