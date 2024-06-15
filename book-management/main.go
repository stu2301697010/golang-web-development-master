package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Book struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Title  string `json:"title"`
	ISBN   string `json:"isbn"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

type CSVBook struct {
	ISBN              string `csv:"ISBN"`
	Title             string `csv:"Book-Title"`
	Author            string `csv:"Book-Author"`
	YearOfPublication int    `csv:"Year-Of-Publication"`
}

var DB *gorm.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	DB.AutoMigrate(&Book{})
}

func LoadCSVData(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Create a custom CSV reader with ';' as the delimiter
	r := csv.NewReader(file)
	r.Comma = ';'

	var csvBooks []CSVBook

	// Custom unmarshalling with the reader
	if err := gocsv.UnmarshalCSV(r, &csvBooks); err != nil {
		log.Fatalf("Failed to parse CSV file: %v", err)
	}

	for _, csvBook := range csvBooks {
		book := Book{
			Title:  csvBook.Title,
			ISBN:   csvBook.ISBN,
			Author: csvBook.Author,
			Year:   csvBook.YearOfPublication,
		}
		DB.Create(&book)
	}
}

func getBooks(c *gin.Context) {
	var books []Book
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	DB.Offset(offset).Limit(limit).Find(&books)

	var totalBooks int64
	DB.Model(&Book{}).Count(&totalBooks)

	c.JSON(http.StatusOK, gin.H{
		"books":       books,
		"page":        page,
		"total_pages": (totalBooks + int64(limit) - 1) / int64(limit),
	})
}

func getBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var book Book

	if err := DB.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

func createBook(c *gin.Context) {
	var newBook Book

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.Create(&newBook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newBook)
}

func updateBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedBook Book

	if err := c.BindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.Model(&Book{}).Where("id = ?", id).Updates(updatedBook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedBook)
}

func deleteBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := DB.Delete(&Book{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book with ID " + strconv.Itoa(id) + " has been deleted"})
}

func main() {
	r := gin.Default()

	InitDB()

	LoadCSVData("books.csv")

	r.GET("/books", getBooks)
	r.GET("/books/:id", getBook)
	r.POST("/books", createBook)
	r.PUT("/books/:id", updateBook)
	r.DELETE("/books/:id", deleteBook)

	r.Run(":8080")
}
