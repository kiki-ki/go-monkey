package lexer

import "github.com/kiki-ki/go-monkey/token"

// 字句解析器

type Lexer struct {
	input        string
	position     int  // 現在の位置
	readPosition int  // 次の文字の位置
	ch           byte // 現在検査してる文字
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpace()

	switch l.ch {
	case ';':
		tok = token.New(token.SEMICOLON, l.ch)
	case '(':
		tok = token.New(token.LPAREN, l.ch)
	case ')':
		tok = token.New(token.RPAREN, l.ch)
	case '{':
		tok = token.New(token.LBRACE, l.ch)
	case '}':
		tok = token.New(token.RBRACE, l.ch)
	case ',':
		tok = token.New(token.COMMA, l.ch)
	case '=':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(token.EQ)
		} else {
			tok = token.New(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(token.NOT_EQ)
		} else {
			tok = token.New(token.BANG, l.ch)
		}
	case '+':
		tok = token.New(token.PLUS, l.ch)
	case '-':
		tok = token.New(token.MINUS, l.ch)
	case '/':
		tok = token.New(token.SLASH, l.ch)
	case '*':
		tok = token.New(token.ASTERISK, l.ch)
	case '<':
		tok = token.New(token.LT, l.ch)
	case '>':
		tok = token.New(token.GT, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookUpIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = token.New(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) makeTwoCharToken(tt token.TokenType) token.Token {
	ch := l.ch
	l.readChar()
	literal := string(ch) + string(l.ch)
	return token.Token{Type: tt, Literal: literal}
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// 識別子,キーワードに利用可能な文字か判断
func isLetter(ch byte) bool {
	return 'a' <= ch && 'z' >= ch || 'A' <= ch && 'Z' >= ch || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
