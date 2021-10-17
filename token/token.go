package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func New(tType TokenType, ch byte) Token {
	return Token{Type: tType, Literal: string(ch)}
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 識別子
	IDENT = "IDENT"

	// リテラル
	INT = "INT"

	// 演算子
	ASSIGN = "="
	PLUS   = "+"

	// デリミタ
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// キーワード
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
