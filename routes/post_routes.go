// routes/post_routes.go

package routes

import (
	"github.com/gin-gonic/gin"
	"blogbackend/controllers" // Make sure to import your controllers package
	"gorm.io/gorm"             // Import GORM if you are using it
)

// SetupPostRoutes sets up routes related to posts
func SetupPostRoutes(router *gin.Engine, db *gorm.DB) {
	postController := controllers.NewPostController(db)

	postRoutes := router.Group("/posts")
	{
		postRoutes.GET("", postController.GetPosts)
		postRoutes.GET("/:id", postController.GetPost)
		postRoutes.POST("/upload", postController.CreatePost)
		// Add additional routes as needed
	}
}
