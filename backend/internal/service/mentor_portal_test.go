package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
)

type mockReportStore struct {
	report *model.GrowthReport
}

func (m *mockReportStore) Submit(_ context.Context, _ int64, _ model.SubmitReportRequest, _ float64) (*model.GrowthReport, error) {
	return m.report, nil
}

func (m *mockReportStore) FindByOrderID(_ context.Context, orderID string) (*model.GrowthReport, error) {
	if m.report == nil || m.report.OrderID != orderID {
		return nil, repository.ErrReportNotFound
	}
	return m.report, nil
}

func TestMentorPortalService_GetReportForUser_Forbidden(t *testing.T) {
	svc := NewMentorPortalService(nil, nil, nil, &mockReportStore{
		report: &model.GrowthReport{OrderID: "ord-1", UserID: 99},
	}, nil, 80, true)

	_, err := svc.GetReportForUser(context.Background(), 1, "ord-1")
	assert.ErrorIs(t, err, repository.ErrOrderForbidden)
}

func TestMentorPortalService_GetReportForUser_OK(t *testing.T) {
	svc := NewMentorPortalService(nil, nil, nil, &mockReportStore{
		report: &model.GrowthReport{OrderID: "ord-1", UserID: 1, Comment: "good"},
	}, nil, 80, true)

	report, err := svc.GetReportForUser(context.Background(), 1, "ord-1")
	require.NoError(t, err)
	assert.Equal(t, "good", report.Comment)
}
