package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
	"github.com/topmentor/backend/pkg/payment"
)

type mockPackageStore struct {
	pkgs []model.LessonPackage
}

func (m *mockPackageStore) ListActive(_ context.Context) ([]model.LessonPackage, error) {
	return m.pkgs, nil
}

func (m *mockPackageStore) FindByID(_ context.Context, id int) (*model.LessonPackage, error) {
	for _, p := range m.pkgs {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, repository.ErrPackageNotFound
}

type mockRechargeStore struct {
	orders map[string]*model.RechargeOrder
}

func (m *mockRechargeStore) CreatePending(_ context.Context, userID int64, pkg model.LessonPackage, channel string) (*model.RechargeOrder, error) {
	if m.orders == nil {
		m.orders = map[string]*model.RechargeOrder{}
	}
	o := &model.RechargeOrder{
		ID:             "R_test_001",
		UserID:         userID,
		PackageID:      pkg.ID,
		AmountCents:    pkg.PriceCents,
		PaymentChannel: channel,
		Status:         "PENDING",
	}
	m.orders[o.ID] = o
	return o, nil
}

func (m *mockRechargeStore) FindByID(_ context.Context, orderID string) (*model.RechargeOrder, error) {
	o, ok := m.orders[orderID]
	if !ok {
		return nil, repository.ErrRechargeNotFound
	}
	return o, nil
}

func (m *mockRechargeStore) CompleteAndCreditLessons(_ context.Context, orderID string, _ int64, lessonCount int, txnID string) (int, error) {
	o, ok := m.orders[orderID]
	if !ok {
		return 0, repository.ErrRechargeNotFound
	}
	if o.Status == "PAID" {
		return 0, repository.ErrRechargeAlreadyPaid
	}
	o.Status = "PAID"
	o.WxTransactionID = txnID
	return lessonCount, nil
}

type mockUserLessonStore struct {
	user *model.User
}

func (m *mockUserLessonStore) FindByID(_ context.Context, id int64) (*model.User, error) {
	if m.user != nil && m.user.ID == id {
		return m.user, nil
	}
	return nil, repository.ErrUserNotFound
}

func TestRechargeService_MockRecharge(t *testing.T) {
	svc := NewRechargeService(
		&mockPackageStore{pkgs: []model.LessonPackage{
			{ID: 1, Name: "体验课 1 节", LessonCount: 1, PriceCents: 9900, IsActive: true},
		}},
		&mockRechargeStore{},
		&mockUserLessonStore{user: &model.User{ID: 1, OpenID: "sms_138", AvailableLessons: 0}},
		payment.NewFactory(payment.MerchantConfig{Mode: "mock"}),
	)

	resp, err := svc.CreateRecharge(context.Background(), 1, model.RechargeRequest{
		PackageID: 1,
		Channel:   "mock",
	}, "127.0.0.1")
	require.NoError(t, err)
	assert.Equal(t, "PAID", resp.Status)
	assert.True(t, resp.MockPaid)
	assert.Equal(t, 1, resp.AvailableLessons)
}

func TestRechargeService_ListPackages(t *testing.T) {
	svc := NewRechargeService(
		&mockPackageStore{pkgs: []model.LessonPackage{{ID: 1, Name: "10节包"}}},
		&mockRechargeStore{},
		&mockUserLessonStore{},
		payment.NewFactory(payment.MerchantConfig{Mode: "mock"}),
	)

	list, err := svc.ListPackages(context.Background())
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestRechargeService_PaymentChannels_MockMode(t *testing.T) {
	svc := NewRechargeService(nil, nil, nil, payment.NewFactory(payment.MerchantConfig{Mode: "mock"}))
	ch := svc.PaymentChannels()
	assert.Equal(t, "mock", ch.Mode)
	assert.Contains(t, ch.Channels, "mock")
}
