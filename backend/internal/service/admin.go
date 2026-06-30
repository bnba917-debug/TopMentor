package service

import (
	"context"
	"errors"
	"strings"

	"github.com/topmentor/backend/internal/model"
	jwtpkg "github.com/topmentor/backend/pkg/jwt"
)

type AdminStore interface {
	ListPendingMentors(ctx context.Context) ([]model.PendingMentorApplication, error)
	ReviewMentor(ctx context.Context, mentorID int64, req model.ReviewMentorRequest) error
	ListCourseware(ctx context.Context) ([]model.Courseware, error)
	CreateCourseware(ctx context.Context, req model.CreateCoursewareRequest) (*model.Courseware, error)
	UpdateCourseware(ctx context.Context, id int64, req model.UpdateCoursewareRequest) (*model.Courseware, error)
	DeleteCourseware(ctx context.Context, id int64) error
	FinanceSummary(ctx context.Context) (*model.FinanceSummary, error)
}

type AdminService struct {
	admin    AdminStore
	jwt      *jwtpkg.Manager
	username string
	password string
}

func NewAdminService(admin AdminStore, jwt *jwtpkg.Manager, username, password string) *AdminService {
	return &AdminService{admin: admin, jwt: jwt, username: username, password: password}
}

var ErrAdminInvalidCredentials = errors.New("invalid admin credentials")
var ErrRejectReasonRequired = errors.New("reject reason required")

func (s *AdminService) Login(_ context.Context, req model.AdminLoginRequest) (*model.AdminLoginResponse, error) {
	if req.Username != s.username || req.Password != s.password {
		return nil, ErrAdminInvalidCredentials
	}
	token, err := s.jwt.IssueAdmin(req.Username)
	if err != nil {
		return nil, err
	}
	return &model.AdminLoginResponse{Token: token, Username: req.Username}, nil
}

func (s *AdminService) ListPendingMentors(ctx context.Context) ([]model.PendingMentorApplication, error) {
	return s.admin.ListPendingMentors(ctx)
}

func (s *AdminService) ReviewMentor(ctx context.Context, mentorID int64, req model.ReviewMentorRequest) error {
	if req.Action == "reject" && strings.TrimSpace(req.RejectReason) == "" {
		return ErrRejectReasonRequired
	}
	return s.admin.ReviewMentor(ctx, mentorID, req)
}

func (s *AdminService) ListCourseware(ctx context.Context) ([]model.Courseware, error) {
	return s.admin.ListCourseware(ctx)
}

func (s *AdminService) CreateCourseware(ctx context.Context, req model.CreateCoursewareRequest) (*model.Courseware, error) {
	return s.admin.CreateCourseware(ctx, req)
}

func (s *AdminService) UpdateCourseware(ctx context.Context, id int64, req model.UpdateCoursewareRequest) (*model.Courseware, error) {
	return s.admin.UpdateCourseware(ctx, id, req)
}

func (s *AdminService) DeleteCourseware(ctx context.Context, id int64) error {
	return s.admin.DeleteCourseware(ctx, id)
}

func (s *AdminService) FinanceSummary(ctx context.Context) (*model.FinanceSummary, error) {
	return s.admin.FinanceSummary(ctx)
}
