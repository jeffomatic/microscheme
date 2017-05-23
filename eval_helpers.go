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
		return nullValue{}, nil
	}

	values, err := mapEval(exprs, env)
	if err != nil {
		return nil, err
	}

	return values[len(values)-1], nil
}

func evalNewProc(paramExprs []expression, body []expression, env *frame) (value, error) {
	var (
		pv   = &procValue{body: body, env: env}
		toks []string
	)

	for _, exp := range paramExprs {
		toks = append(toks, mustExpressionToken(exp))
	}

	if len(toks) >= 2 && toks[len(toks)-2] == "." {
		last := toks[len(toks)-1]

		if !validIdentifier(last) {
			return nil, errInvalidCompoundExpression
		}

		pv.rest = last
		toks = toks[0 : len(toks)-2]
	}

	for _, t := range toks {
		if !validIdentifier(t) {
			return nil, errInvalidCompoundExpression
		}

		pv.formals = append(pv.formals, t)
	}

	return pv, nil
}

func validIdentifier(s string) bool {
	return s != "."
}
