package object

import (
	"bytes"
	"fmt"
	"fungo/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ    = "INTEGER"
	BOOLEAN_OBJ    = "BOOLEAN"
	NULL_OBJ       = "NULL"
	RETURN_VAL_OBJ = "RETURN_VALUE"
	ERROR_OBJ      = "ERROR"
	NOOP_OBJ       = "NOOP"
	FUNCTION_OBJ   = "FUNCTION"
	STRING_OBJ     = "STRING"
	BUILTIN_OBJ    = "BUILTIN"
	ARRAY_OBJ      = "ARRAY"
)

type Object interface {
	Type() ObjectType
	String() string
}

/* ================================= Integer ================================ */
type Integer struct {
	Object
	Value int64
}

func (u Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

/* ================================= Boolean ================================ */
type Boolean struct {
	Object
	Value bool
}

func (b Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

/* ================================== Null ================================== */
type Null struct {
	Object
}

func (n Null) Type() ObjectType {
	return NULL_OBJ
}

func (n Null) String() string {
	return "null"
}

/* =============================== ReturnValue ============================== */
type ReturnValue struct {
	Value Object
}

func (r ReturnValue) Type() ObjectType {
	return RETURN_VAL_OBJ
}

func (r ReturnValue) String() string {
	return r.Value.String()
}

/* ================================== NOOP ================================== */
type Noop struct {
	Object
}

func (n Noop) Type() ObjectType {
	return NOOP_OBJ
}

func (n Noop) String() string {
	return "noop"
}

/* ================================ Function ================================ */
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f Function) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range f.Parameters {
		params = append(params, param.String())
	}

	out.WriteString("fn(" + strings.Join(params, ",") + ") {\n" + f.Body.String() + "\n}")

	return out.String()
}

/* ================================= String ================================= */
type String struct {
	Value string
}

func (s String) Type() ObjectType {
	return STRING_OBJ
}

func (s String) String() string {
	return s.Value
}

/* ================================= BuiltIn ================================ */
type BuiltIn struct {
	Fn     func(args ...Object) Object
	FnName string
}

func (b BuiltIn) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b BuiltIn) String() string {
	return fmt.Sprintf("Built-in function: %s", b.FnName)
}

/* ================================== Array ================================= */
type Array struct {
	Elements []Object
}

func (a Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (a Array) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, element := range a.Elements {
		elements = append(elements, element.String())
	}

	out.WriteString("[" + strings.Join(elements, ", ") + "]")
	return out.String()
}

/* ================================== Error ================================= */
type Error struct {
	Message string
}

func (e Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e Error) String() string {
	return "⛔️ ERROR: " + e.Message
}
