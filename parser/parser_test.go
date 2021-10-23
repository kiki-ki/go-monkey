package parser_test

import (
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

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral is not 'let' got=%q", s.TokenLiteral())
		return false
	}
	letS, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s is not *ast.LetStatement got=%q", s)
		return false
	}
	if letS.Name.Value != name {
		t.Errorf("LetStatement.Name.Value is not '%s' got=%s", name, letS.Name.Value)
		return false
	}
	if letS.Name.TokenLiteral() != name {
		t.Errorf("LetStatement.Name.TokenLiteral() is not '%s' got=%s", name, letS.Name.TokenLiteral())
		return false
	}
	return true
}
