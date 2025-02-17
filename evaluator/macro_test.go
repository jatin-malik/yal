package evaluator

import (
	"fmt"
	"github.com/jatin-malik/yal/ast"
	"github.com/jatin-malik/yal/lexer"
	"github.com/jatin-malik/yal/object"
	"github.com/jatin-malik/yal/parser"
	"testing"
)

func TestMacroDefinitions(t *testing.T) {
	tests := []struct {
		input            string
		expectedProgram  string
		expectedErrorMsg string
	}{
		{`
			let m = macro(x,y){x-y};
			2`,
			"2",
			"",
		},

		{`
			let m = macro(x,y){x-y};
			let n = macro(y,z){z-y};
			4
			let p = macro(a,b){a+b};`,
			"4",
			"",
		},

		{`
		if ( 5 > 2 ) { 
		let x = macro(x,y){x*y};	
		1  
		let y = macro(x,y){x*y};	
		}else{ 
			let x = macro(x,y){x*y};	
			0 
			let y = macro(x,y){x*y};	
		};`,
			"if ( 5 > 2 ){ 1 } else { 0 }",
			"",
		},

		{`
		let add = fn(a,b){
			let x = macro(a,b){a-b};	
			a+b
		};
		add(10,5)`,
			"let add = fn (a, b) { ( a + b ) };add(10, 5)",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			testMacro(t, tt.input, tt.expectedProgram, tt.expectedErrorMsg)
		})
	}
}

func TestMacroExpansion(t *testing.T) {
	tests := []struct {
		input            string
		expectedProgram  string
		expectedErrorMsg string
	}{
		// Basic macro expansion with arithmetic operation
		{`
			let minus = macro(x,y) { quote(unquote(x) - unquote(y)) };
			minus(4, 2)`,
			"( 4 - 2 )",
			"",
		},

		// Macro with conditionals
		{`
			let conditional = macro(a, b) { quote(if (unquote(a) > 0) {unquote(a)} else{ unquote(b)}) };
			conditional(5, 10)`,
			"if ( 5 > 0 ){ 5 } else { 10 }",
			"",
		},

		// Macros inside function calls
		{`
			let callMacro = macro(f, arg) { quote(unquote(f)(unquote(arg))) };
			callMacro(double, 4)`,
			"double(4)",
			"",
		},

		// Macro expanding a function definition
		{`
			let makeAdder = macro(x) { quote(fn(y) { unquote(x) + y }) };
			makeAdder(5);`,
			"fn (y) { ( 5 + y ) }",
			"",
		},

		// Edge case: macro without parameters
		{`
			let constant = macro() { quote(42) };
			constant()`,
			"42",
			"",
		},

		// Edge case: macro with unused parameters
		{`
			let ignoreArg = macro(x) { quote(100) };
			ignoreArg(50)`,
			"100",
			"",
		},

		// Edge case: macro with empty body
		{`
			let empty = macro(x) { quote() };
			empty(5)`,
			"",
			"macro expansion error: quote supports only 1 argument",
		},

		// Edge case: using an undefined macro
		{`
			undefinedMacro(5)`,
			"undefinedMacro(5)", // Should remain as is because it isn't defined
			"",
		},

		{
			`
			let ternary = macro(condition, trueExpr, falseExpr) {
				quote(if (unquote(condition)) { unquote(trueExpr) } else { unquote(falseExpr) })
			};
			ternary(true,1,0)
			ternary(false,1,0)`,
			"if true{ 1 } else { 0 }if false{ 1 } else { 0 }",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			testMacro(t, tt.input, tt.expectedProgram, tt.expectedErrorMsg)
		})
	}
}

func testMacro(t *testing.T, input, expected string, expectedErrorMsg string) {
	l := lexer.New(input)
	p := parser.New(l)
	prg := p.ParseProgram()
	if len(p.Errors) != 0 {
		fmt.Printf("parser errors for input: %q\n", input)
		for _, err := range p.Errors {
			fmt.Println("\t" + err)
		}
		t.FailNow()
	}

	// Macro expansion stage validations
	macroEnv := object.NewEnvironment(nil)
	expandedAST, err := ExpandMacro(prg, macroEnv)
	if err != nil {
		if expectedErrorMsg == "" {
			t.Fatal(err)
		}
		if expectedErrorMsg != err.Error() {
			t.Fatalf("error message does not match.\nexpected: %q\nactual: %q", expectedErrorMsg, err.Error())
		}

		return
	}

	prg = expandedAST.(*ast.Program)
	if prg.String() != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, prg.String())
	}
}
