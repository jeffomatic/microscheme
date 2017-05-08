package main

func interpret(src string) (value, error) {
	exprs, err := parse(tokenize(src))
	if err != nil {
		return nil, err
	}

	return evalSequence(exprs, stdlib.extend())
}
