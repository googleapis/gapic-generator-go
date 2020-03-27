package main

import (
	"testing"
)

func Test_format(t *testing.T) {
	in := []string{
		"gapic: foo bar baz",
		"bazel: foo bar baz",
		"gencli: foo bar baz",
		"samples: foo bar baz",
		"chore(deps): foo bar baz",
		"chore: foo bar baz",
		"unknown: foo bar baz",
		"malformed message",
	}

	want := `# gapic

* foo bar baz

# bazel

* foo bar baz

# gencli

* foo bar baz

# samples

* foo bar baz

# chores

* foo bar baz
* update dependencies (see history)

# other

* unknown: foo bar baz
* malformed message`

	got := format(in)
	if got != want {
		t.Errorf("format(%q)=%q, want %q", in, got, want)
	}
}
