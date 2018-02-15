package models

import (
	logger "github.com/airking05/go-auto-trader/logger"
	"github.com/markcheno/go-talib"

	"gopkg.in/yaml.v2"
	"math"
)

func ConvLogicToYamlString(l Logic) (string, error) {
	d, err := yaml.Marshal(l)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

type Logic interface {
	Execute([]Chart) bool
}

type rsifollow struct {
	period       int          `yaml:"period"`
	positionType PositionType `yaml:"position_type"`
	param        float64      `yaml:"param"`
}
type rsicontrarian struct {
	period       int          `yaml:"period"`
	positionType PositionType `yaml:"position_type"`
	param        float64      `yaml:"param"`
}
type obv struct {
	period       int          `yaml:"period"`
	positionType PositionType `yaml:"position_type"`
	param        float64      `yaml:"param"`
}

type emadif struct {
	period  int     `yaml:"period"`
	difRate float64 `yaml:"difrate"`
}

type smadif struct {
	period  int     `yaml:"period"`
	difRate float64 `yaml:"difrate"`
}

type wmadif struct {
	period  int     `yaml:"period"`
	difRate float64 `yaml:"difrate"`
}

type smaLineCross struct {
	positionType PositionType `yaml:"position_type"`
	shortPeriod  int          `yaml:"short_period"`
	longPeriod   int          `yaml:"long_period"`
	keepPeriod   int          `yaml:"keep_period"`
}

func NewSmaLineCross(positionType PositionType, shortPeriod int, longPeriod int, keepPeriod int) Logic {
	return &smaLineCross{
		positionType: positionType,
		shortPeriod:  shortPeriod,
		longPeriod:   longPeriod,
		keepPeriod:   keepPeriod,
	}
}

func (smaLineCross *smaLineCross) Execute(charts []Chart) bool {
	prices := make([]float64, 0)
	for _, chart := range charts {
		prices = append(prices, chart.Last)
	}
	if len(prices) < smaLineCross.longPeriod {
		logger.Get().Infof("[smaLineCross:Execute]: not enough samples: %d", len(prices))
		return false
	}

	shortSma := talib.Sma(prices, smaLineCross.shortPeriod)
	longSma := talib.Sma(prices, smaLineCross.longPeriod)
	if len(shortSma) < 1 || len(longSma) < 1 {
		logger.Get().Error("[smaLineCross:Execute] there is no ema")
		return false
	}

	logger.Get().Info(longSma)
	logger.Get().Info(shortSma)
	if smaLineCross.positionType == Long {
		for i := 1; i <= smaLineCross.keepPeriod; i++ {
			if shortSma[len(shortSma)-i] < longSma[len(longSma)-i] {
				return false
			}
		}
		if shortSma[len(shortSma)-smaLineCross.keepPeriod-1] >
			longSma[len(longSma)-smaLineCross.keepPeriod-1] {
			return false
		}
		return true
	} else {
		for i := 1; i <= smaLineCross.keepPeriod; i++ {
			if shortSma[len(shortSma)-i] > longSma[len(longSma)-i] {
				return false
			}
		}
		if shortSma[len(shortSma)-smaLineCross.keepPeriod-1] <
			longSma[len(longSma)-smaLineCross.keepPeriod-1] {
			return false
		}
		return true
	}
}

func NewEMADif(period int, difRate float64) Logic {
	return &emadif{
		period:  period,
		difRate: difRate,
	}
}
func (emadif *emadif) Execute(charts []Chart) bool {
	prices := make([]float64, 0)
	for _, chart := range charts {
		prices = append(prices, chart.Last)
	}
	if len(prices) < emadif.period+1 {
		logger.Get().Infof("[emadif:Execute]: not enough samples: %d", len(prices))
		return false
	}
	logger.Get().Infof("[emadif:Execute] len(prices): %d, period: %d", len(prices), emadif.period)

	ema := talib.Ema(prices, emadif.period)
	if len(ema) > 0 {
		logger.Get().Infof("[emadif:Execute] rsi: %v", ema[len(ema)-1])
	} else {
		logger.Get().Error("[emadif:Execute] there is no ema")
		return false
	}
	difRate := math.Abs(prices[len(prices)-1] / ema[len(ema)-1])
	logger.Get().Infof("[emadif:Execute] difRate: %v", difRate)
	if difRate > emadif.difRate {
		return true
	}
	return false
}

func NewSMADif(period int, difRate float64) Logic {
	return &smadif{
		period:  period,
		difRate: difRate,
	}
}

func (smadif *smadif) Execute(charts []Chart) bool {
	prices := make([]float64, 0)
	for _, chart := range charts {
		prices = append(prices, chart.Last)
	}
	if len(prices) < smadif.period+1 {
		logger.Get().Infof("[smadif:Execute]: not enough samples: %d", len(prices))
		return false
	}
	logger.Get().Infof("[smadif:Execute] len(prices): %d, period: %d", len(prices), smadif.period)
	sma := talib.Sma(prices, smadif.period)
	if len(sma) > 0 {
		logger.Get().Infof("[smadif:Execute] sma: %v", sma[len(sma)-1])
	} else {
		logger.Get().Error("[smadif:Execute] there is no sma")
		return false
	}
	difRate := math.Abs(prices[len(prices)-1] / sma[len(sma)-1])
	logger.Get().Infof("[smadif:Execute] difRate: %v", difRate)
	if difRate > smadif.difRate {
		return true
	}
	return false
}

func NewWMADif(period int, difRate float64) Logic {
	return &wmadif{
		period:  period,
		difRate: difRate,
	}
}
func (wmadif *wmadif) Execute(charts []Chart) bool {
	prices := make([]float64, 0)
	for _, chart := range charts {
		prices = append(prices, chart.Last)
	}
	if len(prices) < wmadif.period+1 {
		logger.Get().Infof("[wmadif:Execute]: not enough samples: %d", len(prices))
		return false
	}
	logger.Get().Infof("[wmadif:Execute] len(prices): %d, period: %d", len(prices), wmadif.period)
	wma := talib.Wma(prices, wmadif.period)
	if len(wma) > 0 {
		logger.Get().Infof("[wmadif:Execute] wma: %v", wma[len(wma)-1])
	} else {
		logger.Get().Error("[wmadif:Execute] there is no wma")
		return false
	}
	difRate := math.Abs(prices[len(prices)-1] / wma[len(wma)-1])
	logger.Get().Infof("[wmadif:Execute] difRate: %v", difRate)
	if difRate > wmadif.difRate {
		return true
	}
	return false
}

func NewRSIFollow(positionType PositionType, period int, param float64) Logic {
	return &rsifollow{
		positionType: positionType,
		period:       period,
		param:        param,
	}
}

func NewRSIContrarian(positionType PositionType, period int, param float64) Logic {
	return &rsicontrarian{
		positionType: positionType,
		period:       period,
		param:        param,
	}
}

func NewOBV(positionType PositionType, period int, param float64) Logic {
	return &obv{
		positionType: positionType,
		period:       period,
		param:        param,
	}
}

// We usually use it with 14 days periods.
// It's too much bought if the value is larger than 75%,
// and it's too much sold if the value is smaller than 25%.

func (rsifollow *rsifollow) Execute(charts []Chart) bool {
	prices := make([]float64, 0)
	for _, chart := range charts {
		prices = append(prices, chart.Last)
	}
	if len(prices) < rsifollow.period+1 {
		logger.Get().Infof("[rsifollow:Execute]: not enough samples: %d", len(prices))
		return false
	}
	logger.Get().Infof("[rsifollow:Execute] len(prices): %d, period: %d", len(prices), rsifollow.period)
	rsi := talib.Rsi(prices, rsifollow.period)
	if len(rsi) > 0 {
		logger.Get().Infof("[rsifollow:Execute] rsi: %v", rsi[len(rsi)-1])
	} else {
		logger.Get().Error("[rsifollow:Execute] there is no rsi")
		return false
	}
	if rsifollow.positionType == 1 {
		if rsi[len(rsi)-1] >= rsifollow.param {
			logger.Get().Info("[rsifollow:Execute] judge true")
			return true
		}
	}
	if rsifollow.positionType == 0 {
		if rsi[len(rsi)-1] <= rsifollow.param {
			logger.Get().Info("[rsilogic:Execute] judge true")
			return true
		}
	}
	return false
}

func (rsicontrarian *rsicontrarian) Execute(charts []Chart) bool {
	prices := make([]float64, 0)
	for _, chart := range charts {
		prices = append(prices, chart.Last)
	}
	if len(prices) < rsicontrarian.period+1 {
		logger.Get().Infof("[rsicontrarian:Execute]: not enough samples: %d", len(prices))
		return false
	}
	logger.Get().Infof("[rsicontrarian:Execute] len(prices): %d, period: %d", len(prices), rsicontrarian.period)
	rsi := talib.Rsi(prices, rsicontrarian.period)
	if len(rsi) > 0 {
		logger.Get().Infof("[rsicontrarian:Execute] rsi: %v", rsi[len(rsi)-1])
	} else {
		logger.Get().Error("[rsicontrarian:Execute] there is no rsi")
		return false
	}
	if rsicontrarian.positionType == 1 {
		if rsi[len(rsi)-1] <= rsicontrarian.param {
			logger.Get().Info("[rsicontrarian:Execute] judge true")
			return true
		}
	}
	if rsicontrarian.positionType == 0 {
		if rsi[len(rsi)-1] >= rsicontrarian.param {
			logger.Get().Info("[rsicontrarian:Execute] judge true")
			return true
		}
	}
	return false
}

// Using this logic, you can detect up trend and down trend.
// This manipulates close rate as well as volume.
//
func (obv *obv) Execute(charts []Chart) bool {
	lasts := make([]float64, 0)
	volumes := make([]float64, 0)

	if len(charts) < obv.period+1 {
		logger.Get().Info("[obv:Execute] there is no chart duration")
		return false
	}
	for _, chart := range charts {
		lasts = append(lasts, chart.Last)
		volumes = append(volumes, chart.Volume)
	}

	obvs := talib.Obv(lasts, volumes)
	if len(obvs) > 1 {
		logger.Get().Infof("[obv:Execute] %v", obvs[len(obvs)-1])
	} else {
		logger.Get().Error("[obv:Execute] there is no obv")
		return false
	}
	if obv.positionType == 1 {
		if obvs[len(obvs)-1]/obvs[len(obvs)-2] >= obv.param {
			return true
		}
	} else {
		if obvs[len(obvs)-1]/obvs[len(obvs)-2] <= obv.param {
			return true
		}
	}
	return false
}
