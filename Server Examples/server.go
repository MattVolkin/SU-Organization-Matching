package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"server-example/ent"

	_ "github.com/lib/pq"
)

//host=localhost port=5432 user=dev_user password=testing

func main() {
	dsn := "host=localhost port=5432 user=dev_user password=testing"
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	os.Chdir("../Svelte Examples/plain-svelte-app/dist")
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var submission struct {
				Name    string `json:"name"`
				Club    string `json:"club"`
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
			if submission.Officer == "yes" {
				_, err := client.User.Create().
					SetGoogleID(submission.Name).
					Save(r.Context())
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save user"})
					return
				}
			}
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
