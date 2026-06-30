package service

import (
	"context"
	"errors"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
)

type SlotStore interface {
	ListByMentor(ctx context.Context, q model.SlotListQuery) ([]model.MentorSlot, error)
	FindByID(ctx context.Context, slotID int64) (*model.MentorSlot, error)
}

type BookingStore interface {
	Create(ctx context.Context, userID, slotID int64) (*model.CourseOrder, int, error)
}

type SlotLocker interface {
	TryLock(ctx context.Context, mentorID int64, slotDate, startTime string, userID int64) (bool, error)
	Unlock(ctx context.Context, mentorID int64, slotDate, startTime string) error
}

type BookingService struct {
	slots   SlotStore
	bookings BookingStore
	locker  SlotLocker
}

func NewBookingService(slots SlotStore, bookings BookingStore, locker SlotLocker) *BookingService {
	return &BookingService{slots: slots, bookings: bookings, locker: locker}
}

func (s *BookingService) ListMentorSlots(ctx context.Context, mentorID int64, fromDate, toDate string) ([]model.MentorSlot, error) {
	return s.slots.ListByMentor(ctx, model.SlotListQuery{
		MentorID: mentorID,
		FromDate: fromDate,
		ToDate:   toDate,
	})
}

var ErrSlotLocked = errors.New("slot is being booked by another user")

func (s *BookingService) CreateBooking(ctx context.Context, userID int64, req model.CreateBookingRequest) (*model.BookingResponse, error) {
	slot, err := s.slots.FindByID(ctx, req.SlotID)
	if errors.Is(err, repository.ErrSlotNotFound) {
		return nil, repository.ErrSlotNotFound
	}
	if err != nil {
		return nil, err
	}
	if slot.Status != model.SlotStatusAvailable {
		return nil, repository.ErrSlotUnavailable
	}

	locked, err := s.locker.TryLock(ctx, slot.MentorID, slot.SlotDate, slot.StartTime, userID)
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, ErrSlotLocked
	}
	defer func() {
		_ = s.locker.Unlock(ctx, slot.MentorID, slot.SlotDate, slot.StartTime)
	}()

	order, balance, err := s.bookings.Create(ctx, userID, req.SlotID)
	if err != nil {
		return nil, err
	}

	return &model.BookingResponse{
		Order:            *order,
		AvailableLessons: balance,
	}, nil
}
