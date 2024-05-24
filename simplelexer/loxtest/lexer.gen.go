package loxtest

var _lexerMode0 = []uint32{
	4, 21, 34, 40, 16, 5, 9, 10, 1, 13, 13, 1, 32, 32,
	1, 34, 34, 2, 48, 57, 3, 12, 3, 9, 10, 1, 13, 13,
	1, 32, 32, 1, 4, 0, 5, 0, 1, 1, 5, 0, 6, 1,
	48, 57, 3, 3, 2,
}

var _lexerMode1 = []uint32{
	22, 39, 43, 39, 66, 39, 77, 39, 88, 99, 39, 110, 121, 132,
	143, 149, 160, 171, 182, 193, 204, 215, 16, 5, 32, 33, 10, 34,
	34, 14, 35, 91, 10, 92, 92, 2, 93, 1114111, 10, 3, 0, 5,
	0, 22, 7, 34, 34, 7, 85, 85, 4, 110, 110, 7, 114, 114,
	7, 116, 116, 7, 117, 117, 19, 120, 120, 17, 10, 3, 48, 57,
	21, 65, 70, 21, 97, 102, 21, 10, 3, 48, 57, 1, 65, 70,
	1, 97, 102, 1, 10, 3, 48, 57, 6, 65, 70, 6, 97, 102,
	6, 10, 3, 48, 57, 3, 65, 70, 3, 97, 102, 3, 10, 3,
	48, 57, 8, 65, 70, 8, 97, 102, 8, 10, 3, 48, 57, 9,
	65, 70, 9, 97, 102, 9, 10, 3, 48, 57, 5, 65, 70, 5,
	97, 102, 5, 5, 0, 2, 0, 3, 3, 10, 3, 48, 57, 11,
	65, 70, 11, 97, 102, 11, 10, 3, 48, 57, 12, 65, 70, 12,
	97, 102, 12, 10, 3, 48, 57, 13, 65, 70, 13, 97, 102, 13,
	10, 3, 48, 57, 15, 65, 70, 15, 97, 102, 15, 10, 3, 48,
	57, 16, 65, 70, 16, 97, 102, 16, 10, 3, 48, 57, 18, 65,
	70, 18, 97, 102, 18, 10, 3, 48, 57, 20, 65, 70, 20, 97,
	102, 20,
}

var _lexerModes = [][]uint32{

	_lexerMode0,

	_lexerMode1,
}

const (
	_lexerConsume  = 0
	_lexerAccept   = 1
	_lexerDiscard  = 2
	_lexerTryAgain = 3
	_lexerEOF      = 4
	_lexerError    = -1
)

type _LexerStateMachine struct {
	token     int
	state     int
	mode      []uint32
	modeStack _Stack[[]uint32]
}

func (l *_LexerStateMachine) PushRune(r rune) int {
	if l.mode == nil {
		l.mode = _lexerMode0
	}

	mode := l.mode

	// Find the table row corresponding to state.
	i := int(mode[int(l.state)])
	count := int(mode[i])
	i++
	end := i + count

	// The format of the row is as follows:
	//
	//   gotoCount uint32
	//   [gotoCount]struct{
	//     rangeBegin uint32
	//     rangeEnd   uint32
	//     gotoState  uint32
	//   }
	//   [actionCount]struct {
	//     actionType  uint32
	//     actionParam uint32
	//   }
	//
	// Where 'actionCount' is determined by the amount of uint32 left in the row.

	gotoN := int(mode[i])
	i++

	// Use binary-search to find the next state.
	b := 0
	e := gotoN
	for b < e {
		j := b + (e-b)/2
		k := i + j*3
		switch {
		case r >= rune(mode[k]) && r <= rune(mode[k+1]):
			l.state = int(mode[k+2])
			return _lexerConsume
		case r < rune(mode[k]):
			e = j
		case r > rune(mode[k+1]):
			b = j + 1
		default:
			panic("not reached")
		}
	}

	// Move 'i' to the beginning of the actions section.
	i += gotoN * 3

	for ; i < end; i += 2 {
		switch mode[i] {
		case 1: // PushMode
			modeIndex := int(mode[i+1])
			l.modeStack.Push(mode)
			l.mode = _lexerModes[modeIndex]
		case 2: // PopMode
			l.mode = l.modeStack.Peek(0)
			l.modeStack.Pop(1)
		case 3: // Accept
			l.token = int(mode[i+1])
			l.state = 0
			return _lexerAccept
		case 4: // Discard
			l.state = 0
			return _lexerDiscard
		case 5: // Accum
			l.state = 0
			return _lexerTryAgain
		}
	}

	if l.state == 0 && r == 0 {
		return _lexerEOF
	}

	return _lexerError
}

func (l *_LexerStateMachine) Reset() {
	l.mode = nil
	l.state = 0
}

func (l *_LexerStateMachine) Token() int {
	return l.token
}
