package routes

import (
	"erajaya-interview/constants"
	"erajaya-interview/controller"
	"erajaya-interview/middleware"
	"erajaya-interview/service"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func User(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	userController := do.MustInvoke[controller.UserController](injector)

	routes := route.Group("/api/v1/user")
	{
		// User
		routes.POST("", userController.Register)
		routes.GET("", userController.GetAllUser)
		routes.POST("/login", userController.Login)
		routes.POST("/refresh", userController.Refresh)
		routes.DELETE("", middleware.Authenticate(jwtService), userController.Delete)
		routes.PATCH("", middleware.Authenticate(jwtService), userController.Update)
		routes.GET("/me", middleware.Authenticate(jwtService), userController.Me)
	}
}
