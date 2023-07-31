package parser

import (
	"fmt"
	"fungo/ast"
	"fungo/lexer"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite
}

func TestParserTestSuite(t *testing.T) {
	suite.Run(t, &ParserTestSuite{})
}

func (t *ParserTestSuite) testLetStatement(statement ast.Statement, name string) {
	tokenLiteral := statement.TokenLiteral()
	t.Equal(tokenLiteral, "let")

	letStatement, ok := statement.(*ast.LetStatement)
	t.True(ok)

	letStatementValue := letStatement.Name.Value
	t.Equal(name, letStatementValue)

	letStatementTokenLiteral := letStatement.Name.TokenLiteral()
	t.Equal(name, letStatementTokenLiteral)
}

func (t *ParserTestSuite) testIntegerLiteral(ex ast.Expression, value int64) {
	num, ok := ex.(*ast.IntergerLiteral)

	t.True(ok)
	t.Equal(num.Value, value)
}

func (t *ParserTestSuite) testIdentifier(exp ast.Expression, value string) {
	identifier, ok := exp.(*ast.Identifier)
	t.True(ok)

	t.Equal(identifier.Value, value)
	t.Equal(identifier.TokenLiteral(), value)
}

func (t *ParserTestSuite) testBooleanLiteral(exp ast.Expression, value bool) {
	boolean, ok := exp.(*ast.Boolean)
	t.True(ok)

	t.Equal(boolean.Value, value)
	t.Equal(boolean.TokenLiteral(), fmt.Sprintf("%t", value))
}

func (t *ParserTestSuite) testLiteralExpression(exp ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case int:
		t.testIntegerLiteral(exp, int64(v))
	case int64:
		t.testIntegerLiteral(exp, v)
	case string:
		t.testIdentifier(exp, v)
	case bool:
		t.testBooleanLiteral(exp, v)
	default:
		t.FailNow("type of exp not handled. got=%T", exp)
	}
}

func (t *ParserTestSuite) testInfixExpression(exp ast.Expression, left interface{}, operator string, right interface{}) {
	opExp, ok := exp.(*ast.InfixExpression)

	t.True(ok)
	t.testLiteralExpression(opExp.Left, left)

	t.Equal(opExp.Operator, operator)
	t.testLiteralExpression(opExp.Right, right)
}

func (t *ParserTestSuite) TestParsingInfixExpressions() {
	infixTests := []struct {
		input    string
		lhs      interface{}
		operator string
		rhs      interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, test := range infixTests {
		parser := New(lexer.New(test.input))
		program := parser.ParseProgram()

		t.Empty(parser.Errors())
		t.Len(program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		t.True(ok)

		t.testInfixExpression(statement.Expression, test.lhs, test.operator, test.rhs)
	}
}

func (t *ParserTestSuite) TestOperatorPrecedenceParsing() {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}

	for _, test := range tests {
		parser := New(lexer.New(test.input))
		program := parser.ParseProgram()

		t.Empty(parser.Errors())

		actual := program.String()
		t.Equal(actual, test.expected)
	}
}

func (t *ParserTestSuite) TestParsingPrefixExpressions() {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, test := range prefixTests {
		parser := New(lexer.New(test.input))
		program := parser.ParseProgram()

		t.Empty(parser.Errors())
		t.Len(program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		t.True(ok)

		exp, ok := statement.Expression.(*ast.PrefixExpression)
		t.True(ok)

		t.Equal(exp.Operator, test.operator)
		t.testLiteralExpression(exp.Right, test.value)
	}
}

func (t *ParserTestSuite) TestLetStatements() {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, test := range tests {
		parser := New(lexer.New(test.input))
		program := parser.ParseProgram()

		t.Empty(parser.Errors())
		t.NotNil(program)
		t.Len(program.Statements, 1)

		statement := program.Statements[0]
		t.testLetStatement(statement, test.expectedIdentifier)

		value := statement.(*ast.LetStatement).Value
		t.testLiteralExpression(value, test.expectedValue)
	}
}

func (t *ParserTestSuite) TestReturnStatement() {
	input := `
		return 5;
		return 10;
		return 993322;
	`

	parser := New(lexer.New(input))
	program := parser.ParseProgram()

	t.Empty(parser.Errors())
	t.Len(program.Statements, 3)

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		t.True(ok)

		tokenLiteral := returnStatement.TokenLiteral()
		t.Equal(tokenLiteral, "return")
	}
}

func (t *ParserTestSuite) TestIdentifierExpression() {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, test := range tests {
		parser := New(lexer.New(test.input))
		program := parser.ParseProgram()

		t.Empty(parser.Errors())
		t.Len(program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ReturnStatement)
		t.True(ok)

		t.Equal(statement.TokenLiteral(), "return")
		t.testLiteralExpression(statement.ReturnValue, test.expectedValue)
	}
}

func (t *ParserTestSuite) TestIntergerLiteralExpression() {
	input := "5;"

	parser := New(lexer.New(input))
	program := parser.ParseProgram()

	t.Empty(parser.Errors())
	t.Len(program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok)

	literal, ok := statement.Expression.(*ast.IntergerLiteral)
	t.True(ok)

	t.Equal(literal.Value, int64(5))
	t.Equal(literal.TokenLiteral(), "5")
}

func (t *ParserTestSuite) TestIfExpression() {
	input := `if (x < y) { x }`

	parser := New(lexer.New(input))
	program := parser.ParseProgram()

	t.Empty(parser.Errors())
	t.Len(program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok)

	expression, ok := statement.Expression.(*ast.IfExpression)
	t.True(ok)

	t.testInfixExpression(expression.Condition, "x", "<", "y")

	t.Len(expression.IfCondition.Statements, 1)

	ifCondition, ok := expression.IfCondition.Statements[0].(*ast.ExpressionStatement)
	t.True(ok)

	t.testIdentifier(ifCondition.Expression, "x")
	t.Nil(expression.ElseCondition)
}

func (t *ParserTestSuite) TestIfElseExpression() {
	input := `if (x < y) { x } else { y }`

	parser := New(lexer.New(input))
	program := parser.ParseProgram()

	t.Empty(parser.Errors())
	t.Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	t.True(ok)

	t.testInfixExpression(exp.Condition, "x", "<", "y")
	t.Len(exp.IfCondition.Statements, 1)

	ifCondition, ok := exp.IfCondition.Statements[0].(*ast.ExpressionStatement)
	t.True(ok)

	t.testIdentifier(ifCondition.Expression, "x")
	t.Len(exp.ElseCondition.Statements, 1)

	elseCondition, ok := exp.ElseCondition.Statements[0].(*ast.ExpressionStatement)
	t.True(ok)

	t.testIdentifier(elseCondition.Expression, "y")
}

func (t *ParserTestSuite) TestFunctionLiteralParsing() {
	input := `fn(x, y) { x + y; }`

	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	t.Empty(parser.Errors())

	t.Len(program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok)

	function, ok := statement.Expression.(*ast.FunctionLiteral)
	t.True(ok)
	t.Len(function.Parameters, 2)

	t.testLiteralExpression(function.Parameters[0], "x")
	t.testLiteralExpression(function.Parameters[1], "y")

	t.Len(function.Body.Statements, 1)

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	t.True(ok)

	t.testInfixExpression(bodyStatement.Expression, "x", "+", "y")
}

func (t *ParserTestSuite) TestFunctionParameterParsing() {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, test := range tests {
		parser := New(lexer.New(test.input))
		program := parser.ParseProgram()
		t.Empty(parser.Errors())

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		t.True(ok)

		function, ok := statement.Expression.(*ast.FunctionLiteral)
		t.True(ok)

		t.Equal(len(function.Parameters), len(test.expectedParams))

		for i, identifier := range test.expectedParams {
			t.testLiteralExpression(function.Parameters[i], identifier)
		}
	}
}

func (t *ParserTestSuite) TestExpressionParsing() {
	input := "add(1, 2 * 3, 4 + 5);"

	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	t.Empty(parser.Errors())

	t.Len(program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok)

	expression, ok := statement.Expression.(*ast.CallExpression)
	t.True(ok)

	t.testIdentifier(expression.Function, "add")
	t.Len(expression.Arguments, 3)

	t.testLiteralExpression(expression.Arguments[0], 1)
	t.testInfixExpression(expression.Arguments[1], 2, "*", 3)
	t.testInfixExpression(expression.Arguments[2], 4, "+", 5)
}
