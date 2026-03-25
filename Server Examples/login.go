package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
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
	Email         string
	Token         string
	Created       time.Time
	LastActivity  time.Time
	ProfileFields map[string]string
}

// AuthSessionStore holds active user sessions keyed by token.
type AuthSessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*AuthSession
}

// GoogleProfile contains a minimal user profile returned by Google or loaded
// from the local database.
type GoogleProfile struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

var authSessionStore *AuthSessionStore

const sessionInactivityTimeout = 30 * time.Minute

// NewAuthSessionStore creates an initialized in-memory auth session store.
func NewAuthSessionStore() *AuthSessionStore {
	return &AuthSessionStore{sessions: make(map[string]*AuthSession)}
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
		dbState, userInfo, err := dbClient.Query().Filter("google_id", googleID).FetchUserProfileByGoogleID(r.Context(), "")
		if err != nil {
			fmt.Println("User not found in DB, fetching from Google:", err)
			userInfo, err = fetchGoogleProfile(r.Context(), token)
			if err != nil {
				http.Error(w, "Failed to get user info", http.StatusInternalServerError)
				fmt.Println("User info error:", err)
				return
			}
			dbState, err = dbState.UpsertUserProfileByGoogleID(r.Context(), googleID, userInfo)
			if err != nil {
				http.Error(w, "Failed to persist user info", http.StatusInternalServerError)
				fmt.Println("Persist user info error:", err)
				return
			}
		}
		_ = dbState

		sessionToken := generateAuthSessionToken()
		session := &AuthSession{
			Email:         userInfo.Email,
			Token:         sessionToken,
			Created:       time.Now(),
			LastActivity:  time.Now(),
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

// makeCurrentUserHandler returns the authenticated user's email and session info.
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
