package main

var stdlib *frame

func init() {
	stdlib = newFrame()

	stdlib.set("true", boolValue{true})
	stdlib.set("false", boolValue{false})

	libSrc := map[string]string{
		"+":   "(lambda (a b) (primitive + a b))",
		"-":   "(lambda (a b) (primitive - a b))",
		"*":   "(lambda (a b) (primitive * a b))",
		"/":   "(lambda (a b) (primitive / a b))",
		"not": `(lambda (a) (if a false true))`,
		"or":  `(lambda (a b) (if a true b))`,
		"and": `(lambda (a b) (if a b false))`,
		"xor": `(lambda (a b) (if a (not b) b))`,
		"=":   "(lambda (a b) (primitive = a b))",
		">":   "(lambda (a b) (primitive > a b))",
		">=":  "(lambda (a b) (or (= a b) (> a b)))",
		"<":   "(lambda (a b) (>= b a))",
		"<=":  "(lambda (a b) (> b a))",
	}

	for k, src := range libSrc {
		exprs, err := parse(tokenize(src))
		if err != nil {
			panic("failed to parse " + k + ": " + err.Error())
		}

		if len(exprs) != 1 {
			panic(k + " does not parse to single expression")
		}

		v, err := eval(exprs[0], stdlib)
		if err != nil {
			panic("failed to evaluate " + k + ": " + err.Error())
		}

		stdlib.set(k, v)
	}
}
