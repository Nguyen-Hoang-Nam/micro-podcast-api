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
	Url    string `json:"rss"`
	Update string `json:"update"`
}

func main() {
	port := os.Getenv("PORT")
	env := os.Getenv("ENV")
	client := os.Getenv("CLIENT")

	// Heroku need enviroment PORT
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Create Middleware
	router := gin.New()

	// Add Loger to Middleware
	router.Use(gin.Logger())

	// Add Cors to Middleware
	// Support all origin in development mode
	// Only support https://micro-podcast-svelte.vercel.app in product mode
	if env == "development" {
		router.Use(cors.Default())
	} else {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{client},
			AllowMethods:     []string{"POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept", "Accept-Encoding", "Access-Control-Request-Headers", "Access-Control-Request-Method"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))
	}

	// Subscript rss by link
	router.POST("/rss", func(c *gin.Context) {
		var rss Rss
		c.BindJSON(&rss)

		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(rss.Url)

		if rss.Update == feed.Updated {
			c.JSON(http.StatusOK, gin.H{
				"message": "Already latest",
			})
		} else {
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
		}
	})

	router.Run(":" + port)
}
