package routes

import (
	"erajaya-interview/constants"
	"erajaya-interview/controller"
	"erajaya-interview/middleware"
	"erajaya-interview/service"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func Product(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	productController := do.MustInvoke[controller.ProductController](injector)

	routes := route.Group("/api/v1/product")
	{
		// User
		routes.POST("", middleware.Authenticate(jwtService), productController.CreateProduct)
		routes.GET("", middleware.Authenticate(jwtService), productController.GetAllProducts)
	}
}
