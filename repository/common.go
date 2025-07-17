package repository

import (
	"math"

	"erajaya-interview/dto"

	"gorm.io/gorm"
)

func Paginate(req dto.PaginationRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (req.Page - 1) * req.Limit
		return db.Offset(offset).Limit(req.Limit)
	}
}

func TotalPage(count, perPage int64) int64 {
	totalPage := int64(math.Ceil(float64(count) / float64(perPage)))

	return totalPage
}
