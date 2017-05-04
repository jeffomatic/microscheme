package main

import (
	"reflect"
	"testing"
)

func TestEval(t *testing.T) {
	env := newFrame()
	env.set("testVar", numberValue{1})
	env.set("testProc", &procValue{
		params: []string{"x"},
		body:   []expression{&tokenExpression{"x"}},
	})

	cases := []struct {
		src     string
		want    value
		wantErr error
	}{
		{
			src:  `null`,
			want: nullValue{},
		},
		{
			src:  `1`,
			want: numberValue{1},
		},
		{
			src:  `#t`,
			want: boolValue{true},
		},
		{
			src:  `#f`,
			want: boolValue{false},
		},
		{
			src:  `testVar`,
			want: numberValue{1},
		},
		{
			src:  `(begin)`,
			want: theNullValue,
		},
		{
			src:  `(begin 1 2)`,
			want: numberValue{2},
		},
		{
			src:  `(if #t 1 2)`,
			want: numberValue{1},
		},
		{
			src:  `(if #f 1 2)`,
			want: numberValue{2},
		},
		{
			src: `(lambda () null)`,
			want: &procValue{
				params: nil,
				body:   []expression{&tokenExpression{"null"}},
			},
		},
		{
			src: `(lambda (a b) a b)`,
			want: &procValue{
				params: []string{"a", "b"},
				body: []expression{
					&tokenExpression{"a"},
					&tokenExpression{"b"},
				},
			},
		},
		{
			src:  `(let ((a 1)) a)`,
			want: numberValue{1},
		},
		{
			src: `
(let ((a 1)
	    (b 2))
	a
	b)`,
			want: numberValue{2},
		},
		{
			src:  `(testProc #t)`,
			want: theTrueValue,
		},
	}

	for i, c := range cases {
		t.Logf("Case %d: %v", i, c.src)

		exprs, err := parse(tokenize(c.src))
		if err != nil {
			t.Fatal("parse error:", err)
		}

		if len(exprs) != 1 {
			t.Fatal("should be exactly one top-level expression: ", exprs)
		}

		got, gotErr := eval(exprs[0], env)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("value:\ngot:  %v\nwant: %v", got, c.want)
		}
		if gotErr != c.wantErr {
			t.Errorf("error:\ngot:  %v\nwant: %v", gotErr, c.wantErr)
		}
	}
}
