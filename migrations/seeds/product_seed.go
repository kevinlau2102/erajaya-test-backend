package seeds

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"erajaya-interview/entity"

	"gorm.io/gorm"
)

func ListProductSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("./migrations/json/products.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var listProduct []entity.Product
	if err := json.Unmarshal(jsonData, &listProduct); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.Product{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.Product{}); err != nil {
			return err
		}
	}

	for _, data := range listProduct {
		var product entity.Product
		err := db.Where("name = ?", data.Name).First(&product).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
