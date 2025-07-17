package provider

import (
	"erajaya-interview/config"
	"erajaya-interview/controller"
	"erajaya-interview/repository"
	"erajaya-interview/service"

	"github.com/samber/do"
	"gorm.io/gorm"
)

func ProvideProductDependencies(injector *do.Injector, db *gorm.DB, jwtService service.JWTService, redis config.Redis) {
	// Repository
	productRepository := repository.NewProductRepository(db, redis)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)

	// Service
	productService := service.NewProductService(productRepository, refreshTokenRepository, jwtService, db)

	// Controller
	do.Provide(
		injector, func(i *do.Injector) (controller.ProductController, error) {
			return controller.NewProductController(productService), nil
		},
	)
}
