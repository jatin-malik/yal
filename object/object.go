package object

import "fmt"

type ObjectType string

const (
	IntegerObject ObjectType = "INTEGER"
	BooleanObject ObjectType = "BOOLEAN"
	NullObject    ObjectType = "NULL"
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
