// controllers/post_controller.go

package controllers

import (
    "blogbackend/models"
    "blogbackend/views"
    "database/sql"
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type PostController struct {
	DB *gorm.DB
}

func NewPostController(db *gorm.DB) *PostController {
	return &PostController{DB: db}
}

func (pc *PostController) GetPosts(c *gin.Context) {
	var posts []models.Post
	if err := pc.DB.Order("created_at DESC").Find(&posts).Error; err != nil {
		views.JSON(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	views.JSON(c, http.StatusOK, gin.H{"posts": posts})
}

func (pc *PostController) GetPost(c *gin.Context) {
	var post models.Post
	id := c.Param("id")
	if err := pc.DB.Where("id = ?", id).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			views.JSON(c, http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		views.JSON(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	views.JSON(c, http.StatusOK, post)
}

func (pc *PostController) CreatePost(c *gin.Context) {
    var post models.Post
    if err := c.ShouldBindJSON(&post); err != nil {
        views.JSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Assuming your Image field is a string representing the file path
    // You might need to adjust this based on your actual data structure
    uploadedFileName := fmt.Sprintf("uploads/%s", post.Image)

    // Create a sql.NullString instance
    nullString := sql.NullString{String: uploadedFileName, Valid: true}

    // Assign the sql.NullString instance to the post.Image field
    post.Image = nullString

    if err := pc.DB.Create(&post).Error; err != nil {
        views.JSON(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    views.JSON(c, http.StatusCreated, post)
}
