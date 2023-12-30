package evaluator

import (
	"fmt"

	"fungo/ast"
	"fungo/object"
	"fungo/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

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

func evalIntegerInfixExpression(operator string, leftNode, rightNode object.Object) object.Object {
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

func evalPrefixExpression(prefixExpression *ast.PrefixExpression) object.Object {
	operator, right := prefixExpression.Operator, Eval(prefixExpression.Right)

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

func evalInfixExpression(infixExpression *ast.InfixExpression) object.Object {
	operator, left, right := infixExpression.Operator, Eval(infixExpression.Left), Eval(infixExpression.Right)

	if isError(left) {
		return left
	}

	if isError(right) {
		return right
	}

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
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

func evalIfExpression(expression *ast.IfExpression) object.Object {
	condition := Eval(expression.Condition)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(expression.IfCondition)
	} else if expression.ElseCondition != nil {
		return Eval(expression.ElseCondition)
	} else {
		return NULL
	}
}

func evalReturnExpression(expression *ast.ReturnStatement) object.Object {
	value := Eval(expression.ReturnValue)

	if isError(value) {
		return value
	}

	return &object.ReturnValue{Value: value}
}

func evalIntegerLiteral(integerLiteral *ast.IntegerLiteral) object.Object {
	return &object.Integer{Value: integerLiteral.Value}
}

func newError(format string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.IfExpression:
		return evalIfExpression(node)

	case *ast.ReturnStatement:
		return evalReturnExpression(node)

	// Expressions
	case *ast.PrefixExpression:
		return evalPrefixExpression(node)

	case *ast.InfixExpression:
		return evalInfixExpression(node)

	// Values
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	}

	return nil
}
