package main

// Notes API Backend
// API Documentation: /app/docs/API.md
// This implementation follows the specification defined in /app/docs/API.md

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Note struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
	UserID    int64  `json:"user_id,omitempty"`
}

type CreateNoteRequest struct {
	Text string `json:"text"`
}

type NotesResponse struct {
	Notes []Note `json:"notes"`
}

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Get bot token from environment
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Println("Warning: TELEGRAM_BOT_TOKEN not set, authentication will be disabled")
	}

	// Initialize storage
	storage := NewMemoryStorage()

	// Create router
	router := mux.NewRouter()

	// Health check endpoint (no auth required)
	router.HandleFunc("/health", HealthCheckHandler()).Methods("GET")

	// Create API subrouter
	api := router.PathPrefix("/api").Subrouter()

	// Add middleware
	if botToken != "" {
		api.Use(AuthMiddleware(botToken))
	}

	// Register routes
	api.HandleFunc("/notes", GetNotesHandler(storage)).Methods("GET", "OPTIONS")
	api.HandleFunc("/notes", CreateNoteHandler(storage)).Methods("POST", "OPTIONS")
	api.HandleFunc("/notes/{id}", DeleteNoteHandler(storage)).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/notes", DeleteAllNotesHandler(storage)).Methods("DELETE", "OPTIONS")

	// Setup CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // In production, specify your frontend URL
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Init-Data"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)

	// Start server
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

// GetNotesHandler returns all notes for the user
func GetNotesHandler(storage *MemoryStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserIDFromContext(r)
		log.Printf("[GET /api/notes] User ID: %d - Fetching notes", userID)
		
		notes := storage.GetNotes(userID)
		log.Printf("[GET /api/notes] User ID: %d - Returning %d notes", userID, len(notes))

		response := NotesResponse{Notes: notes}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// CreateNoteHandler creates a new note
func CreateNoteHandler(storage *MemoryStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserIDFromContext(r)
		log.Printf("[POST /api/notes] User ID: %d - Creating new note", userID)
		
		var req CreateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[POST /api/notes] User ID: %d - ERROR: Invalid request body: %v", userID, err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Trim whitespace and validate (per API spec: whitespace-only text is considered empty)
		trimmedText := strings.TrimSpace(req.Text)
		if trimmedText == "" {
			log.Printf("[POST /api/notes] User ID: %d - ERROR: Empty or whitespace-only text field", userID)
			http.Error(w, "Text is required", http.StatusBadRequest)
			return
		}

		note := storage.CreateNote(userID, trimmedText)
		log.Printf("[POST /api/notes] User ID: %d - Note created: ID=%s, Text length=%d", userID, note.ID, len(trimmedText))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	}
}

// DeleteNoteHandler deletes a specific note
func DeleteNoteHandler(storage *MemoryStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		noteID := vars["id"]
		userID := GetUserIDFromContext(r)
		
		log.Printf("[DELETE /api/notes/%s] User ID: %d - Attempting to delete note", noteID, userID)

		if err := storage.DeleteNote(userID, noteID); err != nil {
			log.Printf("[DELETE /api/notes/%s] User ID: %d - ERROR: %v", noteID, userID, err)
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}

		log.Printf("[DELETE /api/notes/%s] User ID: %d - Note deleted successfully", noteID, userID)
		w.WriteHeader(http.StatusNoContent)
	}
}

// DeleteAllNotesHandler deletes all notes for the user
func DeleteAllNotesHandler(storage *MemoryStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserIDFromContext(r)
		log.Printf("[DELETE /api/notes] User ID: %d - Deleting all notes", userID)
		
		// Get count before deletion for logging
		notesBefore := storage.GetNotes(userID)
		count := len(notesBefore)
		
		storage.DeleteAllNotes(userID)
		log.Printf("[DELETE /api/notes] User ID: %d - Deleted %d notes", userID, count)

		w.WriteHeader(http.StatusNoContent)
	}
}

// HealthCheckHandler returns the health status of the service
func HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[GET /health] Health check requested")

		health := map[string]interface{}{
			"status":  "healthy",
			"service": "tma-notes-api",
			"version": "1.0.0",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(health)

		log.Printf("[GET /health] Status: healthy")
	}
}
