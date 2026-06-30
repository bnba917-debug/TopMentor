package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/topmentor/backend/internal/model"
)

var (
	ErrAlreadyVerifiedMentor = errors.New("already verified mentor")
	ErrApplicationPending    = errors.New("application pending review")
)

type mentorApplicationRow struct {
	ID              int64
	MentorID        int64
	IdCardURL       string
	StudentCardURL  string
	EnglishProofURL string
	RejectReason    sql.NullString
	ReviewedAt      sql.NullTime
	CreatedAt       time.Time
}

func (r *MentorRepository) FindByPhoneAny(ctx context.Context, phone string) (*model.Mentor, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, openid, real_name, school_name, major, gender,
		       COALESCE(english_score, ''), COALESCE(avatar_url, ''), COALESCE(bio, ''),
		       COALESCE(intro_video_url, ''), tags,
		       is_verified, balance, created_at
		FROM mentors WHERE phone = $1`, phone)

	m, err := scanMentorRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return m, err
}

func (r *MentorRepository) getLatestApplication(ctx context.Context, mentorID int64) (*mentorApplicationRow, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, mentor_id, id_card_url, student_card_url,
		       COALESCE(english_proof_url, ''), reject_reason, reviewed_at, created_at
		FROM mentor_applications
		WHERE mentor_id = $1
		ORDER BY created_at DESC
		LIMIT 1`, mentorID)

	var app mentorApplicationRow
	err := row.Scan(
		&app.ID, &app.MentorID, &app.IdCardURL, &app.StudentCardURL,
		&app.EnglishProofURL, &app.RejectReason, &app.ReviewedAt, &app.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *MentorRepository) GetApplyStatus(ctx context.Context, phone string) (*model.MentorApplyStatus, error) {
	mentor, err := r.FindByPhoneAny(ctx, phone)
	if err != nil {
		return nil, err
	}
	if mentor == nil {
		return &model.MentorApplyStatus{Status: "none"}, nil
	}

	if mentor.IsVerified == 1 {
		return &model.MentorApplyStatus{
			Status:   "approved",
			MentorID: mentor.ID,
		}, nil
	}

	app, err := r.getLatestApplication(ctx, mentor.ID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return &model.MentorApplyStatus{
			Status:   "none",
			MentorID: mentor.ID,
			Profile:  mentorToApplyDraft(mentor, "", "", ""),
		}, nil
	}

	if !app.ReviewedAt.Valid {
		return &model.MentorApplyStatus{
			Status:        "pending",
			MentorID:      mentor.ID,
			ApplicationID: app.ID,
			AppliedAt:     app.CreatedAt,
			Profile:       mentorToApplyDraft(mentor, app.IdCardURL, app.StudentCardURL, app.EnglishProofURL),
		}, nil
	}

	if app.RejectReason.Valid && app.RejectReason.String != "" {
		return &model.MentorApplyStatus{
			Status:        "rejected",
			MentorID:      mentor.ID,
			ApplicationID: app.ID,
			RejectReason:  app.RejectReason.String,
			AppliedAt:     app.CreatedAt,
			Profile:       mentorToApplyDraft(mentor, app.IdCardURL, app.StudentCardURL, app.EnglishProofURL),
		}, nil
	}

	return &model.MentorApplyStatus{
		Status:   "none",
		MentorID: mentor.ID,
		Profile:  mentorToApplyDraft(mentor, app.IdCardURL, app.StudentCardURL, app.EnglishProofURL),
	}, nil
}

func mentorToApplyDraft(m *model.Mentor, idCard, studentCard, englishProof string) *model.MentorApplyDraft {
	return &model.MentorApplyDraft{
		RealName:        m.RealName,
		SchoolName:      m.SchoolName,
		Major:           m.Major,
		Gender:          m.Gender,
		EnglishScore:    m.EnglishScore,
		Bio:             m.Bio,
		AvatarURL:       m.AvatarURL,
		IntroVideoURL:   m.IntroVideoURL,
		Tags:            m.Tags,
		IdCardURL:       idCard,
		StudentCardURL:  studentCard,
		EnglishProofURL: englishProof,
	}
}

func (r *MentorRepository) SubmitApplication(ctx context.Context, userID int64, phone string, req model.SubmitMentorApplyRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	gender := req.Gender
	if gender == "" {
		gender = "unknown"
	}
	tags := pq.StringArray(req.Tags)

	var mentorID int64
	var isVerified int
	err = tx.QueryRowContext(ctx, `SELECT id, is_verified FROM mentors WHERE phone = $1`, phone).Scan(&mentorID, &isVerified)
	if errors.Is(err, sql.ErrNoRows) {
		openID := fmt.Sprintf("apply_u_%d", userID)
		err = tx.QueryRowContext(ctx, `
			INSERT INTO mentors (openid, phone, real_name, school_name, major, gender, english_score,
			                     avatar_url, bio, intro_video_url, tags, is_verified)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 0)
			RETURNING id`,
			openID, phone, req.RealName, req.SchoolName, req.Major, gender, req.EnglishScore,
			req.AvatarURL, req.Bio, req.IntroVideoURL, tags,
		).Scan(&mentorID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		if isVerified == 1 {
			return ErrAlreadyVerifiedMentor
		}

		var pendingCount int
		if err := tx.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM mentor_applications
			WHERE mentor_id = $1 AND reviewed_at IS NULL`, mentorID).Scan(&pendingCount); err != nil {
			return err
		}
		if pendingCount > 0 {
			return ErrApplicationPending
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE mentors
			SET real_name = $2, school_name = $3, major = $4, gender = $5,
			    english_score = $6, avatar_url = $7, bio = $8, intro_video_url = $9,
			    tags = $10, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1`,
			mentorID, req.RealName, req.SchoolName, req.Major, gender, req.EnglishScore,
			req.AvatarURL, req.Bio, req.IntroVideoURL, tags,
		)
		if err != nil {
			return err
		}
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO mentor_applications (mentor_id, id_card_url, student_card_url, english_proof_url)
		VALUES ($1, $2, $3, NULLIF($4, ''))`,
		mentorID, req.IdCardURL, req.StudentCardURL, req.EnglishProofURL,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
