// controllers/post_controller.go

package controllers

import (
	"blogbackend/models"
	"blogbackend/views"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	DB *sql.DB
}

func NewPostController(db *sql.DB) *PostController {
	return &PostController{DB: db}
}

func (pc *PostController) GetPosts(c *gin.Context) {
	rows, errQuery := pc.DB.Query("SELECT id, title, content, COALESCE(image, ''), created_at FROM posts ORDER BY created_at DESC")
	if errQuery != nil {
		views.JSON(c, http.StatusInternalServerError, gin.H{"error": errQuery.Error()})
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		errScan := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.CreatedAt)
		if errScan != nil {
			views.JSON(c, http.StatusInternalServerError, gin.H{"error": errScan.Error()})
			return
		}
		posts = append(posts, post)
	}

	views.JSON(c, http.StatusOK, gin.H{"posts": posts})
}

func (pc *PostController) GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	errQuery := pc.DB.QueryRow("SELECT id, title, content, COALESCE(image, ''), created_at FROM posts WHERE id = $1", id).
		Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.CreatedAt)

	switch {
	case errQuery == sql.ErrNoRows:
		views.JSON(c, http.StatusNotFound, gin.H{"error": "Post not found"})
	case errQuery != nil:
		views.JSON(c, http.StatusInternalServerError, gin.H{"error": errQuery.Error()})
	default:
		views.JSON(c, http.StatusOK, post)
	}
}
func (pc *PostController) CreatePost(c *gin.Context) {
	// Parse form data, including files
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		views.JSON(c, http.StatusInternalServerError, gin.H{"error": "Error parsing form data"})
		return
	}

	// Retrieve other form data
	title := c.Request.FormValue("title")
	content := c.Request.FormValue("content")

	// Handle file upload
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		views.JSON(c, http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}
	defer file.Close()

	// Save the image to a folder (you might want to use a unique filename)
	uploadedFileName := fmt.Sprintf("uploads/%s", header.Filename)
	err = c.SaveUploadedFile(header, uploadedFileName)
	if err != nil {
		views.JSON(c, http.StatusInternalServerError, gin.H{"error": "Error saving the image"})
		return
	}

	// Create a new Post instance with data
	post := models.Post{
		Title:   title,
		Content: content,
		Image:   sql.NullString{String: uploadedFileName, Valid: true}, // Use NewNullString to handle string to sql.NullString conversion
	}

	// Insert the post into the database
	err = pc.DB.QueryRow("INSERT INTO posts (title, content, image) VALUES ($1, $2, $3) RETURNING id", post.Title, post.Content, post.Image).
		Scan(&post.ID)

	if err != nil {
		views.JSON(c, http.StatusInternalServerError, gin.H{"error": "Error inserting post into the database"})
		return
	}

	views.JSON(c, http.StatusCreated, post)
}
