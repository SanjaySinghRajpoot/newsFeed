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
	userRouter.GET("/", controllers.GetUsers)
	userRouter.GET("/:id", controllers.GetUser)
	userRouter.PUT("/:id/update", controllers.UpdateUser)
	userRouter.DELETE("/:id/delete", controllers.DeleteUser)

	// Post routes
	postRouter := r.Group("/api/posts")
	postRouter.GET("/", controllers.GetPosts)
	postRouter.GET("/:id/show", controllers.ShowPost)
	postRouter.POST("/create", controllers.CreatePost)
	postRouter.PUT("/:id/update", controllers.UpdatePost)
	postRouter.DELETE("/:id/delete", controllers.DeletePost)

	// Friends Route
	friendRouter := r.Group("/api/friend")
	friendRouter.GET("/", controllers.GetFriends)
	friendRouter.POST("/:following_user_id/follow", controllers.FollowRequest)
	friendRouter.DELETE("/:following_user_id/unfollow", controllers.UnfollowRequest)

	// Comment routes
	commentRouter := r.Group("/api/posts/comment")
	commentRouter.POST("/add", controllers.CommentOnPost)
	commentRouter.GET("/:comment_id", controllers.GetComment)
	commentRouter.PUT("/:comment_id/update", controllers.UpdateComment)
	commentRouter.DELETE("/:comment_id/delete", controllers.DeleteComment)

	r.GET("/api/newsfeed", controllers.GetNewsFeed)
}
