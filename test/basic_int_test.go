package test

import (
	"testing"
)

func TestBasicInt(t *testing.T) {
	a := BasicAttrTypesObject{}
	b := BasicAttrTypesObject{}
    
    if !a.Equals(&b) {
        t.Errorf("a and b should be equal")
    }

    a.VInt = 3 
    b.VInt = 5

    if a.Equals(&b) {
        t.Errorf("a and b should not be equal")
    }

	d := a.Diff(&b)
	if len(d.Chgs) != 1 {
		t.Errorf("Found %d differences between a and b, expected 1", len(d.Chgs))
	}

	a.Apply(d)
	if !a.Equals(&b) {
		t.Errorf("a and b should be equal")
	}
}