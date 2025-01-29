package object

import (
	"bytes"
	"fmt"
	"github.com/jatin-malik/make-thy-interpreter/ast"
	"strings"
)

type ObjectType string

const (
	IntegerObject     ObjectType = "INTEGER"
	BooleanObject     ObjectType = "BOOLEAN"
	NullObject        ObjectType = "NULL"
	ReturnValueObject ObjectType = "RETURN_VALUE"
	ErrorValueObject  ObjectType = "ERROR"
	FunctionObject    ObjectType = "FUNCTION"
	StringObject      ObjectType = "STRING"
	ArrayObject       ObjectType = "ARRAY"
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

type Integer struct {
	Value int64
}

func (integer *Integer) Type() ObjectType {
	return IntegerObject
}

func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}

type String struct {
	Value string
}

func (string *String) Type() ObjectType {
	return StringObject
}

func (string *String) Inspect() string {
	return fmt.Sprintf("\"%s\"", string.Value)
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
