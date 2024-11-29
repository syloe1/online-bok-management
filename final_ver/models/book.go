// backend/models/book.go
package models

import (
	"time"
)

type Book struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Author      string    `gorm:"not null" json:"author"`
	Description string    `json:"description"`
	Quantity    int       `gorm:"default:0" json:"quantity"`
	PdfPath     string    `json:"pdf_path"` // 新增字段，用于存储PDF文件路径
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
