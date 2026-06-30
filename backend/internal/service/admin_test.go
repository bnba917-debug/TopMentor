package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/topmentor/backend/internal/model"
	jwtpkg "github.com/topmentor/backend/pkg/jwt"
)

type mockAdminStore struct{}

func (m *mockAdminStore) ListPendingMentors(context.Context) ([]model.PendingMentorApplication, error) {
	return nil, nil
}
func (m *mockAdminStore) ReviewMentor(context.Context, int64, model.ReviewMentorRequest) error {
	return nil
}
func (m *mockAdminStore) ListCourseware(context.Context) ([]model.Courseware, error) {
	return nil, nil
}
func (m *mockAdminStore) CreateCourseware(context.Context, model.CreateCoursewareRequest) (*model.Courseware, error) {
	return nil, nil
}
func (m *mockAdminStore) UpdateCourseware(context.Context, int64, model.UpdateCoursewareRequest) (*model.Courseware, error) {
	return nil, nil
}
func (m *mockAdminStore) DeleteCourseware(context.Context, int64) error { return nil }
func (m *mockAdminStore) FinanceSummary(context.Context) (*model.FinanceSummary, error) {
	return &model.FinanceSummary{}, nil
}

func TestAdminService_Login(t *testing.T) {
	svc := NewAdminService(&mockAdminStore{}, jwtpkg.NewManager("secret", 24), "admin", "pass")

	resp, err := svc.Login(context.Background(), model.AdminLoginRequest{Username: "admin", Password: "pass"})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token)

	_, err = svc.Login(context.Background(), model.AdminLoginRequest{Username: "admin", Password: "wrong"})
	assert.ErrorIs(t, err, ErrAdminInvalidCredentials)
}

func TestAdminService_ReviewRejectRequiresReason(t *testing.T) {
	svc := NewAdminService(&mockAdminStore{}, jwtpkg.NewManager("secret", 24), "admin", "pass")
	err := svc.ReviewMentor(context.Background(), 1, model.ReviewMentorRequest{Action: "reject"})
	assert.ErrorIs(t, err, ErrRejectReasonRequired)
}
