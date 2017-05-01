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
		case '(', ')':
			appendCurrent()
			res = append(res, string(r))
		case ' ', '\t', '\n', '\r':
			appendCurrent()
		default:
			current += string(r)
		}
	}

	appendCurrent()

	return res
}
