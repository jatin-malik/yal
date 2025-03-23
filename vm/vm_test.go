package vm

import (
	"github.com/jatin-malik/yal/compiler"
	"github.com/jatin-malik/yal/lexer"
	"github.com/jatin-malik/yal/parser"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1+2", "3"},
		{"6-2", "4"},
		{"3*4", "12"},
		{"6/3", "2"},

		// Deeply nested expressions
		{"(((1+2)*3)-4)/2", "2"},
		{"(5*(3+(2*2)))", "35"},
		{"(6/(2*(1+2)))", "1"},
		{"((2+3)*(4+(5-2)))", "35"},

		// Edge cases
		{"(((1)))", "1"}, // Extra parentheses should have no effect
		{"(0+((1+2)*3))", "9"},
		{"(100/(10/(2*5)))", "100"},
		{"((8-6)*(3+(4/2)))", "10"},

		// Comparison Operators
		{"1 == 1", "true"},
		{"1 == 2", "false"},
		{"5 != 3", "true"},
		{"4 != 4", "false"},
		{"10 > 5", "true"},
		{"10 > 10", "false"},
		{"3 < 7", "true"},
		{"7 < 3", "false"},

		// Nested Comparisons
		{"(1+2) == (3)", "true"},
		{"(10-5) > (2+2)", "true"},
		{"(2*3) < (10-1)", "true"},
		{"(4/2) != (2-1)", "true"},
		{"(6/3) == (2-0)", "true"},
		{"((2+3)*2) > ((4+1)*2)", "false"},
		{"((5*2)/2) == ((4+1))", "true"},

		// Prefix Operators
		{"!true", "false"},
		{"!false", "true"},
		{"!1", "false"}, // Assuming truthy values like 1 evaluate to true
		{"!0", "false"},
		{"!!true", "true"},
		{"!!false", "false"},
		{"!!1", "true"},
		{"!!0", "true"},

		// Prefix Negative
		{"-5", "-5"},
		{"-(-5)", "5"},
		{"-(3+2)", "-5"},
		{"-1 + 2", "1"}, // Ensure correct order of evaluation
		{"-(2*3)", "-6"},
		{"-(10-5)", "-5"},
		{"-(4/2)", "-2"},
		{"-(-(-3))", "-3"},

		{`if (5 > 3) { 10 } else { 20 }`, "10"},
		{`if (5 > 7) { 10 } else { 20 }`, "20"},
		{`if (5 > 3) { 10 } else { 20 };5`, "5"},
		{`if (5 > 3) { 10 }`, "10"},
		{`if (5 > 3) { 10 };6+1`, "7"},
		{`if (5 > 8) { 10 };2+1`, "3"},
		{`if (5 > 8) { 10 }`, "null"},

		{`let x = 5 ; x`, "5"},
		{`let x = 5 ; x+2`, "7"},
		{`let x = 5 ; let x = 10; x + 4`, "14"},

		// String literals
		{`"hello"`, "hello"},
		{`"world"`, "world"},

		// String concatenation
		{`"hello" + " " + "world"`, "hello world"},
		{`"foo" + "bar"`, "foobar"},
		{`"Go" + "lang"`, "Golang"},

		// Array literals
		{`[1, 2, 3]`, "[1, 2, 3]"},
		{`[10, 20, 30][1]`, "20"},
		{`let arr = [5, 10, 15]; arr[2]`, "15"},
		{`[1 + 1, 2 * 2, 3 - 1]`, "[2, 4, 2]"},
		{`let x = [1, 2, 3]; x[0] + x[2]`, "4"},

		// Hash literals
		{`{"key": "value"}`, "{key:value}"},
		{`{"name": "Alice", "age": 25}["name"]`, "Alice"},
		{`let h = {"a": 1, "b": 2}; h["b"]`, "2"},
		{`{"x": 10, "y": 20}["y"]`, "20"},
		{`let m = {1: "one", 2: "two"}; m[1]`, "one"},

		// Function Calls - Simple cases
		{`let f = fn() { 5 }; f()`, "5"},                             // Explicit return
		{`let f = fn() { return 10; }; f()`, "10"},                   // Explicit return with `return`
		{`let f = fn() { 1+2 }; f()`, "3"},                           // Implicit return
		{`let f = fn() { 2*3 }; f()`, "6"},                           // Implicit return with expression
		{`let f = fn() {}; f()`, "null"},                             // Empty function body should return null
		{`let f = fn() { if (true) { return 42; } }; f()`, "42"},     // Return inside conditional
		{`let f = fn() { if (false) { return 42; } }; f()`, "null"},  // Branch with no return should return null
		{`let f = fn() { 5 }; let g = fn() { f() }; g()`, "5"},       // Nested function calls
		{`let f = fn() { 1+2 }; let g = fn() { f() * 2 }; g()`, "6"}, // Function calls within function

		// Function assignment to multiple variables
		{`let f = fn() { 10 }; let g = f; g()`, "10"},                       // Function aliasing
		{`let f = fn() { 20 }; let g = fn() { f() }; let h = g; h()`, "20"}, // Function aliasing with calls

		// Functions returning functions
		{`let f = fn() { fn() { 99 } }; let g = f(); g()`, "99"}, // Function returning another function
		{`let f = fn() { fn() {} }; let g = f(); g()`, "null"},   // Function returning empty function

		// Calling function multiple times
		{`let f = fn() { 5 }; f(); f(); f()`, "5"}, // Ensure multiple calls work

		// Functions in expressions
		{`let f = fn() { 4 }; f() + 2`, "6"},                           // Function return used in expression
		{`let f = fn() { 10 }; let g = fn() { f() + f() }; g()`, "20"}, // Function calls inside expressions

		// Function reassignments
		{`let f = fn() { 10 }; let f = fn() { 20 }; f()`, "20"}, // Overwriting function

		// Function calls in if-expressions
		{`let f = fn() { 7 }; if (true) { f() }`, "7"},                               // Function inside if-true branch
		{`let f = fn() { 8 }; if (false) { f() } else { 12 }`, "12"},                 // Function inside if-false branch
		{`let f = fn() { if (true) { return 30; } else { return 40; } }; f()`, "30"}, // Function with full if-else

		// Basic Local Variable Usage
		{`let f = fn() { let x = 5; x }; f()`, "5"}, // Local variable is returned
		{`let f = fn() { let x = 5; let x = 10; x }; f()`, "10"},
		{`let f = fn() { let x = 10; return x; }; f()`, "10"}, // Explicit return of local variable
		{`let f = fn() { let x = 2; x * 3 }; f()`, "6"},       // Local variable used in an expression

		// Multiple Local Variables
		{`let f = fn() { let a = 3; let b = 4; a + b }; f()`, "7"}, // Multiple local bindings
		{`let f = fn() { let x = 2; let y = x + 5; y }; f()`, "7"}, // Variable referencing another local variable

		// Shadowing Global Variables
		{`let x = 100; let f = fn() { let x = 5; x }; f()`, "5"},       // Local variable shadows global
		{`let x = 100; let f = fn() { let x = 5; x }; f(); x`, "100"},  // Global variable remains unchanged
		{`let x = 100; let f = fn() { let x = x + 5; x }; f()`, "105"}, // Local variable initializes using global variable
		{`let x = 100; let f = fn() { x + 5 }; f()`, "105"},            // Function accesses global variable
		{`let x = 50; let f = fn() { let x = x + 10; x }; f()`, "60"},  // Function uses outer variable while shadowing

		// Conditionals and Local Variables
		{`let f = fn() { let x = 0; if (true) { let x = 10; x } else { x } }; f()`, "10"},              // Local variable inside if-block
		{`let f = fn() { let x = 0; if (false) { let x = 10; x } else { let x = 20; x } }; f()`, "20"}, // Local variable in else-block

		{`let returnsOneReturner = fn() {
				let returnsOne = fn() { 1; };
				returnsOne;};
				returnsOneReturner()();`, "1"},

		{`let f = fn(x) { x }; f(5)`, "5"},                                 // Single argument
		{`let f = fn(x, y) { x + y }; f(3, 4)`, "7"},                       // Multiple arguments
		{`let f = fn(x, y) { x * y }; f(2, 3)`, "6"},                       // Multiplication with arguments
		{`let f = fn(x) { x * 2 }; f(10)`, "20"},                           // Argument used in expression
		{`let f = fn(x) { x }; f(1 + 2)`, "3"},                             // Expression as argument
		{`let f = fn(x, y) { x + y }; f(2 * 3, 4 + 1)`, "11"},              // Mixed expressions as arguments
		{`let f = fn(x) { x * 2 }; let g = fn(y) { f(y) + 1 }; g(3)`, "7"}, // Function call inside function
		{`let f = fn(x) {}; f(5)`, "null"},                                 // Function ignoring argument returns null
		{`let f = fn() { return 42; }; f(5)`, "42"},                        // Extra arguments ignored
		{`let f = fn(x) { x + 1 }; let g = fn(y) { f(y) * 2 }; g(3)`, "8"}, // Function calling another function with argument
		{`let add = fn(x, y) { x + y }; let square = fn(n) { n * n }; let h = fn(a, b) { square(add(a, b)) }; h(2, 3)`, "25"}, // Nested function evaluation
		{`let f = fn(x) { if (x > 10) { return x; } }; f(5)`, "null"},                                                         // Return in condition (false case)
		{`let f = fn(x) { if (x > 10) { return x; } else { return 0; } }; f(15)`, "15"},                                       // Return in condition (true case)
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			compiler, err := testCompile(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			bytecode := compiler.Output()
			vm := NewStackVM(bytecode.Instructions, bytecode.ConstantPool)
			err = vm.Run()
			if err != nil {
				t.Fatal(err)
			}
			obj := vm.Top()
			if obj.Inspect() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, obj.Inspect())
			}
		})
	}
}

func testCompile(input string) (*compiler.Compiler, error) {
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()

	compiler := compiler.New()
	err := compiler.Compile(program)
	if err != nil {
		return nil, err
	}
	return compiler, nil
}
