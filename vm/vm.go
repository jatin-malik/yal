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
		case bytecode.OpAdd:
			right := svm.pop()
			left := svm.pop()

			//TODO: Assuming these are integers. Generalise.
			res := left.(*object.Integer).Value + right.(*object.Integer).Value
			svm.push(&object.Integer{Value: res})
			ip += 1
		}
	}
	return nil
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
