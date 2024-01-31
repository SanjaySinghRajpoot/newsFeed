package router

import (
	"github.com/SanjaySinghRajpoot/newsFeed/controllers"
	"github.com/SanjaySinghRajpoot/newsFeed/middleware"
	"github.com/gin-gonic/gin"
)

func GetRoute(r *gin.Engine) {
	// User routes
	r.POST("/api/signup", controllers.Signup)
	r.POST("/api/login", controllers.Login)

	r.Use(middleware.RequireAuth)
	r.POST("/api/logout", controllers.Logout)
	userRouter := r.Group("/api/users")
	{
		userRouter.GET("/", controllers.GetUsers)
		userRouter.GET("/:id", controllers.GetUser)
		userRouter.PUT("/:id/update", controllers.UpdateUser)
		userRouter.DELETE("/:id/delete", controllers.DeleteUser)
	}
}
