package vm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jatin-malik/yal/bytecode"
	"github.com/jatin-malik/yal/object"
)

const (
	StackSize   int = 2048
	GlobalsSize int = 65536
	FramesSize  int = 1024
)

var builtInFunctions = []*object.BuiltinFunction{
	object.BuiltinFunctions["len"],
	object.BuiltinFunctions["first"],
	object.BuiltinFunctions["last"],
	object.BuiltinFunctions["rest"],
	object.BuiltinFunctions["push"],
	object.BuiltinFunctions["puts"],
}

// VM mimics a real machine. It emulates the fetch-decode-execute cycle of a real machine and operates upon bytecode.
type VM interface {
	Run() error
}

type Frame struct {
	closure *object.Closure
	ip      int
	bp      int
}

func (frame *Frame) instructions() bytecode.Instructions {
	return frame.closure.Fn.Instructions
}

func NewFrame(closure *object.Closure, bp int) *Frame {
	return &Frame{
		closure: closure,
		bp:      bp,
	}
}

// StackVM is a stack based VM.
type StackVM struct {
	frames         []*Frame
	activeFrameIdx int
	constantPool   []object.Object
	globals        []object.Object

	stack []object.Object
	sp    int // sp always points to the next available slot in stack
}

type StackVMOption func(*StackVM)

// WithGlobals allows setting a custom globals array.
func WithGlobals(globals []object.Object) StackVMOption {
	return func(vm *StackVM) {
		vm.globals = globals
	}
}

func NewStackVM(instructions bytecode.Instructions, constantPool []object.Object, options ...StackVMOption) *StackVM {
	mainClosure := &object.Closure{
		Fn: &object.CompiledFunction{
			Instructions: instructions,
		},
	}
	mainFrame := NewFrame(mainClosure, 0)
	frames := make([]*Frame, FramesSize)
	frames[0] = mainFrame
	vm := &StackVM{
		constantPool: constantPool,
		stack:        make([]object.Object, StackSize),
		frames:       frames,
		globals:      make([]object.Object, GlobalsSize),
	}

	// Apply provided options
	for _, option := range options {
		option(vm)
	}

	return vm
}

func (svm *StackVM) Run() error {
	for svm.frames[svm.activeFrameIdx].ip < len(svm.frames[svm.activeFrameIdx].instructions()) {
		activeFrame := svm.frames[svm.activeFrameIdx]

		opcode := bytecode.OpCode(activeFrame.instructions()[activeFrame.ip]) // Fetch

		switch opcode { // Decode
		case bytecode.OpPush:
			idx := binary.BigEndian.Uint16(activeFrame.instructions()[activeFrame.ip+1:])
			obj := svm.constantPool[idx]
			svm.push(obj)
			activeFrame.ip += 1 + 2
		case bytecode.OpPushTrue:
			svm.push(object.TRUE)
			activeFrame.ip += 1
		case bytecode.OpPushFalse:
			svm.push(object.FALSE)
			activeFrame.ip += 1
		case bytecode.OpPushNull:
			svm.push(object.NULL)
			activeFrame.ip += 1
		case bytecode.OpAdd, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpEqual, bytecode.OpNotEqual, bytecode.OpGT:
			err := svm.executeBinaryOperation(opcode)
			if err != nil {
				return err
			}
			activeFrame.ip += 1
		case bytecode.OpNegateBoolean, bytecode.OpNegateNumber:
			err := svm.executeUnaryOperation(opcode)
			if err != nil {
				return err
			}
			activeFrame.ip += 1
		case bytecode.OpJumpIfFalse:
			jumpTo := binary.BigEndian.Uint16(activeFrame.instructions()[activeFrame.ip+1:])
			if !object.IsTruthy(svm.Top()) {
				activeFrame.ip = int(jumpTo)
			} else {
				activeFrame.ip += 1 + 2
			}

		case bytecode.OpJump:
			jumpTo := binary.BigEndian.Uint16(activeFrame.instructions()[activeFrame.ip+1:])
			activeFrame.ip = int(jumpTo)
		case bytecode.OpSetLocal:
			idx := binary.BigEndian.Uint16(activeFrame.instructions()[activeFrame.ip+1:])
			localBindingsStackIdx := svm.frames[svm.activeFrameIdx].bp + 1 + int(idx)
			svm.stack[localBindingsStackIdx] = svm.pop()
			activeFrame.ip += 1 + 2
		case bytecode.OpSetGlobal:
			idx := binary.BigEndian.Uint16(activeFrame.instructions()[activeFrame.ip+1:])
			svm.globals[idx] = svm.pop()
			activeFrame.ip += 1 + 2
		case bytecode.OpGetLocal:
			idx := binary.BigEndian.Uint16(activeFrame.instructions()[activeFrame.ip+1:])
			localBindingsStackIdx := svm.frames[svm.activeFrameIdx].bp + 1 + int(idx)
			obj := svm.stack[localBindingsStackIdx]
			svm.push(obj)
			activeFrame.ip += 1 + 2
		case bytecode.OpGetGlobal:
			idx := binary.BigEndian.Uint16(activeFrame.instructions()[activeFrame.ip+1:])
			obj := svm.globals[idx]
			svm.push(obj)
			activeFrame.ip += 1 + 2
		case bytecode.OpGetBuiltIn:
			idx := int(activeFrame.instructions()[activeFrame.ip+1])
			obj := builtInFunctions[idx]
			svm.push(obj)
			activeFrame.ip += 2
		case bytecode.OpGetFree:
			idx := int(activeFrame.instructions()[activeFrame.ip+1])
			obj := activeFrame.closure.FreeStore[idx]
			svm.push(obj)
			activeFrame.ip += 2
		case bytecode.OpArray:
			count := binary.BigEndian.Uint16(activeFrame.instructions()[activeFrame.ip+1:])
			arr := svm.buildArray(int(count))
			svm.push(arr)
			activeFrame.ip += 1 + 2
		case bytecode.OpHash:
			count := binary.BigEndian.Uint16(activeFrame.instructions()[activeFrame.ip+1:])
			hash, err := svm.buildHash(int(count))
			if err != nil {
				return err
			}
			svm.push(hash)
			activeFrame.ip += 1 + 2
		case bytecode.OpIndex:
			idx := svm.pop()
			iterable := svm.pop()
			obj, err := evalIndexExpression(iterable, idx)
			if err != nil {
				return err
			}
			svm.push(obj)
			activeFrame.ip += 1
		case bytecode.OpClosure:
			idx := binary.BigEndian.Uint16(activeFrame.instructions()[activeFrame.ip+1:])
			compiledFn := svm.constantPool[idx].(*object.CompiledFunction)
			freeCount := int(activeFrame.instructions()[activeFrame.ip+3])

			freeStore := make([]object.Object, freeCount)
			for i := 0; i < freeCount; i++ {
				freeStore[freeCount-1-i] = svm.pop()
			}
			closure := &object.Closure{
				Fn:        compiledFn,
				FreeStore: freeStore,
			}
			svm.push(closure)
			activeFrame.ip += 1 + 2 + 1
		case bytecode.OpCall:
			argsCount := int(activeFrame.instructions()[activeFrame.ip+1])
			fn := svm.stack[svm.sp-1-argsCount]
			if closure, ok := fn.(*object.Closure); ok {
				svm.pushFrame(closure, svm.sp-1-argsCount)
				svm.sp += closure.Fn.NumLocals
			} else if builtInFn, ok := fn.(*object.BuiltinFunction); ok {
				args := make([]object.Object, argsCount)
				for i := 0; i < argsCount; i++ {
					args[argsCount-1-i] = svm.pop()
				}
				svm.pop() // pops function from stack
				obj := builtInFn.Fn(args...)
				if object.IsErrorValue(obj) {
					return errors.New(obj.(*object.Error).Message)
				}
				svm.push(obj)
			} else {
				return fmt.Errorf("type: %T not a callable object", fn)
			}
			activeFrame.ip += 2
		case bytecode.OpGetCurrentClosure:
			closure := svm.frames[svm.activeFrameIdx].closure
			svm.push(closure)
			activeFrame.ip += 1
		case bytecode.OpReturnValue:
			val := svm.pop()
			svm.sp = svm.frames[svm.activeFrameIdx].bp // clean up activation record
			svm.popFrame()
			svm.push(val)
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
	case object.StringObject:
		l := left.(*object.String)
		r := right.(*object.String)
		svm.push(&object.String{Value: l.Value + r.Value})
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
		if r.Value == 0 {
			return fmt.Errorf("division by zero")
		}
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

func (svm *StackVM) pushFrame(closure *object.Closure, bp int) {
	frame := NewFrame(closure, bp)
	svm.frames[svm.activeFrameIdx+1] = frame
	svm.activeFrameIdx++
}

func (svm *StackVM) popFrame() {
	svm.activeFrameIdx--
}

func (svm *StackVM) buildArray(count int) object.Object {
	objs := make([]object.Object, count)
	for count > 0 {
		objs[count-1] = svm.pop()
		count--
	}
	return &object.Array{Elements: objs}
}

func (svm *StackVM) buildHash(count int) (object.Object, error) {
	pairs := make(map[object.Object]object.Object)

	for count > 0 {
		pairs[svm.pop()] = svm.pop()
		count--
	}

	ho := new(object.Hash)
	elems := make(map[object.HashKey]object.Object)
	for k, v := range pairs {
		// Check if key is hashable
		if key, ok := k.(object.Hashable); !ok {
			return nil, fmt.Errorf("key type %s is not hashable", k.Type())
		} else {
			hashKey := key.HashKey()
			elems[hashKey] = v
		}
	}
	ho.Pairs = elems
	return ho, nil
}

func evalIndexExpression(iterable object.Object, index object.Object) (object.Object, error) {
	switch iterable.Type() {
	case object.ArrayObject:
		return evalArrayIndexExpression(iterable, index)
	case object.HashObject:
		return evalHashIndexExpression(iterable, index)
	default:
		return nil, fmt.Errorf("index expression not supported for type: %s", iterable.Type())
	}
}

func evalArrayIndexExpression(iterable object.Object, index object.Object) (object.Object, error) {
	arr := iterable.(*object.Array)

	// The index has to be an integer
	if i, ok := index.(*object.Integer); !ok {
		return nil, fmt.Errorf("index must be an integer for index expression in arrays")
	} else {
		// Check bounds of the index
		idx := i.Value
		if idx < 0 || idx >= int64(len(arr.Elements)) {
			return nil, fmt.Errorf("index %d out of bounds for arr length %d", idx, len(arr.Elements))
		}

		return arr.Elements[idx], nil
	}

}

func evalHashIndexExpression(iterable object.Object, index object.Object) (object.Object, error) {
	hash := iterable.(*object.Hash)
	// Check if key is hashable
	if key, ok := index.(object.Hashable); !ok {
		return nil, fmt.Errorf("key type %s is not hashable", index.Type())
	} else {
		hashKey := key.HashKey()
		if val, ok := hash.Pairs[hashKey]; ok {
			return val, nil
		} else {
			return object.NULL, nil
		}
	}
}
