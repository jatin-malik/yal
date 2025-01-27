package evaluator

import (
	"github.com/jatin-malik/make-thy-interpreter/ast"
	"github.com/jatin-malik/make-thy-interpreter/object"
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
		}
	case *ast.BlockStatement:
		for _, stmt := range v.Statements {
			result = Eval(stmt, env)
			if isReturnValue(result) {
				return result
			}
		}
	case *ast.ReturnStatement:
		obj := Eval(v.Value, env)
		result = &object.ReturnValue{Value: obj}
	case *ast.ExpressionStatement:
		result = Eval(v.Expr, env)
	case *ast.IntegerLiteral:
		result = &object.Integer{Value: v.Value}
	case *ast.BooleanLiteral:
		result = getBooleanObject(v.Value)
	case *ast.PrefixExpression:
		operandObject := Eval(v.Right, env)
		result = evalPrefixExpression(v.Operator, operandObject)
	case *ast.InfixExpression:
		leftObj := Eval(v.Left, env)
		rightObj := Eval(v.Right, env)
		result = evalInfixExpression(v.Operator, leftObj, rightObj)
	case *ast.IfElseConditional:
		conditionObj := Eval(v.Condition, env)
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
		env.Set(v.Name.Value, rightObj)
	case *ast.Identifier:
		result = env.Get(v.Value)
	default:
		result = object.NULL
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
		return object.NULL
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	// check for type mismatch
	if left.Type() != right.Type() {
		return object.NULL
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
	l, ok1 := left.(*object.Integer)
	r, ok2 := right.(*object.Integer)
	if ok1 && ok2 {
		return &object.Integer{Value: l.Value + r.Value}
	} else {
		return object.NULL
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
		return object.NULL
	}
}

func evalBangPrefixExpression(right object.Object) object.Object {
	if isTruthy(right) {
		return object.FALSE
	}
	return object.TRUE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case object.NULL, object.FALSE:
		return false
	default:
		return true
	}
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
