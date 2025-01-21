// Package repl exposes the REPL functionality for the language.
package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/jatin-malik/make-thy-interpreter/lexer"
	"github.com/jatin-malik/make-thy-interpreter/token"
)

func Start(in io.Reader, out io.Writer) {
	prompt := ">> "
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, prompt)
		// Read
		if !scanner.Scan() {
			if scanner.Err() != nil {
				fmt.Fprintf(out, "scanning errored out")
			}
			return
		}
		input := scanner.Text()

		if strings.ToLower(input) == "quit" {
			return
		}

		// Eval
		l := lexer.New(input)
		for tok := l.NextToken(); tok.Type != token.EOF; {
			fmt.Fprintf(out, "%+v\n", tok)
			tok = l.NextToken()
		}
	}
}
