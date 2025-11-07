# TMA Triplet - Deployment

Docker deployment for Telegram Mini App with Notes API and Echo Bot.

## Quick Start

```bash
cd infra
cp .env.example .env
nano .env  # Add TELEGRAM_BOT_TOKEN
./deploy.sh deploy
```

**Access:** http://localhost

## Services

- **Frontend** (Port 80) - React + Telegram UI
- **Backend** (Port 3000) - Go API with Telegram auth
- **Bot** (Port 3001) - Echo bot (polling/webhook)

All services have health checks and auto-restart.

## Prerequisites

- Docker 20.10+, Docker Compose V2+
- Ports: 80, 3000, 3001

**Install:** [Docker Desktop](https://www.docker.com/products/docker-desktop) or `curl -fsSL https://get.docker.com | sh`

## Configuration

**Environment (.env):**

```bash
TELEGRAM_BOT_TOKEN=your_token       # Required
TELEGRAM_BOT_ID=your_bot_id         # Optional (webapp/init)
VITE_API_BASE_URL=/api              # API path
# WEBHOOK_URL=https://domain.com    # Optional (webhook mode)
```

## Commands

```bash
./deploy.sh deploy    # Build & start
./deploy.sh status    # Service status
./deploy.sh logs      # View logs
./deploy.sh health    # Health checks
./deploy.sh cleanup   # Stop & clean
```

**Direct Docker Compose:**

```bash
docker compose up -d --build
docker compose ps
docker compose logs -f [backend|frontend|bot]
docker compose down
```

## Access Points

- Frontend: http://localhost
- API: http://localhost/api/notes
- Backend health: http://localhost:3000/health
- Bot health: http://localhost:3001/health

## HTTPS Development (ngrok)

Telegram Mini Apps require HTTPS:

```bash
./deploy.sh deploy
ngrok http 80
```

Use ngrok HTTPS URL in @BotFather → Configure Mini App.

**Optional webhook mode:**

```bash
# Update .env:
WEBHOOK_URL=https://abc.ngrok-free.app
docker compose restart bot
```

## Bot Modes

**Polling (default):** No URL needed, works locally  
**Webhook:** Set WEBHOOK_URL, recommended for production

Check logs: `docker compose logs bot`

## Production Setup

```bash
# 1. Install Docker
curl -fsSL https://get.docker.com | sh

# 2. Clone & configure
git clone https://github.com/prohladenn/tma-triplet.git
cd tma-triplet/infra
cp .env.example .env
nano .env  # Set production values

# 3. Deploy
./deploy.sh deploy

# 4. Optional: SSL with certbot
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

## Troubleshooting

**Check health:**

```bash
./deploy.sh health
docker compose ps
```

**View logs:**

```bash
docker compose logs -f [service]
```

**Restart service:**

```bash
docker compose restart [backend|frontend|bot]
```

**Clean rebuild:**

```bash
./deploy.sh cleanup
./deploy.sh deploy
```

**Common issues:**

- Port conflicts: Change ports in docker-compose.yml
- Health checks fail: Check service logs
- Bot webhook errors: Use polling mode (remove WEBHOOK_URL)
