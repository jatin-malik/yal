package object

import (
	"bytes"
	"fmt"
	"github.com/jatin-malik/yal/ast"
	"github.com/jatin-malik/yal/bytecode"
	"strings"
)

type ObjectType string

const (
	IntegerObject          ObjectType = "INTEGER"
	BooleanObject          ObjectType = "BOOLEAN"
	NullObject             ObjectType = "NULL"
	ReturnValueObject      ObjectType = "RETURN_VALUE"
	ErrorValueObject       ObjectType = "ERROR"
	FunctionObject         ObjectType = "FUNCTION"
	CompiledFunctionObject ObjectType = "COMPILED_FUNCTION"
	MacroObject            ObjectType = "MACRO"
	StringObject           ObjectType = "STRING"
	ArrayObject            ObjectType = "ARRAY"
	HashObject             ObjectType = "HASH"
	QuoteObject            ObjectType = "QUOTE"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{true}
	FALSE = &Boolean{false}
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type HashKey struct {
	Type  ObjectType
	Value string
}

type Hashable interface {
	HashKey() HashKey
}

type Integer struct {
	Value int64
}

func (integer *Integer) Type() ObjectType {
	return IntegerObject
}

func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}

func (integer *Integer) HashKey() HashKey {
	return HashKey{Type: integer.Type(), Value: integer.Inspect()}
}

type String struct {
	Value string
}

func (string *String) Type() ObjectType {
	return StringObject
}

func (string *String) Inspect() string {
	return string.Value
}

func (string *String) HashKey() HashKey {
	return HashKey{Type: string.Type(), Value: string.Inspect()}
}

type Boolean struct {
	Value bool
}

func (boolean *Boolean) Type() ObjectType {
	return BooleanObject
}

func (boolean *Boolean) Inspect() string {
	return fmt.Sprintf("%t", boolean.Value)
}

func (boolean *Boolean) HashKey() HashKey {
	return HashKey{Type: boolean.Type(), Value: boolean.Inspect()}
}

type Quote struct {
	Node ast.Node
}

func (quote *Quote) Type() ObjectType {
	return QuoteObject
}

func (quote *Quote) Inspect() string {
	var out bytes.Buffer
	out.WriteString("quote")
	out.WriteString("(")
	out.WriteString(quote.Node.String())
	out.WriteString(")")
	return out.String()
}

type Array struct {
	Elements []Object
}

func (array *Array) Type() ObjectType {
	return ArrayObject
}

func (array *Array) Inspect() string {
	var out bytes.Buffer
	var elements []string
	for _, e := range array.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type Hash struct {
	Pairs map[HashKey]Object
}

func (hash *Hash) Type() ObjectType {
	return HashObject
}

func (hash *Hash) Inspect() string {
	var out bytes.Buffer
	out.WriteString("{")
	count := 0
	for k, v := range hash.Pairs {
		out.WriteString(k.Value)
		out.WriteString(":")
		out.WriteString(v.Inspect())
		if count != len(hash.Pairs)-1 {
			out.WriteString(", ")
		}
		count++
	}
	out.WriteString("}")
	return out.String()
}

type CompiledFunction struct {
	Instructions bytecode.Instructions
}

func (compiledFunction *CompiledFunction) Type() ObjectType {
	return CompiledFunctionObject
}

func (compiledFunction *CompiledFunction) Inspect() string {
	return fmt.Sprintf("COMPILED_FUNCTION(%p)", compiledFunction)
}

type Function struct {
	Env        *Environment
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
}

func (function *Function) Type() ObjectType {
	return FunctionObject
}

func (function *Function) Inspect() string {
	var out bytes.Buffer
	var params []string
	for _, p := range function.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(function.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type Macro struct {
	Env        *Environment
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
}

func (macro *Macro) Type() ObjectType {
	return MacroObject
}

func (macro *Macro) Inspect() string {
	var out bytes.Buffer
	var params []string
	for _, p := range macro.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("macro(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(macro.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// Null is a billion-dollar mistake but sure, why not!
type Null struct {
}

func (null *Null) Type() ObjectType {
	return NullObject
}

func (null *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (returnValue *ReturnValue) Type() ObjectType {
	return ReturnValueObject
}

func (returnValue *ReturnValue) Inspect() string {
	return returnValue.Value.Inspect()
}

type Error struct {
	Message string
}

func (error *Error) Type() ObjectType {
	return ErrorValueObject
}

func (error *Error) Inspect() string {
	return fmt.Sprintf("ERROR: %s", error.Message)
}

// Environment holds the current evaluation context/bindings. Also known as scope.
type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment(outer *Environment) *Environment {
	return &Environment{store: make(map[string]Object), outer: outer}
}

func (env *Environment) Get(name string) Object {
	if obj, ok := env.store[name]; ok {
		return obj
	} else {
		if env.outer != nil {
			return env.outer.Get(name)
		} else if fn, ok := builtinFunctions[name]; ok {
			return fn
		} else {
			msg := fmt.Sprintf("Undefined variable %q", name)
			return NewError(msg)
		}
	}
}

func (env *Environment) Set(name string, value Object) {
	env.store[name] = value
}

func NewError(message string) *Error {
	return &Error{Message: message}
}
