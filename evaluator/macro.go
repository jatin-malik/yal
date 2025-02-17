package evaluator

import (
	"fmt"
	"github.com/jatin-malik/yal/ast"
	"github.com/jatin-malik/yal/object"
)

func ExpandMacro(prg ast.Node, env *object.Environment) (ast.Node, error) {
	return ast.Walker(prg, func(node ast.Node) (ast.Node, error) {
		if isMacroDefinition(node) {
			macroDef := node.(*ast.LetStatement)
			macroLiteral := macroDef.Right.(*ast.MacroLiteral)
			macroObj := &object.Macro{
				Env:        env,
				Parameters: macroLiteral.Parameters,
				Body:       macroLiteral.Body,
			}
			env.Set(macroDef.Name.Value, macroObj) // register macro definition
			return nil, nil                        // remove from AST
		}

		if isMacroCallExpression(node, env) {
			macroCall := node.(*ast.CallExpression)
			obj := Eval(macroCall, env)
			if object.IsErrorValue(obj) {
				errorMsg := obj.(*object.Error).Message
				return nil, fmt.Errorf("macro expansion error: %s", errorMsg)
			}
			if quoted, ok := obj.(*object.Quote); ok {
				return quoted.Node, nil
			}
		}

		return node, nil
	})
}

func isMacroDefinition(node ast.Node) bool {
	if node, ok := node.(*ast.LetStatement); ok {
		if _, ok = node.Right.(*ast.MacroLiteral); ok {
			return true
		}
	}

	return false

}

func isMacroCallExpression(node ast.Node, env *object.Environment) bool {
	if node, ok := node.(*ast.CallExpression); ok {
		if ident, ok := node.Function.(*ast.Identifier); ok {
			obj := env.Get(ident.Value)
			if _, ok := obj.(*object.Macro); ok {
				return true
			}
		}
	}

	return false
}
