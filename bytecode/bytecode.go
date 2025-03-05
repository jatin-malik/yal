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
	OpSub
	OpMul
	OpDiv
	OpPushTrue
	OpPushFalse
	OpEqual
	OpNotEqual
	OpGT
	OpNegateBoolean
	OpNegateNumber
	OpJumpIfFalse
	OpJump
	OpPushNull
	OpSetGlobal
	OpGetGlobal
)

func (op OpCode) String() string {
	switch op {
	case OpPush:
		return "OpPush"
	case OpAdd:
		return "OpAdd"
	case OpSub:
		return "OpSub"
	case OpMul:
		return "OpMul"
	case OpDiv:
		return "OpDiv"
	case OpPushTrue:
		return "OpPushTrue"
	case OpPushFalse:
		return "OpPushFalse"
	case OpEqual:
		return "OpEqual"
	case OpNotEqual:
		return "OpNotEqual"
	case OpGT:
		return "OpGT"
	case OpNegateBoolean:
		return "OpNegateBoolean"
	case OpNegateNumber:
		return "OpNegateNumber"
	case OpJumpIfFalse:
		return "OpJumpIfFalse"
	case OpJump:
		return "OpJump"
	case OpPushNull:
		return "OpPushNull"
	case OpSetGlobal:
		return "OpSetGlobal"
	case OpGetGlobal:
		return "OpGetGlobal"
	default:
		return fmt.Sprintf("OpCode(%d)", op)
	}
}

// Make generates a bytecode instruction from the input opCode and operands. Multibyte operands are encoded in
// BigEndian order.
func Make(opCode OpCode, operands ...int) ([]byte, error) {
	var instructions bytes.Buffer
	instructions.WriteByte(byte(opCode))
	switch opCode {
	case OpPush, OpJumpIfFalse, OpJump, OpSetGlobal, OpGetGlobal:
		if len(operands) != 1 {
			return nil, fmt.Errorf("%s needs one operand", opCode)
		}
		idx := operands[0]
		var operandBytes [2]byte
		binary.BigEndian.PutUint16(operandBytes[:], uint16(idx))
		instructions.Write(operandBytes[:])
	case OpAdd, OpSub, OpMul, OpDiv, OpPushTrue, OpPushFalse, OpEqual, OpNotEqual, OpGT, OpNegateBoolean,
		OpNegateNumber, OpPushNull:
	default:
		return nil, fmt.Errorf("unknown opcode: %d", opCode)
	}

	return instructions.Bytes(), nil
}
