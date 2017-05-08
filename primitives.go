package main

import "errors"

var primitives map[string]func([]expression, *frame) (value, error)

func init() {
	primitives = map[string]func([]expression, *frame) (value, error){
		"+": primitiveAdd,
		"-": primitiveSubtract,
		"*": primitiveMultiply,
		"/": primitiveDivide,
		"=": primitiveEquals,
		">": primitiveGreaterThan,
	}
}

var (
	errInvalidArgumentType = errors.New("bad argument type")
	errDivideByZero        = errors.New("divide by zero")
	errTypeNotOrderable    = errors.New("type is not orderable")
)

func primitiveAdd(argExprs []expression, env *frame) (value, error) {
	args, err := mapEval(argExprs, env)
	if err != nil {
		return nil, err
	}

	var total int
	for _, v := range args {
		n, ok := v.(numberValue)
		if !ok {
			return nil, errInvalidArgumentType
		}
		total += n.underlying
	}

	return numberValue{total}, nil
}

func primitiveSubtract(argExprs []expression, env *frame) (value, error) {
	args, err := mapEval(argExprs, env)
	if err != nil {
		return nil, err
	}

	if len(args) == 0 {
		return numberValue{0}, nil
	}

	first, ok := args[0].(numberValue)
	if !ok {
		return nil, errInvalidArgumentType
	}

	if len(args) == 1 { // Single argument is a special case (negation)
		return numberValue{-first.underlying}, nil
	}

	total := first.underlying
	for _, v := range args[1:] {
		n, ok := v.(numberValue)
		if !ok {
			return nil, errInvalidArgumentType
		}
		total -= n.underlying
	}

	return numberValue{total}, nil
}

func primitiveMultiply(argExprs []expression, env *frame) (value, error) {
	args, err := mapEval(argExprs, env)
	if err != nil {
		return nil, err
	}

	total := 1
	for _, v := range args {
		n, ok := v.(numberValue)
		if !ok {
			return nil, errInvalidArgumentType
		}
		total *= n.underlying
	}

	return numberValue{total}, nil
}

func primitiveDivide(argExprs []expression, env *frame) (value, error) {
	args, err := mapEval(argExprs, env)
	if err != nil {
		return nil, err
	}

	if len(args) == 0 {
		return numberValue{1}, nil
	}

	first, ok := args[0].(numberValue)
	if !ok {
		return nil, errInvalidArgumentType
	}

	res := first.underlying
	for _, v := range args[1:] {
		n, ok := v.(numberValue)
		if !ok {
			return nil, errInvalidArgumentType
		}

		if n.underlying == 0 {
			return nil, errDivideByZero
		}

		res /= n.underlying
	}

	return numberValue{res}, nil
}

func primitiveEquals(argExprs []expression, env *frame) (value, error) {
	if len(argExprs) != 2 {
		return nil, errWrongNumberOfArguments
	}

	args, err := mapEval(argExprs, env)
	if err != nil {
		return nil, err
	}

	res, err := args[0].equals(args[1])
	if err != nil {
		return nil, err
	}

	return boolValue{res}, nil
}

func primitiveGreaterThan(argExprs []expression, env *frame) (value, error) {
	if len(argExprs) != 2 {
		return nil, errWrongNumberOfArguments
	}

	args, err := mapEval(argExprs, env)
	if err != nil {
		return nil, err
	}

	a, ok := args[0].(orderable)
	if !ok {
		return nil, errTypeNotOrderable
	}

	res, err := a.greaterThan(args[1])
	if err != nil {
		return nil, err
	}

	return boolValue{res}, nil
}
