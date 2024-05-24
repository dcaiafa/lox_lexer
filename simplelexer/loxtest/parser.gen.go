package loxtest

import (
	_i1 "github.com/dcaiafa/lox_lexer/simplelexer"
)

var _LHS = []int32{
	0, 1, 2, 2, 2, 3, 3, 4, 4,
}

var _TermCounts = []int32{
	1, 1, 1, 1, 1, 1, 0, 2, 1,
}

var _Actions = []int32{
	9, 18, 27, 36, 39, 48, 57, 60, 69, 8, 0, -6, 1, 1,
	2, 2, 3, 4, 8, 0, -4, 1, -4, 2, -4, 3, -4, 8,
	0, -2, 1, -2, 2, -2, 3, -2, 2, 0, 2147483647, 8, 0, -3,
	1, -3, 2, -3, 3, -3, 8, 0, -8, 1, -8, 2, -8, 3,
	-8, 2, 0, -1, 8, 0, -5, 1, 1, 2, 2, 3, 4, 8,
	0, -7, 1, -7, 2, -7, 3, -7,
}

var _Goto = []int32{
	9, 18, 18, 18, 18, 18, 18, 19, 18, 8, 1, 3, 2, 5,
	3, 6, 4, 7, 0, 2, 2, 8,
}

type _Bounds struct {
	Begin Token
	End   Token
	Empty bool
}

func _cast[T any](v any) T {
	cv, _ := v.(T)
	return cv
}

type Error struct {
	Token Token
}

func _Find(table []int32, y, x int32) (int32, bool) {
	i := int(table[int(y)])
	count := int(table[i])
	i++
	end := i + count
	for ; i < end; i += 2 {
		if table[i] == x {
			return table[i+1], true
		}
	}
	return 0, false
}

type _Lexer interface {
	ReadToken() (Token, int)
}

type lox struct {
	_lex   _Lexer
	_state _Stack[int32]
	_sym   _Stack[any]

	_lookaheadSym any
	_lookahead    int

	_queuedLookahead    int
	_queuedLookaheadSym any
}

func (p *parser) parse(lex _Lexer) bool {
	const accept = 0x7FFFFFFF

	p._lex = lex

	p._state.Push(0)
	p._ReadToken()

	for {
		topState := p._state.Peek(0)
		action, ok := _Find(_Actions, topState, int32(p._lookahead))
		if !ok {
			if !p._Recover() {
				return false
			}
		}
		if action == accept {
			break
		} else if action >= 0 { // shift
			p._state.Push(action)
			p._sym.Push(p._lookaheadSym)
			p._ReadToken()
		} else { // reduce
			prod := -action
			termCount := _TermCounts[int(prod)]
			rule := _LHS[int(prod)]
			res := p._Act(prod)
			p._state.Pop(int(termCount))
			p._sym.Pop(int(termCount))
			topState = p._state.Peek(0)
			nextState, _ := _Find(_Goto, topState, rule)
			p._state.Push(nextState)
			p._sym.Push(res)
		}
	}

	return true
}

func (p *parser) _ReadToken() {
	p._lookaheadSym, p._lookahead = p._lex.ReadToken()
}

func (p *parser) _Recover() bool {
	errSym := p._lookaheadSym

	for {
		for p._lookahead == ERROR {
			p._ReadToken()
		}

		saveState := p._state
		saveSym := p._sym

		for len(p._state) > 1 {
			// Is ERROR a valid lookahead at the top state?
			state := p._state.Peek(0)
			action, ok := _Find(_Actions, state, int32(ERROR))

			// If so, check whether, after processing all reductions, we can reach a
			// state where shifting the current _lookahead is a valid action.
			action2 := action
			for ok && action2 <= 0 {
				prod := -action2
				rule := _LHS[int(prod)]
				state, _ = _Find(_Goto, state, rule)

				action2, ok = _Find(_Actions, state, int32(ERROR))
				if !ok {
					break
				} else if action2 >= 0 {
					_, ok = _Find(_Actions, state, int32(p._lookahead))
					break
				}
			}

			if ok {
				p._queuedLookahead = p._lookahead
				p._queuedLookaheadSym = p._lookaheadSym
				p._lookahead = ERROR
				p._lookaheadSym = err
				return true
			}

			p._state.Pop(1)
			p._sym.Pop(1)
		}

		if p._lookahead == EOF {
			return false
		}

		p._state = saveState
		p._sym = saveSym
		p._ReadToken()
	}
}

func (p *parser) _Act(prod int32) any {
	switch prod {
	case 1:
		return p.on_S(
			_cast[[]_i1.Token](p._sym.Peek(0)),
		)
	case 2:
		return p.on_token(
			_cast[_i1.Token](p._sym.Peek(0)),
		)
	case 3:
		return p.on_token(
			_cast[_i1.Token](p._sym.Peek(0)),
		)
	case 4:
		return p.on_token__err(
			_cast[error](p._sym.Peek(0)),
		)
	case 5: // ZeroOrMore
		return _cast[[]_i1.Token](p._sym.Peek(0))
	case 6: // ZeroOrMore
		{
			var zero []_i1.Token
			return zero
		}
	case 7: // OneOrMore
		return append(
			_cast[[]_i1.Token](p._sym.Peek(1)),
			_cast[_i1.Token](p._sym.Peek(0)),
		)
	case 8: // OneOrMore
		return []_i1.Token{
			_cast[_i1.Token](p._sym.Peek(0)),
		}
	default:
		panic("unreachable")
	}
}
