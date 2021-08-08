package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
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
	// router.Use(cors.Default())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://micro-podcast-svelte.vercel.app/"},
		AllowMethods:     []string{"POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept", "Accept-Encoding"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/rss", func(c *gin.Context) {
		var rss Rss
		c.BindJSON(&rss)

		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(rss.Url)

		if len(feed.Items) > 25 {
			c.JSON(http.StatusOK, gin.H{
				"title":       feed.Title,
				"description": feed.Description,
				"author":      feed.Author,
				"authors":     feed.Authors,
				"update":      feed.Updated,
				"image":       feed.Image,
				"items":       feed.Items[0:25],
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"title":       feed.Title,
				"description": feed.Description,
				"author":      feed.Author,
				"authors":     feed.Authors,
				"update":      feed.Updated,
				"image":       feed.Image,
				"items":       feed.Items,
			})
		}
	})

	router.Run(":" + port)
}
