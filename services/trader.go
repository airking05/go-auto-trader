package services

import (
	"context"
	"sync"
	"time"

	"github.com/airking05/go-auto-trader/logger"
	"github.com/airking05/go-auto-trader/models"
	"github.com/pkg/errors"
)

type TraderBot interface {
	Execute(context.Context, float64) error
	MakePosition(float64, float64) (int, error)
	ClosePosition(float64, float64) (int, error)
	JudgeEntry() bool
	SetPositionMade(uint, int) error
	SetPositionClosed(uint, int, float64) error
	Rate() (float64, error)
	PositionStart() (uint, error)
	GetPhase() models.TradePhase
	SetPosition(*models.Position)
}

type traderBot struct {
	ChartRepository           ChartRepository    `inject:""`
	PositionRepository        PositionRepository `inject:""`
	OrderRepository           OrderRepository    `inject:""`
	ExchangePrivateRepository ExchangePrivateRepository
	Position                  *models.Position
	Phase                     models.TradePhase
	traderConfig              *models.TraderGorm
	entrylogics               models.Logic
}

type TraderServiceImpl struct {
	ChartRepository    ChartRepository    `inject:""`
	PositionRepository PositionRepository `inject:""`
	OrderRepository    OrderRepository    `inject:""`
	TraderRepository   TraderRepository   `inject:""`
	BotMap             BotMap
}

func NewBotMap() BotMap {
	return BotMap{lock: &sync.RWMutex{}, m: make(map[int]context.CancelFunc)}
}

type BotMap struct {
	lock *sync.RWMutex
	m    map[int]context.CancelFunc
}

func (bm *BotMap) Add(id int, cancelFunc context.CancelFunc) {
	bm.lock.Lock()
	defer bm.lock.Unlock()
	bm.m[id] = cancelFunc
}

func (bm *BotMap) Cancel(id int) {
	bm.lock.Lock()
	defer bm.lock.Unlock()
	bm.m[id]()
	delete(bm.m, id)
}

func (bm *BotMap) Delete(id int) {
	bm.lock.Lock()
	defer bm.lock.Unlock()
	delete(bm.m, id)
}

func (t *TraderServiceImpl) Run(traderGorm *models.TraderGorm) error {
	traderId, err := t.TraderRepository.Insert(traderGorm)
	if err != nil {
		return err
	}
	exchangePrivateRepository, err := NewExchangePrivateRepository(traderGorm.ExchangeID, traderGorm.APIKey, traderGorm.SecretKey)
	if err != nil {
		return err
	}
	bot := &traderBot{
		ChartRepository:           t.ChartRepository,
		PositionRepository:        t.PositionRepository,
		OrderRepository:           t.OrderRepository,
		ExchangePrivateRepository: exchangePrivateRepository,
		Phase:        models.Boot,
		traderConfig: traderGorm,
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	t.BotMap.Add(int(traderId), cancelFunc)
	go func() {
		logger.Get().Errorf("%v", bot.Execute(ctx, 1))
		t.BotMap.Delete(int(traderId))
	}()
	return nil
}

func (t *traderBot) JudgeEntry() bool {
	charts, err := t.ChartRepository.FindN(t.traderConfig.ExchangeID, t.traderConfig.Duration, t.traderConfig.Trading, t.traderConfig.Settlement, 50)
	if err != nil {
		return false
	}
	return t.entrylogics.Execute(charts)
}

func (t *traderBot) SetPositionMade(positionID uint, orderID string) error {
	if err := t.PositionRepository.UpdateToMade(positionID); err != nil {
		return err
	}
	var orderType models.OrderType
	if t.traderConfig.PositionType == models.Long {
		orderType = models.Ask
	} else {
		orderType = models.Bid
	}
	orderData := &models.OrderGorm{
		Order: models.Order{
			ExchangeOrderID: orderID,
			Type:            orderType,
			Trading:         t.traderConfig.Trading,
			Settlement:      t.traderConfig.Settlement},
		ExchangeID:  t.traderConfig.ExchangeID,
		TraderID:    (t.traderConfig.ID),
		PositionID:  positionID,
		Status:      true,
		ExcutePrice: t.Position.EntryPrice,
	}
	orderLocalID, err := t.OrderRepository.Insert(orderData)
	if err != nil {
		return err
	}
	if err = t.PositionRepository.UpdateEntryOrder(positionID, orderLocalID); err != nil {
		return err
	}
	return nil
}

func (t *traderBot) SetPositionClosed(positionID uint, orderID string, price float64) error {
	if err := t.PositionRepository.UpdateToClosed(positionID); err != nil {
		return err
	}
	var orderType models.OrderType
	if t.traderConfig.PositionType == models.Long {
		orderType = models.Bid
	} else {
		orderType = models.Ask
	}
	orderData := &models.OrderGorm{
		Order: models.Order{
			ExchangeOrderID: orderID,
			Type:            orderType,
			Trading:         t.traderConfig.Trading,
			Settlement:      t.traderConfig.Settlement},
		ExchangeID:  t.traderConfig.ExchangeID,
		TraderID:    t.traderConfig.ID,
		PositionID:  positionID,
		Status:      true,
		ExcutePrice: price,
	}

	orderLocalID, err := t.OrderRepository.Insert(orderData)
	if err != nil {
		return err
	}

	if err = t.PositionRepository.UpdateExitOrder(positionID, orderLocalID); err != nil {
		return err
	}
	return nil
}

func (t *traderBot) MakePosition(price float64, speedScale float64) (string, error) {
	orderType := models.Bid
	if t.Position.PositionType == models.Long {
		orderType = models.Ask
	}
	balances, err := t.ExchangePrivateRepository.Balances()
	if err != nil {
		return "", errors.Wrap(err, "failed to get balances")
	}
	total := balances[t.Position.Settlement] * t.Position.AssetDistributionRate
	amount := total / price
	t.traderConfig.TradingAmount = amount

	orderExchangeID, err := t.ExchangePrivateRepository.Order(t.Position.Trading, t.Position.Settlement, orderType, price, amount)
	if err != nil {
		logger.Get().Info("order failed:%s", err)
		return "", errors.Wrap(err, "failed to order!")
	}
	elapsedSecond := 0
	for {
		isContracted := true
		activeOrders, err := t.ExchangePrivateRepository.ActiveOrders()
		if err != nil {
			isContracted = false
		} else {
			for _, activeOrder := range activeOrders {
				if activeOrder.ExchangeOrderID == orderExchangeID {
					isContracted = false
				}
			}
		}
		if isContracted == true {
			t.Position.EntryPrice = price
			return orderExchangeID, nil
		}
		if elapsedSecond >= t.Position.WaitLimitSecond {
			err := t.ExchangePrivateRepository.CancelOrder(orderExchangeID, t.Position.Trading+"_"+t.Position.Settlement)
			if err == nil {
				return "", errors.New("failed to complete order")
			}
		}
		time.Sleep(time.Duration(10*speedScale) * time.Second)
		elapsedSecond += 10
	}
}

func (t *traderBot) ClosePosition(price float64, speedScale float64) (string, error) {
	orderType := models.Ask
	if t.Position.PositionType == models.Long {
		orderType = models.Bid
	}

	amount := t.traderConfig.TradingAmount * (1 - 0.0025)
	orderExchangeID, err := t.ExchangePrivateRepository.Order(t.Position.Trading, t.Position.Settlement, orderType, price, amount)
	if err != nil {
		logger.Get().Info("order failed:%s", err)
		return "", errors.Wrap(err, "failed to order!")
	}
	elapsedSecond := 0
	for {
		isContracted := true
		activeOrders, err := t.ExchangePrivateRepository.ActiveOrders()
		if err != nil {
			isContracted = false
		} else {
			for _, activeOrder := range activeOrders {
				if activeOrder.ExchangeOrderID == orderExchangeID {
					isContracted = false
				}
			}
		}
		if isContracted == true {
			return orderExchangeID, nil
		}
		if elapsedSecond >= t.Position.WaitLimitSecond {
			err := t.ExchangePrivateRepository.CancelOrder(orderExchangeID, t.Position.Trading+"_"+t.Position.Settlement)
			if err == nil {
				return "", errors.New("failed to complete order")
			}
		}
		time.Sleep(time.Duration(10*speedScale) * time.Second)
		elapsedSecond += 10
	}
}

func (t *traderBot) judgeEntryLoop(ctx context.Context, speedScale float64) {
	timer := time.NewTicker(time.Millisecond * time.Duration(1000*10*speedScale))

	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			entry := t.JudgeEntry()
			if entry {
				t.Phase = models.MakePosition
				return
			}
		case <-ctx.Done():
			t.Phase = models.ShutDown
			return
		}
	}
}

func (t *traderBot) makePositionLoop(ctx context.Context, positionID uint, speedScale float64) error {
	chart, err := t.ChartRepository.Find(t.traderConfig.ExchangeID, t.traderConfig.Duration, t.traderConfig.Trading, t.traderConfig.Settlement)
	if err != nil {
		return err
	}
	orderID, err := t.MakePosition(chart.Last, speedScale)
	if err != nil {
		t.Phase = models.JudgeEntry
		return err
	}
	if err := t.SetPositionMade(positionID, orderID); err != nil {
		return err
	}

	t.Position.SetPrice(chart.Last)
	t.Phase = models.JudgeExit

	return nil
}

func (t *traderBot) judgeExitLoop(ctx context.Context, speedScale float64) error {

	timer := time.NewTicker(time.Millisecond * time.Duration(1000*10*speedScale))
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			chart, err := t.ChartRepository.Find(t.traderConfig.ExchangeID, t.traderConfig.Duration, t.traderConfig.Trading, t.traderConfig.Settlement)
			if err != nil {
				return err
			}

			if t.Position.IsProfitTakable(chart.Last) || t.Position.IsLossCuttable(chart.Last) {
				t.Phase = models.ClosePosition
				return nil
			}
		case <-ctx.Done():
			t.Phase = models.ShutDown
			return nil
		}
	}
}

func (t *traderBot) closePositionLoop(ctx context.Context, positionID uint, speedScale float64) error {
	chart, err := t.ChartRepository.Find(t.traderConfig.ExchangeID, t.traderConfig.Duration, t.traderConfig.Trading, t.traderConfig.Settlement)
	if err != nil {
		return err
	}

	orderID, err := t.ClosePosition(chart.Last, speedScale)
	if err != nil {
		return err
	}

	if err := t.SetPositionClosed(positionID, orderID, chart.Last); err != nil {
		t.Phase = models.JudgeExit
		return err
	}
	t.Phase = models.Reset
	return nil
}

func (t *traderBot) SetPosition(position *models.Position) {
	t.Position = position
}

func (t *traderBot) PositionStart() (uint, error) {
	position := &models.Position{
		ExchangeID:            t.traderConfig.ExchangeID,
		AssetDistributionRate: t.traderConfig.AssetDistributionRate,
		ProfitTakeRate:        t.traderConfig.ProfitTakeRate,
		LossCutRate:           t.traderConfig.LossCutRate,
		PositionType:          t.traderConfig.PositionType,
		Trading:               t.traderConfig.Trading,
		Settlement:            t.traderConfig.Settlement,
		WaitLimitSecond:       t.traderConfig.WaitLimitSecond,
	}
	positionID, err := t.PositionRepository.Insert(position, t.traderConfig.ID)
	if err != nil {
		return 0, err
	}
	t.Position = position
	return positionID, nil
}

func (t *traderBot) GetPhase() models.TradePhase {
	return t.Phase
}

func (t *traderBot) Execute(ctx context.Context, speedScale float64) error {
	for {
		switch t.Phase {
		case models.Boot:
			logger.Get().Info("boot trader")
			_, err := t.PositionStart()
			if err != nil {
				return err
			}
			t.Phase = models.JudgeEntry
		case models.Resume:
			logger.Get().Info("resume trader")
			t.Phase = models.JudgeExit
		case models.Reset:
			logger.Get().Info("reboot trader")
			_, err := t.PositionStart()
			if err != nil {
				return err
			}
			t.Phase = models.JudgeEntry

		case models.JudgeEntry:
			logger.Get().Info("judging entry")
			t.judgeEntryLoop(ctx, speedScale)

		case models.MakePosition:
			logger.Get().Info("make position")
			if err := t.makePositionLoop(ctx, t.Position.ID, speedScale); err != nil {
				logger.Get().Warn(err)
			}
		case models.JudgeExit:
			logger.Get().Info("judging exit")
			if err := t.judgeExitLoop(ctx, speedScale); err != nil {
				logger.Get().Warn(err)
				return err
			}

		case models.ClosePosition:
			logger.Get().Info("close position")
			if err := t.closePositionLoop(ctx, t.Position.ID, speedScale); err != nil {
				logger.Get().Warn(err)
			}
		case models.ShutDown:
			logger.Get().Info("teardown trader")
			return nil
		}
	}
}
