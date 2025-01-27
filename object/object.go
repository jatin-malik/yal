package object

import "fmt"

type ObjectType string

const (
	IntegerObject     ObjectType = "INTEGER"
	BooleanObject     ObjectType = "BOOLEAN"
	NullObject        ObjectType = "NULL"
	ReturnValueObject ObjectType = "RETURN_VALUE"
	ErrorValueObject  ObjectType = "ERROR"
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

type Boolean struct {
	Value bool
}

func (boolean *Boolean) Type() ObjectType {
	return BooleanObject
}

func (boolean *Boolean) Inspect() string {
	return fmt.Sprintf("%t", boolean.Value)
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
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

func (env *Environment) Get(name string) Object {
	if obj, ok := env.store[name]; ok {
		return obj
	}
	msg := fmt.Sprintf("Undefined variable %q", name)
	return NewError(msg)
}

func (env *Environment) Set(name string, value Object) {
	env.store[name] = value
}

func NewError(message string) *Error {
	return &Error{Message: message}
}
