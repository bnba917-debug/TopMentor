package model

import "time"

type SubmitMentorApplyRequest struct {
	RealName        string   `json:"real_name" binding:"required,max=50"`
	SchoolName      string   `json:"school_name" binding:"required,max=100"`
	Major           string   `json:"major" binding:"required,max=100"`
	Gender          string   `json:"gender" binding:"omitempty,oneof=male female unknown"`
	EnglishScore    string   `json:"english_score" binding:"max=100"`
	Bio             string   `json:"bio" binding:"max=500"`
	Tags            []string `json:"tags" binding:"max=6,dive,max=20"`
	AvatarURL       string   `json:"avatar_url" binding:"required,max=512"`
	IntroVideoURL   string   `json:"intro_video_url" binding:"required,max=512"`
	IdCardURL       string   `json:"id_card_url" binding:"required,max=512"`
	StudentCardURL  string   `json:"student_card_url" binding:"required,max=512"`
	EnglishProofURL string   `json:"english_proof_url" binding:"omitempty,max=512"`
}

type MentorApplyStatus struct {
	Status        string    `json:"status"`
	MentorID      int64     `json:"mentor_id,omitempty"`
	ApplicationID int64     `json:"application_id,omitempty"`
	RejectReason  string    `json:"reject_reason,omitempty"`
	AppliedAt     time.Time `json:"applied_at,omitempty"`
	Profile       *MentorApplyDraft `json:"profile,omitempty"`
}

type MentorApplyDraft struct {
	RealName        string   `json:"real_name"`
	SchoolName      string   `json:"school_name"`
	Major           string   `json:"major"`
	Gender          string   `json:"gender"`
	EnglishScore    string   `json:"english_score"`
	Bio             string   `json:"bio"`
	AvatarURL       string   `json:"avatar_url"`
	IntroVideoURL   string   `json:"intro_video_url"`
	Tags            []string `json:"tags"`
	IdCardURL       string   `json:"id_card_url"`
	StudentCardURL  string   `json:"student_card_url"`
	EnglishProofURL string   `json:"english_proof_url,omitempty"`
}
