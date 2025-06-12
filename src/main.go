package main

import (
	"Weather/internal/db"
	"Weather/internal/models"
	"Weather/internal/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const KrakowLat float64 = 50.049683
const KrakowLong float64 = 19.944544

func init() {
	if os.Getenv("API_KEY") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Can't load env variables: %v", err)
		}
	}
}

func renderHome() (*models.WeatherResponse, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("Can't get API_KEY")
	}
	var query models.Query = models.Query{
		Lat:   KrakowLat,
		Lon:   KrakowLong,
		Appid: apiKey,
	}
	weatherURL := utils.ConstructURL(query)

	fmt.Printf("WeatherURL: %s\n", weatherURL)

	resp, err := http.Get(weatherURL)
	if err != nil {
		log.Fatalf("Error whilst sending get request to %s: %v\n", weatherURL, err)
	}
	defer resp.Body.Close()

	fmt.Printf("Get request status code: %d\n", resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Can't read bytes from response body: %v", err)
	}

	fmt.Printf("Response body: %s\n", string(bodyBytes))

	var weatherResponse *models.WeatherResponse
	err = json.Unmarshal(bodyBytes, &weatherResponse)
	weatherResponse.HotColdThreshold = models.HotColdThreshold
	if err != nil {
		log.Fatalf("Can't unmarshal bodyBytes to weather response struct: %v\n", err)
	}

	return weatherResponse, nil
}

func main() {
	//set up mime extensions to be correct
	utils.Init()

	mux := http.NewServeMux()

	mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		var homeTemplate = template.Must(template.ParseFiles("templates/index.html"))
		weatherResponse, err := renderHome()
		if err != nil {
			log.Fatalf("Error executing renderHome: %v\n", err)
		}

		w.Header().Set("Content-Type", "text/html")

		homeTemplate.Execute(w, weatherResponse)
	})

	mux.HandleFunc("/recommendations", func(w http.ResponseWriter, r *http.Request) {
		var recommendationsTemplate = template.Must(template.ParseFiles("templates/recommendations.html"))
		weatherResponse, err := renderHome()
		if err != nil {
			log.Fatalf("Error executing renderHome: %v\n", err)
		}

		w.Header().Set("Content-Type", "text/html")

		recommendationsTemplate.Execute(w, weatherResponse)
	})
	mux.HandleFunc("/feedback", func(w http.ResponseWriter, r *http.Request) {
		var feedbackTemplate = template.Must(template.ParseFiles("templates/feedback.html"))
		w.Header().Set("Content-Type", "text/html")
		feedbackTemplate.Execute(w, nil)
	})

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir("static/img"))
		w.Header().Set("Cache-Control", "public, max-age=86400")
	})

	mux.HandleFunc("/feedback/share", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		satisfaction := r.FormValue("satisfaction")
		feedback := r.FormValue("feedback")
		fmt.Printf("Received form review. Satisfaction: %s, feedback: %s\n", satisfaction, feedback)

		db, err := db.InitDB()
		if err != nil {
			http.Error(w, "Failed initializing DB", http.StatusInternalServerError)
			log.Printf("DB error: %v\n", err)
			return
		}

		_, err = db.Exec(`
			INSERT INTO feedback (satisfaction, feedback)
			VALUES (?, ?)
		`, satisfaction, feedback)

		if err != nil {
			http.Error(w, "Failed inserting values into DB", http.StatusInternalServerError)
			log.Printf("DB error: %v\n", err)
			return
		}

		w.Write([]byte("Feedback saved successfully!"))
	})

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	fmt.Println("Listening on the 8000 server")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
