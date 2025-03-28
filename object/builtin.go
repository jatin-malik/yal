package object

import "fmt"

const BuiltInFunctionObject ObjectType = "BUILTIN_FUNCTION"

type BuiltInFunc func(args ...Object) Object

type BuiltinFunction struct {
	Fn BuiltInFunc
}

func (builtin *BuiltinFunction) Type() ObjectType {
	return BuiltInFunctionObject
}

func (builtin *BuiltinFunction) Inspect() string {
	return "builtin function"
}

var BuiltinFunctions = map[string]*BuiltinFunction{
	"len":   {builtinLen},
	"first": {builtinFirst},
	"last":  {builtinLast},
	"rest":  {builtinRest},
	"push":  {builtinPush},
	"puts":  {builtinPuts},
}

var (
	builtinLen = func(args ...Object) Object {
		if len(args) != 1 {
			return NewError(fmt.Sprintf("len() requires 1 argument. got %d", len(args)))
		}

		switch arg := args[0].(type) {
		case *String:
			return &Integer{Value: int64(len(arg.Value))}
		case *Array:
			return &Integer{Value: int64(len(arg.Elements))}
		default:
			return NewError(fmt.Sprintf("len(): type %s not supported", arg.Type()))
		}
	}

	builtinFirst = func(args ...Object) Object {
		if len(args) != 1 {
			return NewError(fmt.Sprintf("first() requires 1 argument. got %d", len(args)))
		}

		switch arg := args[0].(type) {
		case *Array:
			if len(arg.Elements) > 0 {
				return arg.Elements[0]
			}
			return NewError("empty array")
		default:
			return NewError(fmt.Sprintf("first(): type %s not supported", arg.Type()))
		}
	}

	builtinLast = func(args ...Object) Object {
		if len(args) != 1 {
			return NewError(fmt.Sprintf("last() requires 1 argument. got %d", len(args)))
		}

		switch arg := args[0].(type) {
		case *Array:
			if len(arg.Elements) > 0 {
				return arg.Elements[len(arg.Elements)-1]
			}
			return NewError("empty array")
		default:
			return NewError(fmt.Sprintf("last(): type %s not supported", arg.Type()))
		}
	}

	builtinRest = func(args ...Object) Object {
		if len(args) != 1 {
			return NewError(fmt.Sprintf("rest() requires 1 argument. got %d", len(args)))
		}

		switch arg := args[0].(type) {
		case *Array:
			if len(arg.Elements) > 0 {
				restArray := make([]Object, len(arg.Elements)-1)
				copy(restArray, arg.Elements[1:])
				return &Array{Elements: restArray}
			}
			return NewError("empty array")
		default:
			return NewError(fmt.Sprintf("rest(): type %s not supported", arg.Type()))
		}
	}

	builtinPush = func(args ...Object) Object {
		if len(args) != 2 {
			return NewError(fmt.Sprintf("push() requires 2 arguments. got %d", len(args)))
		}

		switch arg := args[0].(type) {
		case *Array:
			extArray := make([]Object, len(arg.Elements)+1)
			copy(extArray, arg.Elements)
			extArray[len(arg.Elements)] = args[1]
			return &Array{Elements: extArray}
		default:
			return NewError(fmt.Sprintf("push(): type %s not supported", arg.Type()))
		}
	}

	builtinPuts = func(args ...Object) Object {
		for _, arg := range args {
			fmt.Print(arg.Inspect())
		}
		fmt.Println()
		return NULL
	}
)
