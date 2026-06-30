package model

import "time"

type User struct {
	ID               int64     `json:"id"`
	OpenID           string    `json:"-"`
	Phone            string    `json:"phone"`
	ChildName        string    `json:"child_name"`
	ChildAge         int       `json:"child_age"`
	EnglishLevel     string    `json:"english_level"`
	AvailableLessons int       `json:"available_lessons"`
	LockedLessons    int       `json:"locked_lessons"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Mentor struct {
	ID            int64     `json:"id"`
	OpenID        string    `json:"-"`
	RealName      string    `json:"real_name"`
	SchoolName    string    `json:"school_name"`
	Major         string    `json:"major"`
	Gender        string    `json:"gender"`
	EnglishScore  string    `json:"english_score"`
	AvatarURL     string    `json:"avatar_url"`
	Bio           string    `json:"bio"`
	IntroVideoURL string    `json:"intro_video_url"`
	Tags          []string  `json:"tags"`
	IsVerified    int       `json:"-"`
	Balance       float64   `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
}

type WxLoginRequest struct {
	Code  string `json:"code" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

type LoginResponse struct {
	Token  string         `json:"token"`
	User   User           `json:"user"`
	Mentor *MentorProfile `json:"mentor,omitempty"`
}

type WxLoginResponse = LoginResponse

type SmsSendRequest struct {
	Phone string `json:"phone" binding:"required,len=11"`
}

type SmsSendResponse struct {
	ExpiresIn int    `json:"expires_in"`
	DebugCode string `json:"debug_code,omitempty"`
}

type SmsLoginRequest struct {
	Phone string `json:"phone" binding:"required,len=11"`
	Code  string `json:"code" binding:"required,len=6"`
}

type UpdateProfileRequest struct {
	ChildName    string `json:"child_name" binding:"required"`
	ChildAge     int    `json:"child_age" binding:"required,min=6,max=14"`
	EnglishLevel string `json:"english_level" binding:"required,oneof=beginner intermediate advanced"`
}

type MentorListQuery struct {
	School   string
	Gender   string
	Tag      string
	Page     int
	PageSize int
}

type MentorListResponse struct {
	List  []Mentor `json:"list"`
	Total int      `json:"total"`
}
