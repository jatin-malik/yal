package evaluator

import (
	"fmt"
	"github.com/jatin-malik/make-thy-interpreter/lexer"
	"github.com/jatin-malik/make-thy-interpreter/object"
	"github.com/jatin-malik/make-thy-interpreter/parser"
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
		obj := testEval(tt.input)
		testIntegerObject(t, obj, tt.expected)
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
		obj := testEval(tt.input)
		testBooleanObject(t, obj, tt.expected)
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
		obj := testEval(tt.input)
		testIntegerObject(t, obj, tt.expected)
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
		obj := testEval(tt.input)
		testBooleanObject(t, obj, tt.expected)
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
		{`if (2==1) { 5 }`, "null"},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		if i, ok := tt.expected.(int); ok {
			testIntegerObject(t, obj, int64(i))
		} else {
			testNullObject(t, obj)
		}
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
		{`let a = 5; a;`, 5},   // simple binding, should return 5
		{`let b = 10; b;`, 10}, // simple binding, should return 10

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

		// ================================
		// Undefined variable (edge case)
		// ================================
		{`a;`, "null"}, // undefined variable, should return null
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		if i, ok := tt.expected.(int); ok {
			testIntegerObject(t, obj, int64(i))
		} else {
			testNullObject(t, obj)
		}

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
		{`let a = 10; return a; 5`, 10},   // should return 10 (first return is evaluated)
		{`return 42; 100`, 42},            // should return 42, second part is ignored
		{`let x = 20; return x + 5;`, 25}, // return evaluates an expression, should return 25

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
		obj := testEval(tt.input)
		if i, ok := tt.expected.(int); ok {
			testIntegerObject(t, obj, int64(i))
		} else {
			testNullObject(t, obj)
		}

	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)

	if len(p.Errors) != 0 {
		for _, err := range p.Errors {
			fmt.Println(err)
		}
		return nil
	}
	prg := p.ParseProgram()
	env := object.NewEnvironment()
	obj := Eval(prg, env)
	return obj
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	if i, ok := obj.(*object.Integer); ok {
		if i.Value != expected {
			t.Errorf("expected %d, got %d", expected, i.Value)
		}
	} else {
		t.Errorf("expected *object.Integer, got %T", obj)

	}
}

func testNullObject(t *testing.T, obj object.Object) {
	if !isNull(obj) {
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
