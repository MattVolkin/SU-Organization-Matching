package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"server-example/ent"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// host=localhost port=5432 user=dev_user password=testing

type UserRole int

const (
	Admin UserRole = iota
	Officer
	Member
)

// String converts a UserRole enum value into a stable network-safe label.
func (role UserRole) String() string {
	switch role {
	case Admin:
		return "admin"
	case Officer:
		return "officer"
	case Member:
		return "member"
	default:
		return "unknown"
	}
}

// MarshalJSON serializes UserRole as a string instead of its integer value.
func (role UserRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(role.String())
}

// SurveyResponsePayload is the request body accepted by the /response endpoint.
type SurveyResponsePayload struct {
	QuestionID int  `json:"questionId"`
	Answer     bool `json:"answer"`
}

// DemographicsPayload is the expected payload for the /submit endpoint.
type DemographicsPayload struct {
	Name     string   `json:"name"`
	Gender   string   `json:"gender"`
	Race     []string `json:"race"`
	Religion string   `json:"religion"`
	Major    []string `json:"major"`
}

// OrgJSON defines the reusable wire format for officer organization data.
// This can be marshaled to JSON for responses and decoded from JSON in requests.
type OrgJSON struct {
	ID                   int       `json:"id"`
	ClubName             string    `json:"clubName"`
	Description          string    `json:"description"`
	MeetingTime          string    `json:"meetingTime"`
	ImagePath            string    `json:"imagePath"`
	ExternalLink         string    `json:"externalLink"`
	ContactInfo          string    `json:"contactInfo"`
	IncludeOfficerEmails bool      `json:"includeOfficerEmails"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

var dbClient *DatabaseClient

// main wires up OAuth, database connectivity, middleware-protected routes,
// and static file serving for the Svelte frontend.
func main() {
	// Initialize in-memory session storage used by auth middleware.
	authSessionStore = NewAuthSessionStore()

	// Read OAuth settings from the environment and apply safe defaults.
	googleClientID := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
	googleClientSecret := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET"))
	googleRedirectURL := strings.TrimSpace(os.Getenv("GOOGLE_REDIRECT_URL"))
	if googleRedirectURL == "" {
		googleRedirectURL = "http://localhost:8080/auth/callback"
	}

	if googleClientID == "" {
		log.Fatal("missing GOOGLE_CLIENT_ID environment variable")
	}
	if googleClientSecret == "" {
		log.Fatal("missing GOOGLE_CLIENT_SECRET environment variable")
	}

	log.Printf("Google OAuth redirect URL: %s", googleRedirectURL)
	log.Printf("Google OAuth client ID suffix: %s", redactValueSuffix(googleClientID, 10))

	// Build the OAuth client used by /login and /auth/callback.
	googleOAuth := &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  googleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	oauthCfg := &OAuthRuntimeConfig{
		GoogleOAuth: googleOAuth,
		StateStore: &OAuthStateStore{
			states: make(map[string]OAuthStateRecord),
		},
	}

	// Open the Postgres connection used by Ent queries and mutations.
	dsn := "host=localhost port=5432 user=dev_user password=testing dbname=dev_project_db"
	var err error
	dbClient, err = NewDatabaseClient("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	RegisterDatabaseClient("default", dbClient.WithAlias("default"))
	defer dbClient.Close()

	submitRequiresAuth := readBoolEnv("SUBMIT_REQUIRES_AUTH", true)
	log.Printf("Submit auth required: %t", submitRequiresAuth)

	// Serve static files from the built Svelte distribution directory.
	if err := os.Chdir("../Svelte Examples/plain-svelte-app/dist"); err != nil {
		log.Printf("warning: failed to switch static directory: %v", err)
	}

	// Register API routes and static site handler.
	router := mux.NewRouter()

	router.HandleFunc("/login", makeOAuthLoginHandler(oauthCfg))
	router.HandleFunc("/auth/callback", makeOAuthCallbackHandler(oauthCfg))
	router.HandleFunc("/api/user", makeCurrentUserHandler())
	router.HandleFunc("/logout", makeLogoutHandler())

	router.Handle("/response", requireAuthenticatedSession(http.HandlerFunc(handleSurveyResponseSubmission)))

	submitHandler := requirePostMethod(http.HandlerFunc(handleDemographicsSubmission))
	if submitRequiresAuth {
		submitHandler = requireAuthenticatedSession(submitHandler)
	}
	router.Handle("/submit", submitHandler)

	// All /api/ endpoints require an authenticated session, but only some require specific roles.
	apiRouter := router.PathPrefix("/api/").Subrouter()
	apiRouter.Use(requireAuthenticatedSession)

	apiRouter.HandleFunc("/prefill", makePrefillFieldsHandler())

	apiRouter.Handle("/adjectives", http.HandlerFunc(handleAdjectivesRequest))

	// Officer-only endpoints are nested under /api/officer/ and use additional role-checking middleware.
	officerRouter := apiRouter.PathPrefix("/officer/").Subrouter()

	officerRouter.Use(makeRequireUserRoleHandlerMiddleware(Officer))

	officerRouter.Handle("/orgs", handleOfficerOrgsRequest)

	officerRouter.Handle("/update", handleOfficerUpdateRequest)

	// Admin-only endpoints are nested under /api/admin/ and use additional role-checking middleware.
	adminRouter := apiRouter.PathPrefix("/admin/").Subrouter()
	adminRouter.Use(makeRequireUserRoleHandlerMiddleware(Admin))

	adminRouter.Handle("/orgs", handleAdminOrgsRequest)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))

	// Start HTTP server.
	port := ":8080"
	fmt.Println("Server is running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, router))
}

func handleAdminOrgsRequest(w http.ResponseWriter, r *http.Request) {
	_, clubs, err := dbClient.Query().FetchAllClubs(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch orgs"})
		return
	}

	jsonClubs := orgStructToJSON(clubs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonClubs)
}

func handleOfficerUpdateRequest(w http.ResponseWriter, r *http.Request) {
	var payload OrgJSON
	if !decodeJSONBody(w, r, &payload) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	dbClient.Query().UpdateClubFromJSON(r.Context(), &payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func handleOfficerOrgsRequest(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("X-User-Email")
	_, clubs, err := dbClient.Query().FetchOfficerClubsByUserEmail(r.Context(), email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch officer orgs"})
		return
	}

	jsonClubs := orgStructToJSON(clubs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonClubs)
}

// orgStructToJSON converts Ent club models into the shared JSON wire format.
func orgStructToJSON(clubs []*ent.Club) []OrgJSON {
	result := make([]OrgJSON, 0, len(clubs))
	for _, club := range clubs {
		if club == nil {
			continue
		}
		result = append(result, orgJSONFromEntClub(club))
	}
	return result
}

func orgJSONFromEntClub(club *ent.Club) OrgJSON {
	return OrgJSON{
		ID:                   club.ID,
		ClubName:             club.ClubName,
		Description:          club.Description,
		MeetingTime:          club.MeetingTime,
		ImagePath:            club.ImagePath,
		ExternalLink:         club.ExternalLink,
		ContactInfo:          club.ContactInfo,
		IncludeOfficerEmails: club.IncludeOfficerEmails,
		UpdatedAt:            club.UpdatedAt,
	}
}

// decodeOrgJSONs decodes a request body with the same JSON format used
// by OrgJSON responses so read/write paths can share one contract.
func decodeOrgJSONs(w http.ResponseWriter, r *http.Request) ([]OrgJSON, bool) {
	var payload []OrgJSON
	if !decodeJSONBody(w, r, &payload) {
		return nil, false
	}
	return payload, true
}

func handleAdjectivesRequest(w http.ResponseWriter, r *http.Request) {
	_, adjectives, err := dbClient.Query().FetchAdjectiveQuestionContents(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch adjectives"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adjectives)
}

// requireAuthenticatedSession rejects requests that do not carry a valid session
// token, and forwards the authenticated email to downstream handlers.
func requireAuthenticatedSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractSessionTokenFromRequest(r)
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: missing token"})
			return
		}

		session, ok := authSessionStore.getSession(token)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: invalid or expired token"})
			return
		}

		// Attach authenticated user email for downstream handlers.
		r.Header.Set("X-User-Email", session.Email)
		next.ServeHTTP(w, r)
	})
}

// makeRequireUserRoleHandlerMiddleware returns middleware that enforces a minimum
// user role for access to an endpoint, returning 403 Forbidden when the user's
// role is insufficient.
func makeRequireUserRoleHandlerMiddleware(requiredRole UserRole) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = requiredRole

			userEmail := r.Header.Get("X-User-Email")

			hasRole := dbClient.Query().getUserRole(r.Context(), userEmail) == requiredRole
			if !hasRole {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// requirePostMethod enforces POST-only semantics for endpoints
// expecting JSON form submissions.
func requirePostMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Only POST method is allowed"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// handleSurveyResponseSubmission validates and echoes a response payload from an
// authenticated user.
func handleSurveyResponseSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Only POST method is allowed"})
		return
	}

	var payload SurveyResponsePayload
	if !decodeJSONBody(w, r, &payload) {
		return
	}
	if payload.QuestionID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "questionId must be a positive integer"})
		return
	}

	// Return a confirmation payload so the client can verify what was stored.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"email":      r.Header.Get("X-User-Email"),
		"questionId": payload.QuestionID,
		"answer":     payload.Answer,
		"message":    "Response accepted",
	})
}

// handleDemographicsSubmission validates the demographics payload and returns
// the accepted values as confirmation.
func handleDemographicsSubmission(w http.ResponseWriter, r *http.Request) {
	var submission DemographicsPayload
	if !decodeJSONBody(w, r, &submission) {
		return
	}

	if strings.TrimSpace(submission.Name) == "" ||
		strings.TrimSpace(submission.Gender) == "" ||
		strings.TrimSpace(submission.Religion) == "" ||
		len(submission.Race) == 0 ||
		len(submission.Major) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing one or more required demographics fields"})
		return
	}

	// Echo submitted data as a simple success response for the client.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message":  "Demographics submitted successfully",
		"name":     submission.Name,
		"gender":   submission.Gender,
		"race":     submission.Race,
		"religion": submission.Religion,
		"major":    submission.Major,
		"email":    r.Header.Get("X-User-Email"),
	})
}

// decodeJSONBody decodes a JSON request body into target and writes a 400 response
// if decoding fails.
func decodeJSONBody(w http.ResponseWriter, r *http.Request, target any) bool {
	err := json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return false
	}
	return true
}

// readBoolEnv parses common truthy/falsey values and falls back to a default
// when unset or unrecognized.
func readBoolEnv(key string, defaultValue bool) bool {
	raw := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if raw == "" {
		return defaultValue
	}
	switch raw {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return defaultValue
	}
}
