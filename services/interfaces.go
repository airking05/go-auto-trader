package services

import (
	"github.com/airking05/go-auto-trader/models"
	"time"
)

//go:generate mockery -name=ChartRepository
type ChartRepository interface {
	FindRange(exchangeId models.ExchangeID, duration int, trading string, settlement string, start time.Time, end time.Time) ([]models.Chart, error)
	FindN(exchangeId models.ExchangeID, duration int,
		trading string, settlement string, num int) ([]models.Chart, error)
	Find(exchangeId models.ExchangeID, duration int,
		trading string, settlement string) (models.Chart, error)
	Create(chart *models.Chart) error
	Truncate() error
}

//go:generate mockery -name=OrderRepository
type OrderRepository interface {
	Insert(orderData *models.OrderGorm) (uint, error)
	Find(orderID uint) (models.OrderGorm, error)
	FindNByPositionID(positionID uint) ([]models.OrderGorm, error)
	Truncate() error
}

//go:generate mockery -name=PositionRepository
type PositionRepository interface {
	Insert(position *models.Position, traderID uint) (uint, error)
	FindNByTraderID(traderID uint) ([]models.Position, error)
	FindNUnclosedByTraderID(traderID uint) ([]models.Position, error)
	FindAll() ([]models.Position, error)
	UpdateToMade(positionID uint) error
	UpdateToClosed(positionID uint) error
	UpdateEntryOrder(positionID uint, orderID uint) error
	UpdateExitOrder(positionID uint, orderID uint) error
	Truncate() error
}

//go:generate mockery -name=TraderRepository
type TraderRepository interface {
	Insert(traderConfig *models.TraderGorm) (uint, error)
	Find(traderConfigID uint) (models.TraderGorm, error)
	FindNByStatus(status string, limit int, offset int) ([]models.TraderGorm, error)
	FindAllByStatus(status string) ([]models.TraderGorm, error)
	UpdateStatusByID(traderID uint, status string) error
	FindAll() ([]models.TraderGorm, error)
	Truncate() error
}

//go:generate mockery -name=ExchangePrivateRepository
type ExchangePrivateRepository interface {
	PurchaseFeeRate() (float64, error)
	SellFeeRate() (float64, error)
	TransferFee() (map[string]float64, error)
	Balances() (map[string]float64, error)
	CompleteBalances() (map[string]*models.Balance, error)
	ActiveOrders() ([]*models.Order, error)
	Order(trading string, settlement string,
		ordertype models.OrderType, price float64, amount float64) (string, error)
	Transfer(typ string, addr string,
		amount float64, additionalFee float64) error
	CancelOrder(orderNumber string, productCode string) error
	Address(c string) (string, error)
}
