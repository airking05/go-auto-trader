package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type OrderType int

const (
	Ask OrderType = iota
	Bid
)

type Order struct {
	ExchangeOrderID         string
	Type       OrderType
	Trading    string
	Settlement string
	Price      float64
	Amount     float64
}
type OrderGorm struct {
	gorm.Model
	Order

	ExchangeID ExchangeID `json:"exchange_id"`
	TraderID   uint        `json:"trader_id"`
	PositionID uint
	Pair       string     `json:"currency_pair"`
	OrderID    string     `json:"order_id"`

	Datetime  time.Time `json:"datetime"`
	Status    bool      `json:"status"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	TradeType PositionType `json:"trade_type"`
}