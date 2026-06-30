package model

import "time"

const (
	SlotStatusAvailable = 0
	SlotStatusBooked    = 1
	SlotStatusClosed    = 2
)

type MentorSlot struct {
	ID        int64     `json:"id"`
	MentorID  int64     `json:"mentor_id"`
	SlotDate  string    `json:"slot_date"` // YYYY-MM-DD
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	Status    int       `json:"status"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type SlotListQuery struct {
	MentorID int64
	FromDate string
	ToDate   string
}

type CreateBookingRequest struct {
	SlotID int64 `json:"slot_id" binding:"required"`
}

type CourseOrder struct {
	ID               string    `json:"id"`
	UserID           int64     `json:"user_id"`
	MentorID         int64     `json:"mentor_id"`
	SlotID           int64     `json:"slot_id"`
	Status           string    `json:"status"`
	AgoraChannelName string    `json:"agora_channel_name,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	Slot             *MentorSlot `json:"slot,omitempty"`
}

type BookingResponse struct {
	Order            CourseOrder `json:"order"`
	AvailableLessons int         `json:"available_lessons"`
}
