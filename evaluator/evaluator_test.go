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
	t.True(ok, "*object.Integer")

	t.Equal(expected, result.Value)
}

func (t *EvaluatorTestSuite) testStringObject(expected string, actual object.Object) {
	result, ok := actual.(*object.String)
	t.True(ok, "*object.String")

	t.Equal(expected, result.Value)
}

func (t *EvaluatorTestSuite) testBooleanObject(expected bool, actual object.Object) {
	result, ok := actual.(*object.Boolean)
	t.True(ok, "*object.Boolean")

	t.Equal(expected, result.Value)
}

func (t *EvaluatorTestSuite) testFunctionObject(expectedParams []string, expectedBody string, actual object.Object) {
	result, ok := actual.(*object.Function)
	t.True(ok, "*object.Function")

	t.Equal(len(expectedParams), len(result.Parameters))
	t.Equal(expectedParams, utils.MapString(result.Parameters))
	t.Equal(expectedBody, result.Body.String())
}

func (t *EvaluatorTestSuite) testArrayObject(expected interface{}, actual object.Object) {
	result, ok := actual.(*object.Array)
	t.True(ok, "*object.Array")

	switch expected := expected.(type) {
	case []int:
		t.Equal(len(expected), len(result.Elements))
		t.Equal(expected, utils.MapInt(result.Elements))
	case []string:
		t.Equal(len(expected), len(result.Elements))
		t.Equal(expected, utils.MapString(result.Elements))
	default:
		t.Fail("testArrayObject edge case not handled", expected)
	}
}

func (t *EvaluatorTestSuite) testBuiltInFunctionObject(expected interface{}, actual object.Object) {
	switch expected := expected.(type) {
	case int:
		t.testIntegerObject(int64(expected), actual)
	case nil:
		t.testNullObject(actual)
	case string:
		t.testErrorObject(expected, actual)
	case []int:
	case []string:
		t.testArrayObject(expected, actual)

	default:
		t.Fail("BuiltInFunction edge case not handled", expected)
	}
}

func (t *EvaluatorTestSuite) testArrayLiteral(expected []string, actual object.Object) {
	result, ok := actual.(*object.Array)
	t.True(ok, "*object.Array")

	t.Equal(len(expected), len(result.Elements))
	t.Equal(expected, utils.MapString(result.Elements))
}

func (t *EvaluatorTestSuite) testArrayIndexExpression(expected interface{}, actual object.Object) {
	switch expected := expected.(type) {
	case int:
		t.testIntegerObject(int64(expected), actual)
	default:
		t.testNullObject(actual)
	}
}

func (t *EvaluatorTestSuite) testErrorObject(expected string, actual object.Object) {
	result, ok := actual.(*object.Error)
	t.True(ok, "*object.Error")

	t.Equal(expected, result.Message)
}

func (t *EvaluatorTestSuite) testEval(input string) object.Object {
	parser := parser.NewParser(lexer.NewLexer(input))
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
		{`"Hello" - "World"`, "unknown operator: STRING - STRING"},
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
		input          string
		expectedParams []string
		expectedBody   string
	}{
		{"fn(x) { x + 2; };", []string{"x"}, "(x + 2)"},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testFunctionObject(test.expectedParams, test.expectedBody, result)
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

func (t *EvaluatorTestSuite) TestStringLiteral() {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world"`, "hello world"},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testStringObject(test.expected, result)
	}
}

func (t *EvaluatorTestSuite) TestStringConcatenation() {
	tests := []struct {
		input    string
		expected string
	}{
		{`"Hello" + " " + "World"`, "Hello World"},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testStringObject(test.expected, result)
	}
}

func (t *EvaluatorTestSuite) TestBuiltInFunctions() {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported. got=`INTEGER`"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`first([1, 2, 3])`, 1},
		{`first([])`, nil},
		{`first(1)`, "argument to `first` must be `ARRAY`, got=`INTEGER`"},
		{`last([1, 2, 3])`, 3},
		{`last([])`, nil},
		{`last(1)`, "argument to `last` must be `ARRAY`, got=`INTEGER`"},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, nil},
		{`push([], 1)`, []int{1}},
		{`push([], "hello")`, []string{"hello"}},
		{`push(1, 1)`, "argument to `push` must be `ARRAY`, got=`INTEGER`"},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testBuiltInFunctionObject(test.expected, result)
	}
}

func (t *EvaluatorTestSuite) TestArrayLiterals() {
	tests := []struct {
		input    string
		expected []string
	}{
		{`[1, 2 * 2, 3 + 3]`, []string{"1", "4", "6"}},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testArrayLiteral(test.expected, result)
	}
}

func (t *EvaluatorTestSuite) TestArrayIndexExpressions() {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"let myArray = [1, 2, 3]; myArray[2];", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, test := range tests {
		result := t.testEval(test.input)
		t.testArrayIndexExpression(test.expected, result)
	}
}
