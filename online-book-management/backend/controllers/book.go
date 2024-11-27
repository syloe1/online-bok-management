// backend/controllers/book.go
package controllers

import (
	"log"
	"net/http"

	"online-book-management/backend/config"
	"online-book-management/backend/models"

	"github.com/gin-gonic/gin"
)

// ShowBooksPage 显示图书列表页面
func ShowBooksPage(c *gin.Context) {
	c.HTML(http.StatusOK, "books.html", nil)
}

// ShowAddBookPage 显示添加图书页面
func ShowAddBookPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_book.html", nil)
}

// GetBooks 获取所有图书
func GetBooks(c *gin.Context) {
	log.Println("Received GET /books request")
	var books []models.Book
	result := config.DB.Find(&books)
	if result.Error != nil {
		log.Printf("Error fetching books: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取图书列表"})
		return
	}
	log.Printf("Fetched books: %+v", books)
	c.JSON(http.StatusOK, books)
}

// GetBook 获取单个图书详情
func GetBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		log.Printf("Book not found with ID %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "图书未找到"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// CreateBook 创建新图书
func CreateBook(c *gin.Context) {
	var input struct {
		Title       string `form:"title" binding:"required"`
		Author      string `form:"author" binding:"required"`
		Description string `form:"description"`
		Quantity    int    `form:"quantity" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		log.Printf("CreateBook: binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "所有必填字段不能为空"})
		return
	}

	book := models.Book{
		Title:       input.Title,
		Author:      input.Author,
		Description: input.Description,
		Quantity:    input.Quantity,
	}

	result := config.DB.Create(&book)
	if result.Error != nil {
		log.Printf("CreateBook: error creating book: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建图书"})
		return
	}

	log.Printf("Created book: %+v", book)
	c.JSON(http.StatusOK, book)
}

// UpdateBook 更新图书信息
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		log.Printf("UpdateBook: book not found with ID %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "图书未找到"})
		return
	}

	var input struct {
		Title       string `json:"title"`
		Author      string `json:"author"`
		Description string `json:"description"`
		Quantity    int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("UpdateBook: binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != "" {
		book.Title = input.Title
	}
	if input.Author != "" {
		book.Author = input.Author
	}
	if input.Description != "" {
		book.Description = input.Description
	}
	if input.Quantity >= 0 {
		book.Quantity = input.Quantity
	}

	result := config.DB.Save(&book)
	if result.Error != nil {
		log.Printf("UpdateBook: error updating book: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法更新图书"})
		return
	}

	log.Printf("Updated book: %+v", book)
	c.JSON(http.StatusOK, book)
}

// DeleteBook 删除图书
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		log.Printf("DeleteBook: book not found with ID %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "图书未找到"})
		return
	}

	result := config.DB.Delete(&book)
	if result.Error != nil {
		log.Printf("DeleteBook: error deleting book: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法删除图书"})
		return
	}

	log.Printf("Deleted book with ID %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "图书已删除"})
}
