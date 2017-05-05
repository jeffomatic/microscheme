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
		return evalNumber(expr, env)
	case exprBoolean:
		return boolValue{mustExpressionToken(expr) == "#t"}, nil
	case exprDereference:
		return env.get(mustExpressionToken(expr))
	case exprBegin:
		return evalSequence(mustExpressionChildren(expr)[1:], env)
	case exprIf:
		return evalIf(expr, env)
	case exprLambda:
		return evalLambda(expr, env)
	case exprLet:
		return evalLet(expr, env)
	case exprApplication:
		return evalApplication(expr, env)
	default:
		panic("classified type cannot be evaluated: " + string(t))
	}
}

func evalNumber(expr expression, env *frame) (value, error) {
	num, err := strconv.Atoi(mustExpressionToken(expr))
	if err != nil {
		panic(fmt.Sprintf("value %v should be valid number but error: %v", expr, err))
	}
	return numberValue{num}, nil
}

func evalSequence(exprs []expression, env *frame) (value, error) {
	var (
		last value = theNullValue
		err  error
	)
	for _, c := range exprs {
		last, err = eval(c, env)
		if err != nil {
			return nil, err
		}
	}
	return last, nil
}

func evalIf(expr expression, env *frame) (value, error) {
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
}

func evalLambda(expr expression, env *frame) (value, error) {
	c := mustExpressionChildren(expr)
	paramExprs := c[1]
	body := c[2:]

	var params []string
	for _, p := range mustExpressionChildren(paramExprs) {
		params = append(params, mustExpressionToken(p))
	}

	return &procValue{params: params, body: body}, nil
}

func evalLet(expr expression, env *frame) (value, error) {
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

	return evalSequence(body, nextEnv)
}

func evalApplication(expr expression, env *frame) (value, error) {
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

	return evalSequence(fproc.body, nextEnv)
}
