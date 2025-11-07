/**
 * Notes API Client
 *
 * This implementation follows the API specification defined in /app/docs/API.md
 * All endpoints, request/response formats, and error handling match the documentation.
 */

import { retrieveLaunchParams } from "@tma.js/sdk";

import type { CreateNoteDto, Note, NotesResponse } from "@/types/note";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:3000/api";

class NotesApi {
  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    // Get Telegram init data for authentication
    const launchParams = retrieveLaunchParams();
    const initDataRaw =
      typeof launchParams.initDataRaw === "string"
        ? launchParams.initDataRaw
        : "";

    console.log("[API]", options?.method || "GET", endpoint);
    console.log("[API] Init data length:", initDataRaw.length);
    console.log(
      "[API] Init data preview:",
      initDataRaw.substring(0, 100) + "..."
    );

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        "X-Init-Data": initDataRaw,
        ...options?.headers,
      },
    });

    console.log("[API] Response status:", response.status, response.statusText);

    if (!response.ok) {
      const errorText = await response.text();
      console.error("[API] Request failed:", response.status, errorText);
      throw new Error(`API request failed: ${response.statusText}`);
    }

    // Handle empty responses (204 No Content)
    if (
      response.status === 204 ||
      response.headers.get("content-length") === "0"
    ) {
      return undefined as T;
    }

    return response.json();
  }

  async getNotes(): Promise<Note[]> {
    const data = await this.request<NotesResponse>("/notes");
    return data.notes;
  }

  async createNote(noteData: CreateNoteDto): Promise<Note> {
    return this.request<Note>("/notes", {
      method: "POST",
      body: JSON.stringify(noteData),
    });
  }

  async deleteNote(id: string): Promise<void> {
    await this.request<void>(`/notes/${id}`, {
      method: "DELETE",
    });
  }

  async deleteAllNotes(): Promise<void> {
    await this.request<void>("/notes", {
      method: "DELETE",
    });
  }
}

export const notesApi = new NotesApi();
