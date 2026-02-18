package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	os.Chdir("../Svelte Examples/plain-svelte-app/dist")
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var submission struct {
				Name   string `json:"name"`
				Club   string `json:"club"`
				Contact string `json:"contact"`
				Officer string `json:"officer"`
			}
			err := json.NewDecoder(r.Body).Decode(&submission)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
				return
			}
			fmt.Printf("Received submission: Name=%s, Club=%s, Contact=%s, Officer=%s\n", submission.Name, submission.Club, submission.Contact, submission.Officer)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Form submitted successfully!"})
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Only POST method is allowed"})
		}
	})
	http.Handle("/", http.FileServer(http.Dir(".")))

	port := ":8080"
	fmt.Println("Server is running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
