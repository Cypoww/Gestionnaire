package main

import (
	"encoding/json"
	"fmt"
	"livres/db"
	"strconv"

	//"html/template"
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

type DeleteRequest struct {
	ID string `json:"id"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./html/*.html")

	r.Static("/css", "./assets/css")
	r.Static("/js", "./assets/js")
	r.Static("/img", "./assets/img")

	r.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(200)
				return
			}
			return
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1/")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost/")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	})

	mongo, cancel := db.Connect()
	defer cancel()

	r.GET("/livres", func(c *gin.Context) {
		bookList, err := db.GetLivres(c, mongo)
		if err != nil {
			c.JSON(http.StatusNotFound, "not found")
		}

		res, err := json.Marshal(bookList)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "unable to marshal json")
		}

		c.JSON(http.StatusOK, fmt.Sprintf("%v", string(res)))
	})

	r.GET("/livres/:id", func(c *gin.Context) {
		id, err := c.Params.Get("id")
		if !err {
			c.JSON(http.StatusInternalServerError, "got error")
		}
		res, got := db.GetLivre(c, mongo, id)
		if got == false {
			c.JSON(http.StatusNotFound, "not found")
			return
		}

		c.JSON(http.StatusOK, res)

	})

	r.POST("/livres", func(c *gin.Context) {
		body := c.Request.Body

		fmt.Printf("%v", body)

		var data Book
		err := json.NewDecoder(body).Decode(&data)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}

		err = db.PostLivre(c, mongo, db.Book{
			ID:     data.ID,
			Author: data.Author,
			Title:  data.Title,
			Year:   data.Year,
			Image:  data.Image,
		})
		if err != nil {
			if err.Error() == "already inserted" {
				c.JSON(http.StatusConflict, err.Error())
				return
			}
		}
		c.JSON(http.StatusOK, "ok")
	})

	r.POST("/livres/form", func(c *gin.Context) {
		query := c.Request.URL.Query()

		id := query.Get("id")
		title := query.Get("title")
		author := query.Get("author")
		year := query.Get("year")
		img := query.Get("img")

		integerYear, err := strconv.Atoi(year)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		err = db.PostLivre(c, mongo, db.Book{
			ID:     id,
			Author: author,
			Title:  title,
			Year:   integerYear,
			Image:  img,
		})
		if err != nil {
			if err.Error() == "already inserted" {
				c.JSON(http.StatusConflict, err.Error())
				return
			}
		}
		c.JSON(http.StatusOK, "ok")
	})

	r.DELETE("/livres/:id", func(c *gin.Context) {
		id, got := c.Params.Get("id")
		if !got {
			c.JSON(http.StatusInternalServerError, "got error")
		}

		err := db.DeleteBook(c, id, mongo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, "deleted")

	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/form", func(c *gin.Context) {
		c.HTML(http.StatusOK, "form.html", nil)
	})

	r.Run() // listen and serve on localhost:8080
}
