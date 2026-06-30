package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
)

type mockApplyStore struct {
	status *model.MentorApplyStatus
	err    error
}

func (m *mockApplyStore) GetApplyStatus(_ context.Context, _ string) (*model.MentorApplyStatus, error) {
	return m.status, m.err
}

func (m *mockApplyStore) SubmitApplication(_ context.Context, _ int64, _ string, _ model.SubmitMentorApplyRequest) error {
	return m.err
}

type mockApplyUserStore struct {
	user *model.User
}

func (m *mockApplyUserStore) FindByID(_ context.Context, _ int64) (*model.User, error) {
	return m.user, nil
}

func TestMentorApplyService_GetStatus(t *testing.T) {
	svc := NewMentorApplyService(
		&mockApplyStore{status: &model.MentorApplyStatus{Status: "pending"}},
		&mockApplyUserStore{user: &model.User{ID: 1, Phone: "13900000001"}},
	)
	st, err := svc.GetStatus(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "pending", st.Status)
}

func TestMentorApplyService_Submit_PendingConflict(t *testing.T) {
	svc := NewMentorApplyService(
		&mockApplyStore{err: repository.ErrApplicationPending},
		&mockApplyUserStore{user: &model.User{ID: 1, Phone: "13900000001"}},
	)
	_, err := svc.Submit(context.Background(), 1, model.SubmitMentorApplyRequest{RealName: "测试"})
	assert.ErrorIs(t, err, repository.ErrApplicationPending)
}
