package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"su-organization-matching/matching"
	"su-organization-matching/server/ent"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	gomail "gopkg.in/mail.v2"
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

// SurveyResponsePayload is one item in the /response request payload array.
type SurveyResponsePayload struct {
	QuestionID int  `json:"questionId"`
	Answer     bool `json:"answer"`
}

// DemographicsPayload is the expected payload for the /submit endpoint.
type DemographicsPayload struct {
	Name       string   `json:"name"`
	Gender     string   `json:"gender"`
	Race       []string `json:"race"`
	Religion   string   `json:"religion"`
	Major      []string `json:"major"`
	LGBTQ      string   `json:"lgbtq"`
	Disability string   `json:"disability"`
}

// OrgJSON defines the reusable wire format for officer organization data.
// This can be marshaled to JSON for responses and decoded from JSON in requests.
type OrgJSON struct {
	ID                    int       `json:"id"`
	ClubName              string    `json:"clubName"`
	Description           string    `json:"description"`
	MeetingTime           string    `json:"meetingTime"`
	ImagePath             string    `json:"imagePath"`
	ExternalLink          string    `json:"externalLink"`
	ContactInfo           string    `json:"contactInfo"`
	IncludeOfficerEmails  bool      `json:"includeOfficerEmails"`
	Officers              []string  `json:"officers"`
	Personality           []string  `json:"personality"`
	Activities            []string  `json:"activities"`
	ActivitiesDescription string    `json:"activitiesDescription"`
	Genders               []string  `json:"genders"`
	Ethnicities           []string  `json:"ethnicities"`
	Religions             []string  `json:"religions"`
	StrictGenders         bool      `json:"strict_genders"`
	DedicatedMajors       []string  `json:"dedicated_majors"`
	AssociatedMajors      []string  `json:"associated_majors"`
	Other                 []string  `json:"other"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

// OrgUpdatePayload contains only mutable club fields for PATCH requests.
type OrgUpdatePayload struct {
	ID                    int       `json:"id"`
	ClubName              *string   `json:"clubName,omitempty"`
	Description           *string   `json:"description,omitempty"`
	MeetingTime           *string   `json:"meetingTime,omitempty"`
	ImagePath             *string   `json:"imagePath,omitempty"`
	ExternalLink          *string   `json:"externalLink,omitempty"`
	ContactInfo           *string   `json:"contactInfo,omitempty"`
	IncludeOfficerEmails  *bool     `json:"includeOfficerEmails,omitempty"`
	Officers              *[]string `json:"officers,omitempty"`
	Personality           *[]string `json:"personality,omitempty"`
	Activities            *[]string `json:"activities,omitempty"`
	ActivitiesDescription *string   `json:"activitiesDescription,omitempty"`
	Genders               *[]string `json:"genders,omitempty"`
	Ethnicities           *[]string `json:"ethnicities,omitempty"`
	Religions             *[]string `json:"religions,omitempty"`
	StrictGenders         *bool     `json:"strict_genders,omitempty"`
	DedicatedMajors       *[]string `json:"dedicated_majors,omitempty"`
	AssociatedMajors      *[]string `json:"associated_majors,omitempty"`
	Other                 *[]string `json:"other,omitempty"`
}

type OrgDeletePayload struct {
	ID int `json:"id"`
}

type ResultsPayload struct {
	Results []OrgJSON `json:"results"`
}

type OrgMatchJSON struct {
	ID              int     `json:"id"`
	ClubName        string  `json:"clubName"`
	MatchPercentage float32 `json:"matchPercentage"`
}

type RoleSwitchTestPayload struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type SwipeQuestionInput struct {
	Term         string              `json:"term"`
	Definition   string              `json:"definition"`
	Translations map[string][]string `json:"translations"`
}

type TestSwipeQuestionUpdatePayload struct {
	ID           int                 `json:"id,omitempty"`
	QuestionType string              `json:"question_type"`
	Translations map[string][]string `json:"translations"`
}

// organization defines either an organization or a user with various fields
// corresponding with the questions asked on the quiz.
type Organization struct {
	name              string
	personality       []string
	activities        []string
	genders           []string
	ethnicities       []string
	religions         []string
	strict_genders    bool
	dedicated_majors  []string
	associated_majors []string
	other             []string
	max_score         float32
}

var dbClient *DatabaseClient

var repoRootDir string
var orgPhotosDir string

const testOfficerClubName = "Computer Science Club"
const maxImageUploadSize = 10 << 20
const imageURLPrefix = "/api/images/"

var testRoleSwitchAllowlist = map[string]struct{}{
	"mckallipb@southwestern.edu": {},
	"volkinm@southwestern.edu":   {},
	"kleinj@southwestern.edu":    {},
	"balakrisa@southwestern.edu": {},
}

// main wires up OAuth, database connectivity, middleware-protected routes,
// and static file serving for the Svelte frontend.
func main() {
	// Initialize in-memory session storage used by auth middleware.
	authSessionStore = NewAuthSessionStore()

	serverWorkingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to determine server working directory: %v", err)
	}
	repoRootDir = filepath.Clean(filepath.Join(serverWorkingDir, ".."))
	orgPhotosDir = filepath.Join(repoRootDir, "Components", "OrgPhotos")
	if err := os.MkdirAll(orgPhotosDir, 0o755); err != nil {
		log.Fatalf("failed to ensure image directory exists: %v", err)
	}
	staticDir := filepath.Join(repoRootDir, "Webpages", "dist")

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
	dbClient, err = NewDatabaseClient("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	RegisterDatabaseClient("default", dbClient.WithAlias("default"))
	defer dbClient.Close()

	submitRequiresAuth := readBoolEnv("SUBMIT_REQUIRES_AUTH", true)
	log.Printf("Submit auth required: %t", submitRequiresAuth)

	// Serve static files from the built Svelte distribution directory.
	if err := os.Chdir(staticDir); err != nil {
		log.Printf("warning: failed to switch static directory: %v", err)
	}

	// Register API routes and static site handler.
	router := mux.NewRouter()

	router.HandleFunc("/login", makeOAuthLoginHandler(oauthCfg)).Methods(http.MethodGet)
	router.HandleFunc("/auth/callback", makeOAuthCallbackHandler(oauthCfg)).Methods(http.MethodGet)
	router.HandleFunc("/api/user", makeCurrentUserHandler()).Methods(http.MethodGet)
	router.HandleFunc("/logout", makeLogoutHandler()).Methods(http.MethodPost)
	router.HandleFunc("/api/images/{filename}", handleClubImageFetchRequest).Methods(http.MethodGet)

	router.Handle("/response", requireAuthenticatedSession(http.HandlerFunc(handleSurveyResponseSubmission))).Methods(http.MethodPost)

	submitHandler := http.Handler(http.HandlerFunc(handleDemographicsSubmission))
	if submitRequiresAuth {
		submitHandler = requireAuthenticatedSession(submitHandler)
	}
	router.Handle("/submit", submitHandler).Methods(http.MethodPost)

	// All /api/ endpoints require an authenticated session, but only some require specific roles.
	apiRouter := router.PathPrefix("/api/").Subrouter()
	apiRouter.Use(requireAuthenticatedSession)

	apiRouter.HandleFunc("/prefill", makePrefillFieldsHandler()).Methods(http.MethodGet)

	apiRouter.Handle("/adjectives", http.HandlerFunc(handleAdjectivesRequest)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/test/role", handleTestRoleSwitchRequest).Methods(http.MethodPost)
	apiRouter.HandleFunc("/test/questions", handleTestQuestionBankUpdateRequest).Methods(http.MethodPatch)

	apiRouter.HandleFunc("/results", getUserOrgsHandler).Methods(http.MethodGet)

	apiRouter.HandleFunc("/delete", handleDeleteUserDataRequest).Methods(http.MethodPost)

	// Officer-only endpoints are nested under /api/officer/ and use additional role-checking middleware.
	officerRouter := apiRouter.PathPrefix("/officer/").Subrouter()

	officerRouter.Use(makeRequireUserRoleHandlerMiddleware(Officer))

	officerRouter.HandleFunc("/orgs", handleOfficerOrgsRequest).Methods(http.MethodGet)
	officerRouter.HandleFunc("/orgs", handlePatchOrgRequest).Methods(http.MethodPatch)
	officerRouter.HandleFunc("/orgs/{id:[0-9]+}/image", handleClubImageUploadRequest).Methods(http.MethodPost)
	officerRouter.HandleFunc("/orgs/{id:[0-9]+}/image", handleClubImageDeleteRequest).Methods(http.MethodDelete)

	// Admin-only endpoints are nested under /api/admin/ and use additional role-checking middleware.
	adminRouter := apiRouter.PathPrefix("/admin/").Subrouter()
	adminRouter.Use(makeRequireUserRoleHandlerMiddleware(Admin))

	adminRouter.HandleFunc("/orgs", handleAdminOrgsRequest).Methods(http.MethodGet)
	adminRouter.HandleFunc("/orgs", handleAdminCreateOrgRequest).Methods(http.MethodPost)
	adminRouter.HandleFunc("/orgs", handlePatchOrgRequest).Methods(http.MethodPatch)
	adminRouter.HandleFunc("/orgs", handleAdminDeleteOrgRequest).Methods(http.MethodDelete)
	adminRouter.HandleFunc("/orgs/{id:[0-9]+}/image", handleClubImageUploadRequest).Methods(http.MethodPost)
	adminRouter.HandleFunc("/orgs/{id:[0-9]+}/image", handleClubImageDeleteRequest).Methods(http.MethodDelete)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))

	// Start HTTP server.
	port := ":8080"
	fmt.Println("Server is running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, router))
}

func handleDeleteUserDataRequest(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.Header.Get("X-User-Email"))
	if email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: missing authenticated user email"})
		return
	}
	if _, err := dbClient.Query().DeleteUserDataByEmail(r.Context(), email); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User data deleted successfully"})
}

func getUserOrgsHandler(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.Header.Get("X-User-Email"))
	if email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: missing authenticated user email"})
		return
	}
	_, answers, err := dbClient.Query().FetchUserAnswersByUserEmail(r.Context(), email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "No user answers found"})
		return
	}

	_, userProfile, err := dbClient.Query().FetchUserByEmail(r.Context(), email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch user demographics"})
		return
	}

	_, clubs, err := dbClient.Query().FetchAllClubs(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch orgs"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	includeScores := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("includeScores")), "true") || strings.TrimSpace(r.URL.Query().Get("includeScores")) == "1"
	if includeScores {
		json.NewEncoder(w).Encode(getScoredClubsFromAnswers(answers, userProfile, clubs))
		return
	}

	json.NewEncoder(w).Encode(getClubsFromAnswers(answers, userProfile, clubs))
}

// getClubsFromAnswers ranks clubs for a user by converting stored answers into
// matcher input and sorting all clubs by normalized match score.
func getClubsFromAnswers(answers []DBAnswer, userProfile *ent.User, clubs []*ent.Club) []OrgJSON {
	matchAnswers := make([]matching.Answer, 0, len(answers))
	for _, a := range answers {
		matchAnswers = append(matchAnswers, matching.Answer{
			QuestionID:   a.QuestionID,
			QuestionType: a.QuestionType,
			AnswerText:   a.AnswerText,
			Translations: a.Translations,
		})
	}

	matchClubs := make([]matching.Organization, 0, len(clubs))
	clubsByID := make(map[int]OrgJSON, len(clubs))
	for _, club := range clubs {
		if club == nil {
			continue
		}

		clubsByID[club.ID] = orgJSONFromEntClub(club)
		matchClubs = append(matchClubs, matching.Organization{
			ID:               club.ID,
			Name:             club.ClubName,
			Personality:      club.Personality,
			Activities:       club.Activities,
			Genders:          club.Genders,
			Ethnicities:      club.Ethnicities,
			Religions:        club.Religions,
			StrictGenders:    club.StrictGenders,
			DedicatedMajors:  club.DedicatedMajors,
			AssociatedMajors: club.AssociatedMajors,
			Other:            club.Other,
		})
	}

	ranked := matching.Sort(buildMatchUser(matchAnswers, userProfile), matchClubs)
	ordered := make([]OrgJSON, 0, len(ranked))
	for _, result := range ranked {
		clubJSON, ok := clubsByID[result.Organization.ID]
		if !ok {
			continue
		}
		ordered = append(ordered, clubJSON)
	}

	return ordered
}

func getScoredClubsFromAnswers(answers []DBAnswer, userProfile *ent.User, clubs []*ent.Club) []OrgMatchJSON {
	matchAnswers := make([]matching.Answer, 0, len(answers))
	for _, a := range answers {
		matchAnswers = append(matchAnswers, matching.Answer{
			QuestionID:   a.QuestionID,
			QuestionType: a.QuestionType,
			AnswerText:   a.AnswerText,
			Translations: a.Translations,
		})
	}

	matchClubs := make([]matching.Organization, 0, len(clubs))
	for _, club := range clubs {
		if club == nil {
			continue
		}

		matchClubs = append(matchClubs, matching.Organization{
			ID:               club.ID,
			Name:             club.ClubName,
			Personality:      club.Personality,
			Activities:       club.Activities,
			Genders:          club.Genders,
			Ethnicities:      club.Ethnicities,
			Religions:        club.Religions,
			StrictGenders:    club.StrictGenders,
			DedicatedMajors:  club.DedicatedMajors,
			AssociatedMajors: club.AssociatedMajors,
			Other:            club.Other,
		})
	}

	ranked := matching.Sort(buildMatchUser(matchAnswers, userProfile), matchClubs)
	ordered := make([]OrgMatchJSON, 0, len(ranked))
	for _, result := range ranked {
		ordered = append(ordered, OrgMatchJSON{
			ID:              result.Organization.ID,
			ClubName:        result.Organization.Name,
			MatchPercentage: result.NormalizedScore,
		})
	}

	return ordered
}

// buildMatchUser merges swipe answers with persisted demographics so both
// influence ranking.
func buildMatchUser(answers []matching.Answer, profile *ent.User) matching.UserInfo {
	user := matching.UserFromAnswers(answers)
	if profile == nil {
		return user
	}

	if user.Name == "" {
		user.Name = extractProfileNameTag(profile.Tags)
	}

	user.Genders = mergeUniqueFoldStrings(user.Genders, profile.Genders)
	user.Ethnicities = mergeUniqueFoldStrings(user.Ethnicities, profile.Ethnicities)
	user.Religions = mergeUniqueFoldStrings(user.Religions, profile.Religions)
	user.DedicatedMajors = mergeUniqueFoldStrings(user.DedicatedMajors, profile.DedicatedMajors)
	user.Other = mergeUniqueFoldStrings(user.Other, profile.Other)

	return user
}

func mergeUniqueFoldStrings(base []string, values []string) []string {
	merged := append([]string{}, base...)
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}

		exists := false
		for _, existing := range merged {
			if strings.EqualFold(strings.TrimSpace(existing), trimmed) {
				exists = true
				break
			}
		}
		if exists {
			continue
		}

		merged = append(merged, trimmed)
	}

	return merged
}

func demographicsToOtherSignals(rawLGBTQ string, rawDisability string) []string {
	signals := []string{}
	if normalizeYesNo(rawLGBTQ) == "yes" {
		signals = append(signals, "lgbtq", "lgbtq+")
	}
	if normalizeYesNo(rawDisability) == "yes" {
		signals = append(signals, "disability")
	}
	return signals
}

func normalizeYesNo(raw string) string {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	switch normalized {
	case "yes", "no", "prefer not to say":
		return normalized
	default:
		return ""
	}
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

func handleAdminCreateOrgRequest(w http.ResponseWriter, r *http.Request) {
	var payload OrgJSON
	if !decodeJSONBody(w, r, &payload) {
		return
	}

	_, createdClub, err := dbClient.Query().CreateClubFromJSON(r.Context(), &payload)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// send email
	sendEmailToOfficers(createdClub)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(orgJSONFromEntClub(createdClub))
}

func handlePatchOrgRequest(w http.ResponseWriter, r *http.Request) {
	var payload OrgUpdatePayload
	if !decodeJSONBody(w, r, &payload) {
		return
	}
	if payload.ID <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id must be a positive integer"})
		return
	}

	_, updatedClub, err := dbClient.Query().PatchClubFromJSON(r.Context(), payload.ID, &payload)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orgJSONFromEntClub(updatedClub))
}

func handleAdminDeleteOrgRequest(w http.ResponseWriter, r *http.Request) {
	var payload OrgDeletePayload
	if !decodeJSONBody(w, r, &payload) {
		return
	}
	if payload.ID <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id must be a positive integer"})
		return
	}

	if _, err := dbClient.Query().DeleteClubByID(r.Context(), payload.ID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func handleClubImageUploadRequest(w http.ResponseWriter, r *http.Request) {
	clubID, ok := parseClubIDFromRequest(w, r)
	if !ok {
		return
	}

	userEmail := strings.TrimSpace(r.Header.Get("X-User-Email"))
	clubRecord, ok := fetchWritableClubForUser(w, r, clubID, userEmail)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxImageUploadSize+1024)
	if err := r.ParseMultipartForm(maxImageUploadSize); err != nil {
		writeJSONError(w, http.StatusBadRequest, "request must be multipart/form-data with an image field")
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "missing image file in form field 'image'")
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(io.LimitReader(file, maxImageUploadSize+1))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to read uploaded image")
		return
	}
	if len(imageBytes) == 0 {
		writeJSONError(w, http.StatusBadRequest, "uploaded image is empty")
		return
	}
	if len(imageBytes) > maxImageUploadSize {
		writeJSONError(w, http.StatusRequestEntityTooLarge, "uploaded image exceeds 10 MB limit")
		return
	}

	extension, contentType, ok := detectAllowedImageType(imageBytes)
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "unsupported image type; allowed types are jpeg, png, gif, and webp")
		return
	}

	filename := buildStoredClubImageFilename(clubRecord.ClubName, clubRecord.ID, extension)
	imagePath := filepath.Join(orgPhotosDir, filename)
	tempPath := imagePath + ".tmp"
	if err := os.WriteFile(tempPath, imageBytes, 0o644); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to save uploaded image")
		return
	}
	if err := os.Rename(tempPath, imagePath); err != nil {
		_ = os.Remove(tempPath)
		writeJSONError(w, http.StatusInternalServerError, "failed to finalize uploaded image")
		return
	}

	imageURL := imageURLPrefix + filename
	patch := OrgUpdatePayload{ImagePath: &imageURL}
	_, updatedClub, err := dbClient.Query().PatchClubFromJSON(r.Context(), clubRecord.ID, &patch)
	if err != nil {
		_ = os.Remove(imagePath)
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if oldFilename, managed := managedImageFilename(clubRecord.ImagePath); managed && oldFilename != filename {
		_ = os.Remove(filepath.Join(orgPhotosDir, oldFilename))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message":     "Club image uploaded successfully",
		"filename":    filename,
		"imagePath":   imageURL,
		"contentType": contentType,
		"club":        orgJSONFromEntClub(updatedClub),
	})
}

func handleClubImageDeleteRequest(w http.ResponseWriter, r *http.Request) {
	clubID, ok := parseClubIDFromRequest(w, r)
	if !ok {
		return
	}

	userEmail := strings.TrimSpace(r.Header.Get("X-User-Email"))
	clubRecord, ok := fetchWritableClubForUser(w, r, clubID, userEmail)
	if !ok {
		return
	}

	currentImagePath := strings.TrimSpace(clubRecord.ImagePath)
	if currentImagePath == "" {
		writeJSONError(w, http.StatusNotFound, "club does not have an image to delete")
		return
	}

	var previousImageBytes []byte
	oldFilename, managed := managedImageFilename(currentImagePath)
	if managed {
		oldFilePath := filepath.Join(orgPhotosDir, oldFilename)
		if imageBytes, err := os.ReadFile(oldFilePath); err == nil {
			previousImageBytes = imageBytes
		}
		if err := os.Remove(oldFilePath); err != nil && !os.IsNotExist(err) {
			writeJSONError(w, http.StatusInternalServerError, "failed to delete image from storage")
			return
		}
	}

	emptyImagePath := ""
	patch := OrgUpdatePayload{ImagePath: &emptyImagePath}
	_, updatedClub, err := dbClient.Query().PatchClubFromJSON(r.Context(), clubRecord.ID, &patch)
	if err != nil {
		if managed && len(previousImageBytes) > 0 {
			_ = os.WriteFile(filepath.Join(orgPhotosDir, oldFilename), previousImageBytes, 0o644)
		}
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Club image deleted successfully",
		"club":    orgJSONFromEntClub(updatedClub),
	})
}

func handleClubImageFetchRequest(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimSpace(mux.Vars(r)["filename"])
	if filename == "" || filename != filepath.Base(filename) || filename == "." {
		writeJSONError(w, http.StatusBadRequest, "invalid image filename")
		return
	}

	imagePath := filepath.Join(orgPhotosDir, filename)
	if _, err := os.Stat(imagePath); err != nil {
		if os.IsNotExist(err) {
			writeJSONError(w, http.StatusNotFound, "image not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "failed to access image")
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeFile(w, r, imagePath)
}

func parseClubIDFromRequest(w http.ResponseWriter, r *http.Request) (int, bool) {
	clubID := 0
	if _, err := fmt.Sscanf(mux.Vars(r)["id"], "%d", &clubID); err != nil || clubID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "club id must be a positive integer")
		return 0, false
	}
	return clubID, true
}

func fetchWritableClubForUser(w http.ResponseWriter, r *http.Request, clubID int, userEmail string) (*ent.Club, bool) {
	if userEmail == "" {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized: missing authenticated user email")
		return nil, false
	}

	role := dbClient.Query().getUserRole(r.Context(), userEmail)
	switch role {
	case Admin:
		_, clubRecord, err := dbClient.Query().FetchClubByID(r.Context(), clubID)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, "club not found")
			return nil, false
		}
		return clubRecord, true
	case Officer:
		_, clubs, err := dbClient.Query().FetchOfficerClubsByUserEmail(r.Context(), userEmail)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to load officer clubs")
			return nil, false
		}
		for _, clubRecord := range clubs {
			if clubRecord != nil && clubRecord.ID == clubID {
				return clubRecord, true
			}
		}
		writeJSONError(w, http.StatusForbidden, "Forbidden: officers can only manage images for their own clubs")
		return nil, false
	default:
		writeJSONError(w, http.StatusForbidden, "Forbidden")
		return nil, false
	}
}

func detectAllowedImageType(imageBytes []byte) (string, string, bool) {
	contentType := http.DetectContentType(imageBytes)
	switch contentType {
	case "image/jpeg":
		return ".jpg", contentType, true
	case "image/png":
		return ".png", contentType, true
	case "image/gif":
		return ".gif", contentType, true
	case "image/webp":
		return ".webp", contentType, true
	default:
		return "", contentType, false
	}
}

func buildStoredClubImageFilename(clubName string, clubID int, extension string) string {
	slug := slugifyClubFilename(clubName)
	if slug == "" {
		slug = "club"
	}
	return fmt.Sprintf("%s-%d-%d%s", slug, clubID, time.Now().UnixNano(), extension)
}

func slugifyClubFilename(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return ""
	}

	var builder strings.Builder
	lastWasDash := false
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
			lastWasDash = false
			continue
		}
		if lastWasDash {
			continue
		}
		builder.WriteByte('-')
		lastWasDash = true
	}

	return strings.Trim(builder.String(), "-")
}

func managedImageFilename(imagePath string) (string, bool) {
	trimmed := strings.TrimSpace(imagePath)
	if !strings.HasPrefix(trimmed, imageURLPrefix) {
		return "", false
	}
	filename := strings.TrimPrefix(trimmed, imageURLPrefix)
	if filename == "" || filename != filepath.Base(filename) || filename == "." {
		return "", false
	}
	return filename, true
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
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
		ID:                    club.ID,
		ClubName:              club.ClubName,
		Description:           club.Description,
		MeetingTime:           club.MeetingTime,
		ImagePath:             club.ImagePath,
		ExternalLink:          club.ExternalLink,
		ContactInfo:           club.ContactInfo,
		IncludeOfficerEmails:  club.IncludeOfficerEmails,
		Officers:              officerEmailsFromClub(club),
		Personality:           club.Personality,
		Activities:            club.Activities,
		ActivitiesDescription: club.ActivitiesDescription,
		Genders:               club.Genders,
		Ethnicities:           club.Ethnicities,
		Religions:             club.Religions,
		StrictGenders:         club.StrictGenders,
		DedicatedMajors:       club.DedicatedMajors,
		AssociatedMajors:      club.AssociatedMajors,
		Other:                 club.Other,

		UpdatedAt: club.UpdatedAt,
	}
}

func officerEmailsFromClub(club *ent.Club) []string {
	if club == nil || len(club.Edges.Leaders) == 0 {
		return []string{}
	}

	emails := make([]string, 0, len(club.Edges.Leaders))
	for _, leader := range club.Edges.Leaders {
		if leader == nil {
			continue
		}
		email := strings.TrimSpace(leader.Email)
		if email == "" {
			continue
		}
		emails = append(emails, email)
	}

	return emails
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
	_, adjectives, err := dbClient.Query().FetchSwipeQuestionContents(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch adjectives"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adjectives)
}

func handleTestRoleSwitchRequest(w http.ResponseWriter, r *http.Request) {
	actorEmail := strings.ToLower(strings.TrimSpace(r.Header.Get("X-User-Email")))
	if !isAllowedRoleTestEmail(actorEmail) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden: role switch endpoint is restricted to approved test users"})
		return
	}

	var payload RoleSwitchTestPayload
	if !decodeJSONBody(w, r, &payload) {
		return
	}

	targetEmail := strings.ToLower(strings.TrimSpace(payload.Email))
	if targetEmail == "" {
		targetEmail = actorEmail
	}
	if !isAllowedRoleTestEmail(targetEmail) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "target email is not allowed for role switching"})
		return
	}

	requestedRole := strings.ToLower(strings.TrimSpace(payload.Role))
	switch requestedRole {
	case "admin", "officer", "member":
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "role must be one of: admin, officer, member"})
		return
	}

	if _, err := dbClient.Query().SetUserRoleForTesting(r.Context(), targetEmail, requestedRole, testOfficerClubName); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	effectiveRole := dbClient.Query().getUserRole(r.Context(), targetEmail)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message":         "Role updated for testing",
		"actorEmail":      actorEmail,
		"targetEmail":     targetEmail,
		"requestedRole":   requestedRole,
		"effectiveRole":   effectiveRole,
		"officerClubName": testOfficerClubName,
	})
}

func handleTestQuestionBankUpdateRequest(w http.ResponseWriter, r *http.Request) {
	actorEmail := strings.ToLower(strings.TrimSpace(r.Header.Get("X-User-Email")))
	if !isAllowedRoleTestEmail(actorEmail) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden: test question update endpoint is restricted to approved test users"})
		return
	}

	var payload []TestSwipeQuestionUpdatePayload
	if !decodeJSONBody(w, r, &payload) {
		return
	}

	if len(payload) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "request body must include at least one question"})
		return
	}

	activities := make([]SwipeQuestionInput, 0, len(payload))
	personalityTraits := make([]SwipeQuestionInput, 0, len(payload))
	for i, item := range payload {
		normalizedType := normalizeTestQuestionType(item.QuestionType)
		if normalizedType == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("question at index %d has unsupported question_type", i)})
			return
		}

		if term := pickEnglishTermFromTranslations(item.Translations); term == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("question at index %d must include translations.en[0] for matching", i)})
			return
		}

		question := SwipeQuestionInput{Translations: item.Translations}
		switch normalizedType {
		case "activities":
			activities = append(activities, question)
		case "personality_traits":
			personalityTraits = append(personalityTraits, question)
		}
	}

	if len(activities) == 0 || len(personalityTraits) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "request must include both activities and personality_traits questions"})
		return
	}

	if _, err := dbClient.Query().ReplaceSwipeQuestionsForTesting(r.Context(), activities, personalityTraits); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_, questions, err := dbClient.Query().FetchSwipeQuestionContents(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Questions updated but failed to fetch refreshed question set"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(questions)
}

func normalizeTestQuestionType(raw string) string {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "_", "")
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, " ", "")

	switch normalized {
	case "activity", "activities", "adjective":
		return "activities"
	case "personality", "personalitytrait", "personalitytraits":
		return "personality_traits"
	default:
		return ""
	}
}

func isAllowedRoleTestEmail(email string) bool {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return false
	}
	_, ok := testRoleSwitchAllowlist[normalized]
	return ok
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

// handleSurveyResponseSubmission validates and stores a full replacement set of
// responses from an authenticated user.
func handleSurveyResponseSubmission(w http.ResponseWriter, r *http.Request) {
	var payload []SurveyResponsePayload
	if !decodeJSONBody(w, r, &payload) {
		return
	}

	email := strings.TrimSpace(r.Header.Get("X-User-Email"))
	if _, err := dbClient.Query().ReplaceSurveyResponsesByUserEmail(r.Context(), email, payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Return a confirmation payload so the client can verify what was stored.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"email":         email,
		"responses":     payload,
		"responseCount": len(payload),
		"message":       "Responses accepted and replaced",
	})
}

// handleDemographicsSubmission validates the demographics payload and returns
// the accepted values as confirmation.
func handleDemographicsSubmission(w http.ResponseWriter, r *http.Request) {
	var submission DemographicsPayload
	if !decodeJSONBody(w, r, &submission) {
		return
	}

	lgbtqAnswer := strings.TrimSpace(submission.LGBTQ)
	disabilityAnswer := strings.TrimSpace(submission.Disability)

	if strings.TrimSpace(submission.Gender) == "" ||
		strings.TrimSpace(submission.Religion) == "" ||
		normalizeYesNo(lgbtqAnswer) == "" ||
		normalizeYesNo(disabilityAnswer) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing one or more required demographics fields: gender, religion, lgbtq, and disability are required"})
		return
	}

	email := strings.TrimSpace(r.Header.Get("X-User-Email"))
	if email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: missing authenticated user email"})
		return
	}

	religions := []string{}
	if trimmed := strings.TrimSpace(submission.Religion); trimmed != "" {
		religions = []string{trimmed}
	}

	if _, err := dbClient.Query().UpsertUserDemographicsByEmail(
		r.Context(),
		email,
		submission.Name,
		[]string{submission.Gender},
		submission.Race,
		religions,
		submission.Major,
		demographicsToOtherSignals(lgbtqAnswer, disabilityAnswer),
	); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Echo submitted data and return persistence confirmation for the client.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message":    "Demographics submitted and saved successfully",
		"name":       submission.Name,
		"gender":     submission.Gender,
		"race":       submission.Race,
		"religion":   submission.Religion,
		"major":      submission.Major,
		"lgbtq":      lgbtqAnswer,
		"disability": disabilityAnswer,
		"email":      email,
	})
}

func sendEmailToOfficers(club *ent.Club) {
	if club == nil || len(club.Edges.Leaders) == 0 {
		return
	}

	for _, leader := range club.Edges.Leaders {
		if leader == nil {
			continue
		}
		email := strings.TrimSpace(leader.Email)
		if email == "" {
			continue
		}
		// sendEmail(email, fmt.Sprintf("Your club '%s' has been created/updated", club.ClubName), "Please check the officer portal for details.")
	}
}

func sendEmail(to, subject, body string) {
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", "youremail@email.com")
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	// Set email body
	message.SetBody("text/plain", body)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", "1a2b3c4d5e6f7g")

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
	}
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

// Source - https://stackoverflow.com/a/72498530
// Posted by mtkopone, modified by community. See post 'Timeline' for change history
// Retrieved 2026-03-26, License - CC BY-SA 4.0

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
