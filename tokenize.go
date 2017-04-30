package main

func tokenize(src string) []string {
	var (
		res     []string
		current string
	)

	appendCurrent := func() {
		if len(current) > 0 {
			res = append(res, current)
			current = ""
		}
	}

	for _, r := range src {
		switch r {

		case '(':
			fallthrough
		case ')':
			appendCurrent()
			res = append(res, string(r))

		case ' ':
			fallthrough
		case '\t':
			fallthrough
		case '\n':
			fallthrough
		case '\r':
			appendCurrent()

		default:
			current += string(r)
		}
	}

	appendCurrent()

	return res
}
