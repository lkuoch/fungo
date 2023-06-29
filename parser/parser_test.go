package parser

import (
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
		t.Errorf("letStatementValue not '%s'. got=%s", name, letStatementValue)
	}

	letStatementTokenLiteral := letStatement.Name.TokenLiteral()
	if letStatementTokenLiteral != name {
		t.Errorf("letStatementTokenLiteral not '%s'. got='%s'", name, letStatementTokenLiteral)
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
		t.Errorf("value not %d. got=%s", value, num.TokenLiteral())
		return false
	}

	return true
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
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
			t.Errorf("exp.Operator is not '%s'. got %q", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
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
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
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
		t.Errorf("ident.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}
