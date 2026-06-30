package model

import "time"

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type PendingMentorApplication struct {
	MentorID        int64     `json:"mentor_id"`
	ApplicationID   int64     `json:"application_id"`
	RealName        string    `json:"real_name"`
	SchoolName      string    `json:"school_name"`
	Major           string    `json:"major"`
	Gender          string    `json:"gender"`
	EnglishScore    string    `json:"english_score"`
	IntroVideoURL   string    `json:"intro_video_url"`
	IdCardURL       string    `json:"id_card_url"`
	StudentCardURL  string    `json:"student_card_url"`
	EnglishProofURL string    `json:"english_proof_url,omitempty"`
	AppliedAt       time.Time `json:"applied_at"`
}

type ReviewMentorRequest struct {
	Action       string `json:"action" binding:"required,oneof=approve reject"`
	RejectReason string `json:"reject_reason"`
}

type Courseware struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	CoverURL   string    `json:"cover_url,omitempty"`
	ContentURL string    `json:"content_url"`
	SortOrder  int       `json:"sort_order"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateCoursewareRequest struct {
	Title      string `json:"title" binding:"required"`
	CoverURL   string `json:"cover_url"`
	ContentURL string `json:"content_url" binding:"required"`
	SortOrder  int    `json:"sort_order"`
	IsActive   *bool  `json:"is_active"`
}

type UpdateCoursewareRequest struct {
	Title      *string `json:"title"`
	CoverURL   *string `json:"cover_url"`
	ContentURL *string `json:"content_url"`
	SortOrder  *int    `json:"sort_order"`
	IsActive   *bool   `json:"is_active"`
}

type FinanceSummary struct {
	TotalRechargeYuan   float64 `json:"total_recharge_yuan"`
	TotalWithdrawYuan   float64 `json:"total_withdraw_yuan"`
	TotalMentorEarnYuan float64 `json:"total_mentor_earn_yuan"`
	UnspentLessons      int     `json:"unspent_lessons"`
	CompletedOrders     int     `json:"completed_orders"`
	PendingMentors      int     `json:"pending_mentors"`
	ActiveCourseware    int     `json:"active_courseware"`
}
