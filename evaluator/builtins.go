package evaluator

import "fungo/object"

func builtIn_len(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError("argument to `len` not supported. got=`%s`", args[0].Type())
	}
}

/* ====================== Map of all built in functions ===================== */
var builtInsMap = map[string]*object.BuiltIn{
	"len": {Fn: builtIn_len},
}
