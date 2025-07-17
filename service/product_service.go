package service

import (
	"context"

	"gorm.io/gorm"

	"erajaya-interview/dto"
	"erajaya-interview/entity"
	"erajaya-interview/repository"
)

type (
	ProductService interface {
		CreateProduct(ctx context.Context, req dto.ProductCreateRequest) (dto.ProductResponse, error)
		GetAllProducts(ctx context.Context, req dto.PaginationRequest) (dto.ProductPaginationResponse, error)
	}

	productService struct {
		productRepo      repository.ProductRepository
		refreshTokenRepo repository.RefreshTokenRepository
		jwtService       JWTService
		db               *gorm.DB
	}
)

func NewProductService(
	productRepo repository.ProductRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtService JWTService,
	db *gorm.DB,
) ProductService {
	return &productService{
		productRepo:      productRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		db:               db,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req dto.ProductCreateRequest) (dto.ProductResponse, error) {
	product := entity.Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Quantity:    req.Quantity,
	}

	newProduct, err := s.productRepo.CreateProduct(ctx, nil, product)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return dto.ProductResponse{
		ID:          newProduct.ID.String(),
		Name:        newProduct.Name,
		Price:       newProduct.Price,
		Description: newProduct.Description,
		Quantity:    newProduct.Quantity,
	}, nil
}

func (s *productService) GetAllProducts(
	ctx context.Context,
	req dto.PaginationRequest,
) (dto.ProductPaginationResponse, error) {
	dataWithPaginate, err := s.productRepo.GetAllProducts(ctx, nil, req)
	if err != nil {
		return dto.ProductPaginationResponse{}, err
	}

	var datas []dto.ProductResponse
	for _, product := range dataWithPaginate.Products {
		data := dto.ProductResponse{
			ID:          product.ID.String(),
			Name:        product.Name,
			Price:       product.Price,
			Description: product.Description,
			Quantity:    product.Quantity,
		}

		datas = append(datas, data)
	}

	return dto.ProductPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			Limit:   dataWithPaginate.Limit,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}
