package dto

import "erajaya-interview/entity"

type (
	ProductCreateRequest struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description" binding:"omitempty"`
		Price       float64 `json:"price" binding:"required,gt=0"`
		Quantity    int     `json:"quantity" binding:"required,min=0"`
	}

	ProductResponse struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Quantity    int     `json:"quantity"`
	}

	GetAllProductRepositoryResponse struct {
		Products []entity.Product `json:"products"`
		PaginationResponse
	}

	ProductPaginationResponse struct {
		Data []ProductResponse `json:"data"`
		PaginationResponse
	}
)
