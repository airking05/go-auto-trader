package repositories

import (
	"time"

	"github.com/airking05/go-auto-trader/models"
	"github.com/airking05/go-auto-trader/services"
	"github.com/jinzhu/gorm"
)

func NewChartRepositoryGorm(db *gorm.DB) services.ChartRepository {
	return &ChartRepositoryGorm{
		DB: db,
	}
}

type ChartRepositoryGorm struct {
	DB *gorm.DB `inject:""`
}

func (r *ChartRepositoryGorm) Truncate() error {
	if err := r.DB.Exec("truncate table charts").Error; err != nil {
		return err
	}
	return nil
}

func (r *ChartRepositoryGorm) Create(chart *models.Chart) error {
	if err := r.DB.Create(chart).Error; err != nil {
		return err
	}
	return nil
}

func (r *ChartRepositoryGorm) FindRange(exchangeId models.ExchangeID, duration int,
	trading string, settlement string, start time.Time, end time.Time) ([]models.Chart, error) {
	charts := make([]models.Chart, 0)

	db := r.DB.Raw(`
SELECT * FROM (
  SELECT * FROM charts WHERE (? <= datetime AND datetime <= ?)
) AS S
WHERE
  S.duration = ?
  AND S.exchange_id = ?
  AND S.pair = ?
ORDER BY S.datetime desc`, start, end, duration, int(exchangeId), trading+"_"+settlement).Scan(&charts)
	if db.Error != nil {
		return nil, db.Error
	}

	return charts, nil
}

func (r *ChartRepositoryGorm) FindN(exchangeId models.ExchangeID, duration int,
	trading string, settlement string, num int) ([]models.Chart, error) {
	charts := make([]models.Chart, 0)

	db := r.DB.Raw(`
SELECT * FROM charts
WHERE
  duration = ?
  AND exchange_id = ?
  AND pair = ?
ORDER BY datetime desc limit ?`, duration, int(exchangeId), trading+"_"+settlement, num).Scan(&charts)
	if db.Error != nil {
		return nil, db.Error
	}

	return charts, nil
}

func (r *ChartRepositoryGorm) Find(exchangeId models.ExchangeID, duration int,
	trading string, settlement string) (models.Chart, error) {
	chart := models.Chart{}

	db := r.DB.Raw(`
SELECT * FROM charts
WHERE
  duration = ?
  AND exchange_id = ?
  AND pair = ?
ORDER BY datetime desc limit 1`, duration, int(exchangeId), trading+"_"+settlement).Scan(&chart)
	if db.Error != nil {
		return chart, db.Error
	}

	return chart, nil
}
