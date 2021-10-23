package parser

import (
	"github.com/kiki-ki/go-monkey/ast"
	"github.com/kiki-ki/go-monkey/lexer"
	"github.com/kiki-ki/go-monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// cur, peekがセットされてる状態まで進めておく
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}