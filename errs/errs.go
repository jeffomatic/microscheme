package errs

import (
	"fmt"
	"strings"
)

type wrapper struct {
	underlying error
	message    string
}

func (w wrapper) Error() string {
	if w.message != "" {
		return w.message + ": " + w.underlying.Error()
	}
	return w.underlying.Error()
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

func Wrapf(err error, format string, params ...interface{}) error {
	return Wrap(err, fmt.Sprintf(format, params...))
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
