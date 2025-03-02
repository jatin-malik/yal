package vm

import (
	"encoding/binary"
	"fmt"
	"github.com/jatin-malik/yal/bytecode"
	"github.com/jatin-malik/yal/object"
)

const (
	StackSize int = 2048
)

// VM mimics a real machine. It emulates the fetch-decode-execute cycle of a real machine and operates upon bytecode.
type VM interface {
	Run() error
}

// StackVM is a stack based VM.
type StackVM struct {
	instructions bytecode.Instructions
	constantPool []object.Object

	stack []object.Object
	sp    int // sp always points to the next available slot in stack
}

func NewStackVM(instructions bytecode.Instructions, constantPool []object.Object) *StackVM {
	return &StackVM{
		constantPool: constantPool,
		stack:        make([]object.Object, StackSize),
		instructions: instructions,
	}
}

func (svm *StackVM) Run() error {
	for ip := 0; ip < len(svm.instructions); {
		opcode := bytecode.OpCode(svm.instructions[ip]) // Fetch

		switch opcode { // Decode
		case bytecode.OpPush:
			idx := binary.BigEndian.Uint16(svm.instructions[ip+1:])
			obj := svm.constantPool[idx]
			svm.push(obj)
			ip += 1 + 2
		case bytecode.OpPushTrue:
			svm.push(object.TRUE)
			ip += 1
		case bytecode.OpPushFalse:
			svm.push(object.FALSE)
			ip += 1
		case bytecode.OpPushNull:
			svm.push(object.NULL)
			ip += 1
		case bytecode.OpAdd, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpEqual, bytecode.OpNotEqual, bytecode.OpGT:
			err := svm.executeBinaryOperation(opcode)
			if err != nil {
				return err
			}
			ip += 1
		case bytecode.OpNegateBoolean, bytecode.OpNegateNumber:
			err := svm.executeUnaryOperation(opcode)
			if err != nil {
				return err
			}
			ip += 1
		case bytecode.OpJumpIfFalse:
			jumpTo := binary.BigEndian.Uint16(svm.instructions[ip+1:])
			if !object.IsTruthy(svm.Top()) {
				ip = int(jumpTo)
			} else {
				ip += 1 + 2
			}

		case bytecode.OpJump:
			jumpTo := binary.BigEndian.Uint16(svm.instructions[ip+1:])
			ip = int(jumpTo)
		default:
			return fmt.Errorf("unknown opcode: %d", opcode)
		}
	}
	return nil
}

func (svm *StackVM) executeUnaryOperation(opcode bytecode.OpCode) error {
	operand := svm.pop()

	switch opcode {
	case bytecode.OpNegateNumber:
		return svm.executeNegateNumberUnaryOperation(operand)
	case bytecode.OpNegateBoolean:
		return svm.executeNegateBooleanUnaryOperation(operand)
	}
	return nil
}

func (svm *StackVM) executeNegateNumberUnaryOperation(operand object.Object) error {
	if i, ok := operand.(*object.Integer); ok {
		svm.push(&object.Integer{Value: -i.Value})
		return nil
	} else {
		return fmt.Errorf("invalid type %s with operator '-'", operand.Type())
	}
}

func (svm *StackVM) executeNegateBooleanUnaryOperation(operand object.Object) error {
	if object.IsTruthy(operand) {
		svm.push(object.FALSE)
	} else {
		svm.push(object.TRUE)
	}
	return nil
}

func (svm *StackVM) executeBinaryOperation(opcode bytecode.OpCode) error {
	right := svm.pop()
	left := svm.pop()

	// check for type mismatch
	if left.Type() != right.Type() {
		return fmt.Errorf("incompatible types: %s and %s", left.Type(), right.Type())
	}

	switch opcode {
	case bytecode.OpAdd:
		return svm.executeAddBinaryOperation(left, right)
	case bytecode.OpSub:
		return svm.executeSubBinaryOperation(left, right)
	case bytecode.OpMul:
		return svm.executeMulBinaryOperation(left, right)
	case bytecode.OpDiv:
		return svm.executeDivBinaryOperation(left, right)
	case bytecode.OpEqual:
		return svm.executeEqualsBinaryOperation(left, right)
	case bytecode.OpNotEqual:
		return svm.executeNotEqualsBinaryOperation(left, right)
	case bytecode.OpGT:
		return svm.executeGreaterThanBinaryOperation(left, right)
	}

	return nil
}

func (svm *StackVM) executeAddBinaryOperation(left, right object.Object) error {
	switch left.Type() {
	case object.IntegerObject:
		l := left.(*object.Integer)
		r := right.(*object.Integer)
		svm.push(&object.Integer{Value: l.Value + r.Value})
	default:
		return fmt.Errorf("unsupported operand type %s with '+'", left.Type())
	}
	return nil
}

func (svm *StackVM) executeSubBinaryOperation(left, right object.Object) error {
	switch left.Type() {
	case object.IntegerObject:
		l := left.(*object.Integer)
		r := right.(*object.Integer)
		svm.push(&object.Integer{Value: l.Value - r.Value})
	default:
		return fmt.Errorf("unsupported operand type %s with '-'", left.Type())
	}
	return nil
}

func (svm *StackVM) executeMulBinaryOperation(left, right object.Object) error {
	switch left.Type() {
	case object.IntegerObject:
		l := left.(*object.Integer)
		r := right.(*object.Integer)
		svm.push(&object.Integer{Value: l.Value * r.Value})
	default:
		return fmt.Errorf("unsupported operand type %s with '*'", left.Type())
	}
	return nil
}

func (svm *StackVM) executeDivBinaryOperation(left, right object.Object) error {
	switch left.Type() {
	case object.IntegerObject:
		l := left.(*object.Integer)
		r := right.(*object.Integer)
		svm.push(&object.Integer{Value: l.Value / r.Value})
	default:
		return fmt.Errorf("unsupported operand type %s with '/'", left.Type())
	}
	return nil
}

func (svm *StackVM) executeEqualsBinaryOperation(left, right object.Object) error {
	objType := left.Type()

	switch objType {
	case object.IntegerObject:
		svm.push(getBooleanObject(left.(*object.Integer).Value == right.(*object.Integer).Value))
	case object.StringObject:
		svm.push(getBooleanObject(left.(*object.String).Value == right.(*object.String).Value))
	case object.BooleanObject:
		svm.push(getBooleanObject(left == right)) // pointer comparison
	default:
		return fmt.Errorf("unsupported operand type %s with '=='", left.Type())
	}
	return nil
}

func (svm *StackVM) executeNotEqualsBinaryOperation(left, right object.Object) error {
	objType := left.Type()

	switch objType {
	case object.IntegerObject:
		svm.push(getBooleanObject(left.(*object.Integer).Value != right.(*object.Integer).Value))
	case object.StringObject:
		svm.push(getBooleanObject(left.(*object.String).Value != right.(*object.String).Value))
	case object.BooleanObject:
		svm.push(getBooleanObject(left != right)) // pointer comparison
	default:
		return fmt.Errorf("unsupported operand type %s with '!='", left.Type())
	}
	return nil
}

func (svm *StackVM) executeGreaterThanBinaryOperation(left, right object.Object) error {
	objType := left.Type()

	switch objType {
	case object.IntegerObject:
		svm.push(getBooleanObject(left.(*object.Integer).Value > right.(*object.Integer).Value))
	default:
		return fmt.Errorf("unsupported operand type %s with '>'", left.Type())
	}
	return nil
}

func getBooleanObject(boolValue bool) object.Object {
	if boolValue {
		return object.TRUE
	}
	return object.FALSE
}

func (svm *StackVM) push(obj object.Object) error {
	if svm.sp >= len(svm.stack) {
		return fmt.Errorf("stack overflow")
	}

	svm.stack[svm.sp] = obj
	svm.sp++
	return nil
}

func (svm *StackVM) pop() object.Object {
	if svm.sp == 0 {
		return nil
	}

	obj := svm.stack[svm.sp-1]
	svm.sp--
	return obj
}

func (svm *StackVM) Top() object.Object {
	if svm.sp == 0 {
		return nil
	}

	obj := svm.stack[svm.sp-1]
	return obj
}
