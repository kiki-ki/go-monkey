package lexer_test

import (
	"testing"

	"github.com/kiki-ki/go-interpreter/token"
)

func TextNextToken(t *testing.T) {
	in := `=+(){},;`

	cases := []struct {
		wantType    token.TokenType
		wantLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(in)

	for i, tt := range cases {
		token := l.NextToken()
		if token.Type != tt.wantType {
			t.Fatalf("cases[%d]: token type wrong, want=%q, got=%q", i, tt.wantType, token.Type)
		}
		if token.Literal != tt.wantLiteral {
			t.Fatalf("cases[%d]: token literal wrong, want=%q, got=%q", i, tt.wantLiteral, token.Literal)
		}
	}
}
