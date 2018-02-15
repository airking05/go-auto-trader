package repositories

import (
	"github.com/airking05/go-auto-trader/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type traderGorm struct {
	db *gorm.DB
}

func (t *traderGorm) Insert(traderConfig *models.TraderGorm) (uint, error) {
	if isNew := t.db.NewRecord(traderConfig); isNew == true {
		if err := t.db.Create(&traderConfig).Error; err != nil {

			return 0, err
		}
		return traderConfig.ID, nil
	}
	err := errors.New("failed to insert traderConfig")
	return 0, err
}

func (t *traderGorm) Find(traderID uint) (*models.TraderGorm, error) {

	traderGorm := &models.TraderGorm{}
	traderGorm.ID = traderID
	if err := t.db.Model(traderGorm).Find(&traderGorm).Error; err != nil {
		return traderGorm, errors.Wrap(err, "failed to get traderConfig by id")
	}
	return traderGorm, nil
}

func (t *traderGorm) FindNByStatus(status string, limit int, offset int) ([]models.TraderGorm, error) {
	var traderConfigs []models.TraderGorm

	if err := t.db.Where(&models.TraderGorm{Status: status}).Order("updated_at desc").Offset(offset).Limit(limit).Find(&traderConfigs).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get traderConfig by id")
	}

	return traderConfigs, nil
}

func (t *traderGorm) FindAllByStatus(status string) ([]models.TraderGorm, error) {
	var traderConfigs []models.TraderGorm

	if err := t.db.Where(&models.TraderGorm{Status: status}).Order("updated_at desc").Find(&traderConfigs).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get traderConfig by id")
	}
	return traderConfigs, nil
}

func (t *traderGorm) UpdateStatusByID(traderID uint, status string) error {
	traderGorm := &models.TraderGorm{}
	traderGorm.ID = traderID
	if err := t.db.Model(traderGorm).UpdateColumn("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (t *traderGorm) FindAll() ([]models.TraderGorm, error) {
	var traderConfigs []models.TraderGorm

	if err := t.db.Find(&traderConfigs).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get traderConfig list")
	}

	return traderConfigs, nil
}

func (t *traderGorm) Truncate() error {
	return t.db.Exec("TRUNCATE TABLE trader_configs").Error
}
