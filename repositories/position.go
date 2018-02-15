package repositories

import (
	"github.com/airking05/go-auto-trader/logger"
	"github.com/airking05/go-auto-trader/models"
	"github.com/airking05/go-auto-trader/services"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func NewPosition(exchangeID models.ExchangeID, assetDistributionRate float64, ptrate float64, lcrate float64, ptype models.PositionType, trading string, settlement string, limitsec int) *models.Position {
	return &models.Position{
		ExchangeID:            exchangeID,
		AssetDistributionRate: assetDistributionRate,
		ProfitTakeRate:        ptrate,
		LossCutRate:           lcrate,
		PositionType:          ptype,
		Trading:               trading,
		Settlement:            settlement,

		WaitLimitSecond: limitsec,
	}
}

type positionStorage struct {
	db *gorm.DB
}

func NewPositionStorage(db *gorm.DB) services.PositionRepository {
	return &positionStorage{
		db: db,
	}
}

func (d *positionStorage) FindNUnclosedByTraderID(traderID uint) ([]models.Position, error) {
	var positions []models.Position

	if err := d.db.Where(&models.Position{TraderID: traderID, IsMade: true, IsClosed: false}).Find(&positions).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get positions by trader_id")
	}
	return positions, nil
}

func (d *positionStorage) Insert(position *models.Position, traderID uint) (uint, error) {
	position.TraderID = traderID
	if isNew := d.db.NewRecord(position); isNew == true {
		err := d.db.Create(&position).Error
		if err != nil {
			logger.Get().Error(err)
			return 0, err
		}
		return position.ID, nil
	}
	err := errors.New("failed to insert position to table")
	return 0, err
}

func (d *positionStorage) FindNByTraderID(traderID uint) ([]models.Position, error) {
	var positions []models.Position

	if err := d.db.Where(&models.Position{TraderID: traderID}).Order("created_at desc").Find(&positions).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get positions by trader_id")
	}
	return positions, nil
}

func (d *positionStorage) FindAll() ([]models.Position, error) {
	var positions []models.Position
	if err := d.db.Find(&positions).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get traderConfig list")
	}
	return positions, nil
}

func (t *positionStorage) UpdateToMade(positionID uint) error {
	if err := t.db.Model(&models.Position{ID: positionID}).UpdateColumn("is_made", true).Error; err != nil {
		return err
	}
	return nil
}
func (t *positionStorage) UpdateToClosed(positionID uint) error {
	if err := t.db.Model(&models.Position{ID: positionID}).UpdateColumn("is_closed", true).Error; err != nil {
		return err
	}
	return nil
}
func (t *positionStorage) UpdateEntryOrder(positionID uint, orderID uint) error {
	if err := t.db.Model(&models.Position{ID: positionID}).UpdateColumn("entry_order_id", orderID).Error; err != nil {
		return err
	}
	return nil
}
func (t *positionStorage) UpdateExitOrder(positionID uint, orderID uint) error {
	if err := t.db.Model(&models.Position{ID: positionID}).UpdateColumn("entry_close_id", orderID).Error; err != nil {
		return err
	}
	return nil
}

func (t *positionStorage) Truncate() error {
	return t.db.Exec("TRUNCATE TABLE positions").Error
}
