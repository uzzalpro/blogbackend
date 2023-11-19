// routes/post_routes.go

package routes

import (
	"blogbackend/controllers"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupPostRoutes(router *gin.Engine, db *sql.DB) {
	postController := controllers.NewPostController(db)

	postRoutes := router.Group("/posts")
	{
		postRoutes.GET("", postController.GetPosts)
		postRoutes.GET("/:id", postController.GetPost)
		postRoutes.POST("/upload", postController.CreatePost)
		// Add additional routes as needed
	}
}
