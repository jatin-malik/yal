// Package repl exposes the REPL functionality for the language.
package repl

import (
	"errors"
	"fmt"
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
	constantPool := make([]object.Object, 0)
	globals := make([]object.Object, 100)

	multilineMode := false
	var buffer []string // Stores multi-line input

	for {

		line, err := rl.Readline()

		if errors.Is(err, readline.ErrInterrupt) {
			// Ctrl+C was pressed: Clear input and continue
			fmt.Println("KeyboardInterrupt")
			multilineMode = false
			rl.SetPrompt(prompt) // restore original prompt
			continue
		} else if err == io.EOF {
			// Handle Ctrl+D (EOF)
			fmt.Println("See ya.")
			break
		} else if err != nil {
			// Handle unexpected errors
			fmt.Println("Error:", err)
			break
		}

		if strings.ToLower(line) == "bye" {
			return
		}

		buffer = append(buffer, line)
		if isCompleteStatement(line, multilineMode) {
			multilineMode = false
			rl.SetPrompt(prompt) // restore original prompt
		} else {
			multilineMode = true
			rl.SetPrompt("..") // Change prompt for multi-line
			continue
		}

		input := strings.Join(buffer, "\n")
		buffer = nil // Reset buffer
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

		var obj object.Object
		if engine == "eval" {
			obj = evaluator.Eval(expandedAST, env)
		} else if engine == "vm" {
			compiler := compiler.New(compiler.WithSymbolTable(symTable), compiler.WithConstantPool(constantPool))
			err = compiler.Compile(expandedAST)
			if err != nil {
				_, _ = io.WriteString(out, err.Error()+"\n")
				continue
			}

			bytecode := compiler.Output()
			constantPool = bytecode.ConstantPool // updates shared constant pool
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

func isCompleteStatement(line string, multilineMode bool) bool {
	if !multilineMode {
		lex := lexer.New(line)
		parser := parser.New(lex)
		parser.ParseProgram()
		if len(parser.Errors) != 0 {
			for _, msg := range parser.Errors {
				if strings.HasPrefix(msg, "incomplete") {
					return false
				}
			}
		}
		return true
	}

	// Inside multiline mode
	if line == "" { // Exit multiline mode when user presses RET on empty line
		return true
	}

	return false

}
