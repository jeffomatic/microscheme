package main

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	errNonBooleanPredicate  = errors.New("predicate must evaluate to boolean")
	errApplicationOnNonProc = errors.New("application operator must evaluate to proc")
)

func eval(expr expression, env *frame) (value, error) {
	t, err := classify(expr)
	if err != nil {
		return nil, err
	}

	switch t {
	case exprNull:
		return nullValue{}, nil
	case exprNumber:
		return evalNumber(expr, env)
	case exprBoolean:
		return boolValue{mustExpressionToken(expr) == "#t"}, nil
	case exprString:
		s := mustExpressionToken(expr)
		return stringValue{s[1 : len(s)-1]}, nil
	case exprDereference:
		return env.get(mustExpressionToken(expr))
	case exprDefine:
		return evalDefine(mustExpressionChildren(expr)[1:], env)
	case exprBegin:
		return evalSequence(mustExpressionChildren(expr)[1:], env)
	case exprIf:
		return evalIf(expr, env)
	case exprLambda:
		return evalLambda(expr, env)
	case exprLet:
		return evalLet(expr, env)
	case exprPrimitive:
		return evalPrimitive(expr, env)
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

func evalDefine(exprs []expression, env *frame) (value, error) {
	switch first := exprs[0].(type) {
	case *tokenExpression:
		k := mustExpressionToken(first)

		v, err := eval(exprs[1], env)
		if err != nil {
			return nil, err
		}

		env.set(k, v)
		return nullValue{}, nil

	case *compoundExpression:
		proc, err := evalNewProc(first.children[1:], exprs[1:], env)
		if err != nil {
			return nil, err
		}

		env.set(mustExpressionToken(first.children[0]), proc)
		return nullValue{}, nil

	default:
		panic(fmt.Sprintf("invalid define expression: %v", exprs))
	}
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

	if eq, _ := p.equals(boolValue{true}); eq {
		return eval(consequent, env)
	}

	return eval(alternative, env)
}

func evalLambda(expr expression, env *frame) (value, error) {
	c := mustExpressionChildren(expr)
	return evalNewProc(mustExpressionChildren(c[1]), c[2:], env)
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

func evalPrimitive(expr expression, env *frame) (value, error) {
	c := mustExpressionChildren(expr)
	f, ok := primitives[mustExpressionToken(c[1])]
	if !ok {
		return nil, errInvalidCompoundExpression
	}
	return f(c[2:], env)
}

func evalApplication(expr expression, env *frame) (value, error) {
	c := mustExpressionChildren(expr)
	fexpr := c[0]
	args := c[1:]

	fval, err := eval(fexpr, env)
	if err != nil {
		return nil, err
	}

	proc, ok := fval.(*procValue)
	if !ok {
		return nil, errApplicationOnNonProc
	}
	if len(args) < len(proc.formals) {
		return nil, errWrongNumberOfArguments
	}
	if len(args) > len(proc.formals) && proc.rest == "" {
		return nil, errWrongNumberOfArguments
	}

	nextEnv := proc.env.extend()

	for i, param := range proc.formals {
		argVal, err := eval(args[i], env)
		if err != nil {
			return nil, err
		}

		nextEnv.set(param, argVal)
	}

	if proc.rest != "" {
		var restVals []value
		for i := len(proc.formals); i < len(args); i++ {
			v, err := eval(args[i], env)
			if err != nil {
				return nil, err
			}

			restVals = append(restVals, v)
		}

		nextEnv.set(proc.rest, makeList(restVals))
	}

	return evalSequence(proc.body, nextEnv)
}
