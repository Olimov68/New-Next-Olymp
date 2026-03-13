package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const defaultTTLMinutes = 10

func main() {
	_ = godotenv.Load()

	botToken := strings.TrimSpace(os.Getenv("BOT_TOKEN"))
	if botToken == "" {
		log.Fatal("BOT_TOKEN is required")
	}

	backendURL := strings.TrimSpace(os.Getenv("BACKEND_WEBHOOK_URL"))
	if backendURL == "" {
		backendURL = "http://localhost:8080/api/v1/telegram/webhook"
	}

	ttlMinutes := defaultTTLMinutes
	if raw := strings.TrimSpace(os.Getenv("CODE_TTL_MINUTES")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 {
			ttlMinutes = v
		}
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Telegram bot init failed: %v", err)
	}

	log.Printf("✅ Bot started as @%s (forward to %s, TTL=%dm)", bot.Self.UserName, backendURL, ttlMinutes)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := bot.GetUpdatesChan(u)

	client := &http.Client{Timeout: 5 * time.Second}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := update.Message

		// /start komandasi uchun nomer so'rash
		if msg.IsCommand() && msg.Command() == "start" {
			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButtonContact("📱 Nomerni ulashing"),
				),
			)
			keyboard.OneTimeKeyboard = true
			keyboard.ResizeKeyboard = true

			reply := tgbotapi.NewMessage(msg.Chat.ID,
				"Assalomu alaykum! 👋\n\nNextOlymp'ga kirish uchun telefon raqamingizni ulashing.")
			reply.ReplyMarkup = keyboard
			bot.Send(reply)
			continue
		}

		// Har qanday xabar — backendga forward qil
		payload := buildPayload(update)
		if err := forwardToBackend(client, backendURL, payload); err != nil {
			log.Printf("⚠️  Backend forward error: %v", err)
		} else {
			log.Printf("✉️  Forwarded update %d from @%s to backend", update.UpdateID, func() string {
				if msg.From != nil {
					return msg.From.UserName
				}
				return "unknown"
			}())
		}
	}
}

// buildPayload — Telegram update ni backendga mos JSON formatga o'tkazadi
func buildPayload(update tgbotapi.Update) map[string]interface{} {
	msg := update.Message
	from := map[string]interface{}{
		"id":         msg.From.ID,
		"is_bot":     msg.From.IsBot,
		"first_name": msg.From.FirstName,
		"last_name":  msg.From.LastName,
		"username":   msg.From.UserName,
	}

	message := map[string]interface{}{
		"message_id": msg.MessageID,
		"from":       from,
		"chat": map[string]interface{}{
			"id":   msg.Chat.ID,
			"type": msg.Chat.Type,
		},
		"text": msg.Text,
	}

	// Kontakt (telefon nomer) bo'lsa ham forward qilish
	if msg.Contact != nil {
		message["contact"] = map[string]interface{}{
			"phone_number": msg.Contact.PhoneNumber,
			"first_name":   msg.Contact.FirstName,
			"last_name":    msg.Contact.LastName,
			"user_id":      msg.Contact.UserID,
		}
	}

	return map[string]interface{}{
		"update_id": update.UpdateID,
		"message":   message,
	}
}

// forwardToBackend — backend webhook ga HTTP POST qiladi
func forwardToBackend(client *http.Client, url string, payload map[string]interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("backend returned %d", resp.StatusCode)
	}
	return nil
}
