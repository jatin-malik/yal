package bytecode

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		opCode   OpCode
		operands []int
		expected []byte
		hasError bool
	}{
		// Valid case: OpPush with a 2-byte operand
		{
			opCode:   OpPush,
			operands: []int{0x1234}, // 4660 in decimal
			expected: []byte{byte(OpPush), 0x12, 0x34},
			hasError: false,
		},
		// Edge case: OpPush with the smallest operand (0)
		{
			opCode:   OpPush,
			operands: []int{0},
			expected: []byte{byte(OpPush), 0x00, 0x00},
			hasError: false,
		},
		// Edge case: OpPush with the largest 2-byte operand (0xFFFF)
		{
			opCode:   OpPush,
			operands: []int{0xFFFF}, // 65535 in decimal
			expected: []byte{byte(OpPush), 0xFF, 0xFF},
			hasError: false,
		},
		// Error case: OpPush with no operand
		{
			opCode:   OpPush,
			operands: []int{},
			expected: nil,
			hasError: true,
		},
		// Error case: OpPush with too many operands
		{
			opCode:   OpPush,
			operands: []int{1, 2}, // More than expected
			expected: nil,
			hasError: true,
		},
		// Error case: Unknown opcode
		{
			opCode:   0xFF, // Assuming 0xFF is not a valid opcode
			operands: []int{0x1234},
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("OpCode=%d Operands=%v", tt.opCode, tt.operands), func(t *testing.T) {
			output, err := Make(tt.opCode, tt.operands...)

			if tt.hasError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !bytes.Equal(output, tt.expected) {
					t.Errorf("expected %v, got %v", tt.expected, output)
				}
			}
		})
	}
}
