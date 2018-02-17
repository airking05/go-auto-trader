package models

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var (
	modelList = []interface{}{
		&OrderGorm{}, &TraderGorm{}, &Position{},
	}
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(modelList...).Error; err != nil {
		return errors.Wrap(err, "failed to migrate models")
	}
	return nil
}
