package vm

import (
	"encoding/binary"
	"fmt"
	"github.com/jatin-malik/yal/bytecode"
	"strings"
	"testing"

	"github.com/jatin-malik/yal/compiler"
	"github.com/jatin-malik/yal/lexer"
	"github.com/jatin-malik/yal/object"
	"github.com/jatin-malik/yal/parser"
)

var DEBUG = false // TODO: Control this via configuration

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

func TestClosures(t *testing.T) {
	tests := []struct {
		input, expected string
	}{

		// Simple closure capturing a variable
		{`let adder = fn(x) { fn(y) { x + y } }; let addTwo = adder(2); addTwo(3)`, "5"},

		// Nested closures capturing outer variables
		{`let outer = fn(x) { fn(y) { fn(z) { x + y + z } } }; let mid = outer(1); let inner = mid(2); inner(3)`, "6"},

		// Multiple closures capturing the same outer variable
		{`let makePair = fn(x) { fn() { x } }; let a = makePair(10); let b = makePair(20); a()`, "10"},
		{`let makePair = fn(x) { fn() { x } }; let a = makePair(10); let b = makePair(20); b()`, "20"},

		// Closure capturing a function argument
		{`let apply = fn(f, x) { f(x) }; let mulTwo = fn(x) { x * 2 }; apply(mulTwo, 5)`, "10"},

		// Closure that returns another function
		{`let makeIncrementer = fn(x) { fn(y) { x + y } }; let incFive = makeIncrementer(5); incFive(10)`, "15"},

		// Closure capturing an outer function return value
		{`let outer = fn(x) { let y = x + 2; fn() { y * 2 } }; let f = outer(3); f()`, "10"},

		// Closure as an argument to another function
		{`let twice = fn(f, x) { f(f(x)) }; let addOne = fn(x) { x + 1 }; twice(addOne, 5)`, "7"},

		// Closure with different function scopes
		{`let a = 10; let outer = fn() { let b = 20; fn() { a + b } }; let f = outer(); f()`, "30"},

		// Closure that captures a local variable but doesn't modify it
		{`let outer = fn(x) { let y = x * 2; fn() { y } }; let f = outer(4); f()`, "8"},

		{
			`
			let newAdderOuter = fn(a, b) {
				let c = a + b;
				fn(d) {
					let e = d + c;
					fn(f) { e + f; };
				};
			};

			let newAdderInner = newAdderOuter(1, 2);
			let adder = newAdderInner(3);
			adder(8);`,

			"14",
		},

		{
			`
			let a = 1;
			let newAdderOuter = fn(b) {
			fn(c) {
			fn(d) { a + b + c + d };
			};
			};
			let newAdderInner = newAdderOuter(2);
			let adder = newAdderInner(3);
			adder(8);
			`,
			"14",
		},

		{
			`
			let newClosure = fn(a, b) {
			let one = fn() { a; };
			let two = fn() { b; };
			fn() { one() + two(); };
			};
			let closure = newClosure(9, 90);
			closure();
			`,
			"99",
		},
	}

	runTests(t, tests)
}

func TestNestedLocalBindings(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		// Simple nested let bindings
		{`let f = fn() { let a = 5; let b = 10; a + b }; f()`, "15"},

		// Nested let bindings inside a function
		{`let f = fn() { let a = 2; let g = fn() { let b = 3; a * b }; g() }; f()`, "6"},

		// Nested bindings within multiple function calls
		{`let f = fn(x) { let a = x * 2; let g = fn(y) { let b = y + 3; a + b }; g(4) }; f(5)`, "17"},

		// Deeply nested let bindings
		{`let f = fn() { 
			let a = 2; 
			let g = fn() { 
				let b = 3; 
				let h = fn() { 
					let c = 4; 
					a + b + c 
				}; 
				h() 
			}; 
			g() 
		  }; 
		  f()`, "9"},

		// Shadowing outer variables
		{`let f = fn() { let x = 10; let g = fn() { let x = 20; x }; g() }; f()`, "20"},

		// Shadowing but accessing outer variables
		{`let f = fn() { let x = 10; let g = fn() { let x = 20; let h = fn() { x }; h() }; g() }; f()`, "20"},

		// Nested let bindings with return statements
		{`let f = fn() { 
			let a = 2; 
			let g = fn() { 
				let b = 3; 
				let h = fn() { 
					let c = 4; 
					return a + b + c; 
				}; 
				h() 
			}; 
			return g(); 
		  }; 
		  f()`, "9"},

		// Variables defined inside conditionals inside functions
		{`let f = fn(x) { 
			if (x > 10) { 
				let y = x * 2; 
				y 
			} else { 
				let y = x + 5; 
				y 
			} 
		  }; 
		  f(15)`, "30"},

		// Function inside a function using nested bindings
		{`let outer = fn(x) { 
			let y = x + 1; 
			let inner = fn() { let z = y * 2; z }; 
			inner() 
		  }; 
		  outer(3)`, "8"},

		// Function capturing a local variable and returning another function
		{`let makeAdder = fn(x) { 
			let y = x + 1; 
			fn(z) { y + z } 
		  }; 
		  let addFive = makeAdder(4); 
		  addFive(3)`, "8"},
	}

	runTests(t, tests)
}

func TestRecursiveClosures(t *testing.T) {
	tests := []struct {
		input, expected string
	}{

		{`
		let wrapper = fn() {
			let countDown = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					countDown(x - 1);
				}
			};
			countDown(1);
		};
		wrapper();
		`,
			"0"},

		// Basic recursive function
		{`let factorial = fn(n) { 
			if (n == 0) { return 1; } 
			else { return n * factorial(n - 1); } 
		  }; 
		  factorial(5)`, "120"},

		// Recursive closure function (capturing itself)
		{`let makeFactorial = fn() { 
			fn(f, n) { if (n == 0) { return 1; } else { return n * f(f, n - 1); } } 
		  }; 
		  let fact = makeFactorial(); 
		  fact(fact, 5)`, "120"},

		// Closure inside a recursive function
		{`let factorial = fn(n) { 
			let helper = fn(f, x) { 
				if (x == 0) { return 1; } 
				else { return x * f(f, x - 1); } 
			}; 
			helper(helper, n); 
		  }; 
		  factorial(5)`, "120"},

		// Recursive function calling a closure that captures a variable
		{`let recursiveAdder = fn(x) { 
			let adder = fn(y) { x + y }; 
			if (x == 0) { return 0; } 
			else { return adder(recursiveAdder(x - 1)); } 
		  }; 
		  recursiveAdder(5)`, "15"},

		// Nested recursive closures
		{`let makeRecSum = fn() { 
			fn(f, x) { 
				if (x == 0) { return 0; } 
				else { return x + f(f, x - 1); } 
			} 
		  }; 
		  let sum = makeRecSum(); 
		  sum(sum, 5)`, "15"},

		// Recursive closure capturing an outer variable
		{`let start = 2; 
		  let rec = fn(n) { 
			let inner = fn(f, x) { 
				if (x == 0) { return start; } 
				else { return f(f, x - 1) + 1; } 
			}; 
			inner(inner, n); 
		  }; 
		  rec(3)`, "5"},
	}

	runTests(t, tests)
}

// runTests is a helper that iterates over a list of test cases.
func runTests(t *testing.T, tests []struct {
	input, expected string
}) {
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj, err := testVM(tt.input, DEBUG)
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

func TestRecursiveFibonacci(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{`
		let fibonacci = fn(x) {
			if (x == 0) {
				return 0;
			} else {
				if (x == 1) {
					return 1;
				} else {
					fibonacci(x - 1) + fibonacci(x - 2);
				}
			}
		};
		fibonacci(15);
`,
			"610",
		},
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

func testVM(input string, debug bool) (object.Object, error) {
	compiler, err := testCompile(input)
	if err != nil {
		return nil, err
	}
	code := compiler.Output()

	if debug {
		fmt.Printf("Input - %s\n", input)
		fmt.Printf("Constant Pool - ")
		for _, obj := range code.ConstantPool {
			objType := obj.Type()
			if objType == object.CompiledFunctionObject {
				fmt.Println("Compiled Function -")
				prettyPrintInstructions(obj.(*object.CompiledFunction).Instructions)
			} else {
				fmt.Printf("%s:%v, ", objType, obj)
			}
		}
		fmt.Println()
		prettyPrintInstructions(code.Instructions)
	}

	vm := NewStackVM(code.Instructions, code.ConstantPool)
	if err := vm.Run(); err != nil {
		return nil, err
	}
	return vm.Top(), nil
}

func prettyPrintInstructions(instructions bytecode.Instructions) {
	fmt.Println("==Instructions===")
	for i := 0; i < len(instructions); {
		opCode := bytecode.OpCode(instructions[i])
		fmt.Print(opCode.String() + " ")

		switch opCode {
		case bytecode.OpPush:
			idx := binary.BigEndian.Uint16(instructions[i+1:])
			fmt.Println(idx)
			i += 1 + 2
		case bytecode.OpPushTrue, bytecode.OpPushFalse, bytecode.OpPushNull, bytecode.OpAdd, bytecode.OpSub,
			bytecode.OpMul, bytecode.OpDiv, bytecode.OpEqual, bytecode.OpNotEqual, bytecode.OpGT,
			bytecode.OpNegateBoolean, bytecode.OpNegateNumber, bytecode.OpIndex, bytecode.OpReturnValue,
			bytecode.OpGetCurrentClosure:
			i++
			fmt.Println()
		case bytecode.OpJumpIfFalse:
			jumpTo := binary.BigEndian.Uint16(instructions[i+1:])
			fmt.Println(jumpTo)
			i += 1 + 2
		case bytecode.OpJump:
			jumpTo := binary.BigEndian.Uint16(instructions[i+1:])
			fmt.Println(jumpTo)
			i += 1 + 2
		case bytecode.OpSetLocal:
			idx := binary.BigEndian.Uint16(instructions[i+1:])
			fmt.Println(idx)
			i += 1 + 2
		case bytecode.OpSetGlobal:
			idx := binary.BigEndian.Uint16(instructions[i+1:])
			fmt.Println(idx)
			i += 1 + 2
		case bytecode.OpGetLocal:
			idx := binary.BigEndian.Uint16(instructions[i+1:])
			fmt.Println(idx)
			i += 1 + 2
		case bytecode.OpGetGlobal:
			idx := binary.BigEndian.Uint16(instructions[i+1:])
			fmt.Println(idx)
			i += 1 + 2
		case bytecode.OpGetBuiltIn:
			idx := int(instructions[i+1])
			fmt.Println(idx)
			i += 1 + 1
		case bytecode.OpGetFree:
			idx := int(instructions[i+1])
			fmt.Println(idx)
			i += 1 + 1
		case bytecode.OpArray:
			count := binary.BigEndian.Uint16(instructions[i+1:])
			fmt.Println(count)
			i += 1 + 2
		case bytecode.OpHash:
			count := binary.BigEndian.Uint16(instructions[i+1:])
			fmt.Println(count)
			i += 1 + 2
		case bytecode.OpClosure:
			idx := binary.BigEndian.Uint16(instructions[i+1:])
			freeCount := int(instructions[i+3])
			fmt.Println(idx, " ", freeCount)
			i += 1 + 2 + 1
		case bytecode.OpCall:
			argsCount := int(instructions[i+1])
			fmt.Println(argsCount)
			i += 2
		default:
			fmt.Printf("unknown opcode: %d\n", opCode)
		}

	}
	fmt.Println("=========")
}
