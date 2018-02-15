package repositories

import (
	"github.com/airking05/go-auto-trader/logger"
	"github.com/airking05/go-auto-trader/models"
	"github.com/airking05/go-auto-trader/services"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type orderDataStorage struct {
	db *gorm.DB
}

func NewOrderGormStorage(db *gorm.DB) services.OrderRepository {
	return &orderDataStorage{
		db: db,
	}
}

type OrderGormRepository interface {
	Find(orderID uint) (models.OrderGorm, error)
	FindNByPositionID(positionID uint) ([]models.OrderGorm, error)
	Insert(*models.OrderGorm) (uint, error)
	Truncate() error
}

func (d *orderDataStorage) Insert(orderData *models.OrderGorm) (uint, error) {
	if isNew := d.db.NewRecord(orderData); isNew == true {
		err := d.db.Create(&orderData).Error
		if err != nil {
			logger.Get().Error(err)
			return 0, err
		}
		return orderData.ID, nil
	}
	err := errors.New("failed to insert position to table")
	return 0, err
}

func (d *orderDataStorage) Find(orderID uint) (models.OrderGorm, error) {
	var orderData models.OrderGorm
	if err := d.db.Where(orderID).Find(&orderData).Error; err != nil {
		return orderData, errors.Wrap(err, "failed to get order_data")
	}
	return orderData, nil
}

func (d *orderDataStorage) FindNByPositionID(positionID uint) ([]models.OrderGorm, error) {
	var orderDatas []models.OrderGorm
	if err := d.db.Where(&models.OrderGorm{PositionID: positionID}).Find(&orderDatas).Error; err != nil {
		return orderDatas, errors.Wrap(err, "failed to get order_data by position id")
	}
	return orderDatas, nil
}

func (d *orderDataStorage) Truncate() error {
	return d.db.Exec("TRUNCATE TABLE order_data").Error
}
