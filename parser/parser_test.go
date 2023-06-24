package parser

import (
	"fungo/ast"
	"fungo/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

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
