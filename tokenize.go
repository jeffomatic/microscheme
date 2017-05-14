package main

var escapes = map[rune]string{
	'a': "\a",
	'b': "\b",
	'f': "\f",
	'n': "\n",
	'r': "\r",
	't': "\t",
	'v': "\v",
}

func tokenize(src string) []string {
	var (
		res     []string
		current string

		isString bool
		isEscape bool

		finishCurrent = func() {
			if len(current) > 0 {
				res = append(res, current)
				current = ""
			}
		}
	)

	for _, r := range src {
		if isString && isEscape {
			if s, ok := escapes[r]; ok {
				current += s
			} else {
				current += string(r)
			}
			isEscape = false
		} else if isString {
			switch r {
			case '"':
				current += string(r)
				finishCurrent()
				isString = false
				isEscape = false
			case '\\':
				isEscape = true
			case '\n', '\r':
				panic("multiline string - make this an error")
			default:
				current += string(r)
			}
		} else {
			switch r {
			case '(', ')':
				finishCurrent()
				res = append(res, string(r))
			case ' ', '\t', '\n', '\r':
				finishCurrent()
			case '"':
				finishCurrent()
				current += string(r)
				isString = true
			default:
				current += string(r)
			}
		}
	}

	if isString {
		panic("unclosed string - make this an error")
	}

	finishCurrent()

	return res
}
