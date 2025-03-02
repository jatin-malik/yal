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
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			compiler := testCompile(tt.input)
			bytecode := compiler.Output()
			vm := NewStackVM(bytecode.Instructions, bytecode.ConstantPool)
			err := vm.Run()
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

func testCompile(input string) *compiler.Compiler {
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()

	compiler := compiler.New()
	compiler.Compile(program)
	return compiler
}
