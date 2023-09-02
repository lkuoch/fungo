package evaluator

import (
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

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
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
		return NULL
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
		return NULL
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case token.BANG:
		return evalBangOperatorExpression(right)

	case token.MINUS:
		return evalMinusOperatorExpression(right)

	default:
		return NULL
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == token.EQ:
		return nativeBoolToBooleanObject(left == right)
	case operator == token.NOT_EQ:
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL
	}
}

func evalIfExpression(expression *ast.IfExpression) object.Object {
	condition := Eval(expression.Condition)

	if isTruthy(condition) {
		return Eval(expression.IfCondition)
	} else if expression.ElseCondition != nil {
		return Eval(expression.ElseCondition)
	} else {
		return NULL
	}
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.PrefixExpression:
		return evalPrefixExpression(node.Operator, Eval(node.Right))

	case *ast.InfixExpression:
		return evalInfixExpression(node.Operator, Eval(node.Left), Eval(node.Right))

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.BlockStatement:
		return evalStatements(node.Statements)

	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return nil
}
