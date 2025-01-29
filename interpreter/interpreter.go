package interpreter

import (
	"fmt"
	"github.com/jatin-malik/make-thy-interpreter/evaluator"
	"github.com/jatin-malik/make-thy-interpreter/lexer"
	"github.com/jatin-malik/make-thy-interpreter/object"
	"github.com/jatin-malik/make-thy-interpreter/parser"
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
	env := object.NewEnvironment(nil)
	obj := evaluator.Eval(prg, env)
	if obj != nil {
		fmt.Println(obj.Inspect())
	}
}
