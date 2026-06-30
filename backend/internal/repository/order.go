package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/topmentor/backend/internal/model"
)

var ErrOrderNotFound = errors.New("order not found")
var ErrOrderForbidden = errors.New("order access forbidden")
var ErrOrderNotJoinable = errors.New("order not joinable")

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID int64) ([]model.CourseOrderDetail, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT co.id, co.user_id, co.mentor_id, co.slot_id, co.status,
		       COALESCE(co.agora_channel_name, ''), co.actual_minutes, co.created_at,
		       m.real_name,
		       ms.slot_date::text, ms.start_time::text, ms.end_time::text
		FROM course_orders co
		JOIN mentors m ON m.id = co.mentor_id
		JOIN mentor_slots ms ON ms.id = co.slot_id
		WHERE co.user_id = $1
		ORDER BY co.created_at DESC
		LIMIT 50`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanOrderDetails(rows)
}

func (r *OrderRepository) ListByMentor(ctx context.Context, mentorID int64) ([]model.CourseOrderDetail, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT co.id, co.user_id, co.mentor_id, co.slot_id, co.status,
		       COALESCE(co.agora_channel_name, ''), co.actual_minutes, co.created_at,
		       m.real_name,
		       ms.slot_date::text, ms.start_time::text, ms.end_time::text,
		       COALESCE(u.child_name, ''), COALESCE(u.phone, ''),
		       (gr.id IS NOT NULL)
		FROM course_orders co
		JOIN mentors m ON m.id = co.mentor_id
		JOIN mentor_slots ms ON ms.id = co.slot_id
		JOIN users u ON u.id = co.user_id
		LEFT JOIN growth_reports gr ON gr.order_id = co.id
		WHERE co.mentor_id = $1
		ORDER BY co.created_at DESC
		LIMIT 50`, mentorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.CourseOrderDetail
	for rows.Next() {
		d, err := scanMentorOrderRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *d)
	}
	if list == nil {
		list = []model.CourseOrderDetail{}
	}
	return list, rows.Err()
}

func scanMentorOrderRow(rows *sql.Rows) (*model.CourseOrderDetail, error) {
	var d model.CourseOrderDetail
	var startTime, endTime string
	err := rows.Scan(
		&d.ID, &d.UserID, &d.MentorID, &d.SlotID, &d.Status,
		&d.AgoraChannelName, &d.ActualMinutes, &d.CreatedAt,
		&d.MentorName,
		&d.SlotDate, &startTime, &endTime,
		&d.ChildName, &d.UserPhone,
		&d.HasReport,
	)
	if err != nil {
		return nil, err
	}
	d.StartTime = trimTime(startTime)
	d.EndTime = trimTime(endTime)
	return &d, nil
}

func (r *OrderRepository) FindDetailByID(ctx context.Context, orderID string) (*model.CourseOrderDetail, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT co.id, co.user_id, co.mentor_id, co.slot_id, co.status,
		       COALESCE(co.agora_channel_name, ''), co.actual_minutes, co.created_at,
		       m.real_name,
		       ms.slot_date::text, ms.start_time::text, ms.end_time::text
		FROM course_orders co
		JOIN mentors m ON m.id = co.mentor_id
		JOIN mentor_slots ms ON ms.id = co.slot_id
		WHERE co.id = $1`, orderID)

	d, err := scanOrderDetailRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrOrderNotFound
	}
	return d, err
}

func (r *OrderRepository) Activate(ctx context.Context, orderID string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE course_orders
		SET status = 'ACTIVE', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status IN ('RESERVED', 'ACTIVE')`, orderID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrOrderNotJoinable
	}
	return nil
}

func (r *OrderRepository) UpdateActualMinutes(ctx context.Context, orderID string, minutes int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE course_orders SET actual_minutes = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`, orderID, minutes)
	return err
}

func (r *OrderRepository) Complete(ctx context.Context, orderID string, minutes int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE course_orders
		SET status = 'COMPLETED', actual_minutes = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status IN ('RESERVED', 'ACTIVE')`, orderID, minutes)
	return err
}

type orderRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
	Close() error
}

func scanOrderDetails(rows *sql.Rows) ([]model.CourseOrderDetail, error) {
	defer rows.Close()
	var list []model.CourseOrderDetail
	for rows.Next() {
		d, err := scanOrderFields(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *d)
	}
	if list == nil {
		list = []model.CourseOrderDetail{}
	}
	return list, rows.Err()
}

func scanOrderDetailRow(row *sql.Row) (*model.CourseOrderDetail, error) {
	var d model.CourseOrderDetail
	var startTime, endTime string
	err := row.Scan(
		&d.ID, &d.UserID, &d.MentorID, &d.SlotID, &d.Status,
		&d.AgoraChannelName, &d.ActualMinutes, &d.CreatedAt,
		&d.MentorName,
		&d.SlotDate, &startTime, &endTime,
	)
	if err != nil {
		return nil, err
	}
	d.StartTime = trimTime(startTime)
	d.EndTime = trimTime(endTime)
	return &d, nil
}

func scanOrderFields(s orderRows) (*model.CourseOrderDetail, error) {
	var d model.CourseOrderDetail
	var startTime, endTime string
	err := s.Scan(
		&d.ID, &d.UserID, &d.MentorID, &d.SlotID, &d.Status,
		&d.AgoraChannelName, &d.ActualMinutes, &d.CreatedAt,
		&d.MentorName,
		&d.SlotDate, &startTime, &endTime,
	)
	if err != nil {
		return nil, err
	}
	d.StartTime = trimTime(startTime)
	d.EndTime = trimTime(endTime)
	return &d, nil
}

func SlotEndTime(slotDate, endTime string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s", slotDate, endTime), time.Local)
}
