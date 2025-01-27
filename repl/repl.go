// Package repl exposes the REPL functionality for the language.
package repl

import (
	"bufio"
	"github.com/jatin-malik/make-thy-interpreter/evaluator"
	"github.com/jatin-malik/make-thy-interpreter/object"
	"github.com/jatin-malik/make-thy-interpreter/parser"
	"io"
	"strings"

	"github.com/jatin-malik/make-thy-interpreter/lexer"
)

func Start(in io.Reader, out io.Writer) {
	prompt := ">> "
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment(nil) // shared scope across all REPL statements
	for {
		io.WriteString(out, prompt)
		// Read
		if !scanner.Scan() {
			if scanner.Err() != nil {
				io.WriteString(out, "scanning errored out")
			}
			return
		}
		input := scanner.Text()

		if strings.ToLower(input) == "quit" {
			return
		}

		// Eval
		l := lexer.New(input)
		p := parser.New(l)
		prg := p.ParseProgram()
		if len(p.Errors) != 0 {
			printParserErrors(out, p.Errors)
			continue
		}
		obj := evaluator.Eval(prg, env)
		if obj != nil {
			io.WriteString(out, obj.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		_, err := io.WriteString(out, "\t"+msg+"\n")
		if err != nil {
			return
		}
	}
}
