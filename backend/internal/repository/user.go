package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/topmentor/backend/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByOpenID(ctx context.Context, openID string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, openid, phone, COALESCE(child_name, ''), child_age, english_level,
		       available_lessons, locked_lessons, created_at, updated_at
		FROM users WHERE openid = $1`, openID)

	return scanUser(row)
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, openid, phone, COALESCE(child_name, ''), child_age, english_level,
		       available_lessons, locked_lessons, created_at, updated_at
		FROM users WHERE id = $1`, id)

	user, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return user, err
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, openid, phone, COALESCE(child_name, ''), child_age, english_level,
		       available_lessons, locked_lessons, created_at, updated_at
		FROM users WHERE phone = $1`, phone)

	user, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return user, err
}

func (r *UserRepository) Create(ctx context.Context, openID, phone string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO users (openid, phone) VALUES ($1, $2)
		RETURNING id, openid, phone, COALESCE(child_name, ''), child_age, english_level,
		          available_lessons, locked_lessons, created_at, updated_at`,
		openID, phone)
	return scanUser(row)
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id int64, req model.UpdateProfileRequest) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE users
		SET child_name = $2, child_age = $3, english_level = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, openid, phone, COALESCE(child_name, ''), child_age, english_level,
		          available_lessons, locked_lessons, created_at, updated_at`,
		id, req.ChildName, req.ChildAge, req.EnglishLevel)

	user, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return user, err
}

func scanUser(row *sql.Row) (*model.User, error) {
	var u model.User
	err := row.Scan(
		&u.ID, &u.OpenID, &u.Phone, &u.ChildName, &u.ChildAge, &u.EnglishLevel,
		&u.AvailableLessons, &u.LockedLessons, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

type MentorRepository struct {
	db *sql.DB
}

func NewMentorRepository(db *sql.DB) *MentorRepository {
	return &MentorRepository{db: db}
}

func (r *MentorRepository) List(ctx context.Context, q model.MentorListQuery) ([]model.Mentor, int, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 20
	}
	if q.PageSize > 50 {
		q.PageSize = 50
	}
	offset := (q.Page - 1) * q.PageSize

	where := "WHERE is_verified = 1"
	args := []interface{}{}
	argN := 1

	if q.School != "" {
		where += fmt.Sprintf(" AND school_name ILIKE $%d", argN)
		args = append(args, "%"+q.School+"%")
		argN++
	}
	if q.Gender != "" {
		where += fmt.Sprintf(" AND gender = $%d", argN)
		args = append(args, q.Gender)
		argN++
	}
	if q.Tag != "" {
		where += fmt.Sprintf(" AND $%d = ANY(tags)", argN)
		args = append(args, q.Tag)
		argN++
	}

	countQuery := "SELECT COUNT(*) FROM mentors " + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery := fmt.Sprintf(`
		SELECT id, openid, real_name, school_name, major, gender,
		       COALESCE(english_score, ''), COALESCE(avatar_url, ''), COALESCE(bio, ''),
		       COALESCE(intro_video_url, ''), tags,
		       is_verified, balance, created_at
		FROM mentors %s ORDER BY id DESC LIMIT $%d OFFSET $%d`, where, argN, argN+1)
	args = append(args, q.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []model.Mentor
	for rows.Next() {
		m, err := scanMentor(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, *m)
	}
	if list == nil {
		list = []model.Mentor{}
	}
	return list, total, rows.Err()
}

func (r *MentorRepository) FindByID(ctx context.Context, id int64) (*model.Mentor, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, openid, real_name, school_name, major, gender,
		       COALESCE(english_score, ''), COALESCE(avatar_url, ''), COALESCE(bio, ''),
		       COALESCE(intro_video_url, ''), tags,
		       is_verified, balance, created_at
		FROM mentors WHERE id = $1 AND is_verified = 1`, id)

	m, err := scanMentorRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return m, err
}

func (r *MentorRepository) FindByPhone(ctx context.Context, phone string) (*model.MentorProfile, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, real_name, school_name, is_verified, COALESCE(avatar_url, '')
		FROM mentors WHERE phone = $1 AND is_verified = 1`, phone)

	var p model.MentorProfile
	err := row.Scan(&p.ID, &p.RealName, &p.SchoolName, &p.IsVerified, &p.AvatarURL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *MentorRepository) FindOwnedByID(ctx context.Context, id int64) (*model.MentorPortalProfile, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, COALESCE(phone, ''), real_name, school_name, major, gender,
		       COALESCE(english_score, ''), COALESCE(avatar_url, ''), COALESCE(bio, ''),
		       COALESCE(intro_video_url, ''), tags, is_verified
		FROM mentors WHERE id = $1`, id)

	var p model.MentorPortalProfile
	var tags pq.StringArray
	err := row.Scan(
		&p.ID, &p.Phone, &p.RealName, &p.SchoolName, &p.Major, &p.Gender,
		&p.EnglishScore, &p.AvatarURL, &p.Bio, &p.IntroVideoURL, &tags, &p.IsVerified,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrMentorNotFound
	}
	if err != nil {
		return nil, err
	}
	p.Tags = []string(tags)
	return &p, nil
}

func (r *MentorRepository) UpdateProfile(ctx context.Context, id int64, req model.UpdateMentorProfileRequest) (*model.MentorPortalProfile, error) {
	gender := req.Gender
	if gender == "" {
		gender = "unknown"
	}
	tags := pq.StringArray(req.Tags)
	if tags == nil {
		tags = pq.StringArray{}
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE mentors
		SET real_name = $2, school_name = $3, major = $4, gender = $5,
		    english_score = $6, bio = $7, tags = $8,
		    avatar_url = CASE WHEN $9 <> '' THEN $9 ELSE avatar_url END,
		    intro_video_url = CASE WHEN $10 <> '' THEN $10 ELSE intro_video_url END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`,
		id, req.RealName, req.SchoolName, req.Major, gender, req.EnglishScore,
		req.Bio, tags, req.AvatarURL, req.IntroVideoURL,
	)
	if err != nil {
		return nil, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, ErrMentorNotFound
	}
	return r.FindOwnedByID(ctx, id)
}

func (r *MentorRepository) SetAvatarURL(ctx context.Context, id int64, url string) error {
	return r.setMediaURL(ctx, id, "avatar_url", url)
}

func (r *MentorRepository) SetIntroVideoURL(ctx context.Context, id int64, url string) error {
	return r.setMediaURL(ctx, id, "intro_video_url", url)
}

func (r *MentorRepository) setMediaURL(ctx context.Context, id int64, column, url string) error {
	query := fmt.Sprintf(`UPDATE mentors SET %s = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1`, column)
	res, err := r.db.ExecContext(ctx, query, id, url)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrMentorNotFound
	}
	return nil
}

type mentorScanner interface {
	Scan(dest ...interface{}) error
}

func scanMentorRow(row *sql.Row) (*model.Mentor, error) {
	return scanMentor(row)
}

func scanMentor(s mentorScanner) (*model.Mentor, error) {
	var m model.Mentor
	var tags pq.StringArray
	err := s.Scan(
		&m.ID, &m.OpenID, &m.RealName, &m.SchoolName, &m.Major, &m.Gender,
		&m.EnglishScore, &m.AvatarURL, &m.Bio, &m.IntroVideoURL, &tags,
		&m.IsVerified, &m.Balance, &m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	m.Tags = []string(tags)
	return &m, nil
}

var ErrMentorNotFound = errors.New("mentor not found")
