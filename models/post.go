// models/post.go

package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Image     sql.NullString `json:"image"`
	CreatedAt string         `json:"created_at"`
}
