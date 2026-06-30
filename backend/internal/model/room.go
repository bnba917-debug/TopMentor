package model

const (
	OrderStatusReserved  = "RESERVED"
	OrderStatusActive    = "ACTIVE"
	OrderStatusCompleted = "COMPLETED"
	OrderStatusCancelled = "CANCELLED"
)

type CourseOrderDetail struct {
	CourseOrder
	MentorName    string `json:"mentor_name,omitempty"`
	ChildName     string `json:"child_name,omitempty"`
	UserPhone     string `json:"user_phone,omitempty"`
	ActualMinutes int    `json:"actual_minutes"`
	HasReport     bool   `json:"has_report,omitempty"`
	SlotDate      string `json:"slot_date,omitempty"`
	StartTime     string `json:"start_time,omitempty"`
	EndTime       string `json:"end_time,omitempty"`
}

type JoinRoomRequest struct {
	Role string `json:"role" binding:"required,oneof=user mentor"`
}

type JoinRoomResponse struct {
	AppID           string `json:"app_id"`
	Channel         string `json:"channel"`
	Token           string `json:"token"`
	UID             uint32 `json:"uid"`
	Role            string `json:"role"`
	MockMode        bool   `json:"mock_mode"`
	OrderID         string `json:"order_id"`
	OrderStatus     string `json:"order_status"`
	EndAt           string `json:"end_at"`
	DurationMinutes int    `json:"duration_minutes"`
	LessonStarted   bool   `json:"lesson_started"`
	StartedAt       string `json:"started_at,omitempty"`
	ElapsedMinutes  int    `json:"elapsed_minutes"`
	MentorOnline    bool   `json:"mentor_online"`
	UserOnline      bool   `json:"user_online"`
	MentorName      string `json:"mentor_name,omitempty"`
}

type HeartbeatResponse struct {
	OK             bool   `json:"ok"`
	ActualMinutes  int    `json:"actual_minutes"`
	LessonStarted  bool   `json:"lesson_started"`
	EndAt          string `json:"end_at,omitempty"`
	StartedAt      string `json:"started_at,omitempty"`
	MentorOnline   bool   `json:"mentor_online"`
	UserOnline     bool   `json:"user_online"`
}

type CompleteRoomResponse struct {
	OrderID       string `json:"order_id"`
	Status        string `json:"status"`
	ActualMinutes int    `json:"actual_minutes"`
}
