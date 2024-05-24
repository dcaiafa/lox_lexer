package baselexer

import (
	"bytes"
	gotoken "go/token"
)

const EOF = 0
const ERROR = 1

const (
	actionConsume  int = 0
	actionAccept   int = 1
	actionDiscard  int = 2
	actionTryAgain int = 3
	actionEOF      int = 4
)

// StateMachine is an interface that is satisfied by a state machine
// implementation generated by Lox.
type StateMachine interface {
	// PushRune pushes a rune to the state machine. It can return one of the
	// following:
	// 0 (consume): consume the rune and include it in the token.
	// 1 (accept): accept the current token. The token id is specified by Token().
	// 2 (discard): discard current token.
	// 3 (try-again): Call PushRune again with the same rune.
	// 4 (EOF): We are done.
	// -1 (error): That's an error.
	PushRune(r rune) int

	// Token returns the token-id of the recognized token when PushRune
	// returns 1 (accept).
	Token() int

	// Reset resets the state machine so it can be used again.
	Reset()
}

// Token is the type of tokens produced by Lexer.
type Token struct {
	Type int
	Str  []byte
	Pos  gotoken.Pos
}

// Config configures a Lexer.
type Config struct {
	// StateMachine is the Lox-generated state machine.
	StateMachine StateMachine

	// OnError is called by Lexer when an error is encountered before recovery is
	// applied. Peek() will contain the current rune.
	OnError func(l *Lexer)

	// File is a go/token.File which is used to efficiently tag tokens with file
	// positions.
	File *gotoken.File

	// Input is input to be parsed by the Lexer.
	Input []byte
}

// Lexer is a simple lexer implementation to be used with lexer state machines
// produced by Lox. It produces tokens containing byte slices of the original
// input which may not be appropriate for very large or stream inputs.
type Lexer struct {
	sm          StateMachine
	onError     func(l *Lexer)
	file        *gotoken.File
	input       []byte
	inputReader *bytes.Reader
	offset      int
	charLen     int
	char        rune
	pos         gotoken.Pos
}

// New creates a Lexer using a Lox-generated state machine.
func New(cfg Config) *Lexer {
	l := &Lexer{
		sm:          cfg.StateMachine,
		onError:     cfg.OnError,
		file:        cfg.File,
		input:       cfg.Input,
		inputReader: bytes.NewReader(cfg.Input),
	}
	// Read the first rune.
	l.consume()
	return l
}

// Pos returns the position of the next rune to be consumed by the Lexer.
func (l *Lexer) Pos() gotoken.Pos {
	return l.pos
}

// Peek returns the next input rune to be consumed by the Lexer.
func (l *Lexer) Peek() rune {
	return l.char
}

// ReadToken parses a token from the input. It returns both the token value and
// the token id as required by the generated Lox _Lexer interface. Return token
// id EOF after the end of the input has been reached.
func (l *Lexer) ReadToken() (Token, int) {
	start := -1

	for {
		if start == -1 {
			start = l.offset
			l.pos = l.file.Pos(l.offset)
		}

		action := l.sm.PushRune(l.char)

		switch action {
		case actionConsume:
			l.consume()

		case actionAccept:
			end := l.offset
			t := Token{
				Type: l.sm.Token(),
				Str:  l.input[start:end],
				Pos:  l.pos,
			}
			start = -1
			return t, t.Type

		case actionDiscard:
			start = -1

		case actionTryAgain:
			// fallthrough

		case actionEOF:
			t := Token{
				Type: EOF,
				Pos:  l.pos,
			}
			return t, t.Type

		default: // Error
			t := Token{
				Type: ERROR,
				Pos:  l.pos,
			}

			l.onError(l)

			// Read until the beginning of the next line.
			// TODO: allow custom recovery.
			for l.char != '\n' && l.char != 0 {
				l.consume()
			}
			l.consume()

			l.sm.Reset()

			return t, t.Type
		}
	}
}

func (l *Lexer) consume() {
	if l.char == '\n' {
		l.file.AddLine(l.offset + 1)
	}
	l.offset += l.charLen
	r, charLen, err := l.inputReader.ReadRune()
	if err != nil {
		l.char = 0
		l.charLen = 0
		return
	}
	l.char = r
	l.charLen = charLen
}