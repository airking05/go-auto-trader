package models

import (
	"testing"
)

func TestAnd(t *testing.T) {
	a := NewFakeLogic(true)
	b := NewFakeLogic(true)
	charts := make([]Chart, 0)
	for _, close := range testClose {
		charts = append(charts, Chart{
			Last: close,
		})
	}
	and := And(a, b)
	if and.Execute(charts) == false {
		t.FailNow()
	}
}

func TestOr(t *testing.T) {
	a := NewFakeLogic(true)
	b := NewFakeLogic(false)
	charts := make([]Chart, 0)
	for _, close := range testClose {
		charts = append(charts, Chart{
			Last: close,
		})
	}
	or := Or(a, b)
	if or.Execute(charts) == false {
		t.FailNow()
	}
}

func TestNot(t *testing.T) {
	a := NewFakeLogic(false)
	charts := make([]Chart, 0)
	for _, close := range testClose {
		charts = append(charts, Chart{
			Last: close,
		})
	}
	not := Not(a)
	if not.Execute(charts) == false {
		t.FailNow()
	}
}
