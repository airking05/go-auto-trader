package models

import (
	"github.com/airking05/go-auto-trader/logger"
	"math"

	"github.com/markcheno/go-talib"
)

type goldencross struct {
	period int

	positionType PositionType
	param        float64
}

func NewGoldenCross(positionType PositionType, period int, param float64) Logic {
	return &goldencross{
		positionType: positionType,
		period:       period,
		param:        param,
	}
}

func (goldencross *goldencross) Execute(charts []Chart) bool {
	prices := make([]float64, 0)
	for _, chart := range charts {
		prices = append(prices, chart.Last)
	}
	if len(prices) < goldencross.period+1 {
		logger.Get().Infof("[goldencross:Execute]: not enough samples: %d", len(prices))
		return false
	}
	logger.Get().Infof("[goldencross:Execute] len(prices): %d, period: %d", len(prices), goldencross.period)

	macd, macdsignal, _ := talib.MacdExt(prices, 12, talib.SMA, 26, talib.SMA, 9, talib.SMA)

	if goldencross.positionType == Long {
		if (macd[len(macd)-2] < 0) && (macd[len(macd)-2] < macdsignal[len(macd)-2]) && (macd[len(macd)-1] > macdsignal[len(macd)-1]) {
			m1 := macd[len(macd)-1] - macd[len(macd)-2]
			m2 := macdsignal[len(macdsignal)-1] - macdsignal[len(macdsignal)-2]
			tan := math.Abs((m1 - m2) / (1 + m1*m2))
			if tan > goldencross.param {
				return true
			}
		}
	} else {
		if (macd[len(macd)-2] > 0) && (macd[len(macd)-2] > macdsignal[len(macd)-2]) && (macd[len(macd)-1] < macdsignal[len(macd)-1]) {
			m1 := macd[len(macd)-1] - macd[len(macd)-2]
			m2 := macdsignal[len(macdsignal)-1] - macdsignal[len(macdsignal)-2]
			tan := math.Abs((m1 - m2) / (1 + m1*m2))
			if tan > goldencross.param {
				return true
			}
		}
	}
	return false
}
