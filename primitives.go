package main

import "errors"

var primitives map[string]func([]expression, *frame) (value, error)

func init() {
	primitives = map[string]func([]expression, *frame) (value, error){
		"+": primitiveAdd,
		"-": primitiveSubtract,
		"*": primitiveMultiply,
		"/": primitiveDivide,
	}
}

var errInvalidArgumentType = errors.New("bad argument type")

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
	panic("not implemented")
}

func primitiveMultiply(argExprs []expression, env *frame) (value, error) {
	panic("not implemented")
}

func primitiveDivide(argExprs []expression, env *frame) (value, error) {
	panic("not implemented")
}
