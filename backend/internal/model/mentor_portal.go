package model

import "time"

type MentorProfile struct {
	ID         int64  `json:"id"`
	RealName   string `json:"real_name"`
	SchoolName string `json:"school_name"`
	IsVerified int    `json:"is_verified"`
	AvatarURL  string `json:"avatar_url,omitempty"`
}

type MentorPortalProfile struct {
	ID            int64    `json:"id"`
	Phone         string   `json:"phone,omitempty"`
	RealName      string   `json:"real_name"`
	SchoolName    string   `json:"school_name"`
	Major         string   `json:"major"`
	Gender        string   `json:"gender"`
	EnglishScore  string   `json:"english_score"`
	AvatarURL     string   `json:"avatar_url"`
	Bio           string   `json:"bio"`
	IntroVideoURL string   `json:"intro_video_url"`
	Tags          []string `json:"tags"`
	IsVerified    int      `json:"is_verified"`
}

type UpdateMentorProfileRequest struct {
	RealName      string   `json:"real_name" binding:"required,max=50"`
	SchoolName    string   `json:"school_name" binding:"required,max=100"`
	Major         string   `json:"major" binding:"required,max=100"`
	Gender        string   `json:"gender" binding:"omitempty,oneof=male female unknown"`
	EnglishScore  string   `json:"english_score" binding:"max=100"`
	Bio           string   `json:"bio" binding:"max=500"`
	Tags          []string `json:"tags" binding:"max=6,dive,max=20"`
	AvatarURL     string   `json:"avatar_url" binding:"omitempty,max=512"`
	IntroVideoURL string   `json:"intro_video_url" binding:"omitempty,max=512"`
}

type UploadResult struct {
	URL string `json:"url"`
}

type SubmitReportRequest struct {
	OrderID          string `json:"order_id" binding:"required"`
	SpeakingScore    int    `json:"speaking_score" binding:"required,min=1,max=5"`
	ConfidenceScore  int    `json:"confidence_score" binding:"required,min=1,max=5"`
	VocabularyScore  int    `json:"vocabulary_score" binding:"required,min=1,max=5"`
	Comment          string `json:"comment" binding:"required,min=10"`
}

type GrowthReport struct {
	ID               int64     `json:"id"`
	OrderID          string    `json:"order_id"`
	MentorID         int64     `json:"mentor_id"`
	UserID           int64     `json:"user_id"`
	SpeakingScore    int       `json:"speaking_score"`
	ConfidenceScore  int       `json:"confidence_score"`
	VocabularyScore  int       `json:"vocabulary_score"`
	Comment          string    `json:"comment"`
	MentorName       string    `json:"mentor_name,omitempty"`
	ChildName        string    `json:"child_name,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type WalletSummary struct {
	Balance      float64             `json:"balance"`
	Transactions []WalletTransaction `json:"transactions"`
}

type WalletTransaction struct {
	ID           int64     `json:"id"`
	Amount       float64   `json:"amount"`
	Type         string    `json:"type"`
	BalanceAfter float64   `json:"balance_after"`
	Remark       string    `json:"remark,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type WithdrawRequest struct {
	AmountCents int `json:"amount_cents" binding:"required,min=100"`
}

type WithdrawResponse struct {
	Balance float64 `json:"balance"`
	MockPaid bool   `json:"mock_paid,omitempty"`
}

type SetSlotsRequest struct {
	Slots []SlotToggle `json:"slots" binding:"required,min=1,max=50"`
}

type SlotToggle struct {
	SlotDate  string `json:"slot_date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	Available bool   `json:"available"`
}
