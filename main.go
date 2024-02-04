package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	//"errors"
)

// define a structure
// keep first char of variable name as caps to make it exportable
// for ensuring json serialisation add `json:<variable name in lower case>`
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "7 hobits", Author: "P M", Quantity: 4},
	{ID: "2", Title: "15 habits", Author: "R M", Quantity: 2},
	{ID: "3", Title: "65 rabits", Author: "D M", Quantity: 8},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func findBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func getBookById(c *gin.Context) {
	bookId := c.Param("id")
	book, err := findBookById(bookId)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func addBooks(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missinng id in query parameter."})
		return
	}

	book, err := findBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Insufficient number of books."})
		return
	}
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)

}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missinng id in query parameter."})
		return
	}

	book, err := findBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Please add book first."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)

}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBookById)
	router.PATCH("/checkoutBook", checkoutBook)
	router.PATCH("/returnBook", returnBook)
	router.POST("/books", addBooks)
	router.Run("localhost:8080")
}
