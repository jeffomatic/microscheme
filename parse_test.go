package main

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		src     string
		want    []expression
		wantErr error
	}{
		{
			src:  `foo`,
			want: []expression{&tokenExpression{"foo"}},
		},
		{
			src:  `foo bar`,
			want: []expression{&tokenExpression{"foo"}, &tokenExpression{"bar"}},
		},
		{
			src: `(foo)`,
			want: []expression{&compoundExpression{
				children: []expression{&tokenExpression{"foo"}},
			}},
		},
		{
			src: `(foo bar)`,
			want: []expression{&compoundExpression{
				children: []expression{
					&tokenExpression{"foo"},
					&tokenExpression{"bar"},
				},
			}},
		},
		{
			src: `
				(lambda (foo)
				  (+ foo foo))
			`,
			want: []expression{&compoundExpression{
				children: []expression{
					&tokenExpression{"lambda"},
					&compoundExpression{
						children: []expression{
							&tokenExpression{"foo"},
						},
					},
					&compoundExpression{
						children: []expression{
							&tokenExpression{"+"},
							&tokenExpression{"foo"},
							&tokenExpression{"foo"},
						},
					},
				},
			}},
		},
		{
			src:     `(foo))`,
			wantErr: errInvalidClosingBrace,
		},
		{
			src:     `(foo`,
			wantErr: errUnclosedExpression,
		},
	}

	for _, c := range cases {
		got, gotErr := parse(tokenize(c.src))

		if gotErr != nil || c.wantErr != nil {
			if gotErr != c.wantErr {
				t.Errorf("%s:\ngot error:  %v\nwant error: %v", c.src, gotErr, c.wantErr)
			}
		} else if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s:\ngot:  %v\nwant: %v", c.src, got, c.want)
		}
	}
}
