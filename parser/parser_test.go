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
	t.Equal("let", tokenLiteral)

	letStatement, ok := statement.(*ast.LetStatement)
	t.True(ok, "*ast.LetStatement")

	letStatementValue := letStatement.Name.Value
	t.Equal(letStatementValue, name)

	letStatementTokenLiteral := letStatement.Name.TokenLiteral()
	t.Equal(letStatementTokenLiteral, name)
}

func (t *ParserTestSuite) testIntegerLiteral(ex ast.Expression, value int64) {
	num, ok := ex.(*ast.IntegerLiteral)

	t.True(ok, "*ast.IntegerLiteral")
	t.Equal(value, num.Value)
}

func (t *ParserTestSuite) testIdentifier(exp ast.Expression, value string) {
	identifier, ok := exp.(*ast.Identifier)
	t.True(ok, "*ast.Identifier")

	t.Equal(value, identifier.Value)
	t.Equal(value, identifier.TokenLiteral())
}

func (t *ParserTestSuite) testBooleanLiteral(exp ast.Expression, value bool) {
	boolean, ok := exp.(*ast.Boolean)
	t.True(ok, "*ast.Boolean")

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

	t.True(ok, "*ast.InfixExpression")
	t.testLiteralExpression(opExp.Left, left)

	t.Equal(opExp.Operator, operator)
	t.testLiteralExpression(opExp.Right, right)
}

func (t *ParserTestSuite) TestBooleanExpression() {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, test := range tests {
		parser := NewParser(lexer.NewLexer(test.input))
		program := parser.ParseProgram()

		t.Empty(parser.Errors())
		t.Len(program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		t.True(ok, "*ast.ExpressionStatement")

		boolean, ok := statement.Expression.(*ast.Boolean)
		t.True(ok, "*ast.Boolean")
		t.Equal(test.expected, boolean.Value)
	}
}

func (t *ParserTestSuite) TestParsingInfixExpressions() {
	tests := []struct {
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

	for _, test := range tests {
		parser := NewParser(lexer.NewLexer(test.input))
		program := parser.ParseProgram()

		t.Empty(parser.Errors())
		t.Len(program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		t.True(ok, "*ast.ExpressionStatement")

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
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
	}

	for _, test := range tests {
		parser := NewParser(lexer.NewLexer(test.input))
		program := parser.ParseProgram()

		t.Empty(parser.Errors())

		actual := program.String()
		t.Equal(test.expected, actual)
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
		parser := NewParser(lexer.NewLexer(test.input))
		program := parser.ParseProgram()

		t.Empty(parser.Errors())
		t.Len(program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		t.True(ok, "*ast.ExpressionStatement")

		exp, ok := statement.Expression.(*ast.PrefixExpression)
		t.True(ok, "*ast.PrefixExpression")

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
		parser := NewParser(lexer.NewLexer(test.input))
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

	parser := NewParser(lexer.NewLexer(input))
	program := parser.ParseProgram()

	t.Empty(parser.Errors())
	t.Len(program.Statements, 3)

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		t.True(ok, "*ast.ReturnStatement")

		tokenLiteral := returnStatement.TokenLiteral()
		t.Equal("return", tokenLiteral)
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
		parser := NewParser(lexer.NewLexer(test.input))
		program := parser.ParseProgram()

		t.Empty(parser.Errors())
		t.Len(program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ReturnStatement)
		t.True(ok, "*ast.ReturnStatement")

		t.Equal("return", statement.TokenLiteral())
		t.testLiteralExpression(statement.ReturnValue, test.expectedValue)
	}
}

func (t *ParserTestSuite) TestIntergerLiteralExpression() {
	input := "5;"

	parser := NewParser(lexer.NewLexer(input))
	program := parser.ParseProgram()

	t.Empty(parser.Errors())
	t.Len(program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok, "*ast.ExpressionStatement")

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	t.True(ok, "*ast.IntegerLiteral")

	t.Equal(int64(5), literal.Value)
	t.Equal("5", literal.TokenLiteral())
}

func (t *ParserTestSuite) TestIfExpression() {
	input := `if (x < y) { x }`

	parser := NewParser(lexer.NewLexer(input))
	program := parser.ParseProgram()

	t.Empty(parser.Errors())
	t.Len(program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok, "*ast.ExpressionStatement")

	expression, ok := statement.Expression.(*ast.IfExpression)
	t.True(ok, "*ast.IfExpression")

	t.testInfixExpression(expression.Condition, "x", "<", "y")

	t.Len(expression.IfCondition.Statements, 1)

	ifCondition, ok := expression.IfCondition.Statements[0].(*ast.ExpressionStatement)
	t.True(ok, "*ast.ExpressionStatement")

	t.testIdentifier(ifCondition.Expression, "x")
	t.Nil(expression.ElseCondition)
}

func (t *ParserTestSuite) TestIfElseExpression() {
	input := `if (x < y) { x } else { y }`

	parser := NewParser(lexer.NewLexer(input))
	program := parser.ParseProgram()

	t.Empty(parser.Errors())
	t.Len(program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok, "*ast.ExpressionStatement")

	exp, ok := stmt.Expression.(*ast.IfExpression)
	t.True(ok, "*ast.IfExpression")

	t.testInfixExpression(exp.Condition, "x", "<", "y")
	t.Len(exp.IfCondition.Statements, 1)

	ifCondition, ok := exp.IfCondition.Statements[0].(*ast.ExpressionStatement)
	t.True(ok, "*ast.ExpressionStatement")

	t.testIdentifier(ifCondition.Expression, "x")
	t.Len(exp.ElseCondition.Statements, 1)

	elseCondition, ok := exp.ElseCondition.Statements[0].(*ast.ExpressionStatement)
	t.True(ok, "*ast.ExpressionStatement")

	t.testIdentifier(elseCondition.Expression, "y")
}

func (t *ParserTestSuite) TestFunctionLiteralParsing() {
	input := `fn(x, y) { x + y; }`

	parser := NewParser(lexer.NewLexer(input))
	program := parser.ParseProgram()
	t.Empty(parser.Errors())

	t.Len(program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok, "*ast.ExpressionStatement")

	function, ok := statement.Expression.(*ast.FunctionLiteral)
	t.True(ok, "*ast.FunctionLiteral")
	t.Len(function.Parameters, 2)

	t.testLiteralExpression(function.Parameters[0], "x")
	t.testLiteralExpression(function.Parameters[1], "y")

	t.Len(function.Body.Statements, 1)

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	t.True(ok, "*ast.ExpressionStatement)")

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
		parser := NewParser(lexer.NewLexer(test.input))
		program := parser.ParseProgram()
		t.Empty(parser.Errors())

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		t.True(ok, "*ast.ExpressionStatement")

		function, ok := statement.Expression.(*ast.FunctionLiteral)
		t.True(ok, "*ast.FunctionLiteral")

		t.Equal(len(test.expectedParams), len(function.Parameters))

		for i, identifier := range test.expectedParams {
			t.testLiteralExpression(function.Parameters[i], identifier)
		}
	}
}

func (t *ParserTestSuite) TestExpressionParsing() {
	input := "add(1, 2 * 3, 4 + 5);"

	parser := NewParser(lexer.NewLexer(input))
	program := parser.ParseProgram()
	t.Empty(parser.Errors())

	t.Len(program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok, "*ast.ExpressionStatement")

	expression, ok := statement.Expression.(*ast.CallExpression)
	t.True(ok, "*ast.CallExpression")

	t.testIdentifier(expression.Function, "add")
	t.Len(expression.Arguments, 3)

	t.testLiteralExpression(expression.Arguments[0], 1)
	t.testInfixExpression(expression.Arguments[1], 2, "*", 3)
	t.testInfixExpression(expression.Arguments[2], 4, "+", 5)
}

func (t *ParserTestSuite) TestStringLiteralParsing() {
	input := `"hello world"`

	parser := NewParser(lexer.NewLexer(input))
	program := parser.ParseProgram()
	t.Empty(parser.Errors())

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok, "*ast.ExpressionStatement")

	literal, ok := statement.Expression.(*ast.StringLiteral)
	t.True(ok, "*ast.StringLiteral")

	t.Equal("hello world", literal.Value)
}

func (t *ParserTestSuite) TestParsingArrayLiterals() {
	input := `[1, 2 * 2, 3 + 3]`

	parser := NewParser(lexer.NewLexer(input))
	program := parser.ParseProgram()
	t.Empty(parser.Errors())

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok, "*ast.ExpressionStatement")

	array, ok := statement.Expression.(*ast.ArrayLiteral)
	t.True(ok, "*ast.ArrayLiteral")
	t.Len(array.Elements, 3)

	t.testIntegerLiteral(array.Elements[0], 1)
	t.testInfixExpression(array.Elements[1], 2, "*", 2)
	t.testInfixExpression(array.Elements[2], 3, "+", 3)
}

func (t *ParserTestSuite) TestParsingIndexExpression() {
	input := `myArray[1 + 1]`

	parser := NewParser(lexer.NewLexer(input))
	program := parser.ParseProgram()
	t.Empty(parser.Errors())

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok, "*ast.ExpressionStatement")

	indexExpression, ok := statement.Expression.(*ast.IndexExpression)
	t.True(ok, "*ast.IndexExpression")

	t.Equal("myArray", indexExpression.Ref.String())
	t.testInfixExpression(indexExpression.Index, 1, "+", 1)
}

func (t *ParserTestSuite) TestParsingHashLiteralStringKeys() {
	tests := []struct {
		input    string
		expected map[string]int64
	}{
		{
			input: `{"one": 1, "two": 2, "three": 3}`,
			expected: map[string]int64{
				"one":   1,
				"two":   2,
				"three": 3,
			},
		},
	}

	for _, test := range tests {
		parser := NewParser(lexer.NewLexer(test.input))
		program := parser.ParseProgram()
		t.Empty(parser.Errors())

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		t.True(ok, "*ast.ExpressionStatement")

		result, ok := statement.Expression.(*ast.HashLiteral)
		t.True(ok, "*ast.HashLiteral")

		t.Equal(len(result.Pairs), len(test.expected))

		for key, value := range result.Pairs {
			literal, ok := key.(*ast.StringLiteral)
			t.True(ok, "*ast.StringLiteral")

			t.testIntegerLiteral(value, test.expected[literal.String()])
		}
	}
}

func (t *ParserTestSuite) TestParsingEmptyHashLiteral() {
	input := `{}`

	parser := NewParser(lexer.NewLexer(input))
	program := parser.ParseProgram()
	t.Empty(parser.Errors())

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	t.True(ok, "*ast.ExpressionStatement")

	result, ok := statement.Expression.(*ast.HashLiteral)
	t.True(ok, "*ast.HashLiteral")

	t.Len(result.Pairs, 0)
}
