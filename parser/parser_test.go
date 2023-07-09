package parser

import (
	"fmt"
	"fungo/ast"
	"fungo/lexer"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	tokenLiteral := statement.TokenLiteral()
	if tokenLiteral != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", tokenLiteral)
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.Letstatement. got=%T", statement)
		return false
	}

	letStatementValue := letStatement.Name.Value
	if letStatementValue != name {
		t.Errorf("letStatementValue not '%q'. got=%q", name, letStatementValue)
	}

	letStatementTokenLiteral := letStatement.Name.TokenLiteral()
	if letStatementTokenLiteral != name {
		t.Errorf("letStatementTokenLiteral not '%q'. got='%q'", name, letStatementTokenLiteral)
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, ex ast.Expression, value int64) bool {
	num, ok := ex.(*ast.IntergerLiteral)
	if !ok {
		t.Errorf("ex not *ast.IntegerLiteral. got=%T", ex)
		return false
	}

	if num.Value != value {
		t.Errorf("value not %d. got=%q", value, num.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.IDentifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %q. got=%q", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %q. got=%q", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)

	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value")
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral not %t. got=%q", value, boolean.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	default:
		t.Errorf("type of exp not handled. got=%T", exp)
		return false
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%q)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %q. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
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

	for _, tt := range infixTests {
		parser := New(lexer.New(tt.input))
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain single statments. got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, statement.Expression, tt.lhs, tt.operator, tt.rhs) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
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
	}

	for _, tt := range tests {
		parser := New(lexer.New(tt.input))
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected %q. got=%q", tt.expected, actual)
		}

	}
}

func TestParsingPrefixExpressions(t *testing.T) {
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

	for _, tt := range prefixTests {
		parser := New(lexer.New(tt.input))
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain single statments. got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("statement is not ast.PrefixExpression. got=%T", statement.Expression)
		}

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not '%q'. got %q", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statments. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, token := range tests {
		statement := program.Statements[i]

		if !testLetStatement(t, statement, token.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;
	`

	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	statementsLen := len(program.Statements)
	if statementsLen != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", statementsLen)
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("statement not *ast.ReturnStatement, got = %T", statement)
			continue
		}

		tokenLiteral := returnStatement.TokenLiteral()
		if tokenLiteral != "return" {
			t.Errorf("returnStatement.TokenLiteral not 'return', got=%q", tokenLiteral)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	statementsLen := len(program.Statements)
	if (statementsLen) != 1 {
		t.Errorf("program has not enough statements. got=%d", statementsLen)
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not a ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", statement.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %q. got=%q", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %q. got=%q", "foobar", ident.TokenLiteral())
	}
}

func TestIntergerLiteralExpression(t *testing.T) {
	input := "5;"

	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	statementsLen := len(program.Statements)
	if (statementsLen) != 1 {
		t.Errorf("program has not enough statements. got=%d", statementsLen)
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not a ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntergerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntergerLiteral. got=%T", statement.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("ident.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral not %q. got=%q", "5", literal.TokenLiteral())
	}
}
