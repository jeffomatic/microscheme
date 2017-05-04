package main

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	theNullValue  value
	theTrueValue  value
	theFalseValue value

	errNonBooleanPredicate    = errors.New("predicate must evaluate to boolean")
	errApplicationOnNonProc   = errors.New("application operator must evaluate to proc")
	errWrongNumberOfArguments = errors.New("application with wrong number of arguments")
)

func init() {
	theNullValue = nullValue{}
	theTrueValue = boolValue{true}
	theFalseValue = boolValue{false}
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
		var (
			last value = theNullValue
			err  error
		)
		for _, c := range mustExpressionChildren(expr)[1:] {
			last, err = eval(c, env)
			if err != nil {
				return nil, err
			}
		}
		return last, nil

	case exprIf:
		c := mustExpressionChildren(expr)
		predicate := c[1]
		consequent := c[2]
		alternative := c[3]

		p, err := eval(predicate, env)
		if err != nil {
			return nil, err
		}

		if _, ok := p.(boolValue); !ok {
			return nil, errNonBooleanPredicate
		}

		if eq, _ := p.equals(theTrueValue); eq {
			return eval(consequent, env)
		}

		return eval(alternative, env)

	case exprLambda:
		c := mustExpressionChildren(expr)
		paramExprs := c[1]
		body := c[2:]

		var params []string
		for _, p := range mustExpressionChildren(paramExprs) {
			params = append(params, mustExpressionToken(p))
		}

		return &procValue{params: params, body: body}, nil

	case exprLet:
		c := mustExpressionChildren(expr)
		assignments := c[1]
		body := c[2:]

		nextEnv := env.extend()
		for _, a := range mustExpressionChildren(assignments) {
			aexprs := mustExpressionChildren(a)
			identifier := aexprs[0]
			rvalExpr := aexprs[1]

			rval, err := eval(rvalExpr, env)
			if err != nil {
				return nil, err
			}

			nextEnv.set(mustExpressionToken(identifier), rval)
		}

		var (
			lastVal value
			err     error
		)
		for _, b := range body {
			lastVal, err = eval(b, nextEnv)
			if err != nil {
				return nil, err
			}
		}

		return lastVal, err

	case exprApplication:
		c := mustExpressionChildren(expr)
		fexpr := c[0]
		args := c[1:]

		fval, err := eval(fexpr, env)
		if err != nil {
			return nil, err
		}

		fproc, ok := fval.(*procValue)
		if !ok {
			return nil, errApplicationOnNonProc
		}
		if len(args) != len(fproc.params) {
			return nil, errWrongNumberOfArguments
		}

		nextEnv := env.extend()
		for i, param := range fproc.params {
			argVal, err := eval(args[i], env)
			if err != nil {
				return nil, err
			}

			nextEnv.set(param, argVal)
		}

		var lastVal value
		for _, b := range fproc.body {
			var err error
			lastVal, err = eval(b, nextEnv)
			if err != nil {
				return nil, err
			}
		}

		return lastVal, err

	default:
		panic("classified type cannot be evaluated: " + string(t))
	}
}
