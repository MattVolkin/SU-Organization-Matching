package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"server-example/ent"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// host=localhost port=5432 user=dev_user password=testing

type UserPerms int

const (
	Admin UserPerms = iota
	Officer
	Member
)

type OAuthConfig struct {
	GoogleOAuth *oauth2.Config
	StateStore  *StateStore
}

type StateStore struct {
	mu     sync.RWMutex
	states map[string]OAuthState
}

type OAuthState struct {
	ExpiresAt time.Time
	Popup     bool
}

type UserSession struct {
	Email        string
	Token        string
	Created      time.Time
	LastActivity time.Time
	ProfileFields map[string]string
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*UserSession
}

type ResponsePayload struct {
	QuestionID int  `json:"questionId"`
	Answer     bool `json:"answer"`
}

type GoogleUserInfo struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
}

type DemographicsSubmission struct {
	Name     string   `json:"name"`
	Gender   string   `json:"gender"`
	Race     []string `json:"race"`
	Religion string   `json:"religion"`
	Major    []string `json:"major"`
}

var sessionStore *SessionStore

const sessionInactivityTimeout = 30 * time.Minute

func main() {
	sessionStore = &SessionStore{
		sessions: make(map[string]*UserSession),
	}

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
	log.Printf("Google OAuth client ID suffix: %s", redactSuffix(googleClientID, 10))

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

	oauthCfg := &OAuthConfig{
		GoogleOAuth: googleOAuth,
		StateStore: &StateStore{
			states: make(map[string]OAuthState),
		},
	}

	dsn := "host=localhost port=5432 user=dev_user password=testing dbname=dev_project_db"
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	submitRequiresAuth := boolFromEnv("SUBMIT_REQUIRES_AUTH", true)
	log.Printf("Submit auth required: %t", submitRequiresAuth)

	if err := os.Chdir("../Svelte Examples/plain-svelte-app/dist"); err != nil {
		log.Printf("warning: failed to switch static directory: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/login", makeLoginHandler(oauthCfg))
	mux.HandleFunc("/auth/callback", makeCallbackHandler(oauthCfg))
	mux.HandleFunc("/api/user", makeGetUserHandler())
	mux.HandleFunc("/api/prefill", makeGetPrefillHandler())
	mux.HandleFunc("/logout", makeLogoutHandler())

	mux.Handle("/response", checkUserAuthMiddleware(http.HandlerFunc(handleResponseSubmission)))

	submitHandler := checkSubmissionDataMiddleware(http.HandlerFunc(handleDemographicsSubmission))
	if submitRequiresAuth {
		submitHandler = checkUserAuthMiddleware(submitHandler)
	}
	mux.Handle("/submit", submitHandler)

	mux.Handle("/", http.FileServer(http.Dir(".")))

	port := ":8080"
	fmt.Println("Server is running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, mux))
}

func checkUserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromRequest(r)
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: missing token"})
			return
		}

		session, ok := sessionStore.getSession(token)
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

func checkUserRoleMiddleware(user string, role UserPerms, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hasRole := true
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

func handleResponseSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Only POST method is allowed"})
		return
	}

	var payload ResponsePayload
	if !checkJSON(w, r, &payload) {
		return
	}
	if payload.QuestionID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "questionId must be a positive integer"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"email":      r.Header.Get("X-User-Email"),
		"questionId": payload.QuestionID,
		"answer":     payload.Answer,
		"message":    "Response accepted",
	})
}

func handleDemographicsSubmission(w http.ResponseWriter, r *http.Request) {
	var submission DemographicsSubmission
	if !checkJSON(w, r, &submission) {
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

func checkJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	err := json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return false
	}
	return true
}

func makeLoginHandler(cfg *OAuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := generateStateToken()
		cfg.StateStore.addState(state, r.URL.Query().Get("popup") == "1")

		authURL := cfg.GoogleOAuth.AuthCodeURL(state, oauth2.AccessTypeOffline)
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	}
}

func makeCallbackHandler(cfg *OAuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")

		stateInfo, ok := cfg.StateStore.validateState(state)
		if !ok {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		token, err := cfg.GoogleOAuth.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusUnauthorized)
			fmt.Println("Token exchange error:", err)
			return
		}

		userInfo, err := getGoogleUserInfo(r.Context(), token)
		if err != nil {
			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
			fmt.Println("User info error:", err)
			return
		}

		sessionToken := generateSessionToken()
		session := &UserSession{
			Email:        userInfo.Email,
			Token:        sessionToken,
			Created:      time.Now(),
			LastActivity: time.Now(),
			ProfileFields: buildProfileFields(userInfo),
		}
		sessionStore.addSession(sessionToken, session)
		fmt.Printf("User logged in: %s\n", userInfo.Email)

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		if stateInfo.Popup {
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Successfully authenticated",
			"email":   userInfo.Email,
			"token":   sessionToken,
		})
	}
}

func makeGetUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromRequest(r)
		if token == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		session, ok := sessionStore.getSession(token)
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

func makeLogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromRequest(r)
		if token != "" {
			sessionStore.removeSession(token)
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

func makeGetPrefillHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromRequest(r)
		if token == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		session, ok := sessionStore.getSession(token)
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

func getGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
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

	var userInfo GoogleUserInfo
	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(userInfo.Email) == "" {
		return nil, fmt.Errorf("google userinfo returned empty email")
	}

	return &userInfo, nil
}

func buildProfileFields(userInfo *GoogleUserInfo) map[string]string {
	fields := map[string]string{}
	setIfNotEmpty(fields, "email", userInfo.Email)
	setIfNotEmpty(fields, "name", userInfo.Name)
	setIfNotEmpty(fields, "firstName", userInfo.GivenName)
	setIfNotEmpty(fields, "lastName", userInfo.FamilyName)
	setIfNotEmpty(fields, "avatarUrl", userInfo.Picture)
	setIfNotEmpty(fields, "locale", userInfo.Locale)
	return fields
}

func setIfNotEmpty(target map[string]string, key string, value string) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return
	}
	target[key] = trimmed
}

func tokenFromRequest(r *http.Request) string {
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

func generateStateToken() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func generateSessionToken() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().UnixMicro())
}

func (ss *StateStore) addState(state string, popup bool) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.states[state] = OAuthState{
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Popup:     popup,
	}
}

func (ss *StateStore) validateState(state string) (OAuthState, bool) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	stateInfo, ok := ss.states[state]
	if !ok {
		return OAuthState{}, false
	}
	if time.Now().After(stateInfo.ExpiresAt) {
		delete(ss.states, state)
		return OAuthState{}, false
	}

	delete(ss.states, state)
	return stateInfo, true
}

func (sess *SessionStore) addSession(token string, session *UserSession) {
	sess.mu.Lock()
	defer sess.mu.Unlock()
	sess.sessions[token] = session
}

func (sess *SessionStore) getSession(token string) (*UserSession, bool) {
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

func (sess *SessionStore) removeSession(token string) {
	sess.mu.Lock()
	defer sess.mu.Unlock()
	delete(sess.sessions, token)
}

func redactSuffix(value string, keep int) string {
	if keep <= 0 || len(value) <= keep {
		return "[redacted]"
	}
	return "..." + value[len(value)-keep:]
}

func boolFromEnv(key string, defaultValue bool) bool {
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
