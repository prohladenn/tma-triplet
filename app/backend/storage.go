package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// MemoryStorage provides in-memory storage for notes
type MemoryStorage struct {
	mu    sync.RWMutex
	notes map[int64][]Note // userID -> notes
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		notes: make(map[int64][]Note),
	}
}

// GetNotes returns all notes for a user
func (s *MemoryStorage) GetNotes(userID int64) []Note {
	s.mu.RLock()
	defer s.mu.RUnlock()

	notes, exists := s.notes[userID]
	if !exists {
		return []Note{}
	}

	// Return a copy to avoid external modifications
	result := make([]Note, len(notes))
	copy(result, notes)
	return result
}

// CreateNote creates a new note for a user
func (s *MemoryStorage) CreateNote(userID int64, text string) Note {
	s.mu.Lock()
	defer s.mu.Unlock()

	note := Note{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Text:      text,
		Timestamp: time.Now().UnixMilli(),
		UserID:    userID,
	}

	if _, exists := s.notes[userID]; !exists {
		s.notes[userID] = []Note{}
	}

	s.notes[userID] = append([]Note{note}, s.notes[userID]...)
	return note
}

// DeleteNote deletes a specific note for a user
func (s *MemoryStorage) DeleteNote(userID int64, noteID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	notes, exists := s.notes[userID]
	if !exists {
		return errors.New("note not found")
	}

	for i, note := range notes {
		if note.ID == noteID {
			s.notes[userID] = append(notes[:i], notes[i+1:]...)
			return nil
		}
	}

	return errors.New("note not found")
}

// DeleteAllNotes deletes all notes for a user
func (s *MemoryStorage) DeleteAllNotes(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.notes[userID] = []Note{}
}
