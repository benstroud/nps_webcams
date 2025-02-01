package main

import (
	"log"
	"net/http"

	"html/template"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetHTMLTemplate(template.Must(template.ParseFiles("templates/index.html", "templates/error.html")))
	r.Static("/static", "./static")
	// Fetch all webcam data upfront and log the result
	log.Println("Fetching all webcams...")
	if err := FetchAllWebcams(); err != nil {
		log.Fatalf("Failed to fetch all webcams: %v", err)
	}
	log.Printf("Total webcams fetched: %d", len(allWebcams))

	r.GET("/", func(c *gin.Context) {
		groupedWebcams := GroupWebcamsByPark()
		log.Printf("Webcams grouped by park: %d parks", len(groupedWebcams))
		c.HTML(http.StatusOK, "index.html", gin.H{
			"groupedWebcams": groupedWebcams,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (default)
}
