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
	obj := Eval(prg)
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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) {
	if b, ok := obj.(*object.Boolean); ok {
		if b.Value != expected {
			t.Errorf("expected %v, got %v", expected, b.Value)
		}
	} else {
		t.Errorf("expected *object.Boolean, got %T", obj)

	}
}
