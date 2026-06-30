package service

import (
	"context"
	"errors"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
)

type MentorPortalStore interface {
	FindByPhone(ctx context.Context, phone string) (*model.MentorProfile, error)
	FindOwnedByID(ctx context.Context, id int64) (*model.MentorPortalProfile, error)
	UpdateProfile(ctx context.Context, id int64, req model.UpdateMentorProfileRequest) (*model.MentorPortalProfile, error)
	SetAvatarURL(ctx context.Context, id int64, url string) error
	SetIntroVideoURL(ctx context.Context, id int64, url string) error
}

type ReportStore interface {
	Submit(ctx context.Context, mentorID int64, req model.SubmitReportRequest, earnAmount float64) (*model.GrowthReport, error)
	FindByOrderID(ctx context.Context, orderID string) (*model.GrowthReport, error)
}

type WalletStore interface {
	GetBalance(ctx context.Context, mentorID int64) (float64, error)
	ListTransactions(ctx context.Context, mentorID int64, limit int) ([]model.WalletTransaction, error)
	Withdraw(ctx context.Context, mentorID int64, amount float64, mock bool) (float64, error)
}

type MentorOrderStore interface {
	ListByMentor(ctx context.Context, mentorID int64) ([]model.CourseOrderDetail, error)
}

type MentorSlotWriter interface {
	ListByMentor(ctx context.Context, q model.SlotListQuery) ([]model.MentorSlot, error)
	UpsertSlots(ctx context.Context, mentorID int64, slots []model.SlotToggle) error
}

type MentorPortalService struct {
	mentors   MentorPortalStore
	orders    MentorOrderStore
	slots     MentorSlotWriter
	reports   ReportStore
	wallet    WalletStore
	earnYuan  float64
	withdrawMock bool
}

func NewMentorPortalService(
	mentors MentorPortalStore,
	orders MentorOrderStore,
	slots MentorSlotWriter,
	reports ReportStore,
	wallet WalletStore,
	earnYuan float64,
	withdrawMock bool,
) *MentorPortalService {
	if earnYuan <= 0 {
		earnYuan = 80
	}
	return &MentorPortalService{
		mentors: mentors, orders: orders, slots: slots,
		reports: reports, wallet: wallet,
		earnYuan: earnYuan, withdrawMock: withdrawMock,
	}
}

func (s *MentorPortalService) ResolveMentor(ctx context.Context, phone string) (*model.MentorProfile, error) {
	return s.mentors.FindByPhone(ctx, phone)
}

func (s *MentorPortalService) ListOrders(ctx context.Context, mentorID int64) ([]model.CourseOrderDetail, error) {
	return s.orders.ListByMentor(ctx, mentorID)
}

func (s *MentorPortalService) ListSlots(ctx context.Context, mentorID int64, from, to string) ([]model.MentorSlot, error) {
	return s.slots.ListByMentor(ctx, model.SlotListQuery{MentorID: mentorID, FromDate: from, ToDate: to})
}

func (s *MentorPortalService) SetSlots(ctx context.Context, mentorID int64, req model.SetSlotsRequest) error {
	return s.slots.UpsertSlots(ctx, mentorID, req.Slots)
}

func (s *MentorPortalService) SubmitReport(ctx context.Context, mentorID int64, req model.SubmitReportRequest) (*model.GrowthReport, error) {
	report, err := s.reports.Submit(ctx, mentorID, req, s.earnYuan)
	if errors.Is(err, repository.ErrReportExists) {
		return nil, err
	}
	if errors.Is(err, repository.ErrOrderNotCompleted) {
		return nil, err
	}
	return report, err
}

func (s *MentorPortalService) GetReportForUser(ctx context.Context, userID int64, orderID string) (*model.GrowthReport, error) {
	report, err := s.reports.FindByOrderID(ctx, orderID)
	if errors.Is(err, repository.ErrReportNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if report.UserID != userID {
		return nil, repository.ErrOrderForbidden
	}
	return report, nil
}

func (s *MentorPortalService) Wallet(ctx context.Context, mentorID int64) (*model.WalletSummary, error) {
	balance, err := s.wallet.GetBalance(ctx, mentorID)
	if err != nil {
		return nil, err
	}
	txs, err := s.wallet.ListTransactions(ctx, mentorID, 20)
	if err != nil {
		return nil, err
	}
	return &model.WalletSummary{Balance: balance, Transactions: txs}, nil
}

func (s *MentorPortalService) Withdraw(ctx context.Context, mentorID int64, req model.WithdrawRequest) (*model.WithdrawResponse, error) {
	amount := float64(req.AmountCents) / 100
	balance, err := s.wallet.Withdraw(ctx, mentorID, amount, s.withdrawMock)
	if errors.Is(err, repository.ErrInsufficientBalance) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &model.WithdrawResponse{Balance: balance, MockPaid: s.withdrawMock}, nil
}

func (s *MentorPortalService) GetProfile(ctx context.Context, mentorID int64) (*model.MentorPortalProfile, error) {
	p, err := s.mentors.FindOwnedByID(ctx, mentorID)
	if errors.Is(err, repository.ErrMentorNotFound) {
		return nil, err
	}
	return p, err
}

func (s *MentorPortalService) UpdateProfile(ctx context.Context, mentorID int64, req model.UpdateMentorProfileRequest) (*model.MentorPortalProfile, error) {
	return s.mentors.UpdateProfile(ctx, mentorID, req)
}

func (s *MentorPortalService) ApplyUploadedMedia(ctx context.Context, mentorID int64, kind, url string) error {
	switch kind {
	case "avatar":
		return s.mentors.SetAvatarURL(ctx, mentorID, url)
	case "intro_video":
		return s.mentors.SetIntroVideoURL(ctx, mentorID, url)
	default:
		return ErrInvalidUploadKind
	}
}

var ErrNotMentor = errors.New("not a verified mentor")
var ErrInvalidUploadKind = errors.New("invalid upload kind")
