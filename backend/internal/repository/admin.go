package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/topmentor/backend/internal/model"
)

var ErrMentorApplicationNotFound = errors.New("mentor application not found")
var ErrCoursewareNotFound = errors.New("courseware not found")

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) ListPendingMentors(ctx context.Context) ([]model.PendingMentorApplication, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT m.id, ma.id, m.real_name, m.school_name, m.major, m.gender,
		       COALESCE(m.english_score, ''), COALESCE(m.intro_video_url, ''),
		       ma.id_card_url, ma.student_card_url, COALESCE(ma.english_proof_url, ''), ma.created_at
		FROM mentors m
		JOIN mentor_applications ma ON ma.mentor_id = m.id
		WHERE m.is_verified = 0 AND ma.reviewed_at IS NULL
		ORDER BY ma.created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PendingMentorApplication
	for rows.Next() {
		var item model.PendingMentorApplication
		if err := rows.Scan(
			&item.MentorID, &item.ApplicationID, &item.RealName, &item.SchoolName, &item.Major, &item.Gender,
			&item.EnglishScore, &item.IntroVideoURL,
			&item.IdCardURL, &item.StudentCardURL, &item.EnglishProofURL, &item.AppliedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if list == nil {
		list = []model.PendingMentorApplication{}
	}
	return list, rows.Err()
}

func (r *AdminRepository) ReviewMentor(ctx context.Context, mentorID int64, req model.ReviewMentorRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var appID int64
	err = tx.QueryRowContext(ctx, `
		SELECT ma.id FROM mentor_applications ma
		JOIN mentors m ON m.id = ma.mentor_id
		WHERE m.id = $1 AND m.is_verified = 0 AND ma.reviewed_at IS NULL
		ORDER BY ma.created_at DESC LIMIT 1`, mentorID).Scan(&appID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrMentorApplicationNotFound
	}
	if err != nil {
		return err
	}

	now := time.Now()
	switch req.Action {
	case "approve":
		_, err = tx.ExecContext(ctx,
			`UPDATE mentors SET is_verified = 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`, mentorID)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			UPDATE mentor_applications SET reviewed_at = $2, reject_reason = NULL WHERE id = $1`, appID, now)
	case "reject":
		if strings.TrimSpace(req.RejectReason) == "" {
			return fmt.Errorf("reject reason required")
		}
		_, err = tx.ExecContext(ctx, `
			UPDATE mentor_applications SET reviewed_at = $2, reject_reason = $3 WHERE id = $1`,
			appID, now, req.RejectReason)
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *AdminRepository) ListCourseware(ctx context.Context) ([]model.Courseware, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, COALESCE(cover_url, ''), content_url, sort_order, is_active, created_at, updated_at
		FROM courseware ORDER BY sort_order ASC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Courseware
	for rows.Next() {
		var c model.Courseware
		if err := rows.Scan(&c.ID, &c.Title, &c.CoverURL, &c.ContentURL, &c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	if list == nil {
		list = []model.Courseware{}
	}
	return list, rows.Err()
}

func (r *AdminRepository) CreateCourseware(ctx context.Context, req model.CreateCoursewareRequest) (*model.Courseware, error) {
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO courseware (title, cover_url, content_url, sort_order, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, COALESCE(cover_url, ''), content_url, sort_order, is_active, created_at, updated_at`,
		req.Title, req.CoverURL, req.ContentURL, req.SortOrder, active)

	var c model.Courseware
	err := row.Scan(&c.ID, &c.Title, &c.CoverURL, &c.ContentURL, &c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (r *AdminRepository) UpdateCourseware(ctx context.Context, id int64, req model.UpdateCoursewareRequest) (*model.Courseware, error) {
	current, err := r.FindCoursewareByID(ctx, id)
	if errors.Is(err, ErrCoursewareNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	title := current.Title
	cover := current.CoverURL
	content := current.ContentURL
	sortOrder := current.SortOrder
	active := current.IsActive

	if req.Title != nil {
		title = *req.Title
	}
	if req.CoverURL != nil {
		cover = *req.CoverURL
	}
	if req.ContentURL != nil {
		content = *req.ContentURL
	}
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		active = *req.IsActive
	}

	row := r.db.QueryRowContext(ctx, `
		UPDATE courseware
		SET title = $2, cover_url = $3, content_url = $4, sort_order = $5, is_active = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, title, COALESCE(cover_url, ''), content_url, sort_order, is_active, created_at, updated_at`,
		id, title, cover, content, sortOrder, active)

	var c model.Courseware
	err = row.Scan(&c.ID, &c.Title, &c.CoverURL, &c.ContentURL, &c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCoursewareNotFound
	}
	return &c, err
}

func (r *AdminRepository) FindCoursewareByID(ctx context.Context, id int64) (*model.Courseware, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, title, COALESCE(cover_url, ''), content_url, sort_order, is_active, created_at, updated_at
		FROM courseware WHERE id = $1`, id)

	var c model.Courseware
	err := row.Scan(&c.ID, &c.Title, &c.CoverURL, &c.ContentURL, &c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCoursewareNotFound
	}
	return &c, err
}

func (r *AdminRepository) DeleteCourseware(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM courseware WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrCoursewareNotFound
	}
	return nil
}

func (r *AdminRepository) FinanceSummary(ctx context.Context) (*model.FinanceSummary, error) {
	var summary model.FinanceSummary

	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(amount_cents), 0) FROM recharge_orders WHERE status = 'PAID'`).Scan(&summary.TotalRechargeYuan)
	if err != nil {
		return nil, err
	}
	summary.TotalRechargeYuan /= 100

	err = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(ABS(amount)), 0) FROM wallet_transactions WHERE type = 'WITHDRAW'`).Scan(&summary.TotalWithdrawYuan)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(amount), 0) FROM wallet_transactions WHERE type = 'EARN'`).Scan(&summary.TotalMentorEarnYuan)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(available_lessons), 0) FROM users`).Scan(&summary.UnspentLessons)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM course_orders WHERE status = 'COMPLETED'`).Scan(&summary.CompletedOrders)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM mentors m
		JOIN mentor_applications ma ON ma.mentor_id = m.id
		WHERE m.is_verified = 0 AND ma.reviewed_at IS NULL`).Scan(&summary.PendingMentors)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM courseware WHERE is_active = TRUE`).Scan(&summary.ActiveCourseware)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
