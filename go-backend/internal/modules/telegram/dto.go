package telegram

type VerifyCodeRequest struct {
	Code string `json:"code" binding:"required"`
}

type CheckStatusResponse struct {
	Linked           bool   `json:"linked"`
	TelegramUsername string `json:"telegram_username,omitempty"`
}

// TelegramUpdate - Telegram webhook update (simplified)
type TelegramUpdate struct {
	UpdateID int              `json:"update_id"`
	Message  *TelegramMessage `json:"message,omitempty"`
}

type TelegramMessage struct {
	MessageID int          `json:"message_id"`
	From      *TelegramUser `json:"from,omitempty"`
	Chat      TelegramChat `json:"chat"`
	Text      string       `json:"text"`
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type TelegramChat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}
