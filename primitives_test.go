package main

import (
	"reflect"
	"testing"
)

func TestPrimitiveAdd(t *testing.T) {
	cases := []struct {
		src     string
		want    numberValue
		wantErr error
	}{
		{
			src:  `()`,
			want: numberValue{0},
		},
		{
			src:  `(1)`,
			want: numberValue{1},
		},
		{
			src:  `(1 2)`,
			want: numberValue{3},
		},
		{
			src:  `(1 2 3)`,
			want: numberValue{6},
		},
		{
			src:  `((primitive + 1 2) 3)`,
			want: numberValue{6},
		},
		{
			src:     `(1 #t)`,
			wantErr: errInvalidArgumentType,
		},
	}

	for i, c := range cases {
		t.Logf("Case %d: %s", i, c.src)

		exprs, err := parse(tokenize(c.src))
		if err != nil {
			t.Fatal("parse error", err)
		}

		got, gotErr := primitiveAdd(mustExpressionChildren(exprs[0]), newFrame())

		if gotErr == nil {
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got:  %v\nwant: %v", got, c.want)
			}
		}

		if gotErr != c.wantErr {
			t.Errorf("error:\ngot:  %v\nwant: %v", gotErr, c.wantErr)
		}
	}
}

func TestPrimitiveSubtract(t *testing.T) {
	cases := []struct {
		src     string
		want    numberValue
		wantErr error
	}{
		{
			src:  `()`,
			want: numberValue{0},
		},
		{
			src:  `(1)`,
			want: numberValue{-1},
		},
		{
			src:  `(1 2)`,
			want: numberValue{-1},
		},
		{
			src:  `(1 2 3)`,
			want: numberValue{-4},
		},
		{
			src:  `((primitive - 1 2) 3)`,
			want: numberValue{-4},
		},
		{
			src:     `(1 #t)`,
			wantErr: errInvalidArgumentType,
		},
	}

	for i, c := range cases {
		t.Logf("Case %d: %s", i, c.src)

		exprs, err := parse(tokenize(c.src))
		if err != nil {
			t.Fatal("parse error", err)
		}

		got, gotErr := primitiveSubtract(mustExpressionChildren(exprs[0]), newFrame())

		if gotErr == nil {
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got:  %v\nwant: %v", got, c.want)
			}
		}

		if gotErr != c.wantErr {
			t.Errorf("error:\ngot:  %v\nwant: %v", gotErr, c.wantErr)
		}
	}
}

func TestPrimitiveMultiply(t *testing.T) {
	cases := []struct {
		src     string
		want    numberValue
		wantErr error
	}{
		{
			src:  `()`,
			want: numberValue{1},
		},
		{
			src:  `(1)`,
			want: numberValue{1},
		},
		{
			src:  `(1 2)`,
			want: numberValue{2},
		},
		{
			src:  `(1 2 3)`,
			want: numberValue{6},
		},
		{
			src:  `((primitive * 1 2) 3)`,
			want: numberValue{6},
		},
		{
			src:     `(1 #t)`,
			wantErr: errInvalidArgumentType,
		},
	}

	for i, c := range cases {
		t.Logf("Case %d: %s", i, c.src)

		exprs, err := parse(tokenize(c.src))
		if err != nil {
			t.Fatal("parse error", err)
		}

		got, gotErr := primitiveMultiply(mustExpressionChildren(exprs[0]), newFrame())

		if gotErr == nil {
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got:  %v\nwant: %v", got, c.want)
			}
		}

		if gotErr != c.wantErr {
			t.Errorf("error:\ngot:  %v\nwant: %v", gotErr, c.wantErr)
		}
	}
}

func TestPrimitiveDivide(t *testing.T) {
	cases := []struct {
		src     string
		want    numberValue
		wantErr error
	}{
		{
			src:  `()`,
			want: numberValue{1},
		},
		{
			src:  `(1)`,
			want: numberValue{1},
		},
		{
			src:  `(1 2)`,
			want: numberValue{0},
		},
		{
			src:  `(4 2)`,
			want: numberValue{2},
		},
		{
			src:  `(5 2)`,
			want: numberValue{2},
		},
		{
			src:  `(0 2)`,
			want: numberValue{0},
		},
		{
			src:  `(12 2 3)`,
			want: numberValue{2},
		},
		{
			src:  `((primitive / 12 2) 3)`,
			want: numberValue{2},
		},
		{
			src:     `(1 0)`,
			wantErr: errDivideByZero,
		},
		{
			src:     `(1 #t)`,
			wantErr: errInvalidArgumentType,
		},
	}

	for i, c := range cases {
		t.Logf("Case %d: %s", i, c.src)

		exprs, err := parse(tokenize(c.src))
		if err != nil {
			t.Fatal("parse error", err)
		}

		got, gotErr := primitiveDivide(mustExpressionChildren(exprs[0]), newFrame())

		if gotErr == nil {
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got:  %v\nwant: %v", got, c.want)
			}
		}

		if gotErr != c.wantErr {
			t.Errorf("error:\ngot:  %v\nwant: %v", gotErr, c.wantErr)
		}
	}
}

func TestPrimitiveEqual(t *testing.T) {
	cases := []struct {
		src     string
		want    value
		wantErr error
	}{
		{
			src:     `()`,
			wantErr: errWrongNumberOfArguments,
		},
		{
			src:     `(1)`,
			wantErr: errWrongNumberOfArguments,
		},
		{
			src:     `()`,
			wantErr: errWrongNumberOfArguments,
		},
		{
			src:     `(1 2 3)`,
			wantErr: errWrongNumberOfArguments,
		},
		{
			src:  `(1 1)`,
			want: boolValue{true},
		},
		{
			src:  `(1 2)`,
			want: boolValue{false},
		},
		{
			src:  `(1 (primitive / 2 2))`,
			want: boolValue{true},
		},
	}

	for i, c := range cases {
		t.Logf("Case %d: %s", i, c.src)

		exprs, err := parse(tokenize(c.src))
		if err != nil {
			t.Fatal("parse error", err)
		}

		got, gotErr := primitiveEquals(mustExpressionChildren(exprs[0]), newFrame())

		if gotErr == nil {
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got:  %v\nwant: %v", got, c.want)
			}
		}

		if gotErr != c.wantErr {
			t.Errorf("error:\ngot:  %v\nwant: %v", gotErr, c.wantErr)
		}
	}
}

func TestPrimitiveGreaterThan(t *testing.T) {
	cases := []struct {
		src     string
		want    value
		wantErr error
	}{
		{
			src:     `()`,
			wantErr: errWrongNumberOfArguments,
		},
		{
			src:     `(1)`,
			wantErr: errWrongNumberOfArguments,
		},
		{
			src:     `()`,
			wantErr: errWrongNumberOfArguments,
		},
		{
			src:     `(1 2 3)`,
			wantErr: errWrongNumberOfArguments,
		},
		{
			src:  `(1 1)`,
			want: boolValue{false},
		},
		{
			src:  `(1 2)`,
			want: boolValue{false},
		},
		{
			src:  `(2 1)`,
			want: boolValue{true},
		},
		{
			src:  `(2 (primitive - 2 1))`,
			want: boolValue{true},
		},
		{
			src:     `(#t #f)`,
			wantErr: errTypeNotOrderable,
		},
		{
			src:     `(1 #t)`,
			wantErr: errIncomparableValueTypes,
		},
	}

	for i, c := range cases {
		t.Logf("Case %d: %s", i, c.src)

		exprs, err := parse(tokenize(c.src))
		if err != nil {
			t.Fatal("parse error", err)
		}

		got, gotErr := primitiveGreaterThan(mustExpressionChildren(exprs[0]), newFrame())

		if gotErr == nil {
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got:  %v\nwant: %v", got, c.want)
			}
		}

		if gotErr != c.wantErr {
			t.Errorf("error:\ngot:  %v\nwant: %v", gotErr, c.wantErr)
		}
	}
}
