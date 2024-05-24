package loxtest

import (
	gotoken "go/token"

	baselexer "github.com/dcaiafa/lox_lexer/simplelexer"
)

func Parse(fset *gotoken.FileSet, expr string) ([]Token, error) {
	file := fset.AddFile("expr", -1, len(expr))
	errs := &ErrLogger{
		Fset: fset,
	}

	onError := func(l *baselexer.Lexer) {
		errs.Errorf(l.Pos(), "unexpected character: %c", l.Peek())
	}

	var parser parser
	parser.errLogger = errs
	lex := baselexer.New(baselexer.Config{
		StateMachine: new(_LexerStateMachine),
		OnError:      onError,
		File:         file,
		Input:        []byte(expr),
	})

	_ = parser.parse(lex)
	return parser.result, errs.Err()
}

type Token = baselexer.Token

type parser struct {
	lox
	errLogger *ErrLogger
	result    []Token
}

func (p *parser) on_S(toks []Token) any {
	p.result = toks
	return nil
}

func (p *parser) on_token(tok Token) Token {
	return tok
}

func (p *parser) on_token__err(_ error) Token {
	return Token{
		Type: ERROR,
	}
}

func (p *parser) _onError() {
	if p.errorToken().Type != ERROR {
		p.errLogger.Errorf(p.errorToken().Pos, "unexpected token %v", p.errorToken())
	}
}
