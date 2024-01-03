package evaluator

import (
	"fungo/lexer"
	"fungo/object"
	"fungo/parser"
	"fungo/utils"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EvaluatorTestSuite struct {
	suite.Suite
}

func TestEvaluatorTestSuite(t *testing.T) {
	suite.Run(t, &EvaluatorTestSuite{})
}

func (t *EvaluatorTestSuite) testNullObject(actual object.Object) {
	t.Equal(actual, NULL)
}

func (t *EvaluatorTestSuite) testIntegerObject(expected int64, actual object.Object) {
	result, ok := actual.(*object.Integer)
	t.True(ok)

	t.Equal(expected, result.Value)
}

func (t *EvaluatorTestSuite) testBooleanObject(expected bool, actual object.Object) {
	result, ok := actual.(*object.Boolean)
	t.True(ok)

	t.Equal(expected, result.Value)
}

func (t *EvaluatorTestSuite) testFunctionObject(expectedParamsLen int, expectedParams []string, expectedBody string, actual object.Object) {
	result, ok := actual.(*object.Function)
	t.True(ok)

	t.Len(result.Parameters, expectedParamsLen)
	t.Equal(expectedParams, utils.MapString(result.Parameters))
	t.Equal(expectedBody, result.Body.String())
}

func (t *EvaluatorTestSuite) testErrorObject(expected string, actual object.Object) {
	result, ok := actual.(*object.Error)
	t.True(ok)

	t.Equal(expected, result.Message)
}

func (t *EvaluatorTestSuite) testEval(input string) object.Object {
	parser := parser.New(lexer.New(input))
	program := parser.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func (t *EvaluatorTestSuite) TestEvalIntegerExpression() {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testIntegerObject(test.expected, result)
	}
}

func (t *EvaluatorTestSuite) TestEvalBooleanExpression() {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testBooleanObject(test.expected, result)
	}
}

func (t *EvaluatorTestSuite) TestIfElseExpressions() {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) {10}", 10},
		{"if (false) {10}", nil},
		{"if (1) {10}", 10},
		{"if (1 < 2) {10}", 10},
		{"if (1 > 2) {10}", nil},
		{"if (1 > 2) {10} else {20}", 20},
		{"if (1 < 2) {10} else {20}", 10},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		integer, ok := test.expected.(int)

		if ok {
			t.testIntegerObject(int64(integer), result)
		} else {
			t.testNullObject(result)
		}
	}
}

func (t *EvaluatorTestSuite) TestBangOperator() {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!5", true},
		{"!!false", false},
		{"!!true", true},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testBooleanObject(test.expected, result)
	}
}

func (t *EvaluatorTestSuite) TestIntegerExpression() {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testIntegerObject(test.expected, result)
	}
}

func (t *EvaluatorTestSuite) TestReturnStatements() {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9", 10},
		{"return 2 * 5; 9", 10},
		{"9; return 2 * 5;", 10},
		{`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
			}

			return 1;
		`, 10},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testIntegerObject(test.expected, result)
	}
}

func (t *EvaluatorTestSuite) TestErrorHandling() {
	tests := []struct {
		input    string
		expected string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
			}

			return 1;
		 `, "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testErrorObject(test.expected, result)
	}
}

func (t *EvaluatorTestSuite) TestLetStatement() {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testIntegerObject(test.expected, result)
	}
}

func (t *EvaluatorTestSuite) TestFuncObject() {
	tests := []struct {
		input             string
		expectedParamsLen int
		expectedParams    []string
		expectedBody      string
	}{
		{"fn(x) { x + 2; };", 1, []string{"x"}, "(x + 2)"},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testFunctionObject(test.expectedParamsLen, test.expectedParams, test.expectedBody, result)
	}
}

func (t *EvaluatorTestSuite) TestFunctionApplication() {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { return x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
		{`
			let newAdder = fn(x) {
				fn(y) { x + y };
			};

			let addTwo = newAdder(2);
			addTwo(2);
		 `, 4,
		},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testIntegerObject(test.expected, result)
	}
}
