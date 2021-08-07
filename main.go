package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/mmcdole/gofeed"
)

type Rss struct {
	Url string `json:"rss"`
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.POST("/rss", func(c *gin.Context) {
		var rss Rss
		c.BindJSON(&rss)

		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(rss.Url)

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,X-Requested-With,Content-Type,Accept,content-type,application/json")

		c.JSON(http.StatusOK, gin.H{
			"title":       feed.Title,
			"description": feed.Description,
			"author":      feed.Author,
			"authors":     feed.Authors,
			"update":      feed.Updated,
			"image":       feed.Image,
			"items":       feed.Items,
		})
	})

	router.Run(":" + port)
}
