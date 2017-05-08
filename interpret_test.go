package main

import (
	"reflect"
	"testing"
)

func TestInterpret(t *testing.T) {
	cases := []struct {
		src  string
		want value
	}{
		{
			src:  `(let ((x (+ 1 5)) (y 4)) (+ y x))`,
			want: numberValue{10},
		},
		{
			src:  `((lambda (x y) (+ x y)) 1 2)`,
			want: numberValue{3},
		},
		{
			src: `
				(let ((add1 (lambda (x) (+ x 1)))
				      (add2 (lambda (x) (+ x 2)))
				      (x 5))
				  (add1 (add2 (add1 x))))
			`,
			want: numberValue{9},
		},
		{
			src:  `(if (= 1 2) (+1 2) (+ 3 4))`,
			want: numberValue{7},
		},
		{
			src: `
				(let ((tri (lambda (x f)
				             (if (= x 0)
				                 0
				                 (+ x (f (- x 1) f))))))
				  (tri 100 tri))
			`,
			want: numberValue{5050},
		},
		{
			src: `
				(let ((y (lambda (f)
				           ((lambda (procedure)
				              (f (lambda (arg) ((procedure procedure) arg))))
				            (lambda (procedure)
				              (f (lambda (arg) ((procedure procedure) arg)))))))
				      (fib-pre (lambda (f)
				                 (lambda (n)
				                   (if (= n 0)
				                       0
				                       (if (= n 1)
				                           1
				                           (+ (f (- n 1)) (f (- n 2)))))))))
				  (let ((fib (y fib-pre)))
				    (fib 10)))
			`,
			want: numberValue{55},
		},
	}

	for i, c := range cases {
		t.Logf("Case %d: %v", i, c.src)

		got, err := interpret(c.src)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("value:\ngot:  %v\nwant: %v", got, c.want)
		}
	}
}
