package loxtest

import (
	gotoken "go/token"

	baselexer "github.com/dcaiafa/lox_lexer/simplelexer"
)

func Parse(fset *gotoken.FileSet, expr string) []Token {
	file := fset.AddFile("expr", -1, len(expr))

	var parser parser
	lex := baselexer.New(baselexer.Config{
		StateMachine: new(_LexerStateMachine),
		File:         file,
		Input:        []byte(expr),
	})

	_ = parser.parse(lex)
	return parser.result
}

type Token = baselexer.Token

type parser struct {
	lox
	result []Token
}

func (p *parser) on_S(toks []Token) any {
	p.result = toks
	return nil
}

func (p *parser) on_token(tok Token) Token {
	return tok
}

func (p *parser) on_token__err(err Error) Token {
	return err.Token
}
