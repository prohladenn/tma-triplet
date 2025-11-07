import { notesApi } from "@/services/api/notes.api";
import type { Note } from "@/types/note";
import { useEffect, useState } from "react";

const STORAGE_KEY = "tma-notes";

interface UseNotesReturn {
  notes: Note[];
  isLoading: boolean;
  error: string | null;
  isBackendAvailable: boolean;
  addNote: (text: string) => Promise<void>;
  deleteNote: (id: string) => Promise<void>;
  deleteAllNotes: () => Promise<void>;
  refreshNotes: () => Promise<void>;
}

export const useNotes = (): UseNotesReturn => {
  const [notes, setNotes] = useState<Note[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isBackendAvailable, setIsBackendAvailable] = useState(true);

  // Load notes from backend or fallback to localStorage
  const loadNotes = async () => {
    console.log("[useNotes] Starting to load notes...");
    setIsLoading(true);
    setError(null);

    try {
      console.log("[useNotes] Attempting to fetch notes from API...");
      const notesFromApi = await notesApi.getNotes();
      console.log(
        "[useNotes] ✓ Successfully loaded notes from API:",
        notesFromApi.length,
        "notes"
      );
      setNotes(notesFromApi);
      setIsBackendAvailable(true);
      // Sync to localStorage as backup
      localStorage.setItem(STORAGE_KEY, JSON.stringify(notesFromApi));
      console.log("[useNotes] ✓ Synced notes to localStorage");
    } catch (err) {
      // Fallback to localStorage if API fails
      console.error("[useNotes] ✗ Failed to load notes from API:", err);
      setIsBackendAvailable(false);
      console.warn("Failed to load notes from API, using localStorage", err);
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored) {
        try {
          const parsedNotes = JSON.parse(stored);
          console.log(
            "[useNotes] ✓ Loaded notes from localStorage:",
            parsedNotes.length,
            "notes"
          );
          setNotes(parsedNotes);
        } catch (parseError) {
          console.error("Failed to parse notes from localStorage", parseError);
          setError("Failed to load notes");
        }
      } else {
        console.log("[useNotes] No notes in localStorage");
      }
    } finally {
      setIsLoading(false);
      console.log(
        "[useNotes] Loading complete. Backend available:",
        isBackendAvailable
      );
    }
  };

  // Load notes on mount
  useEffect(() => {
    loadNotes();
  }, []);

  const addNote = async (text: string) => {
    if (!text.trim()) return;

    console.log(
      "[useNotes] Adding note:",
      text.trim().substring(0, 50) + "..."
    );
    setIsLoading(true);
    setError(null);

    try {
      const newNote = await notesApi.createNote({ text: text.trim() });
      console.log("[useNotes] ✓ Note created via API:", newNote.id);
      const updatedNotes = [newNote, ...notes];
      setNotes(updatedNotes);
      localStorage.setItem(STORAGE_KEY, JSON.stringify(updatedNotes));
    } catch (err) {
      // Fallback to localStorage-only mode
      console.error("[useNotes] ✗ Failed to create note via API:", err);
      console.warn(
        "Failed to create note via API, using localStorage only",
        err
      );
      const newNote: Note = {
        id: Date.now().toString(),
        text: text.trim(),
        timestamp: Date.now(),
      };
      console.log("[useNotes] ✓ Note created locally:", newNote.id);
      const updatedNotes = [newNote, ...notes];
      setNotes(updatedNotes);
      localStorage.setItem(STORAGE_KEY, JSON.stringify(updatedNotes));
      setError("Note saved locally only");
    } finally {
      setIsLoading(false);
    }
  };

  const deleteNote = async (id: string) => {
    setIsLoading(true);
    setError(null);

    try {
      await notesApi.deleteNote(id);
      const updatedNotes = notes.filter((note) => note.id !== id);
      setNotes(updatedNotes);
      localStorage.setItem(STORAGE_KEY, JSON.stringify(updatedNotes));
    } catch (err) {
      console.warn(
        "Failed to delete note via API, using localStorage only",
        err
      );
      const updatedNotes = notes.filter((note) => note.id !== id);
      setNotes(updatedNotes);
      localStorage.setItem(STORAGE_KEY, JSON.stringify(updatedNotes));
      setError("Note deleted locally only");
    } finally {
      setIsLoading(false);
    }
  };

  const deleteAllNotes = async () => {
    setIsLoading(true);
    setError(null);

    try {
      await notesApi.deleteAllNotes();
      setNotes([]);
      localStorage.setItem(STORAGE_KEY, JSON.stringify([]));
    } catch (err) {
      console.warn(
        "Failed to delete all notes via API, using localStorage only",
        err
      );
      setNotes([]);
      localStorage.setItem(STORAGE_KEY, JSON.stringify([]));
      setError("Notes deleted locally only");
    } finally {
      setIsLoading(false);
    }
  };

  const refreshNotes = async () => {
    await loadNotes();
  };

  return {
    notes,
    isLoading,
    error,
    isBackendAvailable,
    addNote,
    deleteNote,
    deleteAllNotes,
    refreshNotes,
  };
};
