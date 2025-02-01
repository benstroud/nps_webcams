package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
)

const BASE_URL = "https://developer.nps.gov/api/v1/"

type Crop struct {
	AspectRatio float64 `json:"aspectRatio"`
	URL         string  `json:"url"`
}

type Image struct {
	URL         string `json:"url"`
	Credit      string `json:"credit"`
	AltText     string `json:"altText"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Caption     string `json:"caption"`
	Crops       []Crop `json:"crops"`
}

type RelatedPark struct {
	States      string `json:"states"`
	ParkCode    string `json:"parkCode"`
	Designation string `json:"designation"`
	FullName    string `json:"fullName"`
	URL         string `json:"url"`
	Name        string `json:"name"`
}

type WebcamData struct {
	URL           string        `json:"url"`
	Title         string        `json:"title"`
	GeometryPoiID string        `json:"geometryPoiId"`
	ID            string        `json:"id"`
	Description   string        `json:"description"`
	Images        []Image       `json:"images"`
	RelatedParks  []RelatedPark `json:"relatedParks"`
	Status        string        `json:"status"`
	StatusMessage string        `json:"statusMessage"`
	IsStreaming   bool          `json:"isStreaming"`
	Tags          []string      `json:"tags"`
	Latitude      float64       `json:"latitude"`
	Longitude     float64       `json:"longitude"`
}

type WebcamResponse struct {
	Total string       `json:"total"`
	Data  []WebcamData `json:"data"`
	Limit string       `json:"limit"`
	Start string       `json:"start"`
}

type WebcamDataSlice []WebcamData

func (w WebcamDataSlice) Len() int {
	return len(w)
}

func (w WebcamDataSlice) Less(i, j int) bool {
	return w[i].Title < w[j].Title // Example: sort by title
}

func (w WebcamDataSlice) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

var (
	allWebcams WebcamDataSlice
	mutex      sync.Mutex
)

func GetApiKey() string {
	key := os.Getenv("NPS_API_KEY")
	if key == "" {
		log.Fatal("API key not found. Please set NPS_API_KEY environment variable. Here to register: https://www.nps.gov/subjects/developer/get-started.htm")
	}
	return key
}

func FetchAllWebcams() error {
	mutex.Lock()
	defer mutex.Unlock()

	// Get the total number of webcams from the initial request
	initialResponse, err := GetWebcams(1, 0)
	if err != nil {
		return fmt.Errorf("failed to fetch initial webcams: %w", err)
	}

	total, err := strconv.Atoi(initialResponse.Total)
	if err != nil {
		return fmt.Errorf("error converting total to integer: %w", err)
	}

	start := 0
	limit := 50

	for start < total {
		webcams, err := GetWebcams(limit, start)
		if err != nil {
			return fmt.Errorf("failed to fetch webcams: %w", err)
		}

		allWebcams = append(allWebcams, webcams.Data...)
		log.Printf("Total webcams stored: %d", len(allWebcams))
		start += limit
	}

	sort.Sort(allWebcams) // Sort the webcams by title (or any other criteria)
	log.Printf("Total webcams stored: %d", len(allWebcams))
	return nil
}

func GetWebcams(limit, start int) (WebcamResponse, error) {
	url := fmt.Sprintf("%swebcams?limit=%d&start=%d", BASE_URL, limit, start)
	log.Printf("Fetching data from the API. Url: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return WebcamResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Api-Key", GetApiKey())
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return WebcamResponse{}, fmt.Errorf("error fetching data from the API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WebcamResponse{}, fmt.Errorf("error fetching data from the API. Status code: %d", resp.StatusCode)
	}

	var webcams WebcamResponse

	if err := json.NewDecoder(resp.Body).Decode(&webcams); err != nil {
		return WebcamResponse{}, fmt.Errorf("error decoding JSON response: %w", err)
	}

	log.Printf("Total found: %s", webcams.Total)

	return webcams, nil
}

func GetWebcamsFromMemory(limit, start int) []WebcamData {
	mutex.Lock()
	defer mutex.Unlock()

	end := start + limit
	if end > len(allWebcams) {
		end = len(allWebcams)
	}
	return allWebcams[start:end]
}

func GroupWebcamsByPark() map[string][]WebcamData {
	mutex.Lock()
	defer mutex.Unlock()

	groupedWebcams := make(map[string][]WebcamData)
	for _, webcam := range allWebcams {
		for _, park := range webcam.RelatedParks {
			groupedWebcams[park.FullName] = append(groupedWebcams[park.FullName], webcam)
		}
	}
	log.Printf("grouped entries: %d", len(groupedWebcams))
	return groupedWebcams
}
