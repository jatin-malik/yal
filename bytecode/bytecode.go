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
	OpArray
	OpHash
	OpIndex
	OpCall
	OpReturnValue
	OpSetLocal
	OpGetLocal
	OpGetBuiltIn
	OpGetFree
	OpClosure
	OpGetCurrentClosure
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
	case OpArray:
		return "OpArray"
	case OpHash:
		return "OpHash"
	case OpIndex:
		return "OpIndex"
	case OpCall:
		return "OpCall"
	case OpReturnValue:
		return "OpReturnValue"
	case OpSetLocal:
		return "OpSetLocal"
	case OpGetLocal:
		return "OpGetLocal"
	case OpGetBuiltIn:
		return "OpGetBuiltIn"
	case OpGetFree:
		return "OpGetFree"
	case OpClosure:
		return "OpClosure"
	case OpGetCurrentClosure:
		return "OpGetCurrentClosure"
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
	case OpPush, OpJumpIfFalse, OpJump, OpSetGlobal, OpGetGlobal, OpArray, OpHash, OpSetLocal, OpGetLocal:
		if len(operands) != 1 {
			return nil, fmt.Errorf("%s needs one operand", opCode)
		}
		idx := operands[0]
		var operandBytes [2]byte
		binary.BigEndian.PutUint16(operandBytes[:], uint16(idx))
		instructions.Write(operandBytes[:])
	case OpAdd, OpSub, OpMul, OpDiv, OpPushTrue, OpPushFalse, OpEqual, OpNotEqual, OpGT, OpNegateBoolean,
		OpNegateNumber, OpPushNull, OpIndex, OpReturnValue, OpGetCurrentClosure:
	case OpCall, OpGetBuiltIn, OpGetFree:
		if len(operands) != 1 {
			return nil, fmt.Errorf("%s needs one operand", opCode)
		}
		argsCount := operands[0]
		instructions.WriteByte(byte(argsCount))
	case OpClosure:
		if len(operands) != 2 {
			return nil, fmt.Errorf("%s needs one operand", opCode)
		}
		idx := operands[0]
		var operandBytes [2]byte
		binary.BigEndian.PutUint16(operandBytes[:], uint16(idx))
		instructions.Write(operandBytes[:])

		freeCount := operands[1]
		instructions.WriteByte(byte(freeCount))
	default:
		return nil, fmt.Errorf("unknown opcode: %d", opCode)
	}

	return instructions.Bytes(), nil
}
