package processor

import (
	"fmt"
	"github.com/jatin-malik/yal/compiler"
	"github.com/jatin-malik/yal/evaluator"
	"github.com/jatin-malik/yal/lexer"
	"github.com/jatin-malik/yal/object"
	"github.com/jatin-malik/yal/parser"
	"github.com/jatin-malik/yal/vm"
)

func Process(input string, engine string) {
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

	var obj object.Object
	if engine == "eval" {
		env := object.NewEnvironment(nil)
		obj = evaluator.Eval(expandedAST, env)
	} else if engine == "vm" {
		compiler := compiler.New()
		err = compiler.Compile(expandedAST)
		if err != nil {
			fmt.Println(err)
			return
		}
		bytecode := compiler.Output()
		vm := vm.NewStackVM(bytecode.Instructions, bytecode.ConstantPool)
		err = vm.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		obj = vm.Top()
	}

	if obj != nil {
		fmt.Println(obj.Inspect())
	}
}
