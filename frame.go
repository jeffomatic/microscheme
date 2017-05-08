package main

import (
	"errors"
	"fmt"

	"sgheme/errs"
)

var errBindingNotFound = errors.New("environment does not contain binding")

type frame struct {
	parent *frame
	table  map[string]value
}

func newFrame() *frame {
	f := new(frame)
	f.table = make(map[string]value)
	return f
}

func (f *frame) get(k string) (value, error) {
	if v, ok := f.table[k]; ok {
		return v, nil
	}

	if f.parent == nil {
		return nil, errs.WrapAfterf(errBindingNotFound, "%q", k)
	}

	return f.parent.get(k)
}

func (f *frame) set(k string, v value) {
	f.table[k] = v
}

func (f *frame) extend() *frame {
	res := newFrame()
	res.parent = f
	return res
}

func (f *frame) debug() string {
	res := ""
	i := 0

	for f != nil {
		res += fmt.Sprintf("-- frame %d --\n", i)

		for k, v := range f.table {
			res += fmt.Sprintf("%s: %+v", k, v)
		}

		f = f.parent
		i += 1
	}

	return res
}
