package evaluator

import (
	"fmt"
	"github.com/jatin-malik/yal/ast"
	"github.com/jatin-malik/yal/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	var result object.Object
	switch v := node.(type) {
	case *ast.Program:
		for _, stmt := range v.Statements {
			result = Eval(stmt, env)
			if isReturnValue(result) {
				return result.(*object.ReturnValue).Value // unwrap
			}
			if isErrorValue(result) {
				return result // no unwrap
			}
		}
	case *ast.BlockStatement:
		for _, stmt := range v.Statements {
			result = Eval(stmt, env)
			if isReturnValue(result) || isErrorValue(result) {
				return result
			}
		}
	case *ast.ReturnStatement:
		obj := Eval(v.Value, env)
		if isErrorValue(obj) {
			return obj
		}
		result = &object.ReturnValue{Value: obj}
	case *ast.ExpressionStatement:
		result = Eval(v.Expr, env)
	case *ast.IntegerLiteral:
		result = &object.Integer{Value: v.Value}
	case *ast.StringLiteral:
		result = &object.String{Value: v.Value}
	case *ast.BooleanLiteral:
		result = getBooleanObject(v.Value)
	case *ast.ArrayLiteral:
		var elems []object.Object
		for _, elem := range v.Elements {
			result = Eval(elem, env)
			if isErrorValue(result) {
				return result
			}
			elems = append(elems, result)
		}
		result = &object.Array{Elements: elems}
	case *ast.HashLiteral:
		pairs := make(map[object.Object]object.Object)

		for k, v := range v.Pairs {
			key := Eval(k, env)
			if isErrorValue(key) {
				return key
			}
			val := Eval(v, env)
			if isErrorValue(val) {
				return val
			}
			pairs[key] = val
		}

		result = evalHashLiteral(pairs)

	case *ast.FunctionLiteral:
		result = &object.Function{
			Env:        env,
			Parameters: v.Parameters,
			Body:       v.Body,
		}
	case *ast.CallExpression:
		fn := Eval(v.Function, env)
		if isErrorValue(fn) {
			return fn
		}

		var args []object.Object
		for _, arg := range v.Arguments {
			result = Eval(arg, env)
			if isErrorValue(result) {
				return result
			}
			args = append(args, result)
		}

		result = evalCallExpression(fn, args)
	case *ast.IndexExpression:
		iterable := Eval(v.Left, env)
		if isErrorValue(iterable) {
			return iterable
		}

		idx := Eval(v.Index, env)
		if isErrorValue(idx) {
			return idx
		}

		result = evalIndexExpression(iterable, idx)

	case *ast.PrefixExpression:
		operandObject := Eval(v.Right, env)
		if isErrorValue(operandObject) {
			return operandObject
		}
		result = evalPrefixExpression(v.Operator, operandObject)
	case *ast.InfixExpression:
		leftObj := Eval(v.Left, env)
		if isErrorValue(leftObj) {
			return leftObj
		}
		rightObj := Eval(v.Right, env)
		if isErrorValue(rightObj) {
			return rightObj
		}
		result = evalInfixExpression(v.Operator, leftObj, rightObj)
	case *ast.IfElseConditional:
		conditionObj := Eval(v.Condition, env)
		if isErrorValue(conditionObj) {
			return conditionObj
		}
		if isTruthy(conditionObj) {
			result = Eval(v.Consequence, env)
		} else {
			if v.Alternative != nil {
				result = Eval(v.Alternative, env)
			} else {
				result = object.NULL
			}
		}
	case *ast.LetStatement:
		rightObj := Eval(v.Right, env)
		if isErrorValue(rightObj) {
			return rightObj
		}
		env.Set(v.Name.Value, rightObj)
	case *ast.Identifier:
		result = env.Get(v.Value)
	default:
		msg := fmt.Sprintf("Unknown statement type: %T", v)
		result = object.NewError(msg)
	}

	return result
}

func getBooleanObject(boolValue bool) object.Object {
	if boolValue {
		return object.TRUE
	}
	return object.FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "-":
		return evalMinusPrefixExpression(right)
	case "!":
		return evalBangPrefixExpression(right)
	default:
		msg := fmt.Sprintf("Unknown operator: %s", operator)
		return object.NewError(msg)
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	// check for type mismatch
	if left.Type() != right.Type() {
		errorMsg := fmt.Sprintf("Incompatible types: %s and %s", left.Type(), right.Type())
		return object.NewError(errorMsg)
	}

	switch operator {
	case "+":
		return evalPlusInfixExpression(left, right)
	case "-":
		return evalMinusInfixExpression(left, right)
	case "*":
		return evalAsteriskInfixExpression(left, right)
	case "/":
		return evalSlashInfixExpression(left, right)
	case "==":
		return evalEqualsInfixExpression(left, right)
	case "!=":
		return evalNotEqualsInfixExpression(left, right)
	case "<":
		return evalLTInfixExpression(left, right)
	case ">":
		return evalGTInfixExpression(left, right)
	default:
		return object.NULL
	}
}

func evalPlusInfixExpression(left, right object.Object) object.Object {
	switch left.Type() {
	case object.IntegerObject:
		l := left.(*object.Integer)
		r := right.(*object.Integer)
		return &object.Integer{Value: l.Value + r.Value}
	case object.StringObject:
		l := left.(*object.String)
		r := right.(*object.String)
		return &object.String{Value: l.Value + r.Value}
	default:
		return object.NewError(fmt.Sprintf("unsupported operand type %s with '+'", left.Type()))
	}

}

func evalMinusInfixExpression(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		return &object.Integer{Value: l.Value - r.Value}
	} else {
		return object.NULL
	}
}

func evalAsteriskInfixExpression(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		return &object.Integer{Value: l.Value * r.Value}
	} else {
		return object.NULL
	}
}

func evalSlashInfixExpression(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		if r.Value == 0 {
			return object.NewError("Division by zero")
		}
		return &object.Integer{Value: l.Value / r.Value}
	} else {
		return object.NULL
	}
}

func evalEqualsInfixExpression(left, right object.Object) object.Object {

	objType := left.Type()

	switch objType {
	case object.IntegerObject:
		return getBooleanObject(left.(*object.Integer).Value == right.(*object.Integer).Value)
	case object.StringObject:
		return getBooleanObject(left.(*object.String).Value == right.(*object.String).Value)
	case object.BooleanObject:
		return getBooleanObject(left == right) // no need to unwrap
	default:
		return object.NULL
	}
}

// TODO: Converting x!=y to !(x==y) for reuse. Check for correctness.
func evalNotEqualsInfixExpression(left, right object.Object) object.Object {
	obj := evalEqualsInfixExpression(left, right)
	if isNull(obj) {
		return object.NULL
	}
	return getBooleanObject(!obj.(*object.Boolean).Value)
}

func evalLTInfixExpression(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		return getBooleanObject(l.Value < r.Value)
	} else {
		return object.NULL
	}
}

func evalGTInfixExpression(left, right object.Object) object.Object {
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		return getBooleanObject(l.Value > r.Value)
	} else {
		return object.NULL
	}
}

func evalMinusPrefixExpression(right object.Object) object.Object {
	if i, ok := right.(*object.Integer); ok {
		return &object.Integer{Value: -i.Value}
	} else {
		msg := fmt.Sprintf("Invalid type %s with operator '-'", right.Type())
		return object.NewError(msg)
	}
}

func evalBangPrefixExpression(right object.Object) object.Object {
	if isTruthy(right) {
		return object.FALSE
	}
	return object.TRUE
}

func evalCallExpression(function object.Object, args []object.Object) object.Object {
	switch function.Type() {
	case object.FunctionObject:
		fn := function.(*object.Function)
		// Extend the environment for this function evaluation
		extendedEnv := object.NewEnvironment(fn.Env)

		if len(fn.Parameters) != len(args) {
			msg := fmt.Sprintf("expected %d parameters, got %d args", len(fn.Parameters), len(args))
			return object.NewError(msg)
		}

		for idx, param := range fn.Parameters {
			extendedEnv.Set(param.Value, args[idx])
		}

		result := Eval(fn.Body, extendedEnv)
		if isReturnValue(result) {
			return result.(*object.ReturnValue).Value // unwrap
		}
		return result
	case object.BuiltInFunctionObject:
		fn := function.(*object.BuiltinFunction)
		return fn.Fn(args...)
	default:
		msg := fmt.Sprintf("expected *object.Function, got %s", function.Type())
		return object.NewError(msg)
	}
}

func evalIndexExpression(iterable object.Object, index object.Object) object.Object {
	switch iterable.Type() {
	case object.ArrayObject:
		return evalArrayIndexExpression(iterable, index)
	case object.HashObject:
		return evalHashIndexExpression(iterable, index)
	default:
		msg := fmt.Sprintf("index expression not supported for type: %s", iterable.Type())
		return object.NewError(msg)
	}
}

func evalArrayIndexExpression(iterable object.Object, index object.Object) object.Object {
	arr := iterable.(*object.Array)

	// The index has to be an integer
	if i, ok := index.(*object.Integer); !ok {
		return object.NewError(fmt.Sprintf("index must be an integer for index expression in arrays"))
	} else {
		// Check bounds of the index
		idx := i.Value
		if idx < 0 || idx >= int64(len(arr.Elements)) {
			return object.NewError(fmt.Sprintf("index out of bounds for arr length %d", len(arr.Elements)))
		}

		return arr.Elements[idx]
	}

}

func evalHashIndexExpression(iterable object.Object, index object.Object) object.Object {
	hash := iterable.(*object.Hash)
	// Check if key is hashable
	if key, ok := index.(object.Hashable); !ok {
		return object.NewError(fmt.Sprintf("key type %s is not hashable", index.Type()))
	} else {
		hashKey := key.HashKey()
		if val, ok := hash.Pairs[hashKey]; ok {
			return val
		} else {
			return object.NULL
		}
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case object.NULL, object.FALSE:
		return false
	default:
		return true
	}
}

func evalHashLiteral(pairs map[object.Object]object.Object) object.Object {
	ho := new(object.Hash)
	elems := make(map[object.HashKey]object.Object)
	for k, v := range pairs {
		// Check if key is hashable
		if key, ok := k.(object.Hashable); !ok {
			return object.NewError(fmt.Sprintf("key type %s is not hashable", k.Type()))
		} else {
			hashKey := key.HashKey()
			elems[hashKey] = v
		}
	}
	ho.Pairs = elems
	return ho
}

func isNull(obj object.Object) bool {
	return obj == object.NULL
}

func isReturnValue(obj object.Object) bool {
	if _, ok := obj.(*object.ReturnValue); ok {
		return true
	}
	return false
}

func isErrorValue(obj object.Object) bool {
	if _, ok := obj.(*object.Error); ok {
		return true
	}
	return false
}
