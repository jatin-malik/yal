// Package repl exposes the REPL functionality for the language.
package repl

import (
	"github.com/chzyer/readline"
	"github.com/jatin-malik/yal/compiler"
	"github.com/jatin-malik/yal/evaluator"
	"github.com/jatin-malik/yal/lexer"
	"github.com/jatin-malik/yal/object"
	"github.com/jatin-malik/yal/parser"
	"github.com/jatin-malik/yal/vm"
	"io"
	"strings"
)

func Start(in io.Reader, out io.Writer, engine string) {

	prompt := ">> "
	rl, err := readline.New(prompt)
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	macroEnv := object.NewEnvironment(nil) // shared scope across all macro expansions

	env := object.NewEnvironment(nil) // shared scope across all REPL statements evaluation

	symTable := compiler.NewSymbolTable(nil)
	globals := make([]object.Object, 100)

	for {

		line, err := rl.Readline()
		if err != nil {
			break
		}

		if strings.ToLower(line) == "bye" {
			return
		}

		l := lexer.New(line)
		p := parser.New(l)
		prg := p.ParseProgram()
		if len(p.Errors) != 0 {
			for _, msg := range p.Errors {
				_, _ = io.WriteString(out, msg+"\n")
			}
			continue
		}

		expandedAST, err := evaluator.ExpandMacro(prg, macroEnv)
		if err != nil {
			_, _ = io.WriteString(out, err.Error()+"\n")
			continue
		}

		var obj object.Object
		if engine == "eval" {
			obj = evaluator.Eval(expandedAST, env)
		} else if engine == "vm" {
			compiler := compiler.New(compiler.WithSymbolTable(symTable))
			err = compiler.Compile(expandedAST)
			if err != nil {
				_, _ = io.WriteString(out, err.Error()+"\n")
				continue
			}

			bytecode := compiler.Output()
			vm := vm.NewStackVM(bytecode.Instructions, bytecode.ConstantPool, vm.WithGlobals(globals))
			err = vm.Run()
			if err != nil {
				_, _ = io.WriteString(out, err.Error()+"\n")
				continue
			}
			obj = vm.Top()
		}

		if obj != nil {
			_, _ = io.WriteString(out, obj.Inspect())
			_, _ = io.WriteString(out, "\n")
		}
	}
}
