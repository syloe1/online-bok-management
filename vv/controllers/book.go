// backend/controllers/book.go
package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"online-book-management/backend/config"
	"online-book-management/backend/models"

	"github.com/gin-gonic/gin"
)

// 设定上传文件的目录
const uploadDir = "uploads/"

// ShowBooksPage 显示图书列表页面
func ShowBooksPage(c *gin.Context) {
	user, _ := c.Get("User") // 获取用户信息
	c.HTML(http.StatusOK, "books.html", gin.H{
		"User": user,
	})
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

// CreateBook 创建新图书，并处理PDF文件上传
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

	// 处理文件上传
	var pdfPath string
	file, err := c.FormFile("pdf")
	if err == nil {
		// 检查文件扩展名
		if filepath.Ext(file.Filename) != ".pdf" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持上传PDF文件"})
			return
		}

		// 创建上传目录（如果不存在）
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, os.ModePerm)
		}

		// 生成唯一的文件名，避免覆盖
		filename := filepath.Base(file.Filename)
		pdfPath = filepath.Join(uploadDir, filename)
		if err := c.SaveUploadedFile(file, pdfPath); err != nil {
			log.Printf("CreateBook: 文件上传失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件上传失败"})
			return
		}
	} else {
		log.Println("CreateBook: 未上传PDF文件")
	}

	book := models.Book{
		Title:       input.Title,
		Author:      input.Author,
		Description: input.Description,
		Quantity:    input.Quantity,
		PdfPath:     pdfPath,
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

// UpdateBook 更新图书信息，并处理PDF文件上传
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		log.Printf("UpdateBook: book not found with ID %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "图书未找到"})
		return
	}

	var input struct {
		Title       string `form:"title"`
		Author      string `form:"author"`
		Description string `form:"description"`
		Quantity    int    `form:"quantity"`
	}

	if err := c.ShouldBind(&input); err != nil {
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

	// 处理文件上传
	fileHeader, err := c.FormFile("pdf")
	if err == nil {
		// 检查文件扩展名
		if filepath.Ext(fileHeader.Filename) != ".pdf" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持上传PDF文件"})
			return
		}

		// 创建上传目录（如果不存在）
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, os.ModePerm)
		}

		// 生成唯一的文件名，避免覆盖
		filename := filepath.Base(fileHeader.Filename)
		pdfPath := filepath.Join(uploadDir, filename)
		if err := c.SaveUploadedFile(fileHeader, pdfPath); err != nil {
			log.Printf("UpdateBook: 文件上传失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件上传失败"})
			return
		}
		book.PdfPath = pdfPath
	} else {
		log.Println("UpdateBook: 未上传PDF文件")
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

// UploadBookPDF 单独上传图书的PDF文件
func UploadBookPDF(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		log.Printf("UploadBookPDF: book not found with ID %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "图书未找到"})
		return
	}

	file, err := c.FormFile("pdf")
	if err != nil {
		log.Printf("UploadBookPDF: 获取文件失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "未上传文件"})
		return
	}

	// 检查文件扩展名
	if filepath.Ext(file.Filename) != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持上传PDF文件"})
		return
	}

	// 创建上传目录（如果不存在）
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	// 生成唯一的文件名，避免覆盖
	filename := filepath.Base(file.Filename)
	pdfPath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, pdfPath); err != nil {
		log.Printf("UploadBookPDF: 文件上传失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件上传失败"})
		return
	}

	// 更新图书的PdfPath
	book.PdfPath = pdfPath
	if err := config.DB.Save(&book).Error; err != nil {
		log.Printf("UploadBookPDF: 更新图书PDF路径失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法更新图书PDF路径"})
		return
	}

	log.Printf("Uploaded PDF for book ID %s: %s", id, pdfPath)
	c.JSON(http.StatusOK, gin.H{"message": "PDF文件上传成功", "pdf_path": pdfPath})
}

// SearchBooks 搜索图书
func SearchBooks(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	var books []models.Book
	// 使用LIKE进行模糊搜索，匹配标题或作者
	result := config.DB.Where("title LIKE ? OR author LIKE ?", "%"+query+"%", "%"+query+"%").Find(&books)
	if result.Error != nil {
		log.Printf("SearchBooks: 搜索失败: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}

	c.JSON(http.StatusOK, books)
}

// ShowEditBookPage 显示编辑图书页面
func ShowEditBookPage(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		log.Printf("ShowEditBookPage: book not found with ID %s", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "图书未找到"})
		return
	}
	c.HTML(http.StatusOK, "edit_book.html", gin.H{"book": book})
}
