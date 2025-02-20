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
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			compiler := testCompile(tt.input)
			bytecode := compiler.Emit()
			vm := NewStackVM(bytecode.Instructions, bytecode.ConstantPool)
			vm.Run()
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
