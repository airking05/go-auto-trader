package models

import (
)

type fakeLogic struct {
	Logic
	FakeExecute func([]Chart) bool
}

func (f *fakeLogic) Execute(cs []Chart) bool {
	return f.FakeExecute(cs)
}

func NewFakeLogic(isPassed bool) Logic {
	return &fakeLogic{
		FakeExecute: func([]Chart) bool {
			return isPassed
		},
	}
}
