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
			tokens, err := Parse(fset, input)
			for _, tk := range tokens {
				position := fset.Position(tk.Pos)
				fmt.Fprintf(
					res,
					"%v [%v] %v:%v\n",
					_TokenToString(tk.Type),
					string(tk.Str),
					position.Line, position.Column)
			}
			if err != nil {
				fmt.Fprintln(res, "Error:")
				fmt.Fprintln(res, err.Error())
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
		`, ``)
	/*
			run(`
		1 "\x1" 1
		1
		`, ``)
	*/
}
