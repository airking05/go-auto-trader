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
	ExchangeID               ExchangeID   `json:"exchange_id"`
	Trading                  string `json:"trading"`
	Settlement               string `json:"settlement"`
	LossCutRate              float64             `json:"loss_cut_rate"`
	ProfitTakeRate           float64             `json:"profit_take_rate"`
	TradingAmount            float64             `json:"trading_amount"`
	AssetDistributionRate    float64             `json:"asset_distribution_rate"`
	PositionType             PositionType `json:"position_type"`
	Duration                 int                 `json:"duration"`
	WaitLimitSecond          int                 `json:"wait_limit_second"`
	MakePositionLogicsString string              `json:"make_position_logics_string" sql:"type:text"`
	MakePositionLogicsYaml   LogicYaml     `sql:"-" json:"make_position_logics_yaml"`
	Status                   string              `json:"status"`
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