// Package repl exposes the REPL functionality for the language.
package repl

import (
	"bufio"
	"github.com/jatin-malik/yal/evaluator"
	"github.com/jatin-malik/yal/object"
	"github.com/jatin-malik/yal/parser"
	"io"
	"strings"

	"github.com/jatin-malik/yal/lexer"
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

		if strings.ToLower(input) == "bye" {
			return
		}

		// Eval
		l := lexer.New(input)
		p := parser.New(l)
		prg := p.ParseProgram()
		if len(p.Errors) != 0 {
			for _, msg := range p.Errors {
				io.WriteString(out, msg+"\n")
			}
			continue
		}
		obj := evaluator.Eval(prg, env)
		if obj != nil {
			io.WriteString(out, obj.Inspect())
			io.WriteString(out, "\n")
		}

	}
}
