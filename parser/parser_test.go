package parser

import (
	"github.com/jatin-malik/make-thy-interpreter/ast"
	"testing"

	"github.com/jatin-malik/make-thy-interpreter/lexer"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input                    string
		expectedName             string
		expectedExpressionString string
	}{
		{"let x = -5;", "x", "( -5 )"},
		{"let x = 5;", "x", "5"},
		{"let add = x;", "add", "x"},
		{"let y = !x;", "y", "( !x )"},
		{"let x = 5+1;", "x", "( 5 + 1 )"},
		{"let x = a+b;", "x", "( a + b )"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)

		program := parser.ParseProgram()

		checkParserErrors(parser, t, tt.input)

		if len(program.Statements) != 1 {
			t.Errorf("expected %d statements, got %d\n", 1, len(program.Statements))
		}

		if stmt, ok := program.Statements[0].(*ast.LetStatement); !ok {
			t.Error("expected a let statement")
		} else {
			if stmt.Name.Value != tt.expectedName {
				t.Errorf("expected identifier name = %s, got %s", tt.expectedName, stmt.Name.Value)
			}

			if stmt.Right.String() != tt.expectedExpressionString {
				t.Errorf("expected right expression = %s, got %s", tt.expectedExpressionString, stmt.Right.String())
			}
		}
	}

}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input                    string
		expectedExpressionString string
	}{
		{"return 5;", "5"},
		{"return -5;", "( -5 )"},
		{"return result;", "result"},
		{"return -result;", "( -result )"},
		{"return !result;", "( !result )"},
		{"return 5+3;", "( 5 + 3 )"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)

		program := parser.ParseProgram()

		checkParserErrors(parser, t, tt.input)

		if len(program.Statements) != 1 {
			t.Errorf("expected %d statements, got %d\n", 1, len(program.Statements))
		}

		if stmt, ok := program.Statements[0].(*ast.ReturnStatement); !ok {
			t.Error("expected a return statement")
		} else {
			if stmt.Value.String() != tt.expectedExpressionString {
				t.Errorf("expected right expression = %s, got %s", tt.expectedExpressionString, stmt.Value.String())
			}
		}
	}

}

func TestExpressionParsing(t *testing.T) {
	tests := []struct {
		input                    string
		expectedExpressionString string
	}{
		// Basic Arithmetic Operations
		{"1 + 2", "( 1 + 2 )"},
		{"2 + 3 + 4", "( ( 2 + 3 ) + 4 )"},
		{"0 + 5", "( 0 + 5 )"},
		{"5 - 3", "( 5 - 3 )"},
		{"8 - 2 - 3", "( ( 8 - 2 ) - 3 )"},
		{"3 * 4", "( 3 * 4 )"},
		{"2 * 3 * 5", "( ( 2 * 3 ) * 5 )"},
		{"6 / 2", "( 6 / 2 )"},
		{"10 / 5 / 2", "( ( 10 / 5 ) / 2 )"},

		// Combined Arithmetic Expressions (Operator Precedence)
		{"1 + 2 * 3", "( 1 + ( 2 * 3 ) )"},
		{"3 * 2 + 1", "( ( 3 * 2 ) + 1 )"},
		{"4 + 6 / 2", "( 4 + ( 6 / 2 ) )"},
		{"1 + (2 * 3)", "( 1 + ( 2 * 3 ) )"},
		{"(1 + 2) * 3", "( ( 1 + 2 ) * 3 )"},
		{"1 + (2 + 3) * 4", "( 1 + ( ( 2 + 3 ) * 4 ) )"},
		{"(1 + 2) / (3 - 4)", "( ( 1 + 2 ) / ( 3 - 4 ) )"},

		// Unary Operators
		{"-5", "( -5 )"},
		{"- (3 + 2)", "( -( 3 + 2 ) )"},
		{"!true", "( !true )"},
		{"!false", "( !false )"},

		// Comparison Operators
		{"5 == 2", "( 5 == 2 )"},
		{"3 == 3", "( 3 == 3 )"},
		{"5 != 3", "( 5 != 3 )"},
		{"6 != 6", "( 6 != 6 )"},
		{"5 > 3", "( 5 > 3 )"},
		{"7 > 2", "( 7 > 2 )"},
		{"3 < 5", "( 3 < 5 )"},
		{"4 < 4", "( 4 < 4 )"},

		// Edge Cases
		{"42", "42"},                       // Single Operand
		{"-42", "( -42 )"},                 // Single Negative Operand
		{"()", ""},                         // Empty Parentheses
		{"( ( 1 + 2 ) )", "( 1 + 2 )"},     // Nested Parentheses
		{"( ( ( 1 + 2 ) ) )", "( 1 + 2 )"}, // Deeper Nested Parentheses

		// Long Expressions
		{"1 + 2 + 3 + 4 + 5 + 6", "( ( ( ( ( 1 + 2 ) + 3 ) + 4 ) + 5 ) + 6 )"},
		{"( 1 + 2 ) * ( 3 + 4 ) * ( 5 + 6 )", "( ( ( 1 + 2 ) * ( 3 + 4 ) ) * ( 5 + 6 ) )"},

		// Multiple Spaces and Indentation (if parser is space-sensitive)
		{"  5  +  3  ", "( 5 + 3 )"},
		{" ( 1 + 2 ) *   ( 3 + 4 ) ", "( ( 1 + 2 ) * ( 3 + 4 ) )"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)

		program := parser.ParseProgram()

		checkParserErrors(parser, t, tt.input)

		if len(program.Statements) != 1 {
			t.Errorf("expected %d statements, got %d\n", 1, len(program.Statements))
		}

		if stmt, ok := program.Statements[0].(*ast.ExpressionStatement); !ok {
			t.Error("expected an expression statement")
		} else {
			if stmt.String() != tt.expectedExpressionString {
				t.Errorf("expected expression = %s, got %s", tt.expectedExpressionString, stmt.Expr.String())
			}
		}
	}

}

func checkParserErrors(parser *Parser, t *testing.T, input string) {
	if len(parser.errors) != 0 {
		t.Logf("failed for input %s", input)
		for _, err := range parser.errors {
			t.Log("\t" + err)
		}
		t.FailNow()
	}
}
