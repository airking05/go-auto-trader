package models

import (
	"time"
)

// Long:true False:false
type PositionType int

const (
	Long  PositionType = 1
	Short              = 0
)

// A position is a struct used by Trader.
// It opens and closes positions according to evaluate methods it has.
type Position struct {
	ID                    uint `gorm:"primary_key" json:"position_id"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
	ExchangeID            ExchangeID
	AssetDistributionRate float64
	EntryPrice            float64 `json:"entry_price"`
	ProfitTakeRate        float64
	LossCutRate           float64
	PositionType          PositionType `json:"position_type" sql:"DEFAULT:false"`
	Trading               string
	Settlement            string
	WaitLimitSecond       int
	EntryOrder            OrderGorm
	EntryOrderID          uint `gorm:"type:bigint" json:"entry_order_id"`
	ExitOrder             OrderGorm
	ExitOrderID           uint `gorm:"type:bigint" json:"exit_order_id"`
	IsMade                bool `json:"is_made"`
	IsClosed              bool `json:"is_closed"`
	TraderID              uint `json:"trader_id"`
}

func (p *Position) IsProfitTakable(rate float64) bool {
	if p.PositionType == 1 {
		profitTakePrice := p.EntryPrice * (1 + p.ProfitTakeRate)
		if rate >= profitTakePrice {
			return true
		}
	} else {
		profitTakePrice := p.EntryPrice * (1 - p.ProfitTakeRate)
		if rate <= profitTakePrice {
			return true
		}
	}
	return false
}

func (p *Position) IsLossCuttable(rate float64) bool {
	if p.PositionType == 1 {
		lossCutPrice := p.EntryPrice * (1 - p.LossCutRate)
		if rate <= lossCutPrice {
			return true
		}
	} else {
		lossCutPrice := p.EntryPrice * (1 + p.LossCutRate)
		if rate >= lossCutPrice {
			return true
		}
	}
	return false
}
func (p *Position) SetPrice(price float64) {
	p.EntryPrice = price
}
