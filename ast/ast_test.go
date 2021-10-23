package ast_test

import (
	"testing"

	"github.com/kiki-ki/go-monkey/ast"
	"github.com/kiki-ki/go-monkey/token"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{
					Type: token.LET,
					Literal: "let",
				},
				Name: &ast.Identifier{
					Token: token.Token{
						Type: token.IDENT,
						Literal: "x",
					},
					Value: "x",
				},
				Value: &ast.Identifier{
					Token: token.Token{
						Type: token.IDENT,
						Literal: "y",
					},
					Value: "y",
				},
			},
		},
	}

	want := "let x = y;"
	if program.String() != want {
		t.Errorf("program.String() wrong. want=%q got=%q", want, program.String())
	}
}