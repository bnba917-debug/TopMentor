package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/topmentor/backend/internal/model"
)

var (
	ErrReportExists      = errors.New("report already exists")
	ErrReportNotFound    = errors.New("report not found")
	ErrOrderNotCompleted = errors.New("order not completed")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Submit(
	ctx context.Context,
	mentorID int64,
	req model.SubmitReportRequest,
	earnAmount float64,
) (*model.GrowthReport, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	var orderMentorID int64
	var userID int64
	var status string
	err = tx.QueryRowContext(ctx, `
		SELECT mentor_id, user_id, status FROM course_orders WHERE id = $1 FOR UPDATE`,
		req.OrderID).Scan(&orderMentorID, &userID, &status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	if orderMentorID != mentorID {
		return nil, ErrOrderForbidden
	}
	if status != model.OrderStatusCompleted {
		return nil, ErrOrderNotCompleted
	}

	var exists int
	_ = tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM growth_reports WHERE order_id = $1`, req.OrderID).Scan(&exists)
	if exists > 0 {
		return nil, ErrReportExists
	}

	var report model.GrowthReport
	err = tx.QueryRowContext(ctx, `
		INSERT INTO growth_reports (order_id, mentor_id, user_id, speaking_score, confidence_score, vocabulary_score, comment)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, order_id, mentor_id, user_id, speaking_score, confidence_score, vocabulary_score, comment, created_at`,
		req.OrderID, mentorID, userID, req.SpeakingScore, req.ConfidenceScore, req.VocabularyScore, req.Comment,
	).Scan(
		&report.ID, &report.OrderID, &report.MentorID, &report.UserID,
		&report.SpeakingScore, &report.ConfidenceScore, &report.VocabularyScore, &report.Comment, &report.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	var balance float64
	err = tx.QueryRowContext(ctx,
		`SELECT balance FROM mentors WHERE id = $1 FOR UPDATE`, mentorID).Scan(&balance)
	if err != nil {
		return nil, err
	}

	newBalance := balance + earnAmount
	_, err = tx.ExecContext(ctx,
		`UPDATE mentors SET balance = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1`,
		mentorID, newBalance)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO wallet_transactions (mentor_id, order_id, amount, type, balance_after, remark)
		VALUES ($1, $2, $3, 'EARN', $4, $5)`,
		mentorID, req.OrderID, earnAmount, newBalance, fmt.Sprintf("课时结算 %s", req.OrderID))
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE course_orders SET feedback_submitted_at = $2, mentor_feedback = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`, req.OrderID, time.Now(), req.Comment)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ReportRepository) FindByOrderID(ctx context.Context, orderID string) (*model.GrowthReport, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT gr.id, gr.order_id, gr.mentor_id, gr.user_id,
		       gr.speaking_score, gr.confidence_score, gr.vocabulary_score, gr.comment, gr.created_at,
		       m.real_name, COALESCE(u.child_name, '')
		FROM growth_reports gr
		JOIN mentors m ON m.id = gr.mentor_id
		JOIN users u ON u.id = gr.user_id
		WHERE gr.order_id = $1`, orderID)

	var report model.GrowthReport
	err := row.Scan(
		&report.ID, &report.OrderID, &report.MentorID, &report.UserID,
		&report.SpeakingScore, &report.ConfidenceScore, &report.VocabularyScore, &report.Comment, &report.CreatedAt,
		&report.MentorName, &report.ChildName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrReportNotFound
	}
	return &report, err
}

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) GetBalance(ctx context.Context, mentorID int64) (float64, error) {
	var balance float64
	err := r.db.QueryRowContext(ctx,
		`SELECT balance FROM mentors WHERE id = $1`, mentorID).Scan(&balance)
	return balance, err
}

func (r *WalletRepository) ListTransactions(ctx context.Context, mentorID int64, limit int) ([]model.WalletTransaction, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, amount, type, balance_after, COALESCE(remark, ''), created_at
		FROM wallet_transactions WHERE mentor_id = $1
		ORDER BY created_at DESC LIMIT $2`, mentorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.WalletTransaction
	for rows.Next() {
		var t model.WalletTransaction
		if err := rows.Scan(&t.ID, &t.Amount, &t.Type, &t.BalanceAfter, &t.Remark, &t.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	if list == nil {
		list = []model.WalletTransaction{}
	}
	return list, rows.Err()
}

func (r *WalletRepository) Withdraw(ctx context.Context, mentorID int64, amount float64, mock bool) (float64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	var balance float64
	err = tx.QueryRowContext(ctx,
		`SELECT balance FROM mentors WHERE id = $1 FOR UPDATE`, mentorID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	if balance < amount {
		return 0, ErrInsufficientBalance
	}

	newBalance := balance - amount
	_, err = tx.ExecContext(ctx,
		`UPDATE mentors SET balance = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1`,
		mentorID, newBalance)
	if err != nil {
		return 0, err
	}

	remark := "提现至微信零钱"
	if mock {
		remark = "模拟提现"
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO wallet_transactions (mentor_id, amount, type, balance_after, remark)
		VALUES ($1, $2, 'WITHDRAW', $3, $4)`,
		mentorID, -amount, newBalance, remark)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return newBalance, nil
}

func (r *SlotRepository) UpsertSlots(ctx context.Context, mentorID int64, slots []model.SlotToggle) error {
	for _, s := range slots {
		status := model.SlotStatusClosed
		if s.Available {
			status = model.SlotStatusAvailable
		}
		endTime := addMinutesToTime(s.StartTime, 45)

		_, err := r.db.ExecContext(ctx, `
			INSERT INTO mentor_slots (mentor_id, slot_date, start_time, end_time, status)
			VALUES ($1, $2::date, $3::time, $4::time, $5)
			ON CONFLICT (mentor_id, slot_date, start_time)
			DO UPDATE SET status = EXCLUDED.status, end_time = EXCLUDED.end_time, updated_at = CURRENT_TIMESTAMP
			WHERE mentor_slots.status IN (0, 2)`,
			mentorID, s.SlotDate, s.StartTime, endTime, status)
		if err != nil {
			return err
		}
	}
	return nil
}

func addMinutesToTime(start string, mins int) string {
	t, err := time.Parse("15:04:05", start)
	if err != nil {
		t, _ = time.Parse("15:04", start)
	}
	return t.Add(time.Duration(mins) * time.Minute).Format("15:04:05")
}
