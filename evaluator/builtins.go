package evaluator

import "fungo/object"

func builtIn_len(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to `len` not supported. got=`%s`", args[0].Type())
	}
}

func builtIn_first(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `first` must be `ARRAY`, got=`%s`", args[0].Type())
	}

	array := args[0].(*object.Array)
	if len(array.Elements) > 0 {
		return array.Elements[0]
	}

	return NULL
}

func builtIn_last(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `last` must be `ARRAY`, got=`%s`", args[0].Type())
	}

	array := args[0].(*object.Array)
	length := len(array.Elements)

	if length > 0 {
		return array.Elements[length-1]
	}

	return NULL
}

func builtIn_rest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `rest` must be `ARRAY`, got=%s", args[0].Type())
	}

	array := args[0].(*object.Array)
	length := len(array.Elements)

	if length > 0 {
		elements := make([]object.Object, length-1)
		copy(elements, array.Elements[1:length])
		return &object.Array{Elements: elements}
	}

	return NULL
}

func builtIn_push(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of 1rguments. got=%d, want=2", len(args))
	}

	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `push` must be `ARRAY`, got=`%s`", args[0].Type())
	}

	array := args[0].(*object.Array)
	length := len(array.Elements)

	elements := make([]object.Object, length+1)
	copy(elements, array.Elements)
	elements[length] = args[1]

	return &object.Array{Elements: elements}
}

/* ====================== Map of all built in functions ===================== */
var builtInsMap = map[string]*object.BuiltIn{
	"len": {
		FnName: "len",
		Fn:     builtIn_len,
	},
	"first": {
		FnName: "first",
		Fn:     builtIn_first,
	},
	"last": {
		FnName: "last",
		Fn:     builtIn_last,
	},
	"rest": {
		FnName: "rest",
		Fn:     builtIn_rest,
	},
	"push": {
		FnName: "push",
		Fn:     builtIn_push,
	},
}
