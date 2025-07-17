package provider

import (
	"erajaya-interview/config"
	"erajaya-interview/constants"
	"erajaya-interview/service"
	"os"
	"strconv"

	"github.com/samber/do"
	"gorm.io/gorm"
)

func InitDatabase(injector *do.Injector) {
	do.ProvideNamed(injector, constants.DB, func(i *do.Injector) (*gorm.DB, error) {
		return config.SetUpDatabaseConnection(), nil
	})
}

func RegisterDependencies(injector *do.Injector) {
	InitDatabase(injector)
	InitRedis(injector)

	do.ProvideNamed(injector, constants.JWTService, func(i *do.Injector) (service.JWTService, error) {
		return service.NewJWTService(), nil
	})

	// Initialize
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	redis := do.MustInvokeNamed[config.Redis](injector, constants.RedisClient)

	// Provide Dependencies
	ProvideUserDependencies(injector, db, jwtService, redis)
	ProvideProductDependencies(injector, db, jwtService, redis)
}

func InitRedis(injector *do.Injector) {
	do.ProvideNamed(injector, constants.RedisClient, func(i *do.Injector) (config.Redis, error) {
		addr := os.Getenv("REDIS_ADDRESS")
		password := os.Getenv("REDIS_PASSWORD")
		db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
		return config.NewRedisClient(addr, password, db), nil
	})
}
