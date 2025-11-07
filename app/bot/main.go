package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Get configuration from environment
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	// Create bot with default handler
	opts := []bot.Option{
		bot.WithDefaultHandler(echoHandler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Decide between webhook and polling mode
	if webhookURL != "" {
		// Webhook mode
		log.Printf("Starting in WEBHOOK mode")
		log.Printf("Setting webhook to: %s/webhook", webhookURL)
		_, err = b.SetWebhook(ctx, &bot.SetWebhookParams{
			URL: webhookURL + "/webhook",
		})
		if err != nil {
			log.Fatalf("Failed to set webhook: %v", err)
		}
		log.Printf("Webhook registered successfully")

		// Start HTTP server for webhook
		addr := fmt.Sprintf(":%s", port)
		log.Printf("Starting webhook server on %s", addr)
		
		// Create a mux to handle both webhook and health check
		mux := http.NewServeMux()
		mux.Handle("/webhook", b.WebhookHandler())
		mux.HandleFunc("/health", healthHandler)
		
		go func() {
			if err := http.ListenAndServe(addr, mux); err != nil {
				log.Fatalf("HTTP server error: %v", err)
			}
		}()

		log.Println("Echo bot is running in webhook mode. Press Ctrl+C to stop.")
		b.StartWebhook(ctx)
	} else {
		// Polling mode
		log.Printf("Starting in POLLING mode (no WEBHOOK_URL provided)")
		
		// Delete any existing webhook
		ok, err := b.DeleteWebhook(ctx, &bot.DeleteWebhookParams{})
		if err != nil {
			log.Printf("Warning: Failed to delete existing webhook: %v", err)
		} else if ok {
			log.Printf("Deleted existing webhook")
		}

		// Start health check server
		addr := fmt.Sprintf(":%s", port)
		log.Printf("Starting health check server on %s", addr)
		
		mux := http.NewServeMux()
		mux.HandleFunc("/health", healthHandler)
		
		go func() {
			if err := http.ListenAndServe(addr, mux); err != nil {
				log.Fatalf("HTTP server error: %v", err)
			}
		}()

		log.Println("Echo bot is running in polling mode. Press Ctrl+C to stop.")
		b.Start(ctx)
	}
	
	log.Println("Bot stopped gracefully")
}

// healthHandler responds to health check requests
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"service":"tma-echo-bot","status":"healthy","version":"1.0.0"}`))
}

// echoHandler echoes back any text message received
func echoHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	log.Printf("Received message from user %d: %s", update.Message.From.ID, update.Message.Text)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
