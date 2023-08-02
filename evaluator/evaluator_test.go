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

func (t *EvaluatorTestSuite) testIntegerObject(obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	t.True(ok)

	t.Equal(result.Value, expected)
}

func (t *EvaluatorTestSuite) testBooleanObject(obj object.Object, expected bool) {
	result, ok := obj.(*object.Boolean)
	t.True(ok)

	t.Equal(result.Value, expected)
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
	}

	for _, test := range tests {
		evaluated := t.testEval(test.input)
		t.testBooleanObject(evaluated, test.expected)
	}
}
