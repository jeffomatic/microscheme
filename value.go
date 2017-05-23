package main

import "errors"

var errIncomparableValueTypes = errors.New("cannot compare values of different types")

type value interface {
	valueType()
	equals(v value) (bool, error)
}

type orderable interface {
	greaterThan(v value) (bool, error)
}

type nullValue struct {
}

func (_ nullValue) valueType() {
	// does nothing
}

func (_ nullValue) equals(other value) (bool, error) {
	_, ok := other.(nullValue)
	return ok, nil
}

type numberValue struct {
	underlying int
}

func (_ numberValue) valueType() {
	// does nothing
}

func (v numberValue) equals(other value) (bool, error) {
	switch other := other.(type) {
	case nullValue:
		return false, nil
	case numberValue:
		return v.underlying == other.underlying, nil
	case *numberValue:
		return v.underlying == other.underlying, nil
	default:
		return false, nil
	}
}

func (v numberValue) greaterThan(other value) (bool, error) {
	switch other := other.(type) {
	case numberValue:
		return v.underlying > other.underlying, nil
	case *numberValue:
		return v.underlying > other.underlying, nil
	default:
		return false, errIncomparableValueTypes
	}
}

type boolValue struct {
	underlying bool
}

func (_ boolValue) valueType() {
	// does nothing
}

func (v boolValue) equals(other value) (bool, error) {
	switch other := other.(type) {
	case boolValue:
		return v.underlying == other.underlying, nil
	case *boolValue:
		return v.underlying == other.underlying, nil
	default:
		return false, nil
	}
}

type stringValue struct {
	underlying string
}

func (_ stringValue) valueType() {
	// does nothing
}

func (v stringValue) equals(other value) (bool, error) {
	switch other := other.(type) {
	case stringValue:
		return v.underlying == other.underlying, nil
	case *stringValue:
		return v.underlying == other.underlying, nil
	default:
		return false, nil
	}
}

type pairValue struct {
	car value
	cdr value
}

func (_ pairValue) valueType() {
	// does nothing
}

func (v pairValue) equals(other value) (bool, error) {
	switch other := other.(type) {
	case pairValue:
		return v == other, nil
	case *pairValue:
		if other == nil {
			return false, nil
		}
		return v == *other, nil
	default:
		return false, nil
	}
}

func makeList(vals []value) value {
	var res value = nullValue{}
	for i := len(vals) - 1; i >= 0; i-- {
		res = pairValue{car: vals[i], cdr: res}
	}
	return res
}

type procValue struct {
	formals []string
	rest    string
	body    []expression
	env     *frame
}

func (_ *procValue) valueType() {
	// does nothing
}

func (v *procValue) equals(other value) (bool, error) {
	switch other := other.(type) {
	case *procValue:
		return v == other, nil
	default:
		return false, nil
	}
}
