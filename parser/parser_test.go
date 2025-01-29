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
		{`let name = "elliot";`, "name", `"elliot"`},
		{`let name = "";`, "name", `""`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
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
		})
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
		{`return "abracadbra";`, `"abracadbra"`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
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
		})
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
		{`"hello"`, `"hello"`},
		{`"hello world"`, `"hello world"`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
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
		})
	}

}

func TestFunctionLiteralParsing(t *testing.T) {
	tests := []struct {
		input                    string
		expectedExpressionString string
	}{

		// Simple function with parameters and a return statement
		{"fn (x,y) { return x+y;}", "fn (x, y) { return ( x + y ); }"},

		// Function with a simple arithmetic operation
		{"fn (x) { return x * 2; }", "fn (x) { return ( x * 2 ); }"},

		// Function with a complex expression
		{"fn (x, y) { return x*y + y - x; }",
			"fn (x, y) { return ( ( ( x * y ) + y ) - x ); }"},

		// Function with a variable assignment
		{"fn (x) { let result = x * 2; return result; }",
			"fn (x) { let result = ( x * 2 ); return result; }"},

		// Function with multiple expressions
		{"fn (x, y) { return x + y * 2 - x; }",
			"fn (x, y) { return ( ( x + ( y * 2 ) ) - x ); }"},

		// Function that returns a constant value
		{"fn () { return 42; }", "fn () { return 42; }"},

		// Function with a variable defined inside and optional return
		{"fn (x) { let y = x + 10; y; }",
			"fn (x) { let y = ( x + 10 ); y }"},

		// Simple function with an operation and return
		{"fn (x) { return x - 3; }", "fn (x) { return ( x - 3 ); }"},
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

func TestArrayLiteralParsing(t *testing.T) {
	tests := []struct {
		input                    string
		expectedExpressionString string
	}{
		// ================================
		// Basic Arrays
		// ================================
		{`["hello", 1, true]`, `["hello", 1, true]`},
		{`[]`, `[]`},
		{`[5]`, `[5]`},
		{`["only"]`, `["only"]`},
		{`[false]`, `[false]`},

		// ================================
		// Nested Arrays
		// ================================
		{`[[1, 2], [3, 4]]`, `[[1, 2], [3, 4]]`},
		{`[["a", "b"], ["c", "d"]]`, `[["a", "b"], ["c", "d"]]`},
		{`[[1, "two"], [3, "four"]]`, `[[1, "two"], [3, "four"]]`},
		{`[[]]`, `[[]]`},
		{`[[1, "hello"], [true, 42]]`, `[[1, "hello"], [true, 42]]`},

		// ================================
		// Array Index Expressions
		// ================================
		// Accessing elements using index notation
		{`["hello", 1, true][0]`, `["hello", 1, true][0]`}, // Access the first element (hello)
		{`["hello", 1, true][1]`, `["hello", 1, true][1]`}, // Access the second element (1)
		{`["hello", 1, true][2]`, `["hello", 1, true][2]`}, // Access the third element (true)

		// Accessing out-of-bounds (should result in null or similar)
		{`["hello", 1, true][3]`, `["hello", 1, true][3]`}, // Out of bounds access

		// ================================
		// Nested Array Index Expressions
		// ================================
		// Accessing elements in nested arrays using index notation
		{`[[1, 2], [3, 4]][0][0]`, `[[1, 2], [3, 4]][0][0]`}, // Access first element of first nested array (1)
		{`[[1, 2], [3, 4]][0][1]`, `[[1, 2], [3, 4]][0][1]`}, // Access second element of first nested array (2)
		{`[[1, 2], [3, 4]][1][0]`, `[[1, 2], [3, 4]][1][0]`}, // Access first element of second nested array (3)
		{`[[1, 2], [3, 4]][1][1]`, `[[1, 2], [3, 4]][1][1]`}, // Access second element of second nested array (4)

		// Accessing out-of-bounds nested arrays
		{`[[1, 2], [3, 4]][2]`, `[[1, 2], [3, 4]][2]`},       // Out of bounds (null)
		{`[[1, 2], [3, 4]][0][2]`, `[[1, 2], [3, 4]][0][2]`}, // Out of bounds in a nested array (null)
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
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
		})
	}

}

func TestIfElseConditionalParsing(t *testing.T) {
	tests := []struct {
		input                    string
		expectedExpressionString string
	}{
		// ================================
		// Regular if-else with expressions
		// ================================
		{"if (x > y) { let z = x; } else { let z = y; }",
			"if ( x > y ){ let z = x; } else { let z = y; }"},

		{"if (x > y) { let x = 10; } else { let y = 10; }",
			"if ( x > y ){ let x = 10; } else { let y = 10; }"},

		{"if (x > y) { if (y > 0) { let z = 5; } else { let z = 10; } } else { let z = 0; }",
			"if ( x > y ){ if ( y > 0 ){ let z = 5; } else { let z = 10; } } else { let z = 0; }"},

		{"if (x > y) { let x = x + 1; let y = y + 1; } else { let x = x - 1; let y = y - 1; }",
			"if ( x > y ){ let x = ( x + 1 ); let y = ( y + 1 ); } else { let x = ( x - 1 ); let y = ( y - 1 ); }"},

		// ================================
		// Edge cases: Missing else block
		// ================================
		{"if (x > y) { let z = x; }",
			"if ( x > y ){ let z = x; }"},

		{"if (x > y) { }",
			"if ( x > y ){ }"},

		{"if (x > y) { let z = x + y; }",
			"if ( x > y ){ let z = ( x + y ); }"},

		// ================================
		// Only if block with return statement
		// ================================
		{"if (x > y) { return true; }",
			"if ( x > y ){ return true; }"},

		// ================================
		// Complex conditions and operations
		// ================================
		{"if (x * 2 > y + 10) { let z = x; } else { let z = y; }",
			"if ( ( x * 2 ) > ( y + 10 ) ){ let z = x; } else { let z = y; }"},

		{"if (x + 10 > y) { let z = x * 2; }",
			"if ( ( x + 10 ) > y ){ let z = ( x * 2 ); }"},
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

func TestCallExpressionParsing(t *testing.T) {
	tests := []struct {
		input                    string
		expectedExpressionString string
	}{
		// ================================
		// Basic function call
		// ================================
		{"add(2, 3)", "add(2, 3)"},
		{"subtract(10, 5)", "subtract(10, 5)"},
		{"sqrt(25)", "sqrt(25)"},

		// ================================
		// Function calls with expressions as arguments
		// ================================
		{"add(2 + 3, 4 * 5)", "add(( 2 + 3 ), ( 4 * 5 ))"},
		{"max(x, y + z)", "max(x, ( y + z ))"},

		// ================================
		// Function call with nested function calls
		// ================================
		{"print(add(2, 3))", "print(add(2, 3))"},
		{"sqrt(add(1, 2))", "sqrt(add(1, 2))"},

		// ================================
		// Edge case: Function with no arguments
		// ================================
		{"noop() + noop()", "( noop() + noop() )"},

		// ================================
		// Edge case: Function with many arguments
		// ================================
		{"longFunctionName(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)",
			"longFunctionName(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)"},

		// ================================
		// Edge case: Multiple nested function calls
		// ================================
		{"outer(inner(1, 2), inner(3, 4))", "outer(inner(1, 2), inner(3, 4))"},

		// ================================
		// Edge case: Nested functions with expressions
		// ================================
		{"outer(inner(1 + 2, 3 * 4), 5)", "outer(inner(( 1 + 2 ), ( 3 * 4 )), 5)"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
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
		})

	}

}

func checkParserErrors(parser *Parser, t *testing.T, input string) {
	if len(parser.Errors) != 0 {
		t.Logf("failed for input %s", input)
		for _, err := range parser.Errors {
			t.Log("\t" + err)
		}
		t.FailNow()
	}
}
