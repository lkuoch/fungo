package evaluator

import (
	"fungo/lexer"
	"fungo/object"
	"fungo/parser"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EvaluatorTestSuite struct {
	suite.Suite
}

func TestEvaluatorTestSuite(t *testing.T) {
	suite.Run(t, &EvaluatorTestSuite{})
}

func (t *EvaluatorTestSuite) testNullObject(obj object.Object) {
	t.Equal(obj, NULL)
}

func (t *EvaluatorTestSuite) testIntegerObject(obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	t.True(ok)

	t.Equal(expected, result.Value)
}

func (t *EvaluatorTestSuite) testBooleanObject(obj object.Object, expected bool) {
	result, ok := obj.(*object.Boolean)
	t.True(ok)

	t.Equal(expected, result.Value)
}

func (t *EvaluatorTestSuite) testEval(input string) object.Object {
	parser := parser.New(lexer.New(input))
	program := parser.ParseProgram()

	return Eval(program)
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
		evaluated := t.testEval(test.input)
		t.testIntegerObject(evaluated, test.expected)
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
		evaluated := t.testEval(test.input)
		t.testBooleanObject(evaluated, test.expected)
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
		evaluated := t.testEval(test.input)
		integer, ok := test.expected.(int)

		if ok {
			t.testIntegerObject(evaluated, int64(integer))
		} else {
			t.testNullObject(evaluated)
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
		evaluated := t.testEval(test.input)
		t.testBooleanObject(evaluated, test.expected)
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
		evaluated := t.testEval(test.input)
		t.testIntegerObject(evaluated, test.expected)
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
		evaluated := t.testEval(test.input)
		t.testIntegerObject(evaluated, test.expected)
	}
}
