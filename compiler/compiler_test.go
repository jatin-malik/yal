package compiler

import (
	"bytes"
	"github.com/jatin-malik/yal/bytecode"
	"github.com/jatin-malik/yal/lexer"
	"github.com/jatin-malik/yal/object"
	"github.com/jatin-malik/yal/parser"
	"testing"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		input                string
		expectedByteCode     bytecode.Instructions
		expectedConstantPool []any
	}{
		{
			input: "1 + 2",
			expectedByteCode: bytecode.Instructions{
				0x00,       // OpPush (1)
				0x00, 0x00, // Index 0 (constant pool: 1)
				0x00,       // OpPush (2)
				0x00, 0x01, // Index 1 (constant pool: 2)
				0x01, // OpAdd
			},
			expectedConstantPool: []any{1, 2},
		},

		{
			input: "2 - 1",
			expectedByteCode: bytecode.Instructions{
				0x00,       // OpPush (1)
				0x00, 0x00, // Index 0 (constant pool: 1)
				0x00,       // OpPush (2)
				0x00, 0x01, // Index 1 (constant pool: 2)
				0x02, // OpSub
			},
			expectedConstantPool: []any{2, 1},
		},

		{
			input: "2 * 3",
			expectedByteCode: bytecode.Instructions{
				0x00,       // OpPush (1)
				0x00, 0x00, // Index 0 (constant pool: 1)
				0x00,       // OpPush (2)
				0x00, 0x01, // Index 1 (constant pool: 2)
				0x03, // OpMul
			},
			expectedConstantPool: []any{2, 3},
		},

		{
			input: "4 / 2",
			expectedByteCode: bytecode.Instructions{
				0x00,       // OpPush (1)
				0x00, 0x00, // Index 0 (constant pool: 1)
				0x00,       // OpPush (2)
				0x00, 0x01, // Index 1 (constant pool: 2)
				0x04, // OpDiv
			},
			expectedConstantPool: []any{4, 2},
		},

		{
			input: `if (1) {2} else {3};4`,
			expectedByteCode: bytecode.Instructions{
				// OpPush (1) to check the condition
				0x00,       // OpPush (1)
				0x00, 0x00, // Index 0 (constant pool: 1)

				// OpJumpIfFalse: Jump to the "false" block (else) if condition is false
				0x0C, 0x00, 0x0C, // OpJumpIfFalse with an offset (to false block)

				// (True block) OpPush (2) for the "if" block (if the condition is true)
				0x00,       // OpPush (2)
				0x00, 0x01, // Index 1 (constant pool: 2)

				// OpJump: Skip over the false block
				0x0D, 0x00, 0x0F, // OpJump to skip the false block and jump to the post-else instruction

				// (False block) OpPush (3) for the "else" block (if the condition is false)
				0x00,       // OpPush (3)
				0x00, 0x02, // Index 2 (constant pool: 3)

				// OpPush (1) to check the condition
				0x00,       // OpPush (1)
				0x00, 0x03, // Index 3 (constant pool: 4)
			},
			expectedConstantPool: []any{1, 2, 3, 4},
		},

		{
			input: `if (1) {2};4`,
			expectedByteCode: bytecode.Instructions{
				// OpPush (1) to check the condition
				0x00,       // OpPush (1)
				0x00, 0x00, // Index 0 (constant pool: 1)

				0x0C, 0x00, 0x0C,

				// (True block) OpPush (2) for the "if" block (if the condition is true)
				0x00,       // OpPush (2)
				0x00, 0x01, // Index 1 (constant pool: 2)

				0x0D, 0x00, 0x0D,

				0x0E,

				// OpPush (1) to check the condition
				0x00,       // OpPush (1)
				0x00, 0x02, // Index 2 (constant pool: 3)
			},
			expectedConstantPool: []any{1, 2, 4},
		},

		{
			input: `if (false) {2};4`,
			expectedByteCode: bytecode.Instructions{
				06, // Load false

				0x0C, 0x00, 0x0A,

				// (True block) OpPush (2) for the "if" block (if the condition is true)
				0x00,       // OpPush (2)
				0x00, 0x00, // Index 0 (constant pool: 1)

				0x0D, 0x00, 0x0B,

				0x0E,

				// OpPush (1) to check the condition
				0x00,       // OpPush (1)
				0x00, 0x01, // Index 1 (constant pool: 2)
			},
			expectedConstantPool: []any{2, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			compiler, err := testCompile(tt.input)
			if err != nil {
				t.Fatalf("Compilation failed: %v", err)
			}
			assertBytecode(t, tt.expectedByteCode, compiler.instructions)
			assertConstantPool(t, tt.expectedConstantPool, compiler.constantPool)
		})
	}
}

func testCompile(input string) (*Compiler, error) {
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()

	compiler := New()
	err := compiler.Compile(program)
	if err != nil {
		return nil, err
	}
	return compiler, nil
}

func assertBytecode(t *testing.T, expected, actual bytecode.Instructions) {
	if !bytes.Equal(expected, actual) {
		t.Errorf("Bytecode mismatch:\nExpected:\n%02X\nGot:\n%02X", expected, actual)
	}
}

func assertConstantPool(t *testing.T, expected []any, actual []object.Object) {
	if len(actual) != len(expected) {
		t.Errorf("Constant pool length mismatch: expected %d, got %d", len(expected), len(actual))
		return
	}

	for i, expectedConst := range expected {
		switch expectedConst := expectedConst.(type) {
		case int:
			testIntegerObject(t, actual[i], int64(expectedConst))
		default:
			t.Errorf("Unsupported constant type at index %d: %T", i, expectedConst)
		}
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	if i, ok := obj.(*object.Integer); ok {
		if i.Value != expected {
			t.Errorf("constant: expected %d, got %d", expected, i.Value)
		}
	} else {
		t.Errorf("expected *object.Integer, got %s", obj.Inspect())

	}
}
