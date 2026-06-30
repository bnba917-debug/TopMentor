package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
)

type mockSlotStore struct {
	slots map[int64]*model.MentorSlot
}

func (m *mockSlotStore) ListByMentor(_ context.Context, q model.SlotListQuery) ([]model.MentorSlot, error) {
	var out []model.MentorSlot
	for _, s := range m.slots {
		if s.MentorID == q.MentorID {
			out = append(out, *s)
		}
	}
	return out, nil
}

func (m *mockSlotStore) FindByID(_ context.Context, slotID int64) (*model.MentorSlot, error) {
	s, ok := m.slots[slotID]
	if !ok {
		return nil, repository.ErrSlotNotFound
	}
	return s, nil
}

type mockBookingStore struct {
	orders map[int64]*model.CourseOrder
}

func (m *mockBookingStore) Create(_ context.Context, userID, slotID int64) (*model.CourseOrder, int, error) {
	if slotID != 1 {
		return nil, 0, repository.ErrSlotUnavailable
	}
	o := &model.CourseOrder{
		ID:       "C_test",
		UserID:   userID,
		MentorID: 1,
		SlotID:   slotID,
		Status:   "RESERVED",
	}
	if m.orders == nil {
		m.orders = map[int64]*model.CourseOrder{}
	}
	m.orders[slotID] = o
	return o, 2, nil
}

type mockSlotLocker struct {
	locked map[int64]bool
}

func (m *mockSlotLocker) TryLock(_ context.Context, mentorID int64, _, _ string, userID int64) (bool, error) {
	if m.locked == nil {
		m.locked = map[int64]bool{}
	}
	if m.locked[mentorID] {
		return false, nil
	}
	m.locked[mentorID] = true
	return true, nil
}

func (m *mockSlotLocker) Unlock(_ context.Context, mentorID int64, _, _ string) error {
	delete(m.locked, mentorID)
	return nil
}

func TestBookingService_CreateBooking(t *testing.T) {
	svc := NewBookingService(
		&mockSlotStore{slots: map[int64]*model.MentorSlot{
			1: {ID: 1, MentorID: 1, SlotDate: "2026-07-01", StartTime: "19:00:00", EndTime: "19:45:00", Status: 0},
		}},
		&mockBookingStore{},
		&mockSlotLocker{},
	)

	resp, err := svc.CreateBooking(context.Background(), 10, model.CreateBookingRequest{SlotID: 1})
	require.NoError(t, err)
	assert.Equal(t, "RESERVED", resp.Order.Status)
	assert.Equal(t, 2, resp.AvailableLessons)
}

func TestBookingService_SlotLocked(t *testing.T) {
	locker := &mockSlotLocker{locked: map[int64]bool{1: true}}
	svc := NewBookingService(
		&mockSlotStore{slots: map[int64]*model.MentorSlot{
			1: {ID: 1, MentorID: 1, SlotDate: "2026-07-01", StartTime: "19:00:00", Status: 0},
		}},
		&mockBookingStore{},
		locker,
	)

	_, err := svc.CreateBooking(context.Background(), 10, model.CreateBookingRequest{SlotID: 1})
	assert.ErrorIs(t, err, ErrSlotLocked)
}
