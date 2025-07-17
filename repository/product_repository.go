package repository

import (
	"context"
	"encoding/json"
	"erajaya-interview/config"
	"erajaya-interview/dto"
	"erajaya-interview/entity"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	SORT_BY_CREATED_AT = "created_at"
	SORT_BY_PRICE      = "price"
	SORT_BY_NAME       = "name"

	ORDER_ASC  = "asc"
	ORDER_DESC = "desc"
)

var (
	productListCacheKeyPattern = "product:list:search=%s:sortBy=%s:order=%s:page=%d:limit=%d"
)

type (
	ProductRepository interface {
		CreateProduct(ctx context.Context, tx *gorm.DB, user entity.Product) (entity.Product, error)
		GetAllProducts(
			ctx context.Context,
			tx *gorm.DB,
			req dto.PaginationRequest,
		) (dto.GetAllProductRepositoryResponse, error)
	}

	productRepository struct {
		db    *gorm.DB
		redis config.Redis
	}
)

func NewProductRepository(db *gorm.DB, redis config.Redis) ProductRepository {
	return &productRepository{
		db:    db,
		redis: redis,
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, tx *gorm.DB, product entity.Product) (entity.Product, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&product).Error; err != nil {
		return entity.Product{}, err
	}

	// Invalidate Redis cache for product list
	if err := r.invalidateProductListCache(ctx); err != nil {
		fmt.Printf("Failed to invalidate cache: %v\n", err)
	}

	return product, nil
}

func (r *productRepository) GetAllProducts(
	ctx context.Context,
	tx *gorm.DB,
	req dto.PaginationRequest,
) (dto.GetAllProductRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var products []entity.Product
	var count int64

	req.Default()

	cacheKey := fmt.Sprintf(productListCacheKeyPattern,
		req.Search, req.SortBy, req.Order, req.Page, req.Limit,
	)

	// Check Redis cache first
	cachedData, err := r.redis.Get(ctx, cacheKey)
	if err == nil && cachedData != "" {
		var cachedResponse dto.GetAllProductRepositoryResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err == nil {
			return cachedResponse, nil
		}
	}

	query := tx.WithContext(ctx).Model(&entity.Product{})

	// Filtering
	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	// Count total data
	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllProductRepositoryResponse{}, err
	}

	// Sorting
	sortBy := SORT_BY_CREATED_AT
	if req.SortBy != "" {
		switch req.SortBy {
		case SORT_BY_PRICE, SORT_BY_NAME, SORT_BY_CREATED_AT:
			sortBy = req.SortBy
		}
	}

	order := ORDER_DESC
	if req.Order == ORDER_ASC {
		order = ORDER_ASC
	}

	query = query.Order(fmt.Sprintf("%s %s", sortBy, order))

	// Pagination
	if err := query.Scopes(Paginate(req)).Find(&products).Error; err != nil {
		return dto.GetAllProductRepositoryResponse{}, err
	}

	totalPage := TotalPage(count, int64(req.Limit))
	result := dto.GetAllProductRepositoryResponse{
		Products: products,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			Limit:   req.Limit,
			Count:   count,
			MaxPage: totalPage,
		},
	}

	// Store result to Redis
	bytes, err := json.Marshal(result)
	if err == nil {
		err = r.redis.Set(ctx, cacheKey, string(bytes), 5*time.Minute)
		if err != nil {
			fmt.Printf("Error setting cache: %v\n", err)
		}
	}

	return result, nil
}

func (r *productRepository) invalidateProductListCache(ctx context.Context) error {
	prefix := "product:list:"
	var cursor uint64
	for {
		keys, nextCursor, err := r.redis.Scan(ctx, cursor, prefix+"*", 100)
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := r.redis.Del(ctx, keys...); err != nil {
				return err
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return nil
}
