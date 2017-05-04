package main

import "errors"

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
		return nil, errBindingNotFound
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
