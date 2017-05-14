package main

import "testing"

func TestClassify(t *testing.T) {
	cases := []struct {
		src     string
		want    expressionType
		wantErr error
	}{
		{
			src:  `null`,
			want: exprNull,
		},
		{
			src:  `()`,
			want: exprNull,
		},
		{
			src:  `100`,
			want: exprNumber,
		},
		{
			src:  `-100`,
			want: exprNumber,
		},
		{
			src:  `#t`,
			want: exprBoolean,
		},
		{
			src:  `#f`,
			want: exprBoolean,
		},
		{
			src:  `"foo bar"`,
			want: exprString,
		},
		{
			src:  `foo`,
			want: exprDereference,
		},
		{
			src:  `(define a b)`,
			want: exprDefine,
		},
		{
			src:  `(define a (b))`,
			want: exprDefine,
		},
		{
			src:  `(define (a) b)`,
			want: exprDefine,
		},
		{
			src:  `(define (a) (b))`,
			want: exprDefine,
		},
		{
			src:  `(define (a b) c)`,
			want: exprDefine,
		},
		{
			src:  `(define (a b c) d)`,
			want: exprDefine,
		},
		{
			src:  `(define (a b) c d)`,
			want: exprDefine,
		},
		{
			src:     `(define)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(define a)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(define a b c)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(define () a)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(define ((a)) b)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(define (a (b)) c)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:  `(begin)`,
			want: exprBegin,
		},
		{
			src:  `(begin a b (c))`,
			want: exprBegin,
		},
		{
			src:  `(if a b c)`,
			want: exprIf,
		},
		{
			src:  `(if a b (c))`,
			want: exprIf,
		},
		{
			src:  `(if a b c)`,
			want: exprIf,
		},
		{
			src:     `(if)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(if a)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(if a b)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(if a b c d)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:  `(lambda () c)`,
			want: exprLambda,
		},
		{
			src:  `(lambda (a b) c)`,
			want: exprLambda,
		},
		{
			src:  `(lambda (a) (c))`,
			want: exprLambda,
		},
		{
			src:  `(lambda (a) b c d)`,
			want: exprLambda,
		},
		{
			src:     `(lambda)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(lambda (a))`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(lambda a b)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:  `(let () a)`,
			want: exprLet,
		},
		{
			src:  `(let ((a b)) c)`,
			want: exprLet,
		},
		{
			src:  `(let ((a b) (c d)) e)`,
			want: exprLet,
		},
		{
			src:  `(let ((a b)) c d e)`,
			want: exprLet,
		},
		{
			src:     `(let)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(let a)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(let (a) b)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(let (()) a)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:     `(let ((a)) b)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:  `(primitive +)`,
			want: exprPrimitive,
		},
		{
			src:     `(primitive)`,
			wantErr: errInvalidCompoundExpression,
		},
		{
			src:  `(a)`,
			want: exprApplication,
		},
		{
			src:  `(a b c)`,
			want: exprApplication,
		},
		{
			src:  `(())`,
			want: exprApplication,
		},
		{
			src:  `((a) b)`,
			want: exprApplication,
		},
	}

	for i, c := range cases {
		t.Logf("Case %d: %s", i, c.src)

		exprs, err := parse(tokenize(c.src))
		if err != nil {
			t.Errorf("%s: parse error: %v", c.src, err)
			continue
		}

		if len(exprs) != 1 {
			t.Errorf("%s: invalid parse result; should have one root expression: %v", c.src, exprs)
		}

		got, gotErr := classify(exprs[0])

		if got != c.want {
			t.Errorf("got:  %v\nwant: %v", got, c.want)
		}

		if gotErr != c.wantErr {
			t.Errorf("error:\ngot:  %v\nwant: %v", gotErr, c.wantErr)
		}
	}
}
