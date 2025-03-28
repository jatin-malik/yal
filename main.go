package main

import (
	"flag"
	"fmt"
	"github.com/jatin-malik/yal/processor"
	"github.com/jatin-malik/yal/repl"
	"os"
	"os/user"
)

var engine = flag.String("engine", "vm", "engine to use ( vm or eval )")

func main() {
	flag.Parse()

	if *engine != "vm" && *engine != "eval" {
		fmt.Fprintf(os.Stderr, "Usage: %s [-engine vm|eval]", os.Args[0])
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) > 0 {
		filename := args[0]
		processFile(filename, *engine)
	} else {
		startREPL(*engine)
	}
}

// startREPL starts the REPL for the language with the provided engine mode
func startREPL(engine string) {
	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s. Welcome to the yal language REPL. Executing in %s mode\n",
		currentUser.Username, engine)
	fmt.Printf("To quit the REPL, say bye.\n")

	// Start the REPL
	repl.Start(os.Stdin, os.Stdout, engine)
}

// processFile reads the file and processes it with the provided engine mode
func processFile(filename string, engine string) {
	fmt.Printf("[Processing in %s mode]\n", engine)
	// Read the file contents
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	processor.Process(string(data), engine)
}
