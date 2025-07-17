package migrations

import (
	"erajaya-interview/migrations/seeds"

	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seeds.ListUserSeeder(db); err != nil {
		return err
	}
	if err := seeds.ListProductSeeder(db); err != nil {
		return err
	}

	return nil
}
