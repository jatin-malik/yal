package interpreter

import (
	"fmt"
	"github.com/jatin-malik/yal/evaluator"
	"github.com/jatin-malik/yal/lexer"
	"github.com/jatin-malik/yal/object"
	"github.com/jatin-malik/yal/parser"
)

func Interpret(input string) {
	l := lexer.New(input)
	p := parser.New(l)
	prg := p.ParseProgram()
	if len(p.Errors) != 0 {
		for _, msg := range p.Errors {
			fmt.Println(msg)
		}
		return
	}

	macroEnv := object.NewEnvironment(nil)
	expandedAST, err := evaluator.ExpandMacro(prg, macroEnv)
	if err != nil {
		fmt.Println(err)
		return
	}
	env := object.NewEnvironment(nil)
	obj := evaluator.Eval(expandedAST, env)
	if obj != nil {
		fmt.Println(obj.Inspect())
	}
}
