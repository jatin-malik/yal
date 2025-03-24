package vm

import (
	"strings"
	"testing"

	"github.com/jatin-malik/yal/compiler"
	"github.com/jatin-malik/yal/lexer"
	"github.com/jatin-malik/yal/object"
	"github.com/jatin-malik/yal/parser"
)

// Arithmetic and Nested Expressions
func TestArithmeticExpressions(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// Basic arithmetic
		{"1+2", "3"},
		{"6-2", "4"},
		{"3*4", "12"},
		{"6/3", "2"},

		// Deeply nested expressions
		{"(((1+2)*3)-4)/2", "2"},
		{"(5*(3+(2*2)))", "35"},
		{"(6/(2*(1+2)))", "1"},
		{"((2+3)*(4+(5-2)))", "35"},

		// Extra parentheses (edge case)
		{"(((1)))", "1"},
		{"(0+((1+2)*3))", "9"},
		{"(100/(10/(2*5)))", "100"},
		{"((8-6)*(3+(4/2)))", "10"},
	}

	runTests(t, tests)
}

// Comparison Operators and Nested Comparisons
func TestComparisons(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// Comparison Operators
		{"1 == 1", "true"},
		{"1 == 2", "false"},
		{"5 != 3", "true"},
		{"4 != 4", "false"},
		{"10 > 5", "true"},
		{"10 > 10", "false"},
		{"3 < 7", "true"},
		{"7 < 3", "false"},

		// Nested Comparisons
		{"(1+2) == (3)", "true"},
		{"(10-5) > (2+2)", "true"},
		{"(2*3) < (10-1)", "true"},
		{"(4/2) != (2-1)", "true"},
		{"(6/3) == (2-0)", "true"},
		{"((2+3)*2) > ((4+1)*2)", "false"},
		{"((5*2)/2) == ((4+1))", "true"},
	}

	runTests(t, tests)
}

// Prefix and Negative Expressions
func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// Prefix operators
		{"!true", "false"},
		{"!false", "true"},
		{"!1", "false"},
		{"!0", "false"},
		{"!!true", "true"},
		{"!!false", "false"},
		{"!!1", "true"},
		{"!!0", "true"},

		// Negative expressions
		{"-5", "-5"},
		{"-(-5)", "5"},
		{"-(3+2)", "-5"},
		{"-1 + 2", "1"},
		{"-(2*3)", "-6"},
		{"-(10-5)", "-5"},
		{"-(4/2)", "-2"},
		{"-(-(-3))", "-3"},
	}

	runTests(t, tests)
}

// Conditionals and Let Statements
func TestConditionalsAndLetStatements(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// Conditionals
		{`if (5 > 3) { 10 } else { 20 }`, "10"},
		{`if (5 > 7) { 10 } else { 20 }`, "20"},
		{`if (5 > 3) { 10 } else { 20 };5`, "5"},
		{`if (5 > 3) { 10 }`, "10"},
		{`if (5 > 3) { 10 };6+1`, "7"},
		{`if (5 > 8) { 10 };2+1`, "3"},
		{`if (5 > 8) { 10 }`, "null"},

		// Let statements and variable usage
		{`let x = 5 ; x`, "5"},
		{`let x = 5 ; x+2`, "7"},
		{`let x = 5 ; let x = 10; x + 4`, "14"},
	}

	runTests(t, tests)
}

// String Literals and Concatenation
func TestStrings(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// String literals
		{`"hello"`, "hello"},
		{`"world"`, "world"},

		// String concatenation
		{`"hello" + " " + "world"`, "hello world"},
		{`"foo" + "bar"`, "foobar"},
		{`"Go" + "lang"`, "Golang"},
	}

	runTests(t, tests)
}

// Arrays and Index Expressions
func TestArraysAndIndexExpressions(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// Array literals and indexing
		{`[1, 2, 3]`, "[1, 2, 3]"},
		{`[10, 20, 30][1]`, "20"},
		{`let arr = [5, 10, 15]; arr[2]`, "15"},
		{`[1 + 1, 2 * 2, 3 - 1]`, "[2, 4, 2]"},
		{`let x = [1, 2, 3]; x[0] + x[2]`, "4"},
	}

	runTests(t, tests)
}

// Hash Literals
func TestHashLiterals(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{`{"key": "value"}`, "{key:value}"},
		{`{"name": "Alice", "age": 25}["name"]`, "Alice"},
		{`let h = {"a": 1, "b": 2}; h["b"]`, "2"},
		{`{"x": 10, "y": 20}["y"]`, "20"},
		{`let m = {1: "one", 2: "two"}; m[1]`, "one"},
	}

	runTests(t, tests)
}

// Function Calls and Definitions
func TestFunctionCalls(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// Simple function calls
		{`let f = fn() { 5 }; f()`, "5"},
		{`let f = fn() { return 10; }; f()`, "10"},
		{`let f = fn() { 1+2 }; f()`, "3"},
		{`let f = fn() { 2*3 }; f()`, "6"},
		{`let f = fn() {}; f()`, "null"},
		{`let f = fn() { if (true) { return 42; } }; f()`, "42"},
		{`let f = fn() { if (false) { return 42; } }; f()`, "null"},
		{`let f = fn() { 5 }; let g = fn() { f() }; g()`, "5"},
		{`let f = fn() { 1+2 }; let g = fn() { f() * 2 }; g()`, "6"},

		// Function aliasing
		{`let f = fn() { 10 }; let g = f; g()`, "10"},
		{`let f = fn() { 20 }; let g = fn() { f() }; let h = g; h()`, "20"},

		// Functions returning functions
		{`let f = fn() { fn() { 99 } }; let g = f(); g()`, "99"},
		{`let f = fn() { fn() {} }; let g = f(); g()`, "null"},

		// Multiple function calls (only first result matters)
		{`let f = fn() { 5 }; f(); f(); f()`, "5"},

		// Functions in expressions
		{`let f = fn() { 4 }; f() + 2`, "6"},
		{`let f = fn() { 10 }; let g = fn() { f() + f() }; g()`, "20"},

		// Function reassignments
		{`let f = fn() { 10 }; let f = fn() { 20 }; f()`, "20"},

		// Function calls in if-expressions
		{`let f = fn() { 7 }; if (true) { f() }`, "7"},
		{`let f = fn() { 8 }; if (false) { f() } else { 12 }`, "12"},
		{`let f = fn() { if (true) { return 30; } else { return 40; } }; f()`, "30"},
	}

	runTests(t, tests)
}

// Local Variable and Scoping Tests
func TestLocalVariableScoping(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// Basic local variable usage
		{`let f = fn() { let x = 5; x }; f()`, "5"},
		{`let f = fn() { let x = 5; let x = 10; x }; f()`, "10"},
		{`let f = fn() { let x = 10; return x; }; f()`, "10"},
		{`let f = fn() { let x = 2; x * 3 }; f()`, "6"},

		// Multiple local variables
		{`let f = fn() { let a = 3; let b = 4; a + b }; f()`, "7"},
		{`let f = fn() { let x = 2; let y = x + 5; y }; f()`, "7"},

		// Shadowing global variables
		{`let x = 100; let f = fn() { let x = 5; x }; f()`, "5"},
		{`let x = 100; let f = fn() { let x = 5; x }; f(); x`, "100"},
		{`let x = 100; let f = fn() { let x = x + 5; x }; f()`, "105"},
		{`let x = 100; let f = fn() { x + 5 }; f()`, "105"},
		{`let x = 50; let f = fn() { let x = x + 10; x }; f()`, "60"},

		// Conditionals and local variables
		{`let f = fn() { let x = 0; if (true) { let x = 10; x } else { x } }; f()`, "10"},
		{`let f = fn() { let x = 0; if (false) { let x = 10; x } else { let x = 20; x } }; f()`, "20"},

		// Function returning function (nested scoping)
		{`let returnsOneReturner = fn() {
				let returnsOne = fn() { 1; };
				returnsOne;};
				returnsOneReturner()();`, "1"},
	}

	runTests(t, tests)
}

// Function Arguments Tests
func TestFunctionArguments(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{`let f = fn(x) { x }; f(5)`, "5"},
		{`let f = fn(x, y) { x + y }; f(3, 4)`, "7"},
		{`let f = fn(x, y) { x * y }; f(2, 3)`, "6"},
		{`let f = fn(x) { x * 2 }; f(10)`, "20"},
		{`let f = fn(x) { x }; f(1 + 2)`, "3"},
		{`let f = fn(x, y) { x + y }; f(2 * 3, 4 + 1)`, "11"},
		{`let f = fn(x) { x * 2 }; let g = fn(y) { f(y) + 1 }; g(3)`, "7"},
		{`let f = fn(x) {}; f(5)`, "null"},
		{`let f = fn() { return 42; }; f(5)`, "42"},
		{`let f = fn(x) { x + 1 }; let g = fn(y) { f(y) * 2 }; g(3)`, "8"},
		{`let add = fn(x, y) { x + y }; let square = fn(n) { n * n };
		  let h = fn(a, b) { square(add(a, b)) }; h(2, 3)`, "25"},
		{`let f = fn(x) { if (x > 10) { return x; } }; f(5)`, "null"},
		{`let f = fn(x) { if (x > 10) { return x; } else { return 0; } }; f(15)`, "15"},
	}

	runTests(t, tests)
}

// runTests is a helper that iterates over a list of test cases.
func runTests(t *testing.T, tests []struct {
	input, expected string
}) {
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj, err := testVM(tt.input)
			if err != nil {
				// Check if an error was expected.
				if strings.HasPrefix(tt.expected, "error:") {
					if err.Error() != strings.TrimPrefix(tt.expected, "error: ") {
						t.Errorf("Expected error %s, got %s", tt.expected, err.Error())
					}
					return
				}
				t.Fatal(err)
			}
			if obj.Inspect() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, obj.Inspect())
			}
		})
	}
}

// Built-in functions: len, first, last, rest, push

func TestEvalBuiltInFuncLen(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// Basic String Length
		{`len("hello")`, "5"},
		{`len("world")`, "5"},
		{`len("a")`, "1"},
		{`len("")`, "0"},

		// Arrays (or lists)
		{`len([1, 2, 3, 4])`, "4"},
		{`len([10, 20, 30])`, "3"},
		{`len([true, false])`, "2"},
		{`len([])`, "0"},

		// Mixed Types
		{`len([1, "a", true])`, "3"},
		{`len([true, "string", 100])`, "3"},

		// Nested Arrays
		{`len([ [1, 2], [3, 4] ])`, "2"},
		{`len([ [1, 2, 3], [4, 5] ])`, "2"},
		{`len([["a", "b"], ["c", "d", "e"]])`, "2"},

		// Concatenation with Length
		{`len("hello" + " world")`, "11"},
		{`len("good" + "bye" + "!")`, "8"},
		{`len("a" + "b" + "c")`, "3"},

		// Other Edge Cases
		{`len("a" + "")`, "1"},
		{`len("") + len("test")`, "4"},
	}

	runTests(t, tests)
}

func TestEvalBuiltInFuncFirst(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// Valid Arrays
		{`first([1, 2, 3])`, "1"},
		{`first(["hello", "world"])`, "hello"},
		{`first([true, false])`, "true"},
		{`first([1, "hello", false])`, "1"},

		// Edge Cases
		{`first([])`, "error: empty array"},

		// Single Element Arrays
		{`first([42])`, "42"},
		{`first(["only"])`, "only"},
		{`first([false])`, "false"},

		// Nested Arrays
		{`first([[1, 2], [3, 4]])`, "[1, 2]"},
		{`first([["a", "b"], ["c", "d"]])`, "[a, b]"},

		// Mixed Arrays
		{`first([1, [2, 3], "hello"])`, "1"},
		{`first([true, [false, true]])`, "true"},

		// Invalid Cases
		{`first("hello")`, "error: first(): type STRING not supported"},
		{`first(42)`, "error: first(): type INTEGER not supported"},
		{`first(true)`, "error: first(): type BOOLEAN not supported"},
	}

	runTests(t, tests)
}

func TestEvalBuiltInFuncLast(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// Valid Arrays
		{`last([1, 2, 3])`, "3"},
		{`last(["hello", "world"])`, "world"},
		{`last([true, false])`, "false"},
		{`last([1, "hello", false])`, "false"},

		// Edge Cases
		{`last([])`, "error: empty array"},

		// Single Element Arrays
		{`last([42])`, "42"},
		{`last(["only"])`, "only"},
		{`last([false])`, "false"},

		// Nested Arrays
		{`last([[1, 2], [3, 4]])`, "[3, 4]"},
		{`last([["a", "b"], ["c", "d"]])`, "[c, d]"},

		// Mixed Arrays
		{`last([1, [2, 3], "hello"])`, "hello"},
		{`last([true, [false, true]])`, "[false, true]"},

		// Invalid Cases
		{`last("hello")`, "error: last(): type STRING not supported"},
		{`last(42)`, "error: last(): type INTEGER not supported"},
		{`last(true)`, "error: last(): type BOOLEAN not supported"},
	}

	runTests(t, tests)
}

func TestEvalBuiltInFuncRest(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// Valid Arrays
		{`rest([1, 2, 3])`, "[2, 3]"},
		{`rest(["hello", "world", "foo"])`, "[world, foo]"},
		{`rest([true, false, true])`, "[false, true]"},
		{`rest([1, "hello", false])`, "[hello, false]"},

		// Edge Cases
		{`rest([])`, "error: empty array"},

		// Single Element Arrays
		{`rest([42])`, "[]"},
		{`rest(["only"])`, "[]"},
		{`rest([false])`, "[]"},

		// Nested Arrays
		{`rest([[1, 2], [3, 4], [5, 6]])`, "[[3, 4], [5, 6]]"},
		{`rest([["a", "b"], ["c", "d"], ["e", "f"]])`, "[[c, d], [e, f]]"},

		// Mixed Arrays
		{`rest([1, [2, 3], "hello", true])`, "[[2, 3], hello, true]"},
		{`rest([true, [false, true], 42])`, "[[false, true], 42]"},

		// Invalid Cases
		{`rest("hello")`, "error: rest(): type STRING not supported"},
		{`rest(42)`, "error: rest(): type INTEGER not supported"},
		{`rest(true)`, "error: rest(): type BOOLEAN not supported"},
	}

	runTests(t, tests)
}

func TestEvalBuiltInFuncPush(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// Valid Arrays
		{`push([1, 2], 3)`, "[1, 2, 3]"},
		{`push(["hello", "world"], "test")`, "[hello, world, test]"},
		{`push([true, false], true)`, "[true, false, true]"},
		{`push([1, "hello", false], 100)`, "[1, hello, false, 100]"},

		// Edge Cases
		{`push([], 42)`, "[42]"},

		// Single Element Arrays
		{`push([42], 100)`, "[42, 100]"},
		{`push(["only"], "more")`, "[only, more]"},
		{`push([false], true)`, "[false, true]"},

		// Nested Arrays
		{`push([[1, 2], [3, 4]], [5, 6])`, "[[1, 2], [3, 4], [5, 6]]"},

		// Mixed Arrays
		{`push([1, [2, 3], "hello"], [4, 5])`, "[1, [2, 3], hello, [4, 5]]"},
		{`push([true, [false, true]], ["more"])`, "[true, [false, true], [more]]"},

		// Invalid Cases
		{`push("hello", 42)`, "error: push(): type STRING not supported"},
		{`push(42, 100)`, "error: push(): type INTEGER not supported"},
		{`push(true, false)`, "error: push(): type BOOLEAN not supported"},
	}

	runTests(t, tests)
}

// Helper functions

func testCompile(input string) (*compiler.Compiler, error) {
	lex := lexer.New(input)
	parser := parser.New(lex)
	program := parser.ParseProgram()

	compiler := compiler.New()
	if err := compiler.Compile(program); err != nil {
		return nil, err
	}
	return compiler, nil
}

func testVM(input string) (object.Object, error) {
	compiler, err := testCompile(input)
	if err != nil {
		return nil, err
	}
	bytecode := compiler.Output()
	vm := NewStackVM(bytecode.Instructions, bytecode.ConstantPool)
	if err := vm.Run(); err != nil {
		return nil, err
	}
	return vm.Top(), nil
}
