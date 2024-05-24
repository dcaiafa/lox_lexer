package loxtest

import (
	"errors"
	"fmt"
	gotoken "go/token"
	"strings"
)

type ErrLogger struct {
	Fset *gotoken.FileSet
	str  strings.Builder
}

func (l *ErrLogger) Errorf(pos gotoken.Pos, msg string, args ...any) {
	if l.str.Len() > 0 {
		l.str.WriteString("\n")
	}
	if pos.IsValid() {
		position := l.Fset.Position(pos)
		fmt.Fprint(&l.str, position.String()+": ")
	}
	fmt.Fprintf(&l.str, msg, args...)
}

func (l *ErrLogger) ParserError(err error) {
	var pos gotoken.Pos
	if err, ok := err.(interface{ Pos() Token }); ok {
		pos = err.Pos().Pos
	}
	l.Errorf(pos, "%v", err)
}

func (l *ErrLogger) Err() error {
	if l.str.Len() == 0 {
		return nil
	}
	return errors.New(l.str.String())
}
