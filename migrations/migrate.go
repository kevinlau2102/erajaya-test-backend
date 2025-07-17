package migrations

import (
	"erajaya-interview/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.RefreshToken{},
		&entity.Product{},
	); err != nil {
		return err
	}

	return nil
}
