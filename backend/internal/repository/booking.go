package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/topmentor/backend/internal/model"
)

var (
	ErrSlotNotFound      = errors.New("slot not found")
	ErrSlotUnavailable   = errors.New("slot unavailable")
	ErrInsufficientLessons = errors.New("insufficient lessons")
)

type SlotRepository struct {
	db *sql.DB
}

func NewSlotRepository(db *sql.DB) *SlotRepository {
	return &SlotRepository{db: db}
}

func (r *SlotRepository) ListByMentor(ctx context.Context, q model.SlotListQuery) ([]model.MentorSlot, error) {
	fromDate := q.FromDate
	toDate := q.ToDate
	if fromDate == "" {
		fromDate = time.Now().Format("2006-01-02")
	}
	if toDate == "" {
		toDate = time.Now().AddDate(0, 0, 14).Format("2006-01-02")
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, mentor_id, slot_date::text, start_time::text, end_time::text, status, updated_at
		FROM mentor_slots
		WHERE mentor_id = $1 AND slot_date >= $2::date AND slot_date <= $3::date
		ORDER BY slot_date, start_time`,
		q.MentorID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.MentorSlot
	for rows.Next() {
		s, err := scanSlot(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *s)
	}
	if list == nil {
		list = []model.MentorSlot{}
	}
	return list, rows.Err()
}

func (r *SlotRepository) FindByID(ctx context.Context, slotID int64) (*model.MentorSlot, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, mentor_id, slot_date::text, start_time::text, end_time::text, status, updated_at
		FROM mentor_slots WHERE id = $1`, slotID)

	s, err := scanSlotRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSlotNotFound
	}
	return s, err
}

type slotScanner interface {
	Scan(dest ...interface{}) error
}

func scanSlotRow(row *sql.Row) (*model.MentorSlot, error) {
	return scanSlot(row)
}

func scanSlot(s slotScanner) (*model.MentorSlot, error) {
	var slot model.MentorSlot
	var startTime, endTime string
	err := s.Scan(&slot.ID, &slot.MentorID, &slot.SlotDate, &startTime, &endTime, &slot.Status, &slot.UpdatedAt)
	if err != nil {
		return nil, err
	}
	slot.StartTime = trimTime(startTime)
	slot.EndTime = trimTime(endTime)
	return &slot, nil
}

func trimTime(t string) string {
	if len(t) >= 8 {
		return t[:8]
	}
	return t
}

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func NewCourseOrderID() (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "C" + hex.EncodeToString(b), nil
}

func (r *BookingRepository) Create(ctx context.Context, userID, slotID int64) (*model.CourseOrder, int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = tx.Rollback() }()

	var slot model.MentorSlot
	var startTime, endTime string
	err = tx.QueryRowContext(ctx, `
		SELECT id, mentor_id, slot_date::text, start_time::text, end_time::text, status, updated_at
		FROM mentor_slots WHERE id = $1 FOR UPDATE`, slotID).Scan(
		&slot.ID, &slot.MentorID, &slot.SlotDate, &startTime, &endTime, &slot.Status, &slot.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, ErrSlotNotFound
	}
	if err != nil {
		return nil, 0, err
	}
	slot.StartTime = trimTime(startTime)
	slot.EndTime = trimTime(endTime)

	if slot.Status != model.SlotStatusAvailable {
		return nil, 0, ErrSlotUnavailable
	}

	var available int
	err = tx.QueryRowContext(ctx,
		`SELECT available_lessons FROM users WHERE id = $1 FOR UPDATE`, userID).Scan(&available)
	if err != nil {
		return nil, 0, err
	}
	if available < 1 {
		return nil, 0, ErrInsufficientLessons
	}

	newBalance := available - 1
	_, err = tx.ExecContext(ctx,
		`UPDATE users SET available_lessons = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1`,
		userID, newBalance)
	if err != nil {
		return nil, 0, err
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE mentor_slots SET status = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1`,
		slotID, model.SlotStatusBooked)
	if err != nil {
		return nil, 0, err
	}

	orderID, err := NewCourseOrderID()
	if err != nil {
		return nil, 0, err
	}

	channelName := fmt.Sprintf("tm_%s", orderID)
	var order model.CourseOrder
	err = tx.QueryRowContext(ctx, `
		INSERT INTO course_orders (id, user_id, mentor_id, slot_id, status, agora_channel_name)
		VALUES ($1, $2, $3, $4, 'RESERVED', $5)
		RETURNING id, user_id, mentor_id, slot_id, status, COALESCE(agora_channel_name, ''), created_at`,
		orderID, userID, slot.MentorID, slotID, channelName).Scan(
		&order.ID, &order.UserID, &order.MentorID, &order.SlotID, &order.Status, &order.AgoraChannelName, &order.CreatedAt)
	if err != nil {
		return nil, 0, err
	}

	order.Slot = &slot

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return &order, newBalance, nil
}
