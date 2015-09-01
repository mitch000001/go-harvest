package main

import "testing"

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			"AbcDef",
			"abc_def",
		},
		{
			"abcDef",
			"abc_def",
		},
	}
	for _, test := range tests {
		res := SnakeCase(test.in)

		if test.out != res {
			t.Logf("Expected SnakeCase to return\n%q\n\tgot:\n%q\n", test.out, res)
			t.Fail()
		}
	}
}
