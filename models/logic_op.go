package models

import (
)

func And(a Logic, b Logic) Logic {
	return &logicAnd{a, b}
}

func Or(a Logic, b Logic) Logic {
	return &logicOr{a, b}
}

func Not(a Logic) Logic {
	return &logicNot{a}
}

type logicAnd struct {
	a Logic `yaml:"a"`
	b Logic `yaml:"b"`
}

func (l *logicAnd) Execute(charts []Chart) bool {
	return l.a.Execute(charts) && l.b.Execute(charts)
}

type logicOr struct {
	a Logic `yaml:"a"`
	b Logic `yaml:"b"`
}

func (l *logicOr) Execute(charts []Chart) bool {
	return l.a.Execute(charts) || l.b.Execute(charts)
}

type logicNot struct {
	a Logic `yaml:"a"`
}

func (l *logicNot) Execute(charts []Chart) bool {
	return !l.a.Execute(charts)
}
