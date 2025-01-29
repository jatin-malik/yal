package main

import (
	"fmt"
	"github.com/jatin-malik/make-thy-interpreter/interpreter"
	"github.com/jatin-malik/make-thy-interpreter/repl"
	"os"
	"os/user"
)

func main() {
	if len(os.Args) > 1 {
		// If there's a filename argument, interpret the whole file
		filename := os.Args[1]
		interpretFile(filename)
	} else {
		// Otherwise, start the REPL
		startREPL()
	}
}

// startREPL starts the REPL for the yal language
func startREPL() {
	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s. Welcome to the yal language REPL.\n", currentUser.Username)
	fmt.Printf("To quit the REPL, say bye.\n")

	// Start the REPL
	repl.Start(os.Stdin, os.Stdout)
}

// interpretFile reads the file and interprets it
func interpretFile(filename string) {
	// Read the file contents
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	interpreter.Interpret(string(data))
}
