# Telegram Echo Bot

Telegram bot that echoes messages. Supports polling and webhook modes.

## Features

- **Dual Mode**: Polling (default) or Webhook
- Echoes back any text message
- Containerized with Docker
- Health monitoring at `/health`
- Auto-detects mode from WEBHOOK_URL

## Modes

### Polling (Default)

No public URL needed. Uses `getUpdates` API.

**Pros:** Works locally, no ngrok needed  
**Cons:** Higher latency, more resource usage

**Usage:**
```bash
export TELEGRAM_BOT_TOKEN="your_token"
go run main.go
```

### Webhook (Production)

Requires public HTTPS URL. Telegram pushes updates.

**Pros:** Lower latency, less resources  
**Cons:** Needs public URL

**Usage:**
```bash
export TELEGRAM_BOT_TOKEN="your_token"
export WEBHOOK_URL="https://your-domain.com"
go run main.go
```

## Configuration

**Required:**
```bash
TELEGRAM_BOT_TOKEN=your_token  # From @BotFather
```

**Optional:**
```bash
WEBHOOK_URL=https://domain.com  # Enables webhook mode
PORT=3001                       # Server port (default: 3001)
```

## Docker

**Polling mode:**
```bash
# .env: Only set TELEGRAM_BOT_TOKEN
./deploy.sh start
```

**Webhook mode:**
```bash
# .env: Set both TELEGRAM_BOT_TOKEN and WEBHOOK_URL
./deploy.sh start
```

## Endpoints

- `/health` - Health check (both modes)
- `/webhook` - Webhook endpoint (webhook mode only)

## Switching Modes

Update `WEBHOOK_URL` in `.env` and restart:
```bash
docker compose restart bot
```

**Logs show active mode:**
```
# Polling: "Starting in POLLING mode"
# Webhook: "Starting in WEBHOOK mode"
```

## Logs

```bash
docker compose logs bot
```
