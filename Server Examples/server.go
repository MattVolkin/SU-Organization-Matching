
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"server-example/ent"
	"server-example/ent/user"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
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

// OAuthRuntimeConfig bundles the Google OAuth client configuration together with
// temporary state storage used to defend against CSRF during login.
type OAuthRuntimeConfig struct {
	GoogleOAuth *oauth2.Config
	StateStore  *OAuthStateStore
}

// OAuthStateStore keeps track of short-lived OAuth state tokens created at login
// time and consumed once the callback is received.
type OAuthStateStore struct {
	mu     sync.RWMutex
	states map[string]OAuthStateRecord
}

// OAuthStateRecord stores metadata associated with one OAuth state token.
type OAuthStateRecord struct {
	ExpiresAt time.Time
	Popup     bool
}

// AuthSession is the in-memory session record for an authenticated user.
// It also carries a small set of profile fields for frontend prefill.
type AuthSession struct {
	Email        string
	Token        string
	Created      time.Time
	LastActivity time.Time
	ProfileFields map[string]string
}

// AuthSessionStore holds active user sessions keyed by token.
type AuthSessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*AuthSession
}

// SurveyResponsePayload is the request body accepted by the /response endpoint.
type SurveyResponsePayload struct {
	QuestionID int  `json:"questionId"`
	Answer     bool `json:"answer"`
}

// GoogleProfile contains a minimal user profile returned by Google or loaded
// from the local database.
type GoogleProfile struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
}

// DemographicsPayload is the expected payload for the /submit endpoint.
type DemographicsPayload struct {
	Name     string   `json:"name"`
	Gender   string   `json:"gender"`
	Race     []string `json:"race"`
	Religion string   `json:"religion"`
	Major    []string `json:"major"`
}

var authSessionStore *AuthSessionStore

const sessionInactivityTimeout = 30 * time.Minute

var entClient *ent.Client

// main wires up OAuth, database connectivity, middleware-protected routes,
// and static file serving for the Svelte frontend.
func main() {
	// Initialize in-memory session storage used by auth middleware.
	authSessionStore = &AuthSessionStore{
		sessions: make(map[string]*AuthSession),
	}

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
	entClient, err = ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer entClient.Close()

	submitRequiresAuth := readBoolEnv("SUBMIT_REQUIRES_AUTH", true)
	log.Printf("Submit auth required: %t", submitRequiresAuth)

	// Serve static files from the built Svelte distribution directory.
	if err := os.Chdir("../Svelte Examples/plain-svelte-app/dist"); err != nil {
		log.Printf("warning: failed to switch static directory: %v", err)
	}

	// Register API routes and static site handler.
	mux := http.NewServeMux()

	mux.HandleFunc("/login", makeOAuthLoginHandler(oauthCfg))
	mux.HandleFunc("/auth/callback", makeOAuthCallbackHandler(oauthCfg))
	mux.HandleFunc("/api/user", makeCurrentUserHandler())
	mux.HandleFunc("/api/prefill", makePrefillFieldsHandler())
	mux.HandleFunc("/logout", makeLogoutHandler())

	mux.Handle("/response", requireAuthenticatedSession(http.HandlerFunc(handleSurveyResponseSubmission)))

	submitHandler := requirePostMethod(http.HandlerFunc(handleDemographicsSubmission))
	if submitRequiresAuth {
		submitHandler = requireAuthenticatedSession(submitHandler)
	}
	mux.Handle("/submit", submitHandler)

	mux.Handle("/", http.FileServer(http.Dir(".")))

	// Start HTTP server.
	port := ":8080"
	fmt.Println("Server is running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, mux))
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

// requireUserRole is a placeholder for role-based authorization.
// It currently always allows access.
func requireUserRole(userEmail string, requiredRole UserRole, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = userEmail
		_ = requiredRole
		hasRole := true
		if !hasRole {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden"})
			return
		}
		next.ServeHTTP(w, r)
	})
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

// makeDatabaseErrorResponder returns a reusable error handler that writes a generic
// 500 response while logging the concrete database error server-side.
func makeDatabaseErrorResponder(w http.ResponseWriter) func(error) {
	return func(err error) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		fmt.Println(err)
	}
}

// runDatabaseQuery executes a query callback and routes any error to errFunc,
// returning the zero value for T when an error occurs.
func runDatabaseQuery[T any](queryFunc func() (T, error), errFunc func(error)) T {
	res, err := queryFunc()
	if err != nil {
		errFunc(err)
		var zero T
		return zero
	}
	return res
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

// makeOAuthLoginHandler starts the Google OAuth flow by generating state and
// redirecting the browser to Google's consent page.
func makeOAuthLoginHandler(oauthConfig *OAuthRuntimeConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := generateOAuthStateToken()
		oauthConfig.StateStore.addState(state, r.URL.Query().Get("popup") == "1")

		authURL := oauthConfig.GoogleOAuth.AuthCodeURL(state, oauth2.AccessTypeOffline)
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	}
}

// makeOAuthCallbackHandler completes OAuth login, resolves user identity,
// upserts user data, and creates an application session.
func makeOAuthCallbackHandler(oauthConfig *OAuthRuntimeConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")

		stateInfo, ok := oauthConfig.StateStore.validateState(state)
		if !ok {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		token, err := oauthConfig.GoogleOAuth.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusUnauthorized)
			fmt.Println("Token exchange error:", err)
			return
		}

		googleID, err := extractGoogleSubjectFromToken(token)
		if err != nil {
			http.Error(w, "Failed to identify user", http.StatusUnauthorized)
			fmt.Println("Google token subject error:", err)
			return
		}

		// Prefer existing database data; otherwise bootstrap from Google userinfo.
		userInfo, err := fetchUserProfileByGoogleID(r.Context(), googleID)
		if err != nil {
			fmt.Println("User not found in DB, fetching from Google:", err)
			userInfo, err = fetchGoogleProfile(r.Context(), token)
			if err != nil {
				http.Error(w, "Failed to get user info", http.StatusInternalServerError)
				fmt.Println("User info error:", err)
				return
			}
			if err := upsertUserProfileByGoogleID(r.Context(), googleID, userInfo); err != nil {
				http.Error(w, "Failed to persist user info", http.StatusInternalServerError)
				fmt.Println("Persist user info error:", err)
				return
			}
		}

		sessionToken := generateAuthSessionToken()
		session := &AuthSession{
			Email:        userInfo.Email,
			Token:        sessionToken,
			Created:      time.Now(),
			LastActivity: time.Now(),
			ProfileFields: buildPrefillFields(userInfo),
		}
		authSessionStore.addSession(sessionToken, session)
		fmt.Printf("User logged in: %s\n", userInfo.Email)

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		if stateInfo.Popup {
			// Popup mode reports success to the opener and closes itself.
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `<!doctype html>
<html>
  <body>
    <script>
      (function () {
        const payload = {
          type: "google-auth-success",
          email: %q,
          token: %q
        };
        if (window.opener && !window.opener.closed) {
          window.opener.postMessage(payload, window.location.origin);
          window.close();
          return;
        }
        document.body.textContent = "Login succeeded. You can close this window.";
      })();
    </script>
  </body>
</html>`, userInfo.Email, sessionToken)
			return
		}

		// Non-popup mode returns JSON for clients that call this endpoint directly.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Successfully authenticated",
			"email":   userInfo.Email,
			"token":   sessionToken,
		})
	}
}

// makeCurrentUserHandler returns the authenticated user's email.
func makeCurrentUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractSessionTokenFromRequest(r)
		if token == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		session, ok := authSessionStore.getSession(token)
		if !ok {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"email": session.Email,
		})
	}
}

// makeLogoutHandler clears the in-memory session and expires the auth cookie.
func makeLogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractSessionTokenFromRequest(r)
		if token != "" {
			authSessionStore.removeSession(token)
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
	}
}

// makePrefillFieldsHandler returns session-backed profile fields for client-side
// form prefill.
func makePrefillFieldsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractSessionTokenFromRequest(r)
		if token == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		session, ok := authSessionStore.getSession(token)
		if !ok {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"fields": session.ProfileFields,
		})
	}
}

// fetchUserProfileByGoogleID fetches the stored user profile for a Google
// subject ID.
func fetchUserProfileByGoogleID(ctx context.Context, googleID string) (*GoogleProfile, error) {
	if entClient == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	googleID = strings.TrimSpace(googleID)
	if googleID == "" {
		return nil, fmt.Errorf("google id is required")
	}

	storedUser, err := entClient.User.Query().Where(user.GoogleIDEQ(googleID)).Only(ctx)
	if err != nil {
		return nil, err
	}

	return &GoogleProfile{
		Email: storedUser.Email,
		Name:  extractProfileNameTag(storedUser.Tags),
	}, nil
}

// upsertUserProfileByGoogleID upserts user information using Google subject ID as the
// stable lookup key.
func upsertUserProfileByGoogleID(ctx context.Context, googleID string, userInfo *GoogleProfile) error {
	println("saving user info to database")
	if entClient == nil {
		return fmt.Errorf("database not initialized")
	}
	if userInfo == nil {
		return fmt.Errorf("user info is required")
	}

	googleID = strings.TrimSpace(googleID)
	email := strings.TrimSpace(userInfo.Email)
	if googleID == "" {
		return fmt.Errorf("google id is required")
	}
	if email == "" {
		return fmt.Errorf("email is required")
	}

	tags := mergeProfileNameTag(nil, userInfo.Name)

	// Update existing record when present.
	existing, err := entClient.User.Query().Where(user.GoogleIDEQ(googleID)).Only(ctx)
	if err == nil {
		tags = mergeProfileNameTag(existing.Tags, userInfo.Name)
		return entClient.User.UpdateOneID(existing.ID).
			SetEmail(email).
			SetTags(tags).
			Exec(ctx)
	}

	if !ent.IsNotFound(err) {
		return err
	}

	// Otherwise create a new user record.
	_, createErr := entClient.User.Create().
		SetGoogleID(googleID).
		SetEmail(email).
		SetTags(tags).
		Save(ctx)
	return createErr
}

// fetchGoogleProfile calls Google userinfo endpoint with the OAuth access token
// and returns a minimal profile.
func fetchGoogleProfile(ctx context.Context, token *oauth2.Token) (*GoogleProfile, error) {
	println("Fetching user info from Google API")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo request failed: %s", res.Status)
	}

	var userInfo GoogleProfile
	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(userInfo.Email) == "" {
		return nil, fmt.Errorf("google userinfo returned empty email")
	}

	return &userInfo, nil
}

// buildPrefillFields prepares lightweight key-value fields used by the
// frontend to prefill forms.
func buildPrefillFields(userInfo *GoogleProfile) map[string]string {
	fields := map[string]string{}
	setTrimmedIfNotEmpty(fields, "email", userInfo.Email)
	setTrimmedIfNotEmpty(fields, "name", userInfo.Name)
	return fields
}

// extractGoogleSubjectFromToken extracts the stable Google subject claim (sub) from
// the OAuth id_token payload.
func extractGoogleSubjectFromToken(token *oauth2.Token) (string, error) {
	if token == nil {
		return "", fmt.Errorf("token is nil")
	}
	idTokenRaw, ok := token.Extra("id_token").(string)
	if !ok || strings.TrimSpace(idTokenRaw) == "" {
		return "", fmt.Errorf("id_token missing from oauth token")
	}

	parts := strings.Split(idTokenRaw, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid id_token format")
	}

	// JWT payload is base64url encoded in the second segment.
	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid id_token payload: %w", err)
	}

	var claims struct {
		Sub string `json:"sub"`
	}
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return "", fmt.Errorf("unable to parse id_token claims: %w", err)
	}

	claims.Sub = strings.TrimSpace(claims.Sub)
	if claims.Sub == "" {
		return "", fmt.Errorf("id_token missing sub claim")
	}

	return claims.Sub, nil
}

// extractProfileNameTag reads the first profile_name tag and returns its value.
func extractProfileNameTag(tags []string) string {
	const prefix = "profile_name="
	for _, tag := range tags {
		if strings.HasPrefix(tag, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(tag, prefix))
		}
	}
	return ""
}

// mergeProfileNameTag removes any existing profile_name tag and appends the new
// one when a non-empty name is provided.
func mergeProfileNameTag(existing []string, name string) []string {
	const prefix = "profile_name="
	filtered := make([]string, 0, len(existing)+1)
	for _, tag := range existing {
		if strings.HasPrefix(tag, prefix) {
			continue
		}
		filtered = append(filtered, tag)
	}
	if trimmed := strings.TrimSpace(name); trimmed != "" {
		filtered = append(filtered, prefix+trimmed)
	}
	return filtered
}

// setTrimmedIfNotEmpty writes a trimmed value into a map only when it is not blank.
func setTrimmedIfNotEmpty(target map[string]string, key string, value string) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return
	}
	target[key] = trimmed
}

// extractSessionTokenFromRequest extracts a session token from Authorization header first,
// then falls back to the session_token cookie.
func extractSessionTokenFromRequest(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader != "" {
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			return strings.TrimSpace(authHeader[7:])
		}
		return authHeader
	}

	if cookie, err := r.Cookie("session_token"); err == nil {
		return strings.TrimSpace(cookie.Value)
	}

	return ""
}

// generateOAuthStateToken creates a simple unique value for OAuth state.
func generateOAuthStateToken() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// generateAuthSessionToken creates a unique session token for in-memory sessions.
func generateAuthSessionToken() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().UnixMicro())
}

// addState stores a new OAuth state value with a short expiration window.
func (ss *OAuthStateStore) addState(state string, popup bool) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.states[state] = OAuthStateRecord{
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Popup:     popup,
	}
}

// validateState checks presence and expiration, then consumes the state so it
// cannot be reused.
func (ss *OAuthStateStore) validateState(state string) (OAuthStateRecord, bool) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	stateInfo, ok := ss.states[state]
	if !ok {
		return OAuthStateRecord{}, false
	}
	if time.Now().After(stateInfo.ExpiresAt) {
		delete(ss.states, state)
		return OAuthStateRecord{}, false
	}

	delete(ss.states, state)
	return stateInfo, true
}

// addSession inserts or replaces a session entry by token.
func (sess *AuthSessionStore) addSession(token string, session *AuthSession) {
	sess.mu.Lock()
	defer sess.mu.Unlock()
	sess.sessions[token] = session
}

// getSession returns a live session and updates last-activity timestamp.
// Expired sessions are removed eagerly.
func (sess *AuthSessionStore) getSession(token string) (*AuthSession, bool) {
	sess.mu.Lock()
	defer sess.mu.Unlock()
	session, ok := sess.sessions[token]
	if !ok {
		return nil, false
	}
	if time.Since(session.LastActivity) > sessionInactivityTimeout {
		delete(sess.sessions, token)
		return nil, false
	}
	session.LastActivity = time.Now()
	return session, ok
}

// removeSession deletes a session token from the in-memory store.
func (sess *AuthSessionStore) removeSession(token string) {
	sess.mu.Lock()
	defer sess.mu.Unlock()
	delete(sess.sessions, token)
}

// redactValueSuffix hides sensitive values while preserving a short suffix for
// debugging logs.
func redactValueSuffix(value string, keep int) string {
	if keep <= 0 || len(value) <= keep {
		return "[redacted]"
	}
	return "..." + value[len(value)-keep:]
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
