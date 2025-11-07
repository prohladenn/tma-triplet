# TMA Triplet

Telegram Mini App ecosystem: React frontend, Go API backend, and echo bot.

## 🚀 Quick Start

```bash
cd infra
cp .env.example .env
nano .env  # Add TELEGRAM_BOT_TOKEN
./deploy.sh deploy
```

**Access:**

- Frontend: http://localhost
- API: http://localhost/api/notes
- Bot: Polling mode (no URL needed)

## 📁 Structure

```
app/
├── backend/        # Go API (notes, auth)
├── bot/            # Telegram echo bot
├── docs/           # API documentation
└── frontend/       # React Mini App
infra/              # Docker deployment
```

## 🛠️ Services

**Frontend** (Port 80)

- React 18 + TypeScript + Telegram UI
- Notes CRUD, offline-first with localStorage

**Backend** (Port 3000)

- Go 1.25, in-memory storage
- Dual Telegram auth (standard + webapp/init)

**Bot** (Port 3001)

- Polling (default) OR Webhook mode
- Auto-configures based on WEBHOOK_URL
- No public URL needed in polling mode
- Telegram UI components
- Notes CRUD interface
- Offline-first with localStorage

### Backend (Port 3000)

- Go 1.25
- In-memory storage
- Dual Telegram auth (standard + webapp/init)
- CORS enabled

### Bot (Port 3001)

- Dual mode: Polling (default) OR Webhook
- Echoes messages
- Auto-configures based on WEBHOOK_URL presence
- Works without public URL (polling mode)
- Health monitoring

## 📖 Documentation

- [infra/README.md](infra/README.md) - Deployment guide
- [app/docs/API.md](app/docs/API.md) - API endpoints
- [app/bot/README.md](app/bot/README.md) - Bot modes

## 🔧 Configuration

**Required:**

```bash
TELEGRAM_BOT_TOKEN=your_token_here  # From @BotFather
```

**Optional:**

```bash
WEBHOOK_URL=https://your-domain.com  # Enables webhook mode
TELEGRAM_BOT_ID=your_bot_id          # For webapp/init auth
VITE_API_BASE_URL=/api               # Frontend API path
```

## 🌐 Development (HTTPS)

```bash
./deploy.sh deploy       # Start all services
ngrok http 80            # Expose via HTTPS
```

Use ngrok URL in @BotFather for Mini App configuration.

**Optional webhook mode:**

```bash
WEBHOOK_URL=https://abc.ngrok-free.app
docker compose restart bot
```

## 🔍 Commands

```bash
./deploy.sh deploy    # Build & start all
./deploy.sh status    # Check status
./deploy.sh logs      # View logs
./deploy.sh health    # Health checks
./deploy.sh cleanup   # Stop & clean
```

## 📦 Dependencies

**Backend:** gorilla/mux, rs/cors, telegram-mini-apps/init-data-golang  
**Bot:** go-telegram/bot v1.17.0  
**Frontend:** @telegram-apps/telegram-ui, @tma.js/sdk-react, react-router-dom

## 📝 License

MIT License - see [LICENSE](LICENSE)
