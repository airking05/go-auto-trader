package repositories

import (
	"github.com/airking05/go-auto-trader/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type TraderGormRepository struct {
	DB *gorm.DB `inject:""`
}

func (t *TraderGormRepository) Insert(traderConfig *models.TraderGorm) (uint, error) {
	if isNew := t.DB.NewRecord(traderConfig); isNew == true {
		if err := t.DB.Create(&traderConfig).Error; err != nil {

			return 0, err
		}
		return traderConfig.ID, nil
	}
	err := errors.New("failed to insert traderConfig")
	return 0, err
}

func (t *TraderGormRepository) Find(traderID uint) (*models.TraderGorm, error) {

	traderGorm := &models.TraderGorm{}
	traderGorm.ID = traderID
	if err := t.DB.Model(traderGorm).Find(&traderGorm).Error; err != nil {
		return traderGorm, errors.Wrap(err, "failed to get traderConfig by id")
	}
	return traderGorm, nil
}

func (t *TraderGormRepository) FindNByStatus(status string, limit int, offset int) ([]models.TraderGorm, error) {
	var traderConfigs []models.TraderGorm

	if err := t.DB.Where(&models.TraderGorm{Status: status}).Order("updated_at desc").Offset(offset).Limit(limit).Find(&traderConfigs).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get traderConfig by id")
	}

	return traderConfigs, nil
}

func (t *TraderGormRepository) FindAllByStatus(status string) ([]models.TraderGorm, error) {
	var traderConfigs []models.TraderGorm

	if err := t.DB.Where(&models.TraderGorm{Status: status}).Order("updated_at desc").Find(&traderConfigs).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get traderConfig by id")
	}
	return traderConfigs, nil
}

func (t *TraderGormRepository) UpdateStatusByID(traderID uint, status string) error {
	traderGorm := &models.TraderGorm{}
	traderGorm.ID = traderID
	if err := t.DB.Model(traderGorm).UpdateColumn("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (t *TraderGormRepository) FindAll() ([]models.TraderGorm, error) {
	var traderConfigs []models.TraderGorm

	if err := t.DB.Find(&traderConfigs).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get traderConfig list")
	}

	return traderConfigs, nil
}

func (t *TraderGormRepository) Truncate() error {
	return t.DB.Exec("TRUNCATE TABLE trader_configs").Error
}
