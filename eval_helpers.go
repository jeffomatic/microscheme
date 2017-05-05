package main

func mapEval(exprs []expression, env *frame) ([]value, error) {
	var res []value
	for _, c := range exprs {
		v, err := eval(c, env)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}

func evalSequence(exprs []expression, env *frame) (value, error) {
	if len(exprs) == 0 {
		return theNullValue, nil
	}

	values, err := mapEval(exprs, env)
	if err != nil {
		return nil, err
	}

	return values[len(values)-1], nil
}
