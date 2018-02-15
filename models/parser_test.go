package models

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestParamYamlToLogicAndMap(t *testing.T) {
	str := `
logic: and
a:
  logic: rsifollow
  period: 30
  param: 10
b: 
  logic: rsifollow
  period: 30
  param: 10
`
	p := NewYamlParser()
	_, err := p.Parse(str)
	if err != nil {
		t.Fatal(err)
	}
}

func TestYamlParser(t *testing.T) {
	str := `
logic: and
a:
  logic: rsifollow
  period: 30
  param: 20.0
b:
  logic: rsifollow
  period: 30
  param: 25.0
`
	p := NewYamlParser()
	_, err := p.Parse(str)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogicYamlToJSON(t *testing.T) {
	str := `
logic: and
a:
  logic: rsifollow
  period: 30
  param: 20.0
b:
  logic: rsifollow
  period: 30
  param: 25.0
`

	var yml LogicYaml

	if err := yaml.Unmarshal([]byte(str), &yml); err != nil {
		t.Fatal(err)
	}

	if _, err := json.Marshal(yml); err != nil {
		t.Fatal(err)
	}
}
