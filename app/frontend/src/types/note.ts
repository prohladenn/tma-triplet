export interface Note {
  id: string;
  text: string;
  timestamp: number;
}

export interface CreateNoteDto {
  text: string;
}

export interface UpdateNoteDto {
  text: string;
}

export interface NotesResponse {
  notes: Note[];
}
