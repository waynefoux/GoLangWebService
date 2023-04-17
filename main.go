package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type GreetingResponse struct {
	Message string `json:"message"`
}

type NewGreeting struct {
	Greeting string `json:"greeting`
}

var greetings = []string{"Hello"}
var greetingsMutex sync.RWMutex

func greetingHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}

	// Extract the name from the URL params
	name := strings.TrimPrefix(r.URL.Path, "/greeting/")

	if name == "" {
		name = "world"
	}

	// Generate a random greeting from the list of greetings
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(len(greetings))

	response := GreetingResponse{Message: fmt.Sprintf("%s, %s!", greetings[randomNumber], name)}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

/* Handler to add new greetings to our greeting endpoint */
func addGreetingHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var newGreeting NewGreeting
	err = json.Unmarshal(body, &newGreeting)

	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	greetingsMutex.Lock()
	greetings = append(greetings, newGreeting.Greeting)
	greetingsMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/greeting/", greetingHandler)
	mux.HandleFunc("/greeting/add", addGreetingHandler)
	port := "8080"
	fmt.Printf("Starting web service on port %s...\n", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Printf("Error starting web service: %v\n", err)
	}
}
