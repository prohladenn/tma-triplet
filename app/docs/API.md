# Notes API Documentation

Complete API specification for the Telegram Mini App Notes application.

**Version:** 1.0.0  
**Last Updated:** 2025-11-05

---

## Table of Contents

- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Authentication](#authentication)
- [Data Types](#data-types)
- [API Endpoints](#api-endpoints)
- [Error Handling](#error-handling)
- [Implementation Examples](#implementation-examples)

---

## Quick Start

### Backend (Go)

- Location: `/app/backend/`
- Base URL: `http://localhost:3000/api`
- Start: `cd app/backend && go run .`

### Frontend (TypeScript)

- Location: `/app/frontend/`
- API Client: `src/services/api/notes.api.ts`
- Hook: `src/hooks/useNotes.ts`
- Configure: Set `VITE_API_BASE_URL` in `.env`

---

## Architecture

### Frontend Components

```
src/
├── services/api/
│   └── notes.api.ts          # HTTP client for backend API
├── hooks/
│   └── useNotes.ts           # State management with offline fallback
├── types/
│   └── note.ts               # TypeScript interfaces
└── pages/IndexPage/
    └── IndexPage.tsx         # UI with Telegram UI components
```

**Features:**

- ✅ Offline-first approach with localStorage fallback
- ✅ Backend availability detection
- ✅ Modal notifications for connection status
- ✅ Loading and error states

### Backend Components

```
app/backend/
├── main.go                   # HTTP server and route handlers
├── middleware.go             # Telegram authentication
└── storage.go                # In-memory note storage
```

**Features:**

- ✅ RESTful API endpoints
- ✅ Telegram Mini App authentication
- ✅ User-specific note isolation
- ✅ Thread-safe in-memory storage
- ✅ CORS enabled

---

## Authentication

All API endpoints use Telegram Mini App authentication via init data validation.

### Headers

**Request:**

```
X-Init-Data: <telegram_init_data_string>
```

The init data string contains:

- User information (ID, username, first/last name)
- Query parameters from Telegram
- Authentication hash or signature

### Validation Methods

The backend supports two validation methods:

#### 1. Standard Validation (Bot Token)

For regular Mini Apps launched directly by the bot.

**Configuration:**

```bash
TELEGRAM_BOT_TOKEN=your_bot_token_here
```

**Process:**

1. Backend receives `X-Init-Data` header
2. Validates using bot token via `initdata.Validate()`
3. Verifies HMAC signature
4. Extracts user ID from validated data
5. Adds user ID to request context

#### 2. Third-Party Validation (Bot ID) - webapp/init

For Mini Apps launched via attachment menu or direct links.

**Configuration:**

```bash
TELEGRAM_BOT_ID=your_bot_id_here
```

**Process:**

1. Backend receives `X-Init-Data` header with `signature` parameter
2. Validates using bot ID via `initdata.ValidateThirdParty()`
3. Verifies Ed25519 signature with Telegram's public key
4. Extracts user ID from validated data
5. Adds user ID to request context

**Priority:** If both methods are configured, third-party validation is tried first, falling back to standard validation.

### Development Mode

When neither `TELEGRAM_BOT_TOKEN` nor `TELEGRAM_BOT_ID` is set:

- ⚠️ Authentication is disabled
- Uses default user ID: `12345`
- Warning logged on server start

**Production:** Always set at least one of `TELEGRAM_BOT_TOKEN` or `TELEGRAM_BOT_ID` to enable authentication.

---

## Data Types

### Note

```typescript
{
  id: string; // Unique identifier (UUID format)
  text: string; // Note content
  timestamp: number; // Unix timestamp in milliseconds
}
```

**Go Definition:**

```go
type Note struct {
    ID        string `json:"id"`
    Text      string `json:"text"`
    Timestamp int64  `json:"timestamp"`
    UserID    int64  `json:"user_id,omitempty"` // Internal only
}
```

### CreateNoteDto

```typescript
{
  text: string; // Note content (required, non-empty)
}
```

**Go Definition:**

```go
type CreateNoteRequest struct {
    Text string `json:"text"`
}
```

### NotesResponse

```typescript
{
  notes: Note[];        // Array of notes
}
```

**Go Definition:**

```go
type NotesResponse struct {
    Notes []Note `json:"notes"`
}
```

---

## API Endpoints

### Base URL

`/api`

All endpoints are prefixed with `/api`:

- Development: `http://localhost:3000/api`
- Production: Configure via environment variable

---

### Health Check

Check the health status of the backend service.

**Endpoint:** `GET /health`

**Authentication:** None required

**Request:**

```http
GET /health HTTP/1.1
```

**Response:** `200 OK`

```json
{
  "status": "healthy",
  "service": "tma-notes-api",
  "version": "1.0.0"
}
```

**Notes:**

- No authentication required
- Used by Docker health checks, load balancers, and monitoring systems
- Returns 200 OK when service is operational

---

### 1. Get All Notes

Retrieve all notes for the authenticated user.

**Endpoint:** `GET /api/notes`

**Request:**

```http
GET /api/notes HTTP/1.1
X-Init-Data: <init_data_string>
```

**Response:** `200 OK`

```json
{
  "notes": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "text": "My first note",
      "timestamp": 1730841600000
    },
    {
      "id": "987e6543-e21b-12d3-a456-426614174001",
      "text": "Another note",
      "timestamp": 1730841700000
    }
  ]
}
```

**Notes:**

- Returns empty array `[]` if user has no notes
- Notes are returned in creation order

---

### 2. Create Note

Create a new note for the authenticated user.

**Endpoint:** `POST /api/notes`

**Request:**

```http
POST /api/notes HTTP/1.1
Content-Type: application/json
X-Init-Data: <init_data_string>

{
  "text": "My new note"
}
```

**Response:** `201 Created`

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "text": "My new note",
  "timestamp": 1730841600000
}
```

**Validation:**

- `text` field is required
- `text` must not be empty string
- Whitespace-only text is considered empty

**Error Response:** `400 Bad Request`

```
Text is required
```

or

```
Invalid request body
```

---

### 3. Delete Note

Delete a specific note by ID.

**Endpoint:** `DELETE /api/notes/{id}`

**Request:**

```http
DELETE /api/notes/123e4567-e89b-12d3-a456-426614174000 HTTP/1.1
X-Init-Data: <init_data_string>
```

**Response:** `204 No Content`

```
(empty body)
```

**Error Response:** `404 Not Found`

```
Note not found
```

**Notes:**

- Users can only delete their own notes
- Returns 404 if note doesn't exist OR belongs to another user
- No response body on success

---

### 4. Delete All Notes

Delete all notes for the authenticated user.

**Endpoint:** `DELETE /api/notes`

**Request:**

```http
DELETE /api/notes HTTP/1.1
X-Init-Data: <init_data_string>
```

**Response:** `204 No Content`

```
(empty body)
```

**Notes:**

- Deletes all notes belonging to the authenticated user
- Safe operation - only affects current user's notes
- Returns 204 even if user has no notes
- No response body on success

---

## Error Handling

### HTTP Status Codes

| Code  | Description           | When Used                                      |
| ----- | --------------------- | ---------------------------------------------- |
| `200` | OK                    | Successful GET request                         |
| `201` | Created               | Successful POST request (note created)         |
| `204` | No Content            | Successful DELETE request (no response body)   |
| `400` | Bad Request           | Invalid request data or validation error       |
| `401` | Unauthorized          | Invalid or missing authentication (production) |
| `404` | Not Found             | Resource not found                             |
| `500` | Internal Server Error | Server error                                   |

### Error Response Format

Errors return plain text error messages with appropriate HTTP status codes.

**Examples:**

```
400 Bad Request
Text is required

404 Not Found
Note not found

400 Bad Request
Invalid request body
```

### Frontend Error Handling

The frontend API client handles errors and provides fallback to localStorage:

```typescript
try {
  const notes = await notesApi.getNotes();
  // Use backend data
} catch (error) {
  // Fall back to localStorage
  const localNotes = JSON.parse(localStorage.getItem("notes") || "[]");
}
```

---

## Implementation Examples

### Frontend (TypeScript)

#### API Client Usage

```typescript
import { notesApi } from "@/services/api/notes.api";

// Get all notes
const notes = await notesApi.getNotes();

// Create note
const newNote = await notesApi.createNote({ text: "New note" });

// Delete note
await notesApi.deleteNote(noteId);

// Delete all notes
await notesApi.deleteAllNotes();
```

#### Using the Hook

```typescript
import { useNotes } from "@/hooks/useNotes";

function MyComponent() {
  const {
    notes,
    isLoading,
    error,
    isBackendAvailable,
    addNote,
    deleteNote,
    deleteAllNotes,
  } = useNotes();

  const handleAdd = async () => {
    await addNote("My note text");
  };

  return (
    <div>
      {isBackendAvailable ? "✅ Online" : "⚠️ Offline"}
      {notes.map((note) => (
        <div key={note.id}>{note.text}</div>
      ))}
    </div>
  );
}
```

### Backend (Go)

#### Handler Implementation

```go
// Get notes for user
func GetNotesHandler(storage *MemoryStorage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := GetUserIDFromContext(r)
        notes := storage.GetNotes(userID)

        response := NotesResponse{Notes: notes}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}

// Create note
func CreateNoteHandler(storage *MemoryStorage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req CreateNoteRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        if req.Text == "" {
            http.Error(w, "Text is required", http.StatusBadRequest)
            return
        }

        userID := GetUserIDFromContext(r)
        note := storage.CreateNote(userID, req.Text)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(note)
    }
}
```

---

## CORS Configuration

The API supports Cross-Origin Resource Sharing (CORS):

```go
cors.New(cors.Options{
    AllowedOrigins:   []string{"*"}, // Restrict in production
    AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Init-Data"},
    AllowCredentials: true,
})
```

**Production Recommendation:** Replace `"*"` with specific frontend domain.

---

## Environment Variables

### Backend

| Variable             | Description                                     | Default | Required           |
| -------------------- | ----------------------------------------------- | ------- | ------------------ |
| `PORT`               | Server port                                     | `3000`  | No                 |
| `TELEGRAM_BOT_TOKEN` | Bot token for standard validation               | -       | Yes (production)\* |
| `TELEGRAM_BOT_ID`    | Bot ID for third-party validation (webapp/init) | -       | Yes (production)\* |

\* At least one authentication method (`TELEGRAM_BOT_TOKEN` or `TELEGRAM_BOT_ID`) should be configured in production.

**Authentication Methods:**

- **Standard validation:** Set `TELEGRAM_BOT_TOKEN` for regular Mini Apps
- **Third-party validation:** Set `TELEGRAM_BOT_ID` for apps via attachment menu or direct links
- **Both:** Set both to support all launch methods (third-party is tried first)
- **Dev mode:** Set neither to disable authentication (uses user ID 12345)

**File:** `/app/backend/.env`

### Frontend

| Variable            | Description     | Default                     | Required |
| ------------------- | --------------- | --------------------------- | -------- |
| `VITE_API_BASE_URL` | Backend API URL | `http://localhost:3000/api` | No       |

**File:** `/app/frontend/.env`

---

## Testing

### Manual Testing Checklist

- [ ] GET `/api/notes` returns empty array for new user
- [ ] POST `/api/notes` creates note and returns 201
- [ ] GET `/api/notes` returns created notes
- [ ] DELETE `/api/notes/{id}` removes specific note (204)
- [ ] DELETE `/api/notes` removes all notes (204)
- [ ] POST with empty text returns 400
- [ ] DELETE non-existent note returns 404
- [ ] Frontend shows "Backend Available" modal on connect
- [ ] Frontend falls back to localStorage when offline
- [ ] Frontend shows "Offline Mode" modal when unavailable

### Example cURL Commands

```bash
# Get notes
curl -X GET http://localhost:3000/api/notes \
  -H "X-Init-Data: user=%7B%22id%22%3A12345%7D"

# Create note
curl -X POST http://localhost:3000/api/notes \
  -H "Content-Type: application/json" \
  -H "X-Init-Data: user=%7B%22id%22%3A12345%7D" \
  -d '{"text":"Test note"}'

# Delete note
curl -X DELETE http://localhost:3000/api/notes/some-id \
  -H "X-Init-Data: user=%7B%22id%22%3A12345%7D"

# Delete all notes
curl -X DELETE http://localhost:3000/api/notes \
  -H "X-Init-Data: user=%7B%22id%22%3A12345%7D"
```

---

## Future Improvements

- [ ] Add API versioning (e.g., `/api/v1/notes`)
- [ ] Implement rate limiting
- [ ] Add request/response logging middleware
- [ ] Add PATCH endpoint for note updates
- [ ] Replace in-memory storage with database (PostgreSQL/MongoDB)
- [ ] Restrict CORS origins in production
- [ ] Add API response time metrics
- [ ] Implement pagination for large note lists
- [ ] Add note search/filtering
- [ ] Add note tags/categories

---

## Troubleshooting

### Backend Won't Start

**Error:** `cannot find main module`

- **Solution:** Run from `/app/backend` directory: `cd app/backend && go run .`

**Error:** `port already in use`

- **Solution:** Change port in `.env` or kill process: `lsof -ti:3000 | xargs kill`

### Frontend Can't Connect

**Error:** `Failed to fetch`

- **Check:** Backend is running on correct port
- **Check:** `VITE_API_BASE_URL` matches backend URL
- **Check:** CORS is enabled on backend

### Authentication Fails

**Error:** `401 Unauthorized`

- **Check:** `TELEGRAM_BOT_TOKEN` is set correctly
- **Check:** Init data is valid from Telegram
- **Dev Mode:** Remove `TELEGRAM_BOT_TOKEN` to disable auth

### Empty Response Error

**Error:** `Unexpected end of JSON input`

- **Cause:** Trying to parse 204 No Content response as JSON
- **Solution:** Already fixed - API client handles 204 correctly

---

## Support

For issues or questions:

1. Check this documentation
2. Review backend logs for errors
3. Check browser console for frontend errors
4. Verify environment variables are set correctly

---

**Documentation Path:** `/app/docs/API.md`  
**Project Repository:** [tma-triplet](https://github.com/prohladenn/tma-triplet)
