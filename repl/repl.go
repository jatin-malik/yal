// Package repl exposes the REPL functionality for the language.
package repl

import (
	"bufio"
	"github.com/jatin-malik/yal/compiler"
	"github.com/jatin-malik/yal/vm"
	"io"
	"strings"

	"github.com/jatin-malik/yal/evaluator"
	"github.com/jatin-malik/yal/object"
	"github.com/jatin-malik/yal/parser"

	"github.com/jatin-malik/yal/lexer"
)

func Start(in io.Reader, out io.Writer) {
	prompt := ">> "
	scanner := bufio.NewScanner(in)
	macroEnv := object.NewEnvironment(nil) // shared scope across all macro expansions
	//env := object.NewEnvironment(nil)      // shared scope across all REPL statements evaluation
	for {
		_, _ = io.WriteString(out, prompt)
		// Read
		if !scanner.Scan() {
			if scanner.Err() != nil {
				_, _ = io.WriteString(out, "scanning errored out")
			}
			return
		}
		input := scanner.Text()

		if strings.ToLower(input) == "bye" {
			return
		}

		// Eval
		l := lexer.New(input)
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

		//obj := evaluator.Eval(expandedAST, env)

		compiler := compiler.New()
		err = compiler.Compile(expandedAST)
		if err != nil {
			_, _ = io.WriteString(out, err.Error()+"\n")
			continue
		}

		bytecode := compiler.Output()
		vm := vm.NewStackVM(bytecode.Instructions, bytecode.ConstantPool)
		err = vm.Run()
		if err != nil {
			_, _ = io.WriteString(out, err.Error()+"\n")
			continue
		}
		obj := vm.Top()
		if obj != nil {
			_, _ = io.WriteString(out, obj.Inspect())
			_, _ = io.WriteString(out, "\n")
		}
	}
}
