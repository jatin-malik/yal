package main

import (
	"flag"
	"fmt"
	"github.com/jatin-malik/yal/compiler"
	"github.com/jatin-malik/yal/evaluator"
	"github.com/jatin-malik/yal/lexer"
	"github.com/jatin-malik/yal/object"
	"github.com/jatin-malik/yal/parser"
	"github.com/jatin-malik/yal/vm"
	"time"
)

var engine = flag.String("engine", "", "engine to use ( vm or eval )")

var benchmarkInput = `
	let fibonacci = fn(x) {
		if (x == 0) {
			0
		} else {
			if (x == 1) {
				return 1;
			} else {
				fibonacci(x - 1) + fibonacci(x - 2);
			}
		}
	};
	fibonacci(35);`

func main() {
	flag.Parse()

	if *engine == "" {
		flag.Usage()
		return
	}

	var duration time.Duration

	lexer := lexer.New(benchmarkInput)
	parser := parser.New(lexer)
	prg := parser.ParseProgram()

	if *engine == "eval" {
		// use tree walking interpreter
		env := object.NewEnvironment(nil)
		start := time.Now()
		obj := evaluator.Eval(prg, env)
		duration = time.Since(start)
		fmt.Println(obj.Inspect())
	} else if *engine == "vm" {
		// use bytecode compiler and vm
		compiler := compiler.New()
		start := time.Now()
		compiler.Compile(prg)
		code := compiler.Output()
		vm := vm.NewStackVM(code.Instructions, code.ConstantPool)
		vm.Run()
		duration = time.Since(start)
		fmt.Println(vm.Top().Inspect())
	} else {
		fmt.Println("Unknown engine")
		return
	}

	fmt.Printf("Execution took %d ms.\n", duration.Milliseconds())
}
