package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	validToken    = regexp.MustCompile(`^[a-fA-F0-9]+$`)
	validUsername = regexp.MustCompile(`^[a-z0-9]+$`)
)

type Auth struct {
	*Database
}

// Protect is a middleware that checks the authentication token.
func (auth *Auth) Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from the request.
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		if !isValidToken(token) {
			http.Error(w, "Token not found", http.StatusUnauthorized)
			return
		}

		// Get username that created the token.
		user, err := auth.CheckToken(r.Context(), token)
		if err != nil {
			if err == ErrInvalidUser {
				http.Error(w, "Token not found", http.StatusUnauthorized)
				return
			}
			log.Printf("Database.CheckToken: %v", err)
			http.Error(w, "Cannot read token", http.StatusInternalServerError)
			return
		}

		// User is authenticated.
		next.ServeHTTP(w, r.WithContext(ContextWithUser(r.Context(), user)))
	})
}

type HTTPLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// HTTPLogin handles login form POST requests.
func (auth *Auth) HTTPLogin(w http.ResponseWriter, r *http.Request) {
	// Read request body.
	var req HTTPLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || !isValidUsername(req.Username) {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	// Check username and password.
	user, err := auth.CheckPassword(r.Context(), req.Username, req.Password)
	if err != nil {
		if err == ErrInvalidUser {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}
		log.Printf("Database.CheckPassword: %v", err)
		http.Error(w, "Cannot read password", http.StatusInternalServerError)
		return
	}

	// Reset the token.
	token, err := auth.ResetToken(r.Context(), user.Name)
	if err != nil {
		log.Printf("Database.ResetToken: %v", err)
		http.Error(w, "Cannot reset token", http.StatusInternalServerError)
		return
	}

	// Send token to user.
	res := map[string]string{"token": token}
	json.NewEncoder(w).Encode(res)
}

// HTTPLogout handles logout requests.
func (auth *Auth) HTTPLogout(w http.ResponseWriter, r *http.Request) {
	// Get username from the request.
	user, ok := ContextUser(r.Context())
	if !ok {
		// User is already logged out.
		return
	}

	// Reset the token.
	_, err := auth.ResetToken(r.Context(), user.Name)
	if err != nil {
		log.Printf("Database.ResetToken: %v", err)
		http.Error(w, "Cannot reset token", http.StatusInternalServerError)
		return
	}
}

// ResetToken resets the user token in the database.
func (auth *Auth) ResetToken(ctx context.Context, username string) (string, error) {
	newToken := generateToken()
	return newToken, auth.Database.ResetToken(ctx, username, newToken)
}

// generateToken returns 32 random bytes encoded as a hex string.
func generateToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}

func isValidUsername(username string) bool {
	return len(username) >= 1 && len(username) <= 128 && validUsername.MatchString(username)
}

func isValidToken(token string) bool {
	return len(token) >= 1 && len(token) <= 1024 && validToken.MatchString(token)
}
