package provider

import (
	"erajaya-interview/config"
	"erajaya-interview/controller"
	"erajaya-interview/repository"
	"erajaya-interview/service"

	"github.com/samber/do"
	"gorm.io/gorm"
)

func ProvideUserDependencies(injector *do.Injector, db *gorm.DB, jwtService service.JWTService, redis config.Redis) {
	// Repository
	userRepository := repository.NewUserRepository(db, redis)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)

	// Service
	userService := service.NewUserService(userRepository, refreshTokenRepository, jwtService, db)

	// Controller
	do.Provide(
		injector, func(i *do.Injector) (controller.UserController, error) {
			return controller.NewUserController(userService), nil
		},
	)
}
