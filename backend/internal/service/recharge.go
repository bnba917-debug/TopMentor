package service

import (
	"context"
	"errors"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
	"github.com/topmentor/backend/pkg/payment"
)

type PackageStore interface {
	ListActive(ctx context.Context) ([]model.LessonPackage, error)
	FindByID(ctx context.Context, id int) (*model.LessonPackage, error)
}

type RechargeStore interface {
	CreatePending(ctx context.Context, userID int64, pkg model.LessonPackage, channel string) (*model.RechargeOrder, error)
	FindByID(ctx context.Context, orderID string) (*model.RechargeOrder, error)
	CompleteAndCreditLessons(ctx context.Context, orderID string, userID int64, lessonCount int, providerTxnID string) (int, error)
}

type UserLessonStore interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)
}

type RechargeService struct {
	packages PackageStore
	recharge RechargeStore
	users    UserLessonStore
	pay      *payment.Factory
}

func NewRechargeService(packages PackageStore, recharge RechargeStore, users UserLessonStore, pay *payment.Factory) *RechargeService {
	return &RechargeService{packages: packages, recharge: recharge, users: users, pay: pay}
}

func (s *RechargeService) ListPackages(ctx context.Context) ([]model.LessonPackage, error) {
	return s.packages.ListActive(ctx)
}

func (s *RechargeService) PaymentChannels() model.PaymentChannelsResponse {
	channels := s.pay.ConfiguredChannels()
	names := make([]string, len(channels))
	for i, ch := range channels {
		names[i] = string(ch)
	}
	mode := "live"
	if s.pay.IsMockMode() {
		mode = "mock"
	}
	return model.PaymentChannelsResponse{Mode: mode, Channels: names}
}

func (s *RechargeService) CreateRecharge(ctx context.Context, userID int64, req model.RechargeRequest, clientIP string) (*model.RechargeResponse, error) {
	ch, err := payment.ParseChannel(req.Channel)
	if err != nil {
		return nil, err
	}

	pkg, err := s.packages.FindByID(ctx, req.PackageID)
	if errors.Is(err, repository.ErrPackageNotFound) {
		return nil, repository.ErrPackageNotFound
	}
	if err != nil {
		return nil, err
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	order, err := s.recharge.CreatePending(ctx, userID, *pkg, string(ch))
	if err != nil {
		return nil, err
	}

	provider, err := s.pay.Provider(ch)
	if err != nil {
		return nil, err
	}

	prep, err := provider.Prepare(payment.PrepareRequest{
		OrderID:     order.ID,
		UserID:      userID,
		AmountCents: pkg.PriceCents,
		Description: pkg.Name,
		ClientIP:    clientIP,
		OpenID:      user.OpenID,
	})
	if err != nil {
		return nil, err
	}

	resp := &model.RechargeResponse{
		OrderID:     order.ID,
		Status:      order.Status,
		AmountCents: pkg.PriceCents,
		PackageName: pkg.Name,
		LessonCount: pkg.LessonCount,
		Channel:     string(ch),
		PayURL:      prep.PayURL,
		JsapiParams: prep.JsapiParams,
	}

	if prep.AutoComplete {
		balance, err := s.recharge.CompleteAndCreditLessons(ctx, order.ID, userID, pkg.LessonCount, prep.ProviderTxnID)
		if err != nil {
			return nil, err
		}
		resp.Status = "PAID"
		resp.MockPaid = true
		resp.AvailableLessons = balance
	}

	return resp, nil
}

func (s *RechargeService) GetLessonBalance(ctx context.Context, userID int64) (*model.LessonBalanceResponse, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &model.LessonBalanceResponse{
		AvailableLessons: user.AvailableLessons,
		LockedLessons:    user.LockedLessons,
	}, nil
}

func (s *RechargeService) GetRechargeOrder(ctx context.Context, userID int64, orderID string) (*model.RechargeOrder, error) {
	order, err := s.recharge.FindByID(ctx, orderID)
	if errors.Is(err, repository.ErrRechargeNotFound) {
		return nil, repository.ErrRechargeNotFound
	}
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, repository.ErrRechargeNotFound
	}
	return order, nil
}
