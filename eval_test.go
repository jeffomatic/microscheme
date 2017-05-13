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
			want: nullValue{},
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
				env:    env,
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
				env: env,
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
					b)
			`,
			want: numberValue{2},
		},
		{
			src:  `(primitive + 1 2)`,
			want: numberValue{3},
		},
		{
			src:  `(testProc #t)`,
			want: boolValue{true},
		},
		{
			src: `
				(((lambda (x)
					(lambda (y)
						(primitive + x y)))
				  2)
				 1)
			`,
			want: numberValue{3},
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

// define tests are separate because they modify the environment.
func TestEvalDefine(t *testing.T) {
	env := newFrame()
	env.set("testVar", numberValue{1})
	env.set("testProc", &procValue{
		params: []string{"x"},
		body:   []expression{&tokenExpression{"x"}},
	})

	cases := []struct {
		src       string
		wantBound map[string]value
		wantErr   error
	}{
		{
			src: `(define a 1)`,
			wantBound: map[string]value{
				"a": numberValue{1},
			},
		},
		{
			src: `(define a testVar)`,
			wantBound: map[string]value{
				"a": numberValue{1},
			},
		},
		{
			src: `(define a testVar)`,
			wantBound: map[string]value{
				"a": numberValue{1},
			},
		},
		{
			src: `(define a (testProc 1))`,
			wantBound: map[string]value{
				"a": numberValue{1},
			},
		},
		{
			src: `(define a (lambda (x) x))`,
			wantBound: map[string]value{
				"a": &procValue{
					params: []string{"x"},
					body:   []expression{&tokenExpression{"x"}},
					// env will be set by test harness
				},
			},
		},
		{
			src: `(define (a x) x)`,
			wantBound: map[string]value{
				"a": &procValue{
					params: []string{"x"},
					body:   []expression{&tokenExpression{"x"}},
					// env will be set by test harness
				},
			},
		},
		{
			src: `(define (a) 1)`,
			wantBound: map[string]value{
				"a": &procValue{
					params: nil,
					body:   []expression{&tokenExpression{"1"}},
					// env will be set by test harness
				},
			},
		},
		{
			src: `(define (a x y) (+ x y))`,
			wantBound: map[string]value{
				"a": &procValue{
					params: []string{"x", "y"},
					body: []expression{
						&compoundExpression{
							children: []expression{
								&tokenExpression{"+"},
								&tokenExpression{"x"},
								&tokenExpression{"y"},
							},
						},
					},
					// env will be set by test harness
				},
			},
		},
		{
			src: `(define (a x y) x y)`,
			wantBound: map[string]value{
				"a": &procValue{
					params: []string{"x", "y"},
					body:   []expression{&tokenExpression{"x"}, &tokenExpression{"y"}},
					// env will be set by test harness
				},
			},
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

		caseEnv := env.extend()
		got, gotErr := eval(exprs[0], caseEnv)
		if !reflect.DeepEqual(got, nullValue{}) {
			t.Errorf("value:\ngot:  %v\nwant: %v", got, nullValue{})
		}
		if gotErr != c.wantErr {
			t.Errorf("error:\ngot:  %v\nwant: %v", gotErr, c.wantErr)
		}

		if len(c.wantBound) > 0 {
			for k, v := range c.wantBound {
				// Proc bindings will expect the per-case environment.
				if proc, ok := v.(*procValue); ok {
					proc.env = caseEnv
				}

				got, err := caseEnv.get(k)
				if err != nil {
					t.Errorf("binding %s: not found", k)
				}

				if !reflect.DeepEqual(got, v) {
					t.Errorf("binding:\ngot:  %v\nwant: %v", got, v)
				}
			}
		}
	}
}
