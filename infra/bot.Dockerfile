# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy bot source
COPY app/bot/go.mod app/bot/go.sum ./
RUN go mod download && go mod verify

COPY app/bot/*.go ./

# Build the bot
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o bot .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /app

COPY --from=builder /build/bot .

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT:-3001}/health || exit 1

EXPOSE 3001

CMD ["./bot"]
