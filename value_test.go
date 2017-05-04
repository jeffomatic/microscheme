package main

import "testing"

func TestValueEqual(t *testing.T) {
	var testProc procValue

	cases := []struct {
		a    value
		b    value
		want bool
	}{
		{
			a:    numberValue{underlying: 1},
			b:    numberValue{underlying: 1},
			want: true,
		},
		{
			a:    numberValue{underlying: 1},
			b:    &numberValue{underlying: 1},
			want: true,
		},
		{
			a:    numberValue{underlying: 1},
			b:    numberValue{underlying: 2},
			want: false,
		},
		{
			a:    numberValue{underlying: 1},
			b:    &numberValue{underlying: 2},
			want: false,
		},
		{
			a:    boolValue{underlying: true},
			b:    boolValue{underlying: true},
			want: true,
		},
		{
			a:    boolValue{underlying: true},
			b:    &boolValue{underlying: true},
			want: true,
		},
		{
			a:    boolValue{underlying: true},
			b:    boolValue{underlying: false},
			want: false,
		},
		{
			a:    boolValue{underlying: true},
			b:    &boolValue{underlying: false},
			want: false,
		},
		{
			a:    &testProc,
			b:    &testProc,
			want: true,
		},
		{
			a:    new(procValue),
			b:    new(procValue),
			want: false,
		},
	}

	for i, c := range cases {
		t.Logf("Case %d: %v.equals(%v)", i, c.a, c.b)

		got, err := c.a.equals(c.b)
		if err != nil {
			t.Errorf("error:\ngot:  %v\nwant: %v", err, nil)
			continue
		}

		if got != c.want {
			t.Errorf("value:\ngot:  %v\nwant: %v", got, c.want)
		}
	}
}

func TestValueIncomparable(t *testing.T) {
	vals := []value{
		numberValue{},
		boolValue{},
		new(procValue),
	}

	for i, v1 := range vals {
		if i == len(vals)-1 {
			break
		}

		for _, v2 := range vals[i+1:] {
			var err error

			_, err = v1.equals(v2)
			if err != errIncomparableValueTypes {
				t.Errorf("%v and %v should be incomparable", v1, v2)
			}

			_, err = v2.equals(v1)
			if err != errIncomparableValueTypes {
				t.Errorf("%v and %v should be incomparable", v2, v1)
			}
		}
	}
}
