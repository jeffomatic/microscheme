package main

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	cases := []struct {
		src  string
		want []string
	}{
		{
			src:  `foo`,
			want: []string{"foo"},
		},
		{
			src:  `foo bar`,
			want: []string{"foo", "bar"},
		},
		{
			src:  `(foo)`,
			want: []string{"(", "foo", ")"},
		},
		{
			src:  `(foo bar)`,
			want: []string{"(", "foo", "bar", ")"},
		},
		{
			src: `
				foo
				bar
				baz
			`,
			want: []string{"foo", "bar", "baz"},
		},
		{
			src: `
				(lambda (foo)
				  (+ foo foo))
			`,
			want: []string{"(", "lambda", "(", "foo", ")", "(", "+", "foo", "foo", ")", ")"},
		},
	}

	for _, c := range cases {
		got := tokenize(c.src)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s:\ngot:  %v\nwant: %v", c.src, got, c.want)
		}
	}
}
