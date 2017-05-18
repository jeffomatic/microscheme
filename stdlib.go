package main

import "fmt"

var stdlib *frame

const src = `
(define true #t)
(define false #f)
(define (+ a b) (primitive + a b))
(define (- a b) (primitive - a b))
(define (* a b) (primitive * a b))
(define (/ a b) (primitive / a b))
(define (not a) (if a false true))
(define (or a b) (if a true b))
(define (and a b) (if a b false))
(define (xor a b) (if a (not b) b))
(define (= a b) (primitive = a b))
(define (> a b) (primitive > a b))
(define (>= a b) (or (= a b) (> a b)))
(define (< a b) (>= b a))
(define (<= a b) (> b a))
(define (cons a b) (primitive cons a b))
(define (car a) (primitive car a))
(define (cdr a) (primitive cdr a))
`

func init() {
	stdlib = newFrame()

	exprs, err := parse(tokenize(src))
	if err != nil {
		panic("failed to parse stdlib source: " + err.Error())
	}

	for _, expr := range exprs {
		_, err = eval(expr, stdlib)
		if err != nil {
			panic(fmt.Sprintf("failed to evaluate: %v:\nerror: %v", expr, err))
		}
	}
}
