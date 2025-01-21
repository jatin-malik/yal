package parser

import (
	"github.com/jatin-malik/make-thy-interpreter/ast"
	"testing"

	"github.com/jatin-malik/make-thy-interpreter/lexer"
)

func TestLetStatement(t *testing.T) {
	input := "let x = 5;"
	l := lexer.New(input)
	parser := New(l)

	program := parser.ParseProgram()

	checkParserErrors(parser, t)

	if len(program.Statements) != 1 {
		t.Errorf("expected %d statements, got %d\n", 1, len(program.Statements))
	}

	if stmt, ok := program.Statements[0].(*ast.LetStatement); !ok {
		t.Error("expected a let statement")
	} else {
		if stmt.Name.Value != "x" {
			t.Errorf("expected identifier name = %s, got %s", "x", stmt.Name.Value)
		}

		if il, ok := stmt.Right.(*ast.IntegerLiteral); !ok {
			t.Error("expected an integer literal")
		} else {
			if il.Value != 5 {
				t.Errorf("expected integer value = %d, got %d", 5, il.Value)
			}
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := "return 5;"
	l := lexer.New(input)
	parser := New(l)

	program := parser.ParseProgram()

	checkParserErrors(parser, t)

	if len(program.Statements) != 1 {
		t.Errorf("expected %d statements, got %d\n", 1, len(program.Statements))
	}

	if stmt, ok := program.Statements[0].(*ast.ReturnStatement); !ok {
		t.Error("expected a return statement")
	} else {
		if il, ok := stmt.Value.(*ast.IntegerLiteral); !ok {
			t.Error("expected an integer literal")
		} else {
			if il.Value != 5 {
				t.Errorf("expected integer value = %d, got %d", 5, il.Value)
			}
		}
	}
}

func checkParserErrors(parser *Parser, t *testing.T) {
	if len(parser.errors) != 0 {
		for _, err := range parser.errors {
			t.Log(err)
		}
		t.FailNow()
	}
}
