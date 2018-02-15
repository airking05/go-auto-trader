package models

import (
	"encoding/json"
	errors "github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func convLogicYaml(m map[interface{}]interface{}) (LogicYaml, error) {
	nm := make(LogicYaml)
	for k, v := range m {
		if sk, ok := k.(string); !ok {
			return nil, errors.Errorf("cannot convert %v to string", k)
		} else {
			nm[sk] = v
		}
	}
	return nm, nil
}

type LogicYaml map[string]interface{}

func (y LogicYaml) MarshalJSON() ([]byte, error) {
	newmap := make(map[string]interface{})
	for k, v := range y {
		if x, ok := v.(map[interface{}]interface{}); ok {
			xm, err := convLogicYaml(x)
			if err != nil {
				return nil, errors.Wrap(err, "cannot convert map to LogicYaml")
			}
			newmap[k] = xm
		} else {
			newmap[k] = v
		}
	}
	return json.Marshal(newmap)
}

func (y LogicYaml) toLogic() Logic {
	name, ok := y["logic"]
	if !ok {
		panic(errors.Errorf("logic name isn't set"))
	}

	switch name {
	case "and":
		a := y.getAsLogicYaml("a").toLogic()
		b := y.getAsLogicYaml("b").toLogic()
		return And(a, b)
	case "or":
		a := y.getAsLogicYaml("a").toLogic()
		b := y.getAsLogicYaml("b").toLogic()
		return Or(a, b)
	case "not":
		a := y.getAsLogicYaml("a").toLogic()
		return Not(a)
	case "wmadif":
		period := y.getAsInt("period")
		difrate := y.getAsFloat64("difrate")
		return NewWMADif(period, difrate)
	case "smadif":
		period := y.getAsInt("period")
		difrate := y.getAsFloat64("difrate")
		return NewSMADif(period, difrate)
	case "emadif":
		period := y.getAsInt("period")
		difrate := y.getAsFloat64("difrate")
		return NewEMADif(period, difrate)
	case "obv":
		period := y.getAsInt("period")
		param := y.getAsFloat64("param")
		return NewOBV(Long, period, param)
	case "rsifollow":
		period := y.getAsInt("period")
		param := y.getAsFloat64("param")
		return NewRSIFollow(Long, period, param)
	case "rsicontrarian":
		period := y.getAsInt("period")
		param := y.getAsFloat64("param")
		return NewRSIContrarian(Long, period, param)
	case "goldencross":
		period := y.getAsInt("period")
		param := y.getAsFloat64("param")
		return NewGoldenCross(Long, period, param)
	case "smaLineCross":
		longPeriod := y.getAsInt("long_period")
		shortPeriod := y.getAsInt("short_period")
		keepPeriod := y.getAsInt("keep_period")
		return NewSmaLineCross(Long, shortPeriod, longPeriod, keepPeriod)
	}

	panic(errors.Errorf("unknown logic name: %s", name))
}

func (y LogicYaml) getAsInterface(key string) interface{} {
	if x, ok := y[key]; ok {
		return x
	}
	panic(errors.Errorf("key '%s' is not found", key))
}

func (y LogicYaml) getAsFloat64(key string) float64 {
	x := y.getAsInterface(key)
	if ret, ok := x.(float64); ok {
		return ret
	}
	x = y.getAsInt(key)
	if ret, ok := x.(int); ok {
		return float64(ret)
	}
	panic(errors.Errorf("key '%s'(value=%v) can't parse into float64", key, x))
}

func (y LogicYaml) getAsInt(key string) int {
	x := y.getAsInterface(key)
	if ret, ok := x.(int); ok {
		return ret
	}
	panic(errors.Errorf("key '%s'(value=%v) can't parse into int", key, x))
}

func (y LogicYaml) getAsString(key string) string {
	x := y.getAsInterface(key)
	if ret, ok := x.(string); ok {
		return ret
	}
	panic(errors.Errorf("key '%s'(value=%v) can't parse into string", key, x))
}

func (y LogicYaml) getAsLogicYaml(key string) LogicYaml {
	if yml, ok := y[key]; ok {
		if logicYml, ok := yml.(LogicYaml); ok {
			return logicYml
		}
		if yml, ok := yml.(map[interface{}]interface{}); ok {
			logicYml := make(LogicYaml)
			for k, v := range yml {
				if key, ok := k.(string); ok {
					logicYml[key] = v
				} else {
					panic(errors.Errorf("key '%s' is not string", k))
				}
			}
			return logicYml
		}
		panic(errors.Errorf("key '%s' is not logic", key))
	}
	panic(errors.Errorf("key '%s' is not found", key))
}

type YamlParser struct {
}

func NewYamlParser() *YamlParser {
	return &YamlParser{}
}

func (p *YamlParser) ParseLogicYaml(yml LogicYaml) (logic Logic, err error) {
	defer func() {
		if r := recover(); r != nil {
			logic = nil
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	return yml.toLogic(), nil
}

func (p *YamlParser) Parse(src string) (logic Logic, err error) {
	var yml LogicYaml
	if err := yaml.Unmarshal([]byte(src), &yml); err != nil {
		return nil, err
	}

	return p.ParseLogicYaml(yml)
}
