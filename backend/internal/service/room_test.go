package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
	"github.com/topmentor/backend/pkg/agora"
)

type mockOrderStore struct {
	order *model.CourseOrderDetail
}

func (m *mockOrderStore) ListByUser(_ context.Context, userID int64) ([]model.CourseOrderDetail, error) {
	if m.order != nil && m.order.UserID == userID {
		return []model.CourseOrderDetail{*m.order}, nil
	}
	return nil, nil
}

func (m *mockOrderStore) FindDetailByID(_ context.Context, orderID string) (*model.CourseOrderDetail, error) {
	if m.order != nil && m.order.ID == orderID {
		return m.order, nil
	}
	return nil, repository.ErrOrderNotFound
}

func (m *mockOrderStore) Activate(_ context.Context, _ string) error { return nil }

func (m *mockOrderStore) UpdateActualMinutes(_ context.Context, _ string, minutes int) error {
	if m.order != nil {
		m.order.ActualMinutes = minutes
	}
	return nil
}

func (m *mockOrderStore) Complete(_ context.Context, orderID string, minutes int) error {
	if m.order != nil {
		m.order.Status = model.OrderStatusCompleted
		m.order.ActualMinutes = minutes
	}
	return nil
}

type mockRoomSession struct {
	markedActive bool
	activeAt     time.Time
}

func (m *mockRoomSession) MarkActive(_ context.Context, _ string) error {
	if !m.markedActive {
		m.markedActive = true
		m.activeAt = time.Now()
	}
	return nil
}

func (m *mockRoomSession) ActiveAt(_ context.Context, _ string) (time.Time, bool, error) {
	if m.markedActive {
		return m.activeAt, true, nil
	}
	return time.Time{}, false, nil
}

func (m *mockRoomSession) RecordHeartbeat(_ context.Context, _, _ string) error { return nil }

func (m *mockRoomSession) IsOnline(_ context.Context, _ string, _ string) (bool, error) {
	return false, nil
}

func (m *mockRoomSession) ElapsedMinutes(_ context.Context, _ string) (int, error) {
	if !m.markedActive {
		return 0, nil
	}
	return 5, nil
}

func TestRoomService_Join(t *testing.T) {
	order := &model.CourseOrderDetail{
		CourseOrder: model.CourseOrder{
			ID:               "C_test",
			UserID:           1,
			MentorID:         2,
			Status:           model.OrderStatusReserved,
			AgoraChannelName: "tm_C_test",
		},
		MentorName: "张同学",
		SlotDate:   "2026-12-01",
		StartTime:  "19:00:00",
		EndTime:    "19:45:00",
	}

	svc := NewRoomService(
		&mockOrderStore{order: order},
		&mockRoomSession{},
		agora.NewTokenService(agora.Config{MockMode: true}),
		45,
	)

	resp, err := svc.Join(context.Background(), 1, 0, "C_test", "user")
	require.NoError(t, err)
	assert.Equal(t, "tm_C_test", resp.Channel)
	assert.True(t, resp.MockMode)
	assert.Equal(t, uint32(1), resp.UID)
	assert.False(t, resp.LessonStarted)
}

func TestRoomService_Join_MentorStartsLesson(t *testing.T) {
	order := &model.CourseOrderDetail{
		CourseOrder: model.CourseOrder{
			ID:               "C_test",
			UserID:           1,
			MentorID:         2,
			Status:           model.OrderStatusReserved,
			AgoraChannelName: "tm_C_test",
		},
		SlotDate: "2026-12-01", EndTime: "19:45:00",
	}
	svc := NewRoomService(
		&mockOrderStore{order: order},
		&mockRoomSession{},
		agora.NewTokenService(agora.Config{MockMode: true}),
		45,
	)

	resp, err := svc.Join(context.Background(), 99, 2, "C_test", "mentor")
	require.NoError(t, err)
	assert.True(t, resp.LessonStarted)
	assert.Equal(t, model.OrderStatusActive, resp.OrderStatus)
}

func TestRoomService_Join_Mentor(t *testing.T) {
	order := &model.CourseOrderDetail{
		CourseOrder: model.CourseOrder{
			ID:               "C_test",
			UserID:           1,
			MentorID:         2,
			Status:           model.OrderStatusReserved,
			AgoraChannelName: "tm_C_test",
		},
		SlotDate: "2026-12-01", EndTime: "19:45:00",
	}
	svc := NewRoomService(
		&mockOrderStore{order: order},
		&mockRoomSession{},
		agora.NewTokenService(agora.Config{MockMode: true}),
		45,
	)

	resp, err := svc.Join(context.Background(), 99, 2, "C_test", "mentor")
	require.NoError(t, err)
	assert.Equal(t, uint32(100002), resp.UID)
}

func TestRoomService_Join_MentorWrongAccount(t *testing.T) {
	order := &model.CourseOrderDetail{
		CourseOrder: model.CourseOrder{
			ID: "C_test", UserID: 1, MentorID: 2, Status: model.OrderStatusReserved, AgoraChannelName: "ch",
		},
		SlotDate: "2026-12-01", EndTime: "19:45:00",
	}
	svc := NewRoomService(&mockOrderStore{order: order}, &mockRoomSession{}, agora.NewTokenService(agora.Config{MockMode: true}), 45)

	_, err := svc.Join(context.Background(), 99, 3, "C_test", "mentor")
	assert.ErrorIs(t, err, repository.ErrOrderForbidden)
}

func TestRoomService_Join_Forbidden(t *testing.T) {
	order := &model.CourseOrderDetail{
		CourseOrder: model.CourseOrder{ID: "C_test", UserID: 1, Status: model.OrderStatusReserved, AgoraChannelName: "ch"},
		SlotDate: "2026-12-01", EndTime: "19:45:00",
	}
	svc := NewRoomService(&mockOrderStore{order: order}, &mockRoomSession{}, agora.NewTokenService(agora.Config{MockMode: true}), 45)

	_, err := svc.Join(context.Background(), 999, 0, "C_test", "user")
	assert.ErrorIs(t, err, repository.ErrOrderForbidden)
}

func TestRoomService_Heartbeat(t *testing.T) {
	order := &model.CourseOrderDetail{
		CourseOrder: model.CourseOrder{ID: "C_test", UserID: 1, Status: model.OrderStatusActive, AgoraChannelName: "ch"},
		SlotDate: "2026-12-01", EndTime: "19:45:00",
	}
	svc := NewRoomService(&mockOrderStore{order: order}, &mockRoomSession{}, agora.NewTokenService(agora.Config{MockMode: true}), 45)

	resp, err := svc.Heartbeat(context.Background(), 1, 0, "C_test", "user")
	require.NoError(t, err)
	assert.True(t, resp.OK)
	assert.Equal(t, 0, resp.ActualMinutes)
	assert.False(t, resp.LessonStarted)
}

func TestRoomService_Complete_Mentor(t *testing.T) {
	order := &model.CourseOrderDetail{
		CourseOrder: model.CourseOrder{
			ID: "C_test", UserID: 1, MentorID: 2, Status: model.OrderStatusActive, AgoraChannelName: "ch",
		},
	}
	sessions := &mockRoomSession{markedActive: true, activeAt: time.Now()}
	svc := NewRoomService(&mockOrderStore{order: order}, sessions, agora.NewTokenService(agora.Config{MockMode: true}), 45)

	resp, err := svc.Complete(context.Background(), 2, "C_test")
	require.NoError(t, err)
	assert.Equal(t, model.OrderStatusCompleted, resp.Status)
}

func TestRoomService_Complete_ParentForbidden(t *testing.T) {
	order := &model.CourseOrderDetail{
		CourseOrder: model.CourseOrder{
			ID: "C_test", UserID: 1, MentorID: 2, Status: model.OrderStatusActive, AgoraChannelName: "ch",
		},
	}
	svc := NewRoomService(&mockOrderStore{order: order}, &mockRoomSession{}, agora.NewTokenService(agora.Config{MockMode: true}), 45)

	_, err := svc.Complete(context.Background(), 0, "C_test")
	assert.ErrorIs(t, err, repository.ErrOrderForbidden)
}
