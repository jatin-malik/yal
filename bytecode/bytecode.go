// Package bytecode mimics machine code for the YAL virtual machine.
package bytecode

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte
type OpCode byte

const (
	OpPush OpCode = iota
	OpAdd
)

// Make generates a bytecode instruction from the input opCode and operands. Multibyte operands are encoded in
// BigEndian order.
func Make(opCode OpCode, operands ...int) ([]byte, error) {
	var instructions bytes.Buffer
	instructions.WriteByte(byte(opCode))
	switch opCode {
	case OpPush:
		// OpPush expects one operand, the index to the constant in the constant pool. The index is 2 bytes wide.
		if len(operands) != 1 {
			return nil, fmt.Errorf("OpPush needs one operand")
		}
		idx := operands[0]
		var operandBytes [2]byte
		binary.BigEndian.PutUint16(operandBytes[:], uint16(idx))
		instructions.Write(operandBytes[:])
	case OpAdd:
	default:
		return nil, fmt.Errorf("unknown opcode: %d", opCode)
	}

	return instructions.Bytes(), nil
}
