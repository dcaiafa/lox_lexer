package loxtest

import (
	"fmt"
	gotoken "go/token"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	run := func(input string, output string) {
		t.Run("test", func(t *testing.T) {
			t.Helper()
			fset := gotoken.NewFileSet()
			res := new(strings.Builder)
			tokens := Parse(fset, input)
			for _, tk := range tokens {
				position := fset.Position(tk.Pos)
				fmt.Fprintf(
					res,
					"%v [%v] %v:%v\n",
					_TokenToString(tk.Type),
					string(tk.Str),
					position.Line, position.Column)
			}

			output = strings.TrimSpace(output)
			resStr := strings.TrimSpace(res.String())

			if output != resStr {
				t.Log("Input:\n", input)
				t.Log("Expected ouput:\n", output)
				t.Log("Actual output:\n", resStr)
				t.Fatal("Unexpected output")
			}
		})
	}

	run(`
123 "foo"
 "hello\nworld" 987 "The \"crazy\" bear!"
"this is good" "this is bad \x1" "this will be discarded"
"life goes on"
    `, `
NUM [123] 2:1
STR ["foo"] 2:5
STR ["hello\nworld"] 3:2
NUM [987] 3:17
STR ["The \"crazy\" bear!"] 3:21
STR ["this is good"] 4:1
ERROR [] 4:16
STR ["life goes on"] 5:1
		`)
	/*
			run(`
		1 "\x1" 1
		1
		`, ``)
	*/
}
