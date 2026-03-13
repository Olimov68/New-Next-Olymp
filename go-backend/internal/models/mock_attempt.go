package models

import "time"

// MockAttempt — foydalanuvchining test topshirish urinishi
type MockAttempt struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	MockTestID uint       `gorm:"not null;index" json:"mock_test_id"`
	StartedAt  time.Time  `gorm:"not null" json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	Score      float64    `gorm:"default:0" json:"score"`
	MaxScore   float64    `gorm:"default:0" json:"max_score"`
	Correct    int        `gorm:"default:0" json:"correct"`
	Wrong      int        `gorm:"default:0" json:"wrong"`
	Unanswered int        `gorm:"default:0" json:"unanswered"`
	Percentage float64    `gorm:"default:0" json:"percentage"`
	TimeTaken  int        `gorm:"default:0" json:"time_taken"` // sekundlarda
	Status     string     `gorm:"size:20;default:in_progress;not null" json:"status"` // in_progress | completed | timed_out | abandoned
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	MockTest *MockTest `gorm:"foreignKey:MockTestID" json:"mock_test,omitempty"`
	Answers  []MockAttemptAnswer `gorm:"foreignKey:AttemptID;constraint:OnDelete:CASCADE" json:"answers,omitempty"`
}

// MockAttemptAnswer — urinishdagi har bir javob
type MockAttemptAnswer struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	AttemptID        uint      `gorm:"not null;index" json:"attempt_id"`
	QuestionID       uint      `gorm:"not null;index" json:"question_id"`
	SelectedOptionID *uint     `json:"selected_option_id"`
	IsCorrect        bool      `gorm:"default:false" json:"is_correct"`
	AnsweredAt       time.Time `gorm:"autoCreateTime" json:"answered_at"`

	// Relations
	Question       *Question       `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
	SelectedOption *QuestionOption `gorm:"foreignKey:SelectedOptionID" json:"selected_option,omitempty"`
}
