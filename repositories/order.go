package repositories

import (
	"github.com/airking05/go-auto-trader/logger"
	"github.com/airking05/go-auto-trader/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type OrderDataStorage struct {
	DB *gorm.DB `inject:""`
}

type OrderGormRepository interface {
	Find(orderID uint) (*models.OrderGorm, error)
	FindNByPositionID(positionID uint) ([]models.OrderGorm, error)
	Insert(*models.OrderGorm) (uint, error)
	Truncate() error
}

func (d *OrderDataStorage) Insert(orderData *models.OrderGorm) (uint, error) {
	if isNew := d.DB.NewRecord(orderData); isNew == true {
		err := d.DB.Create(orderData).Error
		if err != nil {
			logger.Get().Error(err)
			return 0, err
		}
		return orderData.ID, nil
	}
	err := errors.New("failed to insert position to table")
	return 0, err
}

func (d *OrderDataStorage) Find(orderID uint) (*models.OrderGorm, error) {
	var orderData *models.OrderGorm
	if err := d.DB.Where(orderID).Find(orderData).Error; err != nil {
		return orderData, errors.Wrap(err, "failed to get order_data")
	}
	return orderData, nil
}

func (d *OrderDataStorage) FindNByPositionID(positionID uint) ([]*models.OrderGorm, error) {
	var orderDatas []*models.OrderGorm
	if err := d.DB.Where(&models.OrderGorm{PositionID: positionID}).Find(orderDatas).Error; err != nil {
		return orderDatas, errors.Wrap(err, "failed to get order_data by position id")
	}
	return orderDatas, nil
}

func (d *OrderDataStorage) Truncate() error {
	return d.DB.Exec("TRUNCATE TABLE order_data").Error
}
