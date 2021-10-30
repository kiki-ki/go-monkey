package parser_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/kiki-ki/go-monkey/ast"
	"github.com/kiki-ki/go-monkey/lexer"
	"github.com/kiki-ki/go-monkey/parser"
)

func TestLetStatements(t *testing.T) {
	cases := []struct {
		input     string
		wantIdent string
		wantVal   interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foo = y;", "foo", "y"},
	}

	for _, tt := range cases {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		fmt.Println(program.Statements)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}
		s := program.Statements[0]
		if !testLetStatement(t, s, tt.wantIdent) {
			return
		}
		val := s.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.wantVal) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral does not 'let' got=%q", s.TokenLiteral())
		return false
	}
	letS, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s does not *ast.LetStatement got=%q", s)
		return false
	}
	if letS.Name.Value != name {
		t.Errorf("LetStatement.Name.Value does not '%s' got=%s", name, letS.Name.Value)
		return false
	}
	if letS.Name.TokenLiteral() != name {
		t.Errorf("LetStatement.Name.TokenLiteral() does not '%s' got=%s", name, letS.Name.TokenLiteral())
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	cases := []struct {
		input   string
		wantVal interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foo;", "foo"},
	}

	for _, tt := range cases {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}
		s := program.Statements[0]
		returnS, ok := s.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("s not *ast.returnStatement. got=%T", s)
		}
		if returnS.TokenLiteral() != "return" {
			t.Fatalf("returnS.TokenLiteral not 'return', got %q", returnS.TokenLiteral())
		}
		if testLiteralExpression(t, returnS.ReturnValue, tt.wantVal) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	in := "foo;"
	l := lexer.New(in)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	ident, ok := s.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("s.Expression is not ast.Identifier. got=%T", s.Expression)
	}
	if ident.Value != "foo" {
		t.Fatalf("ident.Value is not %q. got=%q", "foo", ident.Value)
	}
	if ident.TokenLiteral() != "foo" {
		t.Fatalf("ident.Value is not %q. got=%q", "foo", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	in := "5;"
	l := lexer.New(in)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	literal, ok := s.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("s.Expression is not ast.IntegerLiteral. got=%T", s.Expression)
	}
	if literal.Value != 5 {
		t.Fatalf("literal.Value is not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Fatalf("literal.Value is not %q. got=%q", "5", literal.TokenLiteral())
	}
}

func TestBooleanExpression(t *testing.T) {
	cases := []struct {
		input    string
		wantBool bool
	}{
		{"true;", true},
		{"false;", false},
	}
	for _, tt := range cases {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}
		s, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		b, ok := s.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("s.Expression is not ast.Boolean. got=%T", s.Expression)
		}
		if b.Value != tt.wantBool {
			t.Fatalf("b.Value is not %t. got=%t", tt.wantBool, b.Value)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	is, ok := s.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("s.Expression is not ast.IfExpression. got=%T", s.Expression)
	}
	if !testInfixExpression(t, is.Condition, "x", "<", "y") {
		return
	}
	if len(is.Consequence.Statements) != 1 {
		t.Fatalf("consequence is not 1 statement. got=%d", len(is.Consequence.Statements))
	}
	con, ok := is.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", is.Consequence.Statements[0])
	}
	if !testIdentifier(t, con.Expression, "x") {
		return
	}
	if is.Alternative != nil {
		t.Fatalf("is.Alternative.Statements was not nil. got=%+v", is.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	is, ok := s.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("s.Expression is not ast.IfExpression. got=%T", s.Expression)
	}
	if !testInfixExpression(t, is.Condition, "x", "<", "y") {
		return
	}
	if len(is.Consequence.Statements) != 1 {
		t.Fatalf("consequence is not 1 statement. got=%d", len(is.Consequence.Statements))
	}
	con, ok := is.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", is.Consequence.Statements[0])
	}
	if !testIdentifier(t, con.Expression, "x") {
		return
	}
	if len(is.Alternative.Statements) != 1 {
		t.Fatalf("alternative is not 1 statement. got=%d", len(is.Alternative.Statements))
	}
	alt, ok := is.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", is.Alternative.Statements[0])
	}
	if !testIdentifier(t, alt.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	fl, ok := s.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("s.Expression is not ast.FunctionLiteral. got=%T", s.Expression)
	}
	if len(fl.Parameters) != 2 {
		t.Fatalf("parameters length is wrong. got=%d", len(fl.Parameters))
	}
	if !testLiteralExpression(t, fl.Parameters[0], "x") {
		return
	}
	if !testLiteralExpression(t, fl.Parameters[1], "y") {
		return
	}
	if len(fl.Body.Statements) != 1 {
		t.Fatalf("Body.Statements length is wrong. got=%d", len(fl.Body.Statements))
	}
	body, ok := fl.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", fl.Body.Statements[0])
	}
	testInfixExpression(t, body.Expression, "x", "+", "y")
}

func TestFunctionParametersParsing(t *testing.T) {
	cases := []struct {
		input      string
		wantParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) { x };", []string{"x"}},
		{"fn(x, y) { x + y };", []string{"x", "y"}},
	}
	for _, tt := range cases {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		s := program.Statements[0].(*ast.ExpressionStatement)
		fn := s.Expression.(*ast.FunctionLiteral)
		if len(fn.Parameters) != len(tt.wantParams) {
			t.Fatalf("parameter length is wrong. want=%d got=%d", len(tt.wantParams), len(fn.Parameters))
		}
		for i, p := range tt.wantParams {
			testLiteralExpression(t, fn.Parameters[i], p)
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	cases := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-5;", "-", 5},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range cases {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not contain 1 statements. got=%d", len(program.Statements))
		}
		s, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		exp, ok := s.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("s.Expression is not ast.PrefixExpression. got=%T", s.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %q. got=%q", tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	cases := []struct {
		input    string
		leftVal  interface{}
		operator string
		rightVal interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
	}

	for _, tt := range cases {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not contain 1 statements. got=%d", len(program.Statements))
		}
		s, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		if !testInfixExpression(t, s.Expression, tt.leftVal, tt.operator, tt.rightVal) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	cases := []struct {
		input string
		want  string
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
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"5 < 4 == 3 > 4", "((5 < 4) == (3 > 4))"},
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
		{"add(a + b + c * d / e + f)", "add((((a + b) + ((c * d) / e)) + f))"},
	}

	for _, tt := range cases {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		got := program.String()
		if got != tt.want {
			t.Errorf("want=%q, got=%q", tt.want, got)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}
	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("s is not  ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok := s.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("exp is not  ast.CallExpression. got=%T", s.Expression)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}
	if !testLiteralExpression(t, exp.Arguments[0], 1) {
		return
	}
	if !testInfixExpression(t, exp.Arguments[1], 2, "*", 3) {
		return
	}
	if !testInfixExpression(t, exp.Arguments[2], 4, "+", 5) {
		return
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, wantVal int64) bool {
	i, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if i.Value != wantVal {
		t.Errorf("i.Value is not %d. got=%d", wantVal, i.Value)
		return false
	}
	if i.TokenLiteral() != strconv.Itoa(int(wantVal)) {
		t.Errorf("i.TokenLiteral is not %q. got=%q", strconv.Itoa(int(wantVal)), i.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, wantVal bool) bool {
	b, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp is not *ast.IntegerLiteral. got=%T", exp)
		return false
	}
	if b.Value != wantVal {
		t.Errorf("b.Value is not %t. got=%t", wantVal, b.Value)
		return false
	}
	if b.TokenLiteral() != fmt.Sprintf("%t", wantVal) {
		t.Errorf("b.TokenLiteral is not %t. got=%s", wantVal, b.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, want interface{}) bool {
	switch v := want.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, int64(v))
	case bool:
		return testBooleanLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp is not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, op string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != op {
		t.Errorf("exp.Operator is not %q. got=%q", op, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
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
