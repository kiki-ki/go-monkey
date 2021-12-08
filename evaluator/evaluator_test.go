package evaluator_test

import (
	"testing"

	"github.com/kiki-ki/go-monkey/evaluator"
	"github.com/kiki-ki/go-monkey/lexer"
	"github.com/kiki-ki/go-monkey/object"
	"github.com/kiki-ki/go-monkey/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	cases := []struct {
		desc  string
		want  int64
		input string
	}{
		{
			desc:  "five",
			want:  5,
			input: "5",
		},
		{
			desc:  "ten",
			want:  10,
			input: "10",
		},
	}
	for _, c := range cases {
		evaluated := testEval(c.input)
		testIntegerObject(t, evaluated, c.want)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return evaluator.Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, want int64) bool {
	res, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if res.Value != want {
		t.Errorf("object has wrong value. got=%d, want=%d", res.Value, want)
		return false
	}
	return true
}
