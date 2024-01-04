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
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

/* ================================= Integer ================================ */
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

/* ================================= Boolean ================================ */
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

/* ================================== Null ================================== */
type Null struct {
	Object
}

func (n Null) Type() ObjectType {
	return NULL_OBJ
}

func (n Null) Inspect() string {
	return "null"
}

/* =============================== ReturnValue ============================== */
type ReturnValue struct {
	Value Object
}

func (r ReturnValue) Type() ObjectType {
	return RETURN_VAL_OBJ
}

func (r ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

/* ================================== NOOP ================================== */
type Noop struct {
	Object
}

func (n Noop) Type() ObjectType {
	return NOOP_OBJ
}

func (n Noop) Inspect() string {
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

func (f Function) Inspect() string {
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

func (s String) Inspect() string {
	return s.Value
}

/* ================================= BuiltIn ================================ */
type BuiltInFunction func(args ...Object) Object

type BuiltIn struct {
	Fn BuiltInFunction
}

func (b BuiltIn) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b BuiltIn) Inspect() string {
	return "built-in function"
}

/* ================================== Error ================================= */
type Error struct {
	Message string
}

func (e Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e Error) Inspect() string {
	return "⛔️ ERROR: " + e.Message
}
