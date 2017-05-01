package main

import "errors"

var (
	errInvalidClosingBrace = errors.New("invalid closing brace")
	errUnclosedExpression  = errors.New("unclosed expression")
)

type expression interface {
	expression()
}

type tokenExpression struct {
	token string
}

func (_ *tokenExpression) expression() {
	// does nothing
}

type compoundExpression struct {
	children []expression
}

func (_ *compoundExpression) expression() {
	// does nothing
}

func parse(tokens []string) ([]expression, error) {
	var (
		res   []expression
		stack []*compoundExpression
	)

	for _, t := range tokens {
		switch t {
		case "(":
			n := new(compoundExpression)
			if len(stack) == 0 {
				res = append(res, n)
				stack = append(stack, n)
			} else {
				parent := stack[len(stack)-1]
				parent.children = append(parent.children, n)
				stack = append(stack, n)
			}

		case ")":
			if len(stack) == 0 {
				return res, errInvalidClosingBrace
			}
			stack = stack[0 : len(stack)-1]

		default:
			n := new(tokenExpression)
			n.token = t
			if len(stack) == 0 {
				res = append(res, n)
			} else {
				parent := stack[len(stack)-1]
				parent.children = append(parent.children, n)
			}
		}
	}

	if len(stack) != 0 {
		return res, errUnclosedExpression
	}

	return res, nil
}
