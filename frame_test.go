package main

import "testing"

func TestFrame(t *testing.T) {
	f1 := newFrame()
	f2 := f1.extend()

	f1.set("number", numberValue{})
	f2.set("bool", boolValue{})

	var (
		v   value
		err error
	)

	v, err = f1.get("number")
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if ok, _ := v.equals(numberValue{}); !ok {
		t.Errorf("unequal value")
	}

	v, err = f1.get("bool")
	if err != errBindingNotFound {
		t.Error("error:\ngot:  %v\nwant: %v", err, errBindingNotFound)
	}

	v, err = f2.get("number")
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if ok, _ := v.equals(numberValue{}); !ok {
		t.Errorf("unequal value")
	}

	v, err = f2.get("bool")
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if ok, _ := v.equals(boolValue{}); !ok {
		t.Errorf("unequal value")
	}

	v, err = f1.get("invalid")
	if err != errBindingNotFound {
		t.Error("error:\ngot:  %v\nwant: %v", err, errBindingNotFound)
	}
}