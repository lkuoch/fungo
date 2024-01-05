package evaluator

import (
	"fungo/ast"
	"fungo/object"
	"fungo/token"
)

var (
	NULL  = &object.Null{}
	NOOP  = &object.Noop{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			resultType := result.Type()
			if resultType == object.RETURN_VAL_OBJ || resultType == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value, ok := right.(*object.Integer)

	if !ok {
		return NULL
	}

	return &object.Integer{Value: -value.Value}
}

func evalIntegerInfixExpression(operator string, leftNode object.Object, rightNode object.Object) object.Object {
	left, ok := leftNode.(*object.Integer)
	if !ok {
		return NULL
	}

	right, ok := rightNode.(*object.Integer)
	if !ok {
		return NULL
	}

	switch operator {
	case token.PLUS:
		return &object.Integer{Value: left.Value + right.Value}

	case token.MINUS:
		return &object.Integer{Value: left.Value - right.Value}

	case token.ASTERISK:
		return &object.Integer{Value: left.Value * right.Value}

	case token.SLASH:
		return &object.Integer{Value: left.Value / right.Value}

	case token.LT:
		return nativeBoolToBooleanObject(left.Value < right.Value)

	case token.GT:
		return nativeBoolToBooleanObject(left.Value > right.Value)

	case token.EQ:
		return nativeBoolToBooleanObject(left.Value == right.Value)

	case token.NOT_EQ:
		return nativeBoolToBooleanObject(left.Value != right.Value)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, leftNode object.Object, rightNode object.Object) object.Object {
	left, ok := leftNode.(*object.String)
	if !ok {
		return NULL
	}

	right, ok := rightNode.(*object.String)
	if !ok {
		return NULL
	}

	switch operator {
	case token.PLUS:
		return &object.String{Value: left.Value + right.Value}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(prefixExpression *ast.PrefixExpression, env *object.Environment) object.Object {
	operator, right := prefixExpression.Operator, Eval(prefixExpression.Right, env)

	if isError(right) {
		return right
	}

	switch operator {
	case token.BANG:
		return evalBangOperatorExpression(right)

	case token.MINUS:
		return evalMinusOperatorExpression(right)

	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(infixExpression *ast.InfixExpression, env *object.Environment) object.Object {
	operator, left, right := infixExpression.Operator, Eval(infixExpression.Left, env), Eval(infixExpression.Right, env)

	if isError(left) {
		return left
	}

	if isError(right) {
		return right
	}

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == token.EQ:
		return nativeBoolToBooleanObject(left == right)
	case operator == token.NOT_EQ:
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(expression *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(expression.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(expression.IfCondition, env)
	} else if expression.ElseCondition != nil {
		return Eval(expression.ElseCondition, env)
	} else {
		return NULL
	}
}

func evalReturnExpression(expression *ast.ReturnStatement, env *object.Environment) object.Object {
	value := Eval(expression.ReturnValue, env)

	if isError(value) {
		return value
	}

	return &object.ReturnValue{Value: value}
}

func evalIntegerLiteral(integerLiteral *ast.IntegerLiteral) object.Object {
	return &object.Integer{Value: integerLiteral.Value}
}

func evalIdentifier(identifier *ast.Identifier, env *object.Environment) object.Object {
	if value, ok := env.Get(identifier.Value); ok {
		return value
	}

	if builtIn, ok := builtInsMap[identifier.Value]; ok {
		return builtIn
	}

	return newError("identifier not found: " + identifier.Value)
}

func evalLetStatement(statement *ast.LetStatement, env *object.Environment) object.Object {
	value := Eval(statement.Value, env)
	if isError(value) {
		return value
	}

	env.Set(statement.Name.Value, value)

	return NOOP
}

func evalFunctionLiteral(fn *ast.FunctionLiteral, env *object.Environment) object.Object {
	return &object.Function{
		Parameters: fn.Parameters,
		Body:       fn.Body,
		Env:        env,
	}
}

func evalStringLiteral(str *ast.StringLiteral, _env *object.Environment) object.Object {
	return &object.String{
		Value: str.Value,
	}
}

func evalArrayLiteral(array *ast.ArrayLiteral, env *object.Environment) object.Object {
	elements := evalExpressions(array.Elements, env)

	if len(elements) == 1 && isError(elements[0]) {
		return elements[0]
	}

	return &object.Array{
		Elements: elements,
	}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, exp := range exps {
		evaluated := Eval(exp, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func evalCallExpression(exp *ast.CallExpression, env *object.Environment) object.Object {
	fn := Eval(exp.Function, env)
	if isError(fn) {
		return fn
	}
	args := evalExpressions(exp.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	return applyFunction(fn, args)
}

func evalArrayIndexExpression(ref *object.Array, index *object.Integer) object.Object {
	idx := index.Value
	array := ref.Elements
	max := int64(len(array) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return array[idx]
}

func evalHashIndexExpression(hash *object.Hash, index object.Object) object.Object {
	key, ok := index.(object.Hashable)

	if !ok {
		return newError("unusable as hash key: `%s`", index.Type())
	}

	pair, ok := hash.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalIndexExpression(exp *ast.IndexExpression, env *object.Environment) object.Object {
	ref := Eval(exp.Ref, env)
	if isError(ref) {
		return ref
	}

	index := Eval(exp.Index, env)
	if isError(index) {
		return index
	}

	switch {
	case ref.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(ref.(*object.Array), index.(*object.Integer))
	case ref.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(ref.(*object.Hash), index)
	default:
		return newError("index operator not supported: %s", ref.Type())
	}
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		return evalReturnExpression(node, env)

	case *ast.LetStatement:
		evalLetStatement(node, env)

	// Expressions
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)

	case *ast.InfixExpression:
		return evalInfixExpression(node, env)

	case *ast.CallExpression:
		return evalCallExpression(node, env)

	case *ast.IndexExpression:
		return evalIndexExpression(node, env)

	// Values
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, env)

	case *ast.StringLiteral:
		return evalStringLiteral(node, env)

	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	}

	return nil
}
