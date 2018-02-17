package models

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Status string

const (
	Running  = "running"
	Paused   = "paused"
	Resuming = "resuming"
	Removed  = "removed"
	Waiting  = "waiting"
)

type TradePhase int

const (
	Boot TradePhase = iota + 1
	Resume
	JudgeEntry
	MakePosition
	JudgeExit
	ClosePosition
	Reset

	ShutDown
)

type TraderGorm struct {
	gorm.Model
	ExchangeID               ExchangeID `yaml:"exchange_id"`
	APIKey                   string     `yaml:"api_key"`
	SecretKey                string     `yaml:"secret_key"`
	Trading                  string     `yaml:"trading"`
	Settlement               string     `yaml:"settlement"`
	LossCutRate              float64    `yaml:"loss_cut_rate"`
	ProfitTakeRate           float64    `yaml:"profit_take_rate"`
	AssetDistributionRate    float64    `yaml:"asset_distribution_rate"`
	TradingAmount            float64
	PositionType             PositionType `yaml:"position_type"`
	Duration                 int          `yaml:"duration"`
	WaitLimitSecond          int          `yaml:"wait_limit_second"`
	MakePositionLogicsString string       `yaml:"make_position_logics" sql:"type:text"`
	MakePositionLogicsYaml   LogicYaml    `sql:"-" yaml:"make_position_logics_yaml"`
	Status                   string       `yaml:"status"`
}

func (t *TraderGorm) BeforeSave(scope *gorm.Scope) (err error) {
	yml, err := yaml.Marshal(t.MakePositionLogicsYaml)
	if err != nil {
		return errors.Wrap(err, "couldn't marshal makePositionLogics to yaml")
	}

	if err := scope.SetColumn("MakePositionLogicsString", string(yml)); err != nil {
		return errors.Wrap(err, "failed to setColumn makePositionLogicsString")
	}
	return nil
}

func (t *TraderGorm) AfterFind() (err error) {
	var yml LogicYaml
	if err := yaml.Unmarshal([]byte(t.MakePositionLogicsString), &yml); err != nil {
		return errors.Wrap(err, "cannot unmarshal makePositionLogic as Yaml")
	}
	t.MakePositionLogicsYaml = yml
	return nil
}
