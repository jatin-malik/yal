package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/jatin-malik/make-thy-interpreter/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s. Welcome to the Monkey language REPL\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
