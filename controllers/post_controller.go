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
	// Try to bind JSON data
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		// If JSON binding fails, try form data
		if err := c.ShouldBind(&post); err != nil {
			views.JSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Handle file upload
	_, header, err := c.Request.FormFile("image")
	if err == nil {
		// Save the image to a folder (you might want to use a unique filename)
		uploadedFileName := fmt.Sprintf("uploads/%s", header.Filename)
		err = c.SaveUploadedFile(header, uploadedFileName)
		if err != nil {
			views.JSON(c, http.StatusInternalServerError, gin.H{"error": "Error saving the image"})
			return
		}

		// Set Image to a valid value
		post.Image = sql.NullString{String: uploadedFileName, Valid: true}
	} else {
		// No image provided, set Image to NULL
		post.Image = sql.NullString{Valid: false}
	}

	// Insert the post into the database
	if err := pc.DB.Create(&post).Error; err != nil {
		views.JSON(c, http.StatusInternalServerError, gin.H{"error": "Error inserting post into the database"})
		return
	}

	views.JSON(c, http.StatusCreated, post)
}
