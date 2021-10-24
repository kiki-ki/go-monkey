package parser_test

import (
	"strconv"
	"testing"

	"github.com/kiki-ki/go-monkey/ast"
	"github.com/kiki-ki/go-monkey/lexer"
	"github.com/kiki-ki/go-monkey/parser"
)

func TestLetStatements(t *testing.T) {
	in := `
	let x = 5;
	let y = 10;
	let foo = 898989;
	`

	l := lexer.New(in)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	cases := []struct {
		wantIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range cases {
		s := program.Statements[i]
		if !testLetStatement(t, s, tt.wantIdentifier) {
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
	in := `
	return 5;
	return 10;
	return 898989;
	`

	l := lexer.New(in)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, s := range program.Statements {
		returnS, ok := s.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("s does not *ast.ReturnStatement. got=%q", s)
			continue
		}
		if returnS.TokenLiteral() != "return" {
			t.Fatalf("returnS.TokenLiteral() does not 'return'. got=%s", returnS.TokenLiteral())
			continue
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

func TestParsingPrefixExpressions(t *testing.T) {
	cases := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-5;", "-", 5},
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
		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	cases := []struct {
		input    string
		leftVal  int64
		operator string
		rightVal int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
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
