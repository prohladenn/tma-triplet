# Telegram Mini App Notes Backend

A simple Go backend for the Telegram Mini App notes application with authentication using Telegram init data.

## 📖 API Documentation

**Complete API specification: [/app/docs/API.md](/app/docs/API.md)**

The API documentation includes:

- All endpoints with request/response formats
- Authentication details
- Error handling
- CORS configuration
- Example requests
- Implementation examples

## Features

## Features

- ✅ RESTful API for notes CRUD operations
- ✅ Telegram Mini App authentication using init-data-golang
  - Standard validation (bot token) for regular Mini Apps
  - Third-party validation (bot ID) for webapp/init (attachment menu/direct links)
- ✅ User-specific notes storage
- ✅ In-memory storage (easy to replace with database)
- ✅ CORS enabled
- ✅ Clean architecture
- ✅ Comprehensive logging for all endpoints

## Prerequisites

- Go 1.21 or higher
- Telegram Bot Token or Bot ID (get from [@BotFather](https://t.me/BotFather))

## Installation

1. Install dependencies:

```bash
cd app/backend
go mod download
```

2. Create `.env` file:

```bash
cp .env.example .env
```

3. Edit `.env` and configure authentication:

**For standard Mini Apps:**

```bash
PORT=3000
TELEGRAM_BOT_TOKEN=your_bot_token_here
```

**For apps via attachment menu or direct links (webapp/init):**

```bash
PORT=3000
TELEGRAM_BOT_ID=your_bot_id_here
```

**For both methods:**

```bash
PORT=3000
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_BOT_ID=your_bot_id_here
```

TELEGRAM_BOT_TOKEN=your_actual_bot_token

````

## Running

### Development

```bash
go run .
````

### Production Build

```bash
go build -o notes-server
./notes-server
```

## API Endpoints

All endpoints require `X-Init-Data` header with Telegram init data (except in development mode).

### GET /api/notes

Get all notes for the authenticated user.

**Response:**

```json
{
  "notes": [
    {
      "id": "1730822400000",
      "text": "My note",
      "timestamp": 1730822400000
    }
  ]
}
```

### POST /api/notes

Create a new note.

**Request:**

```json
{
  "text": "My new note"
}
```

**Response:**

```json
{
  "id": "1730822500000",
  "text": "My new note",
  "timestamp": 1730822500000
}
```

### DELETE /api/notes/:id

Delete a specific note.

**Response:** 204 No Content

### DELETE /api/notes

Delete all notes for the authenticated user.

**Response:** 204 No Content

## Authentication

The backend uses Telegram Mini App init data for authentication. The init data is sent in the `X-Init-Data` header and validated using the init-data-golang library.

In development mode (when `TELEGRAM_BOT_TOKEN` is not set), authentication is disabled and all requests use a default user (ID: 0).

## Project Structure

```
backend/
├── main.go         # Server setup and route handlers
├── middleware.go   # Authentication middleware
├── storage.go      # In-memory storage implementation
├── go.mod          # Go module definition
└── .env.example    # Environment variables template
```

## Environment Variables

- `PORT` - Server port (default: 3000)
- `TELEGRAM_BOT_TOKEN` - Your Telegram bot token for init data validation

## Security Notes

- Init data is validated with 24-hour expiration
- Each user can only access their own notes
- CORS is configured (update allowed origins in production)

## Future Enhancements

- [ ] Add PostgreSQL/MongoDB database
- [ ] Add note update endpoint
- [ ] Add pagination for notes list
- [ ] Add note search functionality
- [ ] Add rate limiting
- [ ] Add logging middleware
- [ ] Add unit tests
