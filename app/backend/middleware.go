package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type contextKey string

const userIDKey contextKey = "userID"

// AuthMiddleware validates Telegram init data and extracts user information
// Supports both standard validation (with bot token) and third-party validation (with bot ID)
func AuthMiddleware(botToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for OPTIONS requests (CORS preflight)
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			// Get init data from header
			initDataRaw := r.Header.Get("X-Init-Data")
			log.Printf("[Auth] Processing request: %s %s", r.Method, r.URL.Path)
			log.Printf("[Auth] Init data present: %v (length: %d)", initDataRaw != "", len(initDataRaw))
			
			if initDataRaw == "" {
				// For development: allow requests without init data
				log.Println("[Auth] Warning: No X-Init-Data header, using default user ID: 12345")
				ctx := context.WithValue(r.Context(), userIDKey, int64(12345))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Try to validate using the appropriate method
			var userID int64
			var validationErr error

			// First, try third-party validation (for webapp/init)
			// This is for apps launched via attachment menu or direct links
			botIDStr := os.Getenv("TELEGRAM_BOT_ID")
			if botIDStr != "" {
				log.Printf("[Auth] Attempting third-party validation with bot ID: %s", botIDStr)
				botID, err := strconv.ParseInt(botIDStr, 10, 64)
				if err == nil {
					expIn := 24 * time.Hour
					validationErr = initdata.ValidateThirdParty(initDataRaw, botID, expIn)
					if validationErr == nil {
						log.Println("[Auth] ✓ Third-party validation successful")
						// Parse to get user ID
						data, parseErr := initdata.Parse(initDataRaw)
						if parseErr == nil && data.User.ID != 0 {
							userID = data.User.ID
							log.Printf("[Auth] ✓ Extracted user ID from third-party validation: %d", userID)
						} else {
							log.Printf("[Auth] ✗ Failed to parse user ID: %v", parseErr)
						}
					} else {
						log.Printf("[Auth] ✗ Third-party validation failed: %v", validationErr)
					}
				} else {
					log.Printf("[Auth] ✗ Failed to parse bot ID: %v", err)
				}
			}

			// If third-party validation failed or wasn't configured, try standard validation
			if userID == 0 && botToken != "" {
				log.Printf("[Auth] Attempting standard validation with bot token")
				expIn := 24 * time.Hour
				validationErr = initdata.Validate(initDataRaw, botToken, expIn)
				if validationErr == nil {
					log.Println("[Auth] ✓ Standard validation successful")
					// Parse to get user ID
					data, parseErr := initdata.Parse(initDataRaw)
					if parseErr == nil && data.User.ID != 0 {
						userID = data.User.ID
						log.Printf("[Auth] ✓ Extracted user ID from standard validation: %d", userID)
					} else {
						log.Printf("[Auth] ✗ Failed to parse user ID: %v", parseErr)
					}
				} else {
					log.Printf("[Auth] ✗ Standard validation failed: %v", validationErr)
				}
			}

			// Check validation result
			if validationErr != nil {
				log.Printf("[Auth] ✗ All validation methods failed: %v", validationErr)
				http.Error(w, "Unauthorized: Invalid init data", http.StatusUnauthorized)
				return
			}

			if userID == 0 {
				log.Println("[Auth] ✗ No user ID found in init data")
				http.Error(w, "Unauthorized: No user data", http.StatusUnauthorized)
				return
			}

			log.Printf("[Auth] ✓ Request authenticated for user ID: %d", userID)
			// Add user ID to context
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(r *http.Request) int64 {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		return 12345 // Default user for development/unauthenticated requests
	}
	return userID
}
