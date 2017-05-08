package errs

import (
	"fmt"
	"strings"
)

type wrapper struct {
	underlying error
	message    string
	suffix     bool
}

func (w wrapper) Error() string {
	if w.message == "" {
		return w.underlying.Error()
	}

	var toks []string
	if w.suffix {
		toks = []string{w.underlying.Error(), w.message}
	} else {
		toks = []string{w.message, w.underlying.Error()}
	}

	return strings.Join(toks, ": ")
}

func Wrap(err error, messages ...string) error {
	if err == nil {
		return nil
	}

	var message string
	if len(messages) > 0 {
		message = strings.Join(messages, ": ")
	}

	return wrapper{underlying: err, message: message}
}

func WrapAfter(err error, messages ...string) error {
	if err == nil {
		return nil
	}

	w := Wrap(err, messages...).(wrapper)
	w.suffix = true // w will not be nil

	return w
}

func Wrapf(err error, format string, params ...interface{}) error {
	return Wrap(err, fmt.Sprintf(format, params...))
}

func WrapAfterf(err error, format string, params ...interface{}) error {
	return WrapAfter(err, fmt.Sprintf(format, params...))
}

func Root(err error) error {
	if err == nil {
		return nil
	}

	if w, ok := err.(wrapper); ok {
		return Root(w.underlying)
	}

	return err
}
