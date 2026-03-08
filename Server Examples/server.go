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

type UserPerms int

const (
	Admin UserPerms = iota
	Officer
	Member
)

func main() {
	dsn := "host=localhost port=5432 user=dev_user password=testing dbname=dev_project_db"
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	os.Chdir("../Svelte Examples/plain-svelte-app/dist")

	mux := http.NewServeMux()

	mux.Handle("/submit", checkSubmissionDataMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbError := makeGenericDBError(w)

		var submission struct {
			Name    string `json:"name"`
			Email   string `json:"email"`
			Club    string `json:"club"`
			Contact string `json:"contact"`
			Officer string `json:"officer"`
			Login   string `json:"login"`
		}
		checkJson(w, r, &submission)

		fmt.Printf("Received submission: Name=%s, Club=%s, Contact=%s, Officer=%s, Email=%s, Login=%s\n", submission.Name, submission.Club, submission.Contact, submission.Officer, submission.Email, submission.Login)
		if submission.Officer == "yes" {
			query := func() (any, error) {
				return client.User.Create().
					SetGoogleID(submission.Name).
					SetEmail(submission.Email).
					Save(r.Context())
			}
			databaseQuery(query, dbError)
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Form submitted successfully!"})
	})))
	mux.Handle("/", http.FileServer(http.Dir(".")))

	port := ":8080"
	fmt.Println("Server is running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, mux))
}

func checkUserAuthMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authenticated (this is just a placeholder)
		authenticated := true // Replace with actual authentication logic

		if !authenticated {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func checkUserRoleMiddleware(user string, role UserPerms, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user has the required role (this is just a placeholder)
		hasRole := true // Replace with actual role-checking logic

		if !hasRole {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func checkSubmissionDataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Only POST method is allowed"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func makeGenericDBError(w http.ResponseWriter) func(error) {
	return func(err error) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		fmt.Println(err)
	}
}

func databaseQuery[T any](queryFunc func() (T, error), errFunc func(error)) T {
	res, err := queryFunc()
	if err != nil {
		errFunc(err)
		var zero T
		return zero
	}
	return res
}

func checkJson(w http.ResponseWriter, r *http.Request, target any) {
	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}
}
