package evaluator

import (
	"errors"
	"fmt"
	"github.com/jatin-malik/yal/lexer"
	"github.com/jatin-malik/yal/object"
	"github.com/jatin-malik/yal/parser"
	"testing"
)

func TestEvalIntegerLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-100", -100},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			testIntegerObject(t, obj, tt.expected)
		})
	}
}

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"quote(4+4)", "( 4 + 4 )"},
		{"quote(2 * (3 + 5))", "( 2 * ( 3 + 5 ) )"},
		{"quote(a + b)", "( a + b )"},
		{"quote(fn(x) { x + 1 })", "fn (x) { ( x + 1 ) }"},
		{"quote(quote(1 + 2))", "quote(( 1 + 2 ))"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			testQuoteObject(t, obj, tt.expected)
		})
	}
}

func TestUnquote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"quote(unquote(4+4))", "8"},
		{"quote(unquote(2 * (3 + 5)))", "16"},
		{"quote(1 + unquote(2 + 3))", "( 1 + 5 )"},
		{"quote(fn(x) { unquote(2 + 2) })", "fn (x) { 4 }"},
		{"quote(unquote(quote(1 + 2)))", "( 1 + 2 )"}, // Unquote should only evaluate the outermost level
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			testQuoteObject(t, obj, tt.expected)
		})
	}
}

func TestEvalBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!false", true},
		{"!true", false},
		{"!!false", false},
		{"!!true", true},
		{"!!!false", true},
		{"!5", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			testBooleanObject(t, obj, tt.expected)
		})
	}
}

func TestEvalArithmeticInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		// ================================
		// Basic operations
		// ================================
		{"5+2", 7},          // Simple addition
		{"10-3", 7},         // Simple subtraction
		{"4*5", 20},         // Simple multiplication
		{"20/4", 5},         // Simple division
		{"5*3+2", 17},       // Multiplication before addition
		{"10-(2+3)", 5},     // Parentheses with addition inside subtraction
		{"(3+2)*2", 10},     // Parentheses with multiplication
		{"(10-5)/(5+5)", 0}, // Parentheses with division

		// ================================
		// Operations with precedence
		// ================================
		{"2+3*5", 17},         // Multiplication before addition
		{"(2+3)*5", 25},       // Parentheses with addition first
		{"(10-3)*(4+2)", 42},  // Parentheses and multiplication
		{"(2+3)*(4+2)-6", 24}, // Mixed operations inside parentheses
		{"2+3*5-6/2", 14},     // Mixed operations with precedence

		// ================================
		// Negative numbers
		// ================================
		{"-5+2", -3},      // Negative number with addition
		{"-5-2", -7},      // Negative number with subtraction
		{"-5*2", -10},     // Negative number with multiplication
		{"-10/2", -5},     // Negative number with division
		{"(5-10)*2", -10}, // Parentheses with negative result
		{"-5+(3*4)", 7},   // Negative number with multiplication in parentheses

		// ================================
		// Large numbers
		// ================================
		{"1000000+2000000", 3000000},       // Large addition
		{"5000000-1000000", 4000000},       // Large subtraction
		{"1000000*1000000", 1000000000000}, // Large multiplication
		{"1000000000/100000", 10000},       // Large division
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			testIntegerObject(t, obj, tt.expected)
		})

	}
}

func TestEvalComparisonInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// ================================
		// Equal (==)
		// ================================
		{"5==5", true},
		{"10==5", false},
		{"0==0", true},
		{"-1==1", false},
		{"true==true", true},
		{"false==false", true},
		{"true==false", false},
		{"\"apple\"==\"apple\"", true},             // Simple equal strings
		{"\"apple\"==\"banana\"", false},           // Different strings
		{"\"hello\"==\"hello\"", true},             // Equal strings with exact same value
		{"\"Hello\"==\"hello\"", false},            // Case-sensitive check (should fail)
		{"\"abc\"==\"abcd\"", false},               // String of different lengths
		{"\"true\"==\"true\"", true},               // Strings with boolean values as strings
		{"\"a\"==\"A\"", false},                    // Case-sensitive comparison
		{"\"123\"==\"123\"", true},                 // Numeric strings (equal)
		{"\"\"==\"\"", true},                       // Empty string comparison
		{"\"hello world\"==\"hello world\"", true}, // Long string comparison

		// ================================
		// Not Equal (!=)
		// ================================
		{"5!=5", false},
		{"10!=5", true},
		{"0!=0", false},
		{"-1!=1", true},
		{"true!=true", false},
		{"false!=false", false},
		{"true!=false", true},
		{"\"apple\"!=\"apple\"", false},  // Identical strings (should not be unequal)
		{"\"apple\"!=\"banana\"", true},  // Different strings
		{"\"hello\"!=\"hello\"", false},  // Identical strings
		{"\"Hello\"!=\"hello\"", true},   // Case-sensitive check (should be unequal)
		{"\"abc\"!=\"abcd\"", true},      // Strings of different lengths
		{"\"true\"!=\"false\"", true},    // Different strings
		{"\"123\"!=\"1234\"", true},      // Numeric strings (unequal)
		{"\"\"!=\"hello\"", true},        // Empty string compared to non-empty
		{"\"goodbye\"!=\"hello\"", true}, // Completely different strings

		// ================================
		// Less Than (<)
		// ================================
		{"5<5", false},
		{"5<10", true},
		{"0<5", true},
		{"-1<0", true},
		{"10<5", false},

		// ================================
		// Greater Than (>)
		// ================================
		{"5>5", false},
		{"10>5", true},
		{"0>5", false},
		{"-1>0", false},
		{"10>5", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			testBooleanObject(t, obj, tt.expected)
		})
	}
}

func TestEvalIfElseConditional(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// ================================
		// Basic if-else conditions
		// ================================
		{`if (2 > 1) { 1 } else { 0 }`, 1},
		{`if (1 > 2) { 1 } else { 0 }`, 0},
		{`if (5 == 5) { 10 } else { 0 }`, 10},
		{`if (0 != 1) { 100 } else { 0 }`, 100},

		// ================================
		// if with nested conditions
		// ================================
		{`if (5 > 2) { if (3 > 1) { 10 } else { 20 } } else { 30 }`, 10},
		{`if (3 < 2) { if (4 > 5) { 10 } else { 20 } } else { 30 }`, 30},

		// Others
		{`if (5 > 2) { 1 }`, 1},
		{`if (0 == 0) { 1 } else { 0 }`, 1},
		{`if (0 == 1) { 1 } else { 0 }`, 0},
		{`if (5 == 5) { if (2 > 1) { 15 } else { 0 } } else { 100 }`, 15},
		{`if (false) { 5 } else { 0 }`, 0},
		{`if (2) { 5 } else { 0 }`, 5},
		{`if (2==1) { 5 }`, nil},

		// ================================
		// Basic if-else conditions with strings
		// ================================
		{`if ("hello" == "hello") { "equal" } else { "not equal" }`, "equal"},      // String equality
		{`if ("apple" == "banana") { "equal" } else { "not equal" }`, "not equal"}, // String inequality
		{`if ("cat" != "dog") { "different" } else { "same" }`, "different"},       // String inequality

		// ================================
		// if with nested conditions using strings
		// ================================
		{`if ("apple" == "apple") { if ("hello" == "hello") { "both equal" } else { "second not equal" } } else { "first not equal" }`, "both equal"}, // Nested string comparison (true in both)
		{`if ("apple" != "banana") { if ("cat" == "dog") { "nested true" } else { "nested false" } } else { "first false" }`, "nested false"},         // Nested false condition with strings

		// ================================
		// Handling strings with boolean or null results
		// ================================
		{`if ("hello" == "world") { "yes" }`, nil},                                                // False condition with no else branch (expected null)
		{`if ("apple" == "apple") { "yes" } else { "no" }`, "yes"},                                // True condition with else
		{`if ("orange" == "apple") { "same" } else { "different" }`, "different"},                 // String inequality with else
		{`if (false) { "not executed" } else { "else branch executed" }`, "else branch executed"}, // False condition

	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			if i, ok := tt.expected.(int); ok {
				testIntegerObject(t, obj, int64(i))
			} else if i, ok := tt.expected.(string); ok {
				testStringObject(t, obj, i)
			} else {
				testNullObject(t, obj)
			}
		})
	}
}

func TestEvalVariableBinding(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// ================================
		// Basic variable binding
		// ================================
		{`let a = 5; a;`, 5},               // simple binding, should return 5
		{`let b = 10; b;`, 10},             // simple binding, should return 10
		{`let b = "elliot"; b;`, "elliot"}, // simple binding, should return 10

		// ================================
		// Variable binding with expressions
		// ================================
		{`let x = 5 + 3; x;`, 8},   // binding with expression, should return 8
		{`let y = 10 * 2; y;`, 20}, // binding with multiplication, should return 20
		{`let z = 5; z * 2;`, 10},  // binding followed by an expression, should return 10

		// ================================
		// Rebinding variables (new assignment)
		// ================================
		{`let a = 5; let a = 10; a;`, 10},                 // rebinding variable a, should return 10
		{`let a = 5; let b = 10; let a = 15; a + b;`, 25}, // rebinding a and adding with b, should return 25

		// ================================
		// Variable binding with conditional
		// ================================
		{`let a = 5; if (a == 5) { let b = 10; b } else { let b = 20; b }`, 10}, // true condition, should return 10
		{`let a = 5; if (a != 5) { let b = 10; b } else { let b = 20; b }`, 20}, // false condition, should return 20

		// ================================
		// Edge cases
		// ================================
		{`let a = 0; let b = a + 5; b;`, 5},    // simple binding with a sum, should return 5
		{`let a = -10; a;`, -10},               // negative value, should return -10
		{`let a = 0; let b = a - 3; b;`, -3},   // subtraction binding, should return -3
		{`let a = 100; let b = a / 2; b;`, 50}, // division, should return 50
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			if i, ok := tt.expected.(int); ok {
				testIntegerObject(t, obj, int64(i))
			} else if i, ok := tt.expected.(string); ok {
				testStringObject(t, obj, i)
			}
		})
	}
}

func TestEvalReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// ================================
		// Basic return statement
		// ================================
		{`let a = "elliot"; return a; 5`, "elliot"}, // should return 10 (first return is evaluated)
		{`return 42; 100`, 42},                      // should return 42, second part is ignored
		{`let x = 20; return x + 5;`, 25},           // return evaluates an expression, should return 25

		// ================================
		// Multiple return statements
		// ================================
		{`return 5; return 10;`, 5},      // should return 5, second return is ignored
		{`return 5 + 5; return 10;`, 10}, // return evaluates first expression, should return 10
		{`return 5 * 2; return 0;`, 10},  // evaluates the first return, should return 10

		// ================================
		// Return after assignments
		// ================================
		{`let a = 3; let b = 5; return a + b;`, 8},   // should return 8 (a + b)
		{`let a = 3; return a * 2;`, 6},              // should return 6 (a * 2)
		{`let x = 5; let y = 10; return x + y;`, 15}, // should return 15 (x + y)

		// ================================
		// Edge cases with return
		// ================================
		{`let a = 1; if (a == 1) { return a; } else { return 0; }`, 1}, // returns 1 since a == 1
		{`let a = 0; if (a == 1) { return a; } else { return 0; }`, 0}, // returns 0 since a != 1

		// ================================
		// Unreachable code after return
		// ================================
		{`return 1; let x = 2; x;`, 1},           // `x` should never be evaluated after the return, should return 1
		{`return 10; let a = 20; return a;`, 10}, // first return is evaluated, second return is ignored, should return 10
		{`return 100; return 200;`, 100},         // second return is ignored, should return 100

	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			if i, ok := tt.expected.(int); ok {
				testIntegerObject(t, obj, int64(i))
			}
		})
	}
}

func TestEvalFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// ================================
		// Basic Function Call
		// ================================
		{`let add = fn(x, y) { x + y }; add(2, 3);`, 5},                    // simple addition function
		{`let multiply = fn(x, y) { x * y }; multiply(4, 5);`, 20},         // simple multiplication function
		{`fn(x) { return x; }("elliot");`, "elliot"},                       // direct function invocation
		{`let subtract = fn(x, y) { return x - y; }; subtract(10, 3);`, 7}, // function with subtraction

		// ================================
		// Edge Cases
		// ================================
		{`let noReturn = fn() {}; noReturn();`, nil},      // function without a return value
		{`let noArg = fn() { return 42; }; noArg();`, 42}, // function with no arguments

		// ================================
		// Function with Boolean Return
		// ================================
		{`let isBig = fn(x) { if (x > 1000) { return true; } else { return false; } }; isBig(1001);`, true}, // checking even number
		{`let isBig = fn(x) { if (x > 1000) { return true; } else { return false; } }; isBig(100);`, false}, // checking odd number

		// ================================
		// Function with Closures
		// ================================
		{`let outer = fn(x) { let inner = fn(y) { return x + y; }; return inner; }; let closure = outer(5); closure(3);`, 8}, // closure capturing `x`
		{`let outer = fn(x) { let inner = fn(y) { return x * y; }; return inner; }; let closure = outer(2); closure(4);`, 8}, // closure with multiplication

		// ================================
		// Closure with Different Scopes
		// ================================
		{`let outer = fn(x) { let inner = fn(y) { return x + y; }; return inner; }; let closure_a = outer(10); let closure_b = outer(20); closure_a(5);`, 15}, // different closure instances, same outer function
		{`let outer = fn(x) { let inner = fn(y) { return x - y; }; return inner; }; let closure = outer(10); closure(4);`, 6},                                 // closure with subtraction

		// ================================
		// Closures and Variable Capturing
		// ================================
		{`let x = 10; let closure = fn() { return x; }; let x = 20; closure();`, 20}, // closure captures the latest value of x (20)
		{`let x = 5; let closure = fn() { return x; }; let x = 15; closure();`, 15},  // closure captures the latest value of x (15)

		// ================================
		// Nested Function Calls
		// ================================
		{`let add = fn(x, y) { return x + y; }; let multiply = fn(x, y) { return x * y; }; multiply(add(2, 3), 4);`, 20}, // nested function calls (add + multiply)
		{`let square = fn(x) { return x * x; }; square(square(3));`, 81},                                                 // function calling itself

		// ================================
		// Functions Returning Functions (Higher-Order Functions)
		// ================================
		{`let multiplyBy = fn(x) { return fn(y) { return x * y; }; }; let multiplyByTwo = multiplyBy(2); multiplyByTwo(3);`, 6}, // higher-order function returning another function
		{`let applyFn = fn(f, x) { return f(x); }; applyFn(fn(x) { return x + 1; }, 5);`, 6},                                    // higher-order function that accepts another function
	}

	// Running each test
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			switch expected := tt.expected.(type) {
			case int:
				testIntegerObject(t, obj, int64(expected))
			case bool:
				testBooleanObject(t, obj, expected)
			}
		})
	}
}

func TestEvalArrays(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// ================================
		// Basic Array Creation
		// ================================
		{`[1, 2, 3]`, []interface{}{int64(1), int64(2), int64(3)}}, // Array with integers
		{`["a", "b", "c"]`, []interface{}{"a", "b", "c"}},          // Array with strings

		// ================================
		// Empty Array
		// ================================
		{`[]`, []interface{}{}}, // Empty array

		// ================================
		// Array Access (Indexing)
		// ================================
		{`[1, 2, 3][0]`, int64(1)}, // Access first element
		{`[1, 2, 3][1]`, int64(2)}, // Access second element
		{`[1, 2, 3][2]`, int64(3)}, // Access third element
		{`[1, 2, 3][3]`, errors.New("index out of bounds for arr length 3")}, // Access out of bounds (null)

		{`[[1, 2], [3, 4]][0][1]`, int64(2)}, // Access second element of the first nested array

		// ================================
		// Array with Mixed Types
		// ================================
		{`[1, "a", true]`, []interface{}{int64(1), "a", true}}, // Mixed types in the array
		{`[1, "a", true][1]`, "a"},                             // Access string element
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case int64:
				testIntegerObject(t, obj, expected)
			case bool:
				testBooleanObject(t, obj, expected)
			case string:
				testStringObject(t, obj, expected)
			case error:
				testErrorObject(t, obj, expected.Error())
			case []interface{}:
				testArrayObject(t, obj, expected)
			}
		})
	}
}

func TestEvalHashAccess(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// ================================
		// Basic Hash Access with String Keys
		// ================================
		{`{"name": "Alice", "age": 30}["name"]`, "Alice"},     // Access string key "name"
		{`{"name": "Alice", "age": 30}["age"]`, int64(30)},    // Access string key "age"
		{`{"score": 95, "passed": true}["score"]`, int64(95)}, // Access string key "score"
		{`{"score": 95, "passed": true}["passed"]`, true},     // Access string key "passed"

		// ================================
		// Access with Integer Keys
		// ================================
		{`{1: "one", 2: "two", 3: "three"}[1]`, "one"},   // Access integer key 1
		{`{1: "one", 2: "two", 3: "three"}[2]`, "two"},   // Access integer key 2
		{`{1: "one", 2: "two", 3: "three"}[3]`, "three"}, // Access integer key 3

		// ================================
		// Access with Boolean Keys
		// ================================
		{`{true: "yes", false: "no"}[true]`, "yes"}, // Access boolean key true
		{`{true: "yes", false: "no"}[false]`, "no"}, // Access boolean key false

		// ================================
		// Missing Keys (Should return null)
		// ================================
		{`{"name": "Alice"}["age"]`, nil}, // Accessing non-existent key "age"
		{`{"score": 95}["passed"]`, nil},  // Accessing non-existent key "passed"

		// ================================
		// Invalid Key Types (e.g., arrays as keys)
		// ================================
		{`{[1,2]: "array"}[2]`, errors.New("key type ARRAY is not hashable")},      // Invalid key type: array as a key
		{`{"name": "Alice"}[[1,2]]`, errors.New("key type ARRAY is not hashable")}, // Invalid key type: array as a key
		{`{true: "yes"}[1]`, nil},

		// ================================
		// Nested Hash Access
		// ================================
		{`{"person": {"name": "Alice", "age": 30}}["person"]["name"]`, "Alice"},    // Access nested hash key "name"
		{`{"person": {"name": "Alice", "age": 30}}["person"]["age"]`, int64(30)},   // Access nested hash key "age"
		{`{"team": {"leader": "Bob", "members": 5}}["team"]["leader"]`, "Bob"},     // Access nested hash key "leader"
		{`{"team": {"leader": "Bob", "members": 5}}["team"]["members"]`, int64(5)}, // Access nested hash key "members"

		// ================================
		// Nested Hash Access with Missing Keys
		// ================================
		{`{"team": {"leader": "Bob"}}["team"]["members"]`, nil},  // Accessing non-existent nested key "members"
		{`{"team": {"leader": "Bob"}}["team"]["location"]`, nil}, // Accessing non-existent nested key "location"
		{`{"person": {"name": "Alice"}}["person"]["age"]`, nil},  // Accessing non-existent nested key "age"

		// ================================
		// Nested Hash with Mixed Keys
		// ================================
		{`{"user": {1: "Bob", true: "Alice"}}["user"][1]`, "Bob"},      // Mixed key types in nested hash
		{`{"user": {1: "Bob", true: "Alice"}}["user"][true]`, "Alice"}, // Mixed key types in nested hash

		// ================================
		// Out-of-Bounds Access (Accessing Hash Key of Another Hash/Array)
		// ================================
		{`{"nested": {"inner": {"key": "value"}}}["nested"]["inner"]["key"]`, "value"}, // Access nested keys
		{`{"nested": {"inner": {"key": "value"}}}["nested"]["inner"][0]`, nil},         // Trying to access array index on a hash

		// ================================
		// Edge Cases with Nested Hashes and Missing Keys
		// ================================
		{`{"outer": {"inner": {"key": "value"}}}["outer"]["inner"]["missing"]`, nil},                                                       // Accessing non-existent key in nested hash
		{`{"outer": {"inner": {"key": "value"}}}["outer"]["missing"]["key"]`, errors.New("index expression not supported for type: NULL")}, // Accessing non-existent outer key

		// ================================
		// Multiple Hashes with the Same Key
		// ================================
		{`{"user": {"name": "Bob"}}["user"]["name"]`, "Bob"},     // Access name from "user"
		{`{"user": {"name": "Alice"}}["user"]["name"]`, "Alice"}, // Access name from "user"
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case int64:
				testIntegerObject(t, obj, expected)
			case string:
				testStringObject(t, obj, expected)
			case bool:
				testBooleanObject(t, obj, expected)
			case nil:
				testNullObject(t, obj)
			case error:
				testErrorObject(t, obj, expected.Error())
			default:
				t.Errorf("unexpected type %T for expected value: %v", expected, expected)
			}
		})
	}
}

func TestEvalStringConcatenation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// ================================
		// Basic Concatenation
		// ================================
		{`"Hello" + " " + "World"`, "Hello World"}, // Basic string concatenation
		{`"foo" + "bar"`, "foobar"},                // Simple concatenation of two strings
		{`"apple" + "pie"`, "applepie"},            // Simple concatenation of two words
		{`"a" + "b" + "c"`, "abc"},                 // Concatenation of multiple single characters

		// ================================
		// Concatenation with variables (string variables)
		// ================================
		{`let greeting = "Hello"; let name = "Alice"; greeting + " " + name`, "Hello Alice"},   // Using variables in concatenation
		{`let part_a = "Good"; let part_b = "Morning"; part_a + " " + part_b`, "Good Morning"}, // Concatenating variables holding strings

		// ================================
		// Concatenation with empty strings
		// ================================
		{`"" + "test"`, "test"}, // Empty string and a non-empty string
		{`"test" + ""`, "test"}, // Non-empty string and an empty string
		{`"" + ""`, ""},         // Concatenation of two empty strings

		// ================================
		// Concatenation with large strings
		// ================================
		{`"a" + "b" + "c" + "d" + "e" + "f" + "g" + "h" + "i" + "j"`, "abcdefghij"}, // Concatenating many characters
		{`"Lorem ipsum dolor sit amet, consectetur adipiscing elit. " + "Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."`,
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."}, // Long string concatenation
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			testStringObject(t, obj, tt.expected)
		})
	}
}

func TestEvalWithMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{

		// 1. Basic arithmetic macro evaluation
		{`
			let minus = macro(a, b) { quote(unquote(a) - unquote(b)) };
			minus(10, 2)`,
			8,
		},
		{`
			let add = macro(a, b) { quote(unquote(a) + unquote(b)) };
			add(3, 7)`,
			10,
		},
		{`
			let mul = macro(a, b) { quote(unquote(a) * unquote(b)) };
			mul(4, 5)`,
			20,
		},
		{`
			let div = macro(a, b) { quote(unquote(a) / unquote(b)) };
			div(10, 2)`,
			5,
		},

		// 2. Macro with variable substitution
		{`
			let x = 10;
			let y = 5;
			let sub = macro(a, b) { quote(unquote(a) - unquote(b)) };
			sub(x, y)`,
			5,
		},

		// 4. Macros inside conditionals
		{`
			let conditional = macro(a, b) { quote(if (unquote(a) > 0) {unquote(a)} else {unquote(b)}) };
			conditional(5, 10)`,
			5,
		},
		{`
			let conditional = macro(a, b) { quote(if (unquote(a) > 0) {unquote(a)} else {unquote(b)}) };
			conditional(-5, 10)`,
			10,
		},

		// 5. Macros with functions
		{`
			let makeAdder = macro(x) { quote(fn(y) { unquote(x) + y }) };
			let addFive = makeAdder(5);
			addFive(3)`,
			8,
		},

		{`
		let unless = macro(condition, consequence, alternative){ 
			quote(if (!(unquote(condition))) { unquote(consequence); }
 			else { unquote(alternative); }); };
		unless(10 > 5, 2, 3);`,
			3},

		// 6. Edge Cases

		// 6.1 Unused parameter
		{`
			let ignoreArg = macro(x) { quote(100) };
			ignoreArg(50)`,
			100,
		},

		// 6.3 Macro with boolean expressions
		{`
			let boolCheck = macro(a) { quote(unquote(a) == true) };
			boolCheck(true)`,
			true,
		},
		{`
			let boolCheck = macro(a) { quote(unquote(a) == false) };
			boolCheck(true)`,
			false,
		},

		// 6.4 Error case: Division by zero inside macro
		{`
			let divZero = macro(a) { quote(unquote(a) / 0) };
			divZero(10)`,
			fmt.Errorf("Division by zero"),
		},

		// 6.5 Error case: Invalid macro argument
		{`
			let invalid = macro(a) { quote(unquote(a) + 1) };
			invalid("hello")`,
			fmt.Errorf("Incompatible types: STRING and INTEGER"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case int64:
				testIntegerObject(t, obj, expected)
			case int:
				testIntegerObject(t, obj, int64(expected))
			case bool:
				testBooleanObject(t, obj, expected)
			case string:
				testStringObject(t, obj, expected)
			case error:
				testErrorObject(t, obj, expected.Error())
			case []interface{}:
				testArrayObject(t, obj, expected)
			}
		})
	}
}

func TestEvalBuiltInFuncLen(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		// ================================
		// Basic String Length
		// ================================
		{`len("hello")`, 5}, // Basic string length
		{`len("world")`, 5}, // Another basic string length
		{`len("a")`, 1},     // Length of a single character
		{`len("")`, 0},      // Length of an empty string

		// ================================
		// Arrays (or lists)
		// ================================
		{`len([1, 2, 3, 4])`, 4},  // Length of an array (or list) of integers
		{`len([10, 20, 30])`, 3},  // Another array length
		{`len([true, false])`, 2}, // Array with boolean values
		{`len([])`, 0},            // Empty array

		// ================================
		// Mixed Types (if supported)
		// ================================
		{`len([1, "a", true])`, 3},        // Mixed array with number, string, and boolean
		{`len([true, "string", 100])`, 3}, // Mixed array with booleans, string, and number

		// ================================
		// Nested Arrays
		// ================================
		{`len([ [1, 2], [3, 4] ])`, 2},            // Nested arrays (should count top level)
		{`len([ [1, 2, 3], [4, 5] ])`, 2},         // Nested arrays with varying lengths
		{`len([["a", "b"], ["c", "d", "e"]])`, 2}, // Nested arrays with mixed sizes

		// ================================
		// Concatenation with Length
		// ================================
		{`len("hello" + " world")`, 11},  // Length of concatenated string
		{`len("good" + "bye" + "!")`, 8}, // Length of concatenated string with multiple parts
		{`len("a" + "b" + "c")`, 3},      // Length of concatenation of multiple characters

		// ================================
		// Other Edge Cases
		// ================================
		{`len("a" + "")`, 1},         // Concatenation with an empty string
		{`len("") + len("test")`, 4}, // Adding lengths of two strings
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			testIntegerObject(t, obj, int64(tt.expected))
		})
	}
}

func TestEvalBuiltInFuncFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// ================================
		// Valid Arrays (Non-Empty)
		// ================================
		{`first([1, 2, 3])`, int64(1)},           // First element of the array is an integer
		{`first(["hello", "world"])`, "hello"},   // First element of the array is a string
		{`first([true, false])`, true},           // First element of the array is a boolean
		{`first([1, "hello", false])`, int64(1)}, // Mixed array, first element is integer

		// ================================
		// Edge Cases (Empty Arrays)
		// ================================
		{`first([])`, errors.New("empty array")}, // Empty array, should return null

		// ================================
		// Single Element Arrays
		// ================================
		{`first([42])`, int64(42)},  // Single integer in array
		{`first(["only"])`, "only"}, // Single string in array
		{`first([false])`, false},   // Single boolean in array

		// ================================
		// Nested Arrays
		// ================================
		{`first([[1, 2], [3, 4]])`, []interface{}{int64(1), int64(2)}}, // First element is an array
		{`first([["a", "b"], ["c", "d"]])`, []interface{}{"a", "b"}},   // First element is a string array

		// ================================
		// Mixed Arrays with Nested Elements
		// ================================
		{`first([1, [2, 3], "hello"])`, int64(1)}, // First element is a number
		{`first([true, [false, true]])`, true},    // First element is a boolean

		// ================================
		// Invalid Cases (non-array arguments)
		// ================================
		{`first("hello")`, errors.New("first(): type STRING not supported")},
		{`first(42)`, errors.New("first(): type INTEGER not supported")},
		{`first(true)`, errors.New("first(): type BOOLEAN not supported")},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case int64:
				testIntegerObject(t, obj, expected)
			case string:
				testStringObject(t, obj, expected)
			case bool:
				testBooleanObject(t, obj, expected)
			case []interface{}:
				// Check for nested arrays
				testArrayObject(t, obj, expected)
			case error:
				testErrorObject(t, obj, expected.Error())
			default:
				t.Errorf("unexpected type for expected value: %T", expected)
			}
		})
	}
}

func TestEvalBuiltInFuncLast(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// ================================
		// Valid Arrays (Non-Empty)
		// ================================
		{`last([1, 2, 3])`, int64(3)},         // Last element of the array is an integer
		{`last(["hello", "world"])`, "world"}, // Last element of the array is a string
		{`last([true, false])`, false},        // Last element of the array is a boolean
		{`last([1, "hello", false])`, false},  // Mixed array, last element is a boolean

		// ================================
		// Edge Cases (Empty Arrays)
		// ================================
		{`last([])`, errors.New("empty array")}, // Empty array, should return an error

		// ================================
		// Single Element Arrays
		// ================================
		{`last([42])`, int64(42)},  // Single integer in array
		{`last(["only"])`, "only"}, // Single string in array
		{`last([false])`, false},   // Single boolean in array

		// ================================
		// Nested Arrays
		// ================================
		{`last([[1, 2], [3, 4]])`, []interface{}{int64(3), int64(4)}}, // Last element is an array
		{`last([["a", "b"], ["c", "d"]])`, []interface{}{"c", "d"}},   // Last element is a string array

		// ================================
		// Mixed Arrays with Nested Elements
		// ================================
		{`last([1, [2, 3], "hello"])`, "hello"},                     // Last element is a string
		{`last([true, [false, true]])`, []interface{}{false, true}}, // Last element is a nested array

		// ================================
		// Invalid Cases (non-array arguments)
		// ================================
		{`last("hello")`, errors.New("last(): type STRING not supported")}, // Invalid, not an array
		{`last(42)`, errors.New("last(): type INTEGER not supported")},     // Invalid, not an array
		{`last(true)`, errors.New("last(): type BOOLEAN not supported")},   // Invalid, not an array
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			switch expected := tt.expected.(type) {
			case int64:
				testIntegerObject(t, obj, expected)
			case string:
				testStringObject(t, obj, expected)
			case bool:
				testBooleanObject(t, obj, expected)
			case []interface{}:
				// Check for nested arrays
				testArrayObject(t, obj, expected)
			case error:
				testErrorObject(t, obj, expected.Error())
			default:
				t.Errorf("unexpected type for expected value: %T", expected)
			}
		})
	}
}

func TestEvalBuiltInFuncRest(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// ================================
		// Valid Arrays (Non-Empty)
		// ================================
		{`rest([1, 2, 3])`, []interface{}{int64(2), int64(3)}},             // Rest of the array excluding the first element
		{`rest(["hello", "world", "foo"])`, []interface{}{"world", "foo"}}, // Rest of the array with strings
		{`rest([true, false, true])`, []interface{}{false, true}},          // Rest of the array with booleans
		{`rest([1, "hello", false])`, []interface{}{"hello", false}},       // Mixed array, excluding the first element

		// ================================
		// Edge Cases (Empty Arrays)
		// ================================
		{`rest([])`, errors.New("empty array")}, // Empty array, should return an error or "null"

		// ================================
		// Single Element Arrays
		// ================================
		{`rest([42])`, []interface{}{}},     // Single element array, result is an empty array
		{`rest(["only"])`, []interface{}{}}, // Single string element, result is an empty array
		{`rest([false])`, []interface{}{}},  // Single boolean element, result is an empty array

		// ================================
		// Nested Arrays
		// ================================
		{`rest([[1, 2], [3, 4], [5, 6]])`, []interface{}{[]interface{}{int64(3), int64(4)}, []interface{}{int64(5), int64(6)}}}, // Nested arrays, excluding first
		{`rest([["a", "b"], ["c", "d"], ["e", "f"]])`, []interface{}{[]interface{}{"c", "d"}, []interface{}{"e", "f"}}},         // Nested arrays with strings

		// ================================
		// Mixed Arrays with Nested Elements
		// ================================
		{`rest([1, [2, 3], "hello", true])`, []interface{}{[]interface{}{int64(2), int64(3)}, "hello", true}}, // Mixed array with a nested array and others
		{`rest([true, [false, true], 42])`, []interface{}{[]interface{}{false, true}, int64(42)}},             // Boolean and nested array mixed

		// ================================
		// Invalid Cases (non-array arguments)
		// ================================
		{`rest("hello")`, errors.New("rest(): type STRING not supported")}, // Invalid, not an array
		{`rest(42)`, errors.New("rest(): type INTEGER not supported")},     // Invalid, not an array
		{`rest(true)`, errors.New("rest(): type BOOLEAN not supported")},   // Invalid, not an array
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case []interface{}:
				testArrayObject(t, obj, expected)
			case error:
				testErrorObject(t, obj, expected.Error())
			default:
				t.Errorf("unexpected type for expected value: %T", expected)
			}
		})
	}
}

func TestEvalBuiltInFuncPush(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// ================================
		// Valid Arrays (Non-Empty)
		// ================================
		{`push([1, 2], 3)`, []interface{}{int64(1), int64(2), int64(3)}},                        // Add an integer to the end
		{`push(["hello", "world"], "test")`, []interface{}{"hello", "world", "test"}},           // Add a string
		{`push([true, false], true)`, []interface{}{true, false, true}},                         // Add a boolean
		{`push([1, "hello", false], 100)`, []interface{}{int64(1), "hello", false, int64(100)}}, // Mixed types

		// ================================
		// Edge Cases (Empty Arrays)
		// ================================
		{`push([], 42)`, []interface{}{int64(42)}}, // Adding to an empty array

		// ================================
		// Single Element Arrays
		// ================================
		{`push([42], 100)`, []interface{}{int64(42), int64(100)}}, // Add to a single-element array
		{`push(["only"], "more")`, []interface{}{"only", "more"}}, // Add to a single-element string array
		{`push([false], true)`, []interface{}{false, true}},       // Add to a single-element boolean array

		// ================================
		// Nested Arrays
		// ================================
		{`push([[1, 2], [3, 4]], [5, 6])`, []interface{}{[]interface{}{int64(1), int64(2)}, []interface{}{int64(3), int64(4)}, []interface{}{int64(5), int64(6)}}}, // Add a nested array
		{`push([["a", "b"], ["c", "d"]], ["e", "f"])`, []interface{}{[]interface{}{"a", "b"}, []interface{}{"c", "d"}, []interface{}{"e", "f"}}},                   // Add another nested array

		// ================================
		// Mixed Arrays with Nested Elements
		// ================================
		{`push([1, [2, 3], "hello"], [4, 5])`, []interface{}{int64(1), []interface{}{int64(2), int64(3)}, "hello", []interface{}{int64(4), int64(5)}}}, // Add nested array to mixed array
		{`push([true, [false, true]], ["more"])`, []interface{}{true, []interface{}{false, true}, []interface{}{"more"}}},                              // Add a nested array of strings

		// ================================
		// Invalid Cases (non-array arguments)
		// ================================
		{`push("hello", 42)`, errors.New("push(): type STRING not supported")},  // Invalid, not an array
		{`push(42, 100)`, errors.New("push(): type INTEGER not supported")},     // Invalid, not an array
		{`push(true, false)`, errors.New("push(): type BOOLEAN not supported")}, // Invalid, not an array
	}

	// Running each test
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case []interface{}:

				testArrayObject(t, obj, expected)
			case error:

				testErrorObject(t, obj, expected.Error())
			default:
				t.Errorf("unexpected type for expected value: %T", expected)
			}
		})
	}
}

// TODO: Put these scenarios with their respective normal cases in other test functions.
func TestEvalErrorHandling(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		// ================================
		// Undefined Variables
		// ================================
		{"return x;", `Undefined variable "x"`},                           // using undefined variable 'x'
		{"let a = 5; return b;", `Undefined variable "b"`},                // 'b' is not defined
		{"let a = 5; let b = a + c; return a;", `Undefined variable "c"`}, // 'c' is undefined

		// ================================
		// Division by Zero
		// ================================
		{"let x = 5; let y = 0; return x / y;", "Division by zero"}, // division by zero
		{"return 10 / 0;", "Division by zero"},                      // division by zero

		// ================================
		// Invalid Operations
		// ================================
		//{"return 'string' + 5;", "Runtime Error: Invalid operation between 'string' and 'number'."}, // invalid type operation
		{"let a = true; let b = 10; return a + b;", "Incompatible types: BOOLEAN and INTEGER"},   // invalid type operation
		{`let a = 10; let b = "hello"; return a - b;`, "Incompatible types: INTEGER and STRING"}, // invalid type operation
		{"-true", "Invalid type BOOLEAN with operator '-'"},
		{"!(true+2)", "Incompatible types: BOOLEAN and INTEGER"},

		// ================================
		// Invalid Condition Expressions
		// ================================
		{"if (x > 5) { return 1; } else { return 0; }", `Undefined variable "x"`}, // undefined variable 'x'
		{"if (10 / 0) { return 1; } else { return 0; }", "Division by zero"},      // division by zero in condition

		// ================================
		// Type Mismatch in Conditionals
		// ================================
		{"if (true == 2) { return 1; } else { return 0; }", "Incompatible types: BOOLEAN and INTEGER"},
		{"if (false > 5) { return 1; } else { return 0; }", "Incompatible types: BOOLEAN and INTEGER"},

		// ================================
		// Functions
		// ================================
		{`let func = fn(x) { return x; }; func(1, 2);`, "expected 1 parameters, got 2 args"},           // too many arguments
		{`let func = fn(x,y) { return x+y; }; func(10);`, "expected 2 parameters, got 1 args"},         // too few arguments
		{`let func = fn(x) { return x + 5; }; func(true);`, "Incompatible types: BOOLEAN and INTEGER"}, // invalid argument type
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			obj := testEval(tt.input)
			if obj != nil {
				testErrorObject(t, obj, tt.expectedError)
			}
		})
	}
}

// TODO: error handling for parser and expansion errors.
func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	prg := p.ParseProgram()
	if len(p.Errors) != 0 {
		fmt.Printf("parser errors for input: %q\n", input)
		for _, err := range p.Errors {
			fmt.Println("\t" + err)
		}
		return nil
	}
	macroEnv := object.NewEnvironment(nil)
	expandedAST, err := ExpandMacro(prg, macroEnv)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	env := object.NewEnvironment(nil)
	obj := Eval(expandedAST, env)
	return obj
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	if i, ok := obj.(*object.Integer); ok {
		if i.Value != expected {
			t.Errorf("expected %d, got %d", expected, i.Value)
		}
	} else {
		t.Errorf("expected *object.Integer, got %s", obj.Inspect())

	}
}

func testQuoteObject(t *testing.T, obj object.Object, expected string) {
	if i, ok := obj.(*object.Quote); ok {
		if i.Node.String() != expected {
			t.Errorf("expected %s, got %s", expected, i.Node.String())
		}
	} else {
		t.Errorf("expected *object.Quote, got %s", obj.Inspect())

	}
}

func testStringObject(t *testing.T, obj object.Object, expected string) {
	if i, ok := obj.(*object.String); ok {
		if i.Value != expected {
			t.Errorf("expected %s, got %s", expected, i.Value)
		}
	} else {
		t.Errorf("expected *object.String, got %s", obj)

	}
}

func testNullObject(t *testing.T, obj object.Object) {
	if !object.IsNull(obj) {
		t.Errorf("expected null, got %v", obj)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) {
	if b, ok := obj.(*object.Boolean); ok {
		if b.Value != expected {
			t.Errorf("expected %v, got %v", expected, b.Value)
		}
	} else {
		t.Errorf("expected *object.Boolean, got %T", obj)

	}
}

func testErrorObject(t *testing.T, obj object.Object, expectedError string) {
	if err, ok := obj.(*object.Error); ok {
		if err.Message != expectedError {
			t.Errorf("expected error msg %q, got %q", expectedError, err.Message)
		}
	} else {
		t.Errorf("expected *object.Error, got %T", obj)

	}
}

// Helper function to test arrays
func testArrayObject(t *testing.T, obj object.Object, expected []interface{}) {
	arr, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("expected an array, got %s", obj.Type())
	}

	if len(arr.Elements) != len(expected) {
		t.Errorf("expected array of length %d, got %d", len(expected), len(arr.Elements))
	}

	// Iterate through the array elements and compare them to the expected values
	for i, expectedElem := range expected {
		switch expectedElem := expectedElem.(type) {
		case int64:
			testIntegerObject(t, arr.Elements[i], expectedElem)
		case string:
			testStringObject(t, arr.Elements[i], expectedElem)
		case bool:
			testBooleanObject(t, arr.Elements[i], expectedElem)
		case []interface{}:
			testArrayObject(t, arr.Elements[i], expectedElem) // Recursively check nested arrays
		default:
			t.Errorf("unsupported element type in array at index %d: %T", i, expectedElem)
		}
	}
}

//func testHashObject(t *testing.T, obj object.Object, expected map[interface{}]interface{}) {
//	// Ensure the object is of type Hash
//	hash, ok := obj.(*object.Hash)
//	if !ok {
//		t.Fatalf("expected *Hash, got %T", obj)
//	}
//
//	// Check the number of keys in the returned hash
//	if len(hash.Pairs) != len(expected) {
//		t.Fatalf("expected %d pairs, got %d", len(expected), len(hash.Pairs))
//	}
//
//	// Iterate over the expected key-value pairs and compare
//	for key, expectedValue := range expected {
//		expectedKey := toHashableKey(key) // Convert the expected key into the correct type for comparison
//		hashValue, ok := hash.Pairs[expectedKey]
//		if !ok {
//			t.Errorf("expected key %v not found in hash", key)
//			continue
//		}
//
//		// Compare the values (use deep equality depending on type)
//		switch expectedValue := expectedValue.(type) {
//		case int64:
//			testIntegerObject(t, hashValue, expectedValue)
//		case string:
//			testStringObject(t, hashValue, expectedValue)
//		case bool:
//			testBooleanObject(t, hashValue, expectedValue)
//		case nil:
//			testNullObject(t, hashValue)
//		default:
//			t.Errorf("unexpected value type: %T", expectedValue)
//		}
//	}
//}
