package model

import "time"

type LessonPackage struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	LessonCount int       `json:"lesson_count"`
	PriceCents  int       `json:"price_cents"`
	IsTrial     bool      `json:"is_trial"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type RechargeOrder struct {
	ID              string    `json:"id"`
	UserID          int64     `json:"user_id"`
	PackageID       int       `json:"package_id"`
	AmountCents     int       `json:"amount_cents"`
	PaymentChannel  string    `json:"payment_channel"`
	WxTransactionID string    `json:"wx_transaction_id,omitempty"`
	Status          string    `json:"status"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type RechargeRequest struct {
	PackageID int    `json:"package_id" binding:"required"`
	Channel   string `json:"channel" binding:"required"`
}

type RechargeResponse struct {
	OrderID          string            `json:"order_id"`
	Status           string            `json:"status"`
	AmountCents      int               `json:"amount_cents"`
	PackageName      string            `json:"package_name"`
	LessonCount      int               `json:"lesson_count"`
	Channel          string            `json:"channel"`
	PayURL           string            `json:"pay_url,omitempty"`
	JsapiParams      map[string]string `json:"jsapi_params,omitempty"`
	MockPaid         bool              `json:"mock_paid,omitempty"`
	AvailableLessons int               `json:"available_lessons,omitempty"`
}

type LessonBalanceResponse struct {
	AvailableLessons int `json:"available_lessons"`
	LockedLessons    int `json:"locked_lessons"`
}

type PaymentChannelsResponse struct {
	Mode     string   `json:"mode"`
	Channels []string `json:"channels"`
}
