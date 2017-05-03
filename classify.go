package main

import (
	"errors"
	"regexp"
)

type expressionType int

var (
	errInvalidExpressionType     = errors.New("invalid expression type")
	errInvalidCompoundExpression = errors.New("invalid compound expression")
)

const (
	exprInvalid expressionType = iota

	exprNull = iota

	// token expression types
	exprNumber      = iota
	exprBoolean     = iota
	exprDereference = iota

	// compound expression types
	exprSequence    = iota
	exprIf          = iota
	exprLambda      = iota
	exprLet         = iota
	exprApplication = iota
)

func classify(expr expression) (expressionType, error) {
	switch e := expr.(type) {
	case *tokenExpression:
		return classifyToken(e)
	case *compoundExpression:
		return classifyCompound(e)
	default:
		return exprInvalid, errInvalidExpressionType
	}
}

func classifyToken(expr *tokenExpression) (expressionType, error) {
	numberRegexp := regexp.MustCompile(`^-?\d+$`)

	switch {
	case expr.token == "null":
		return exprNull, nil
	case expr.token == "#t":
		return exprBoolean, nil
	case expr.token == "#f":
		return exprBoolean, nil
	case numberRegexp.Match([]byte(expr.token)):
		return exprNumber, nil
	default:
		return exprDereference, nil
	}
}

func classifyCompound(expr *compoundExpression) (expressionType, error) {
	if len(expr.children) == 0 {
		return exprNull, nil
	}

	if isCompoundExpression(expr.children[0]) {
		return exprApplication, nil
	}

	c, ok := expr.children[0].(*tokenExpression)
	if !ok {
		return exprInvalid, errInvalidCompoundExpression
	}

	switch c.token {
	case "begin":
		return exprSequence, nil
	case "if":
		if len(expr.children) != 4 {
			return exprInvalid, errInvalidCompoundExpression
		}
		return exprIf, nil
	case "lambda":
		if len(expr.children) < 3 || !isCompoundExpression(expr.children[1]) {
			return exprInvalid, errInvalidCompoundExpression
		}
		return exprLambda, nil
	case "let":
		if len(expr.children) < 3 || !isCompoundExpression(expr.children[1]) {
			return exprInvalid, errInvalidCompoundExpression
		}

		switch assignments := expr.children[1].(type) {
		case *compoundExpression:
			for _, c := range assignments.children {
				switch assign := c.(type) {
				case *compoundExpression:
					if len(assign.children) != 2 {
						return exprInvalid, errInvalidCompoundExpression
					}

					if !isTokenExpression(assign.children[0]) {
						return exprInvalid, errInvalidCompoundExpression
					}
				default:
					return exprInvalid, errInvalidCompoundExpression
				}
			}
		default:
			return exprInvalid, errInvalidCompoundExpression
		}

		return exprLet, nil
	default:
		return exprApplication, nil
	}
}
