// Package repl exposes the REPL functionality for the language.
package repl

import (
	"bufio"
	"github.com/jatin-malik/make-thy-interpreter/parser"
	"io"
	"strings"

	"github.com/jatin-malik/make-thy-interpreter/lexer"
)

func Start(in io.Reader, out io.Writer) {
	prompt := ">> "
	scanner := bufio.NewScanner(in)
	for {
		_, err := io.WriteString(out, prompt)
		if err != nil {
			return
		}
		// Read
		if !scanner.Scan() {
			if scanner.Err() != nil {
				_, err := io.WriteString(out, "scanning errored out")
				if err != nil {
					return
				}
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
		_, err = io.WriteString(out, prg.String()+"\n")
		if err != nil {
			return
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
