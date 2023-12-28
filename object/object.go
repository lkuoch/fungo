package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ    = "INTEGER"
	BOOLEAN_OBJ    = "BOOLEAN"
	NULL_OBJ       = "NULL"
	RETURN_VAL_OBJ = "RETURN_VALUE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Specifications
type Integer struct {
	Object
	Value int64
}

func (u Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Object
	Value bool
}

func (b Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct {
	Object
}

func (n Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (r ReturnValue) Type() ObjectType {
	return RETURN_VAL_OBJ
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}
