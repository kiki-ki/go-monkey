package lexer_test

import (
	"testing"

	"github.com/kiki-ki/go-monkey/lexer"
	"github.com/kiki-ki/go-monkey/token"
)

func TestNextToken(t *testing.T) {
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

	l := lexer.New(in)

	for i, tt := range cases {
		tok := l.NextToken()
		if tok.Type != tt.wantType {
			t.Fatalf("cases[%d]: token type wrong, want=%q, got=%q", i, tt.wantType, tok.Type)
		}
		if tok.Literal != tt.wantLiteral {
			t.Fatalf("cases[%d]: token literal wrong, want=%q, got=%q", i, tt.wantLiteral, tok.Literal)
		}
	}
}
