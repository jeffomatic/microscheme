package main

import (
	"fmt"
	"strconv"
)

var theNullValue value

func init() {
	theNullValue = nullValue{}
}

func eval(expr expression, env *frame) (value, error) {
	t, err := classify(expr)
	if err != nil {
		return nil, err
	}

	switch t {
	case exprNull:
		return theNullValue, nil
	case exprNumber:
		num, err := strconv.Atoi(mustExpressionToken(expr))
		if err != nil {
			panic(fmt.Sprintf("value %v should be valid number but error: %v", expr, err))
		}
		return numberValue{num}, nil
	case exprBoolean:
		return boolValue{mustExpressionToken(expr) == "#t"}, nil
	case exprDereference:
		return env.get(mustExpressionToken(expr))
	case exprSequence:
		return theNullValue, nil
	case exprIf:
		return theNullValue, nil
	case exprLambda:
		return theNullValue, nil
	case exprLet:
		return theNullValue, nil
	case exprApplication:
		return theNullValue, nil
	default:
		panic("classified type cannot be evaluated: " + string(t))
	}
}
