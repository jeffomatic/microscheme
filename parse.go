package main

import "errors"

var (
	errInvalidClosingBrace = errors.New("invalid closing brace")
	errUnclosedExpression  = errors.New("unclosed expression")
)

type expression interface {
	expressionType()
}

type tokenExpression struct {
	token string
}

func (_ *tokenExpression) expressionType() {
	// does nothing
}

type compoundExpression struct {
	children []expression
}

func (_ *compoundExpression) expressionType() {
	// does nothing
}

func isTokenExpression(expr expression) bool {
	_, ok := expr.(*tokenExpression)
	return ok
}

func isCompoundExpression(expr expression) bool {
	_, ok := expr.(*compoundExpression)
	return ok
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
