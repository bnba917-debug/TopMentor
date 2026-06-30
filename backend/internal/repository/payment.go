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

var ErrPackageNotFound = errors.New("package not found")
var ErrRechargeNotFound = errors.New("recharge order not found")
var ErrRechargeAlreadyPaid = errors.New("recharge order already paid")

type PackageRepository struct {
	db *sql.DB
}

func NewPackageRepository(db *sql.DB) *PackageRepository {
	return &PackageRepository{db: db}
}

func (r *PackageRepository) ListActive(ctx context.Context) ([]model.LessonPackage, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, lesson_count, price_cents, is_trial, is_active, created_at
		FROM lesson_packages WHERE is_active = TRUE ORDER BY price_cents ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.LessonPackage
	for rows.Next() {
		var p model.LessonPackage
		if err := rows.Scan(&p.ID, &p.Name, &p.LessonCount, &p.PriceCents, &p.IsTrial, &p.IsActive, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	if list == nil {
		list = []model.LessonPackage{}
	}
	return list, rows.Err()
}

func (r *PackageRepository) FindByID(ctx context.Context, id int) (*model.LessonPackage, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, lesson_count, price_cents, is_trial, is_active, created_at
		FROM lesson_packages WHERE id = $1 AND is_active = TRUE`, id)

	var p model.LessonPackage
	err := row.Scan(&p.ID, &p.Name, &p.LessonCount, &p.PriceCents, &p.IsTrial, &p.IsActive, &p.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPackageNotFound
	}
	return &p, err
}

type RechargeRepository struct {
	db *sql.DB
}

func NewRechargeRepository(db *sql.DB) *RechargeRepository {
	return &RechargeRepository{db: db}
}

func NewRechargeOrderID() (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "R" + hex.EncodeToString(b), nil
}

func (r *RechargeRepository) CreatePending(ctx context.Context, userID int64, pkg model.LessonPackage, channel string) (*model.RechargeOrder, error) {
	orderID, err := NewRechargeOrderID()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO recharge_orders (id, user_id, package_id, amount_cents, payment_channel, status)
		VALUES ($1, $2, $3, $4, $5, 'PENDING')
		RETURNING id, user_id, package_id, amount_cents, payment_channel,
		          COALESCE(wx_transaction_id, ''), status, paid_at, created_at`,
		orderID, userID, pkg.ID, pkg.PriceCents, channel)

	return scanRechargeOrder(row)
}

func (r *RechargeRepository) FindByID(ctx context.Context, orderID string) (*model.RechargeOrder, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, package_id, amount_cents, payment_channel,
		       COALESCE(wx_transaction_id, ''), status, paid_at, created_at
		FROM recharge_orders WHERE id = $1`, orderID)

	order, err := scanRechargeOrder(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRechargeNotFound
	}
	return order, err
}

func (r *RechargeRepository) CompleteAndCreditLessons(
	ctx context.Context,
	orderID string,
	userID int64,
	lessonCount int,
	providerTxnID string,
) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	var status string
	err = tx.QueryRowContext(ctx,
		`SELECT status FROM recharge_orders WHERE id = $1 FOR UPDATE`, orderID).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrRechargeNotFound
	}
	if err != nil {
		return 0, err
	}
	if status == "PAID" {
		return 0, ErrRechargeAlreadyPaid
	}

	var available int
	err = tx.QueryRowContext(ctx,
		`SELECT available_lessons FROM users WHERE id = $1 FOR UPDATE`, userID).Scan(&available)
	if err != nil {
		return 0, err
	}

	newBalance := available + lessonCount
	_, err = tx.ExecContext(ctx,
		`UPDATE users SET available_lessons = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1`,
		userID, newBalance)
	if err != nil {
		return 0, err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE recharge_orders
		SET status = 'PAID', paid_at = $2, wx_transaction_id = $3
		WHERE id = $1 AND status = 'PENDING'`,
		orderID, time.Now(), providerTxnID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return newBalance, nil
}

func scanRechargeOrder(row *sql.Row) (*model.RechargeOrder, error) {
	var o model.RechargeOrder
	var paidAt sql.NullTime
	err := row.Scan(
		&o.ID, &o.UserID, &o.PackageID, &o.AmountCents, &o.PaymentChannel,
		&o.WxTransactionID, &o.Status, &paidAt, &o.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if paidAt.Valid {
		t := paidAt.Time
		o.PaidAt = &t
	}
	return &o, nil
}

func FormatPriceYuan(cents int) string {
	return fmt.Sprintf("%.2f", float64(cents)/100)
}
