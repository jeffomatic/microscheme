package errs

import (
	"errors"
	"testing"
)

func TestErrs(t *testing.T) {
	e := errors.New("root")
	w1 := Wrap(e, "w1")
	w2 := Wrapf(w1, "w%d", 2)
	w3 := WrapAfterf(w2, "%s3", "w")

	if w1.Error() != "w1: root" {
		t.Errorf("w1.Error():\ngot:  %v\nwant: %v", w1.Error(), "w1: root")
	}

	if w2.Error() != "w2: w1: root" {
		t.Errorf("w2.Error():\ngot:  %v\nwant: %v", w2.Error(), "w2: w1: root")
	}

	if w3.Error() != "w2: w1: root: w3" {
		t.Errorf("w3.Error():\ngot:  %v\nwant: %v", w3.Error(), "w2: w1: root: w3")
	}

	if Root(w1) != e {
		t.Errorf("Root(w1):\ngot:  %v\nwant: %v", Root(w1), e)
	}

	if Root(w2) != e {
		t.Errorf("Root(w2):\ngot:  %v\nwant: %v", Root(w2), e)
	}

	if Root(w3) != e {
		t.Errorf("Root(w3):\ngot:  %v\nwant: %v", Root(w3), e)
	}
}
