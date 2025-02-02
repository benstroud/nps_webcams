package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"html/template"

	"github.com/gin-gonic/gin"
)

const UPDATE_DATA_INTERVAL = 30 * time.Minute

var gatherDataTaskRunning bool
var lastUpdated = time.Now()
var mu sync.Mutex

// Fetch all webcams data
func gatherData() {
	log.Println("Fetching all webcams...")
	if err := FetchAllWebcams(); err != nil {
		log.Fatalf("Failed to fetch all webcams: %v", err)
	}
	lastUpdated = time.Now()
	log.Printf("Total webcams fetched: %d. Last updated: %v", len(allWebcams), lastUpdated)
}

// Background task to fetch data periodically
func gatherDataBackgroundTask() {
	mu.Lock()
	if gatherDataTaskRunning {
		mu.Unlock()
		log.Println("Skipping gather data task: Already running")
		return
	}
	gatherDataTaskRunning = true
	mu.Unlock()

	log.Println("Running gather data background task...")
	gatherData()

	mu.Lock()
	gatherDataTaskRunning = false
	mu.Unlock()
}

// Start the background task to fetch data periodically
func startGatherDataBackgroundTask() {
	ticker := time.NewTicker(UPDATE_DATA_INTERVAL)
	defer ticker.Stop()

	for range ticker.C {
		gatherDataBackgroundTask()
	}
}

func getLastUpdatedMinutes() int {
	mu.Lock()
	defer mu.Unlock()
	return int(time.Since(lastUpdated).Minutes())
}

func main() {
	go startGatherDataBackgroundTask()

	r := gin.Default()
	r.SetHTMLTemplate(template.Must(template.ParseFiles("templates/index.html", "templates/error.html")))
	r.Static("/static", "./static")

	// Initial data fetch
	gatherData()

	r.GET("/", func(c *gin.Context) {
		groupedWebcams := GroupWebcamsByPark()
		log.Printf("Webcams grouped by park: %d parks", len(groupedWebcams))
		c.HTML(http.StatusOK, "index.html", gin.H{
			"groupedWebcams":     groupedWebcams,
			"lastUpdatedMinutes": getLastUpdatedMinutes(),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (default)
}
