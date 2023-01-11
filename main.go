package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	// BLOB
	Image string `json:"img"`
}

func Add(b Book) {
	books = append(books, Book{Title: b.Title, Author: b.Author, Year: b.Year})
}

var books []Book

func init() {
	books = []Book{
		{ID: "1", Title: "Le Seigneur des Anneaux", Author: "J.R.R. Tolkien", Year: 1954},
		{ID: "2", Title: "Le Hobbit", Author: "J.R.R. Tolkien", Year: 1937},
		{ID: "3", Title: "Le Trône de Fer", Author: "George R.R. Martin", Year: 1996},
		{ID: "4", Image: ""},
	}
}

func getBooks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": books})
}

func getBook(c *gin.Context) {
	// Récupérez l'ID du livre à partir de la route
	id := c.Param("id")

	// Recherchez le livre correspondant dans la liste
	for _, book := range books {
		if book.ID == id {
			c.JSON(http.StatusOK, gin.H{"data": book})
			return
		}
	}

	// Si le livre n'est pas trouvé, renvoyez un code d'erreur 404
	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

func Delete(id string) error {
	for i, b := range books {
		if b.ID == id {
			books[i] = books[len(books)-1]
			books = books[:len(books)-1]
			return nil
		}
	}
	return fmt.Errorf("Book with ID %d not found", id)
}

//var templates = template.Must(template.ParseFiles("index.html"))

func main() {
	r := gin.Default()

	r.GET("/livres", func(c *gin.Context) {
		c.JSON(http.StatusOK, fmt.Sprintf("%v", books))
	})

	r.GET("/livres/:id", func(c *gin.Context) {
		getBook(c)
	})

	r.POST("/livres", func(c *gin.Context) {
		body := c.Request.Body

		var data Book
		err := json.NewDecoder(body).Decode(&data)
		if err != nil {
			panic("marche pas")
		}

		Add(data)
		c.JSON(http.StatusOK, fmt.Sprintf("%v", data.ID))
	})

	r.DELETE("/livres/:id", func(c *gin.Context) {
		Delete(c.Param("id"))
	})

	//r.GET("/", homeHandler)

	r.Run() // listen and serve on localhost:8080

}

//func homeHandler(c *gin.Context) {
	//err := templates.ExecuteTemplate(c.Writer, "index.html", p)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
//	}
//}
