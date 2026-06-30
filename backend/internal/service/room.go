package service

import (
	"context"
	"errors"
	"time"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
	"github.com/topmentor/backend/pkg/agora"
)

type OrderStore interface {
	ListByUser(ctx context.Context, userID int64) ([]model.CourseOrderDetail, error)
	FindDetailByID(ctx context.Context, orderID string) (*model.CourseOrderDetail, error)
	Activate(ctx context.Context, orderID string) error
	UpdateActualMinutes(ctx context.Context, orderID string, minutes int) error
	Complete(ctx context.Context, orderID string, minutes int) error
}

type RoomSessionStore interface {
	MarkActive(ctx context.Context, orderID string) error
	ActiveAt(ctx context.Context, orderID string) (time.Time, bool, error)
	RecordHeartbeat(ctx context.Context, orderID, role string) error
	IsOnline(ctx context.Context, orderID, role string) (bool, error)
	ElapsedMinutes(ctx context.Context, orderID string) (int, error)
}

type RoomService struct {
	orders   OrderStore
	sessions RoomSessionStore
	agora    *agora.TokenService
	lessonMin int
}

func NewRoomService(orders OrderStore, sessions RoomSessionStore, agoraSvc *agora.TokenService, lessonMinutes int) *RoomService {
	if lessonMinutes <= 0 {
		lessonMinutes = 45
	}
	return &RoomService{
		orders:    orders,
		sessions:  sessions,
		agora:     agoraSvc,
		lessonMin: lessonMinutes,
	}
}

func (s *RoomService) ListUserOrders(ctx context.Context, userID int64) ([]model.CourseOrderDetail, error) {
	return s.orders.ListByUser(ctx, userID)
}

func (s *RoomService) GetOrder(ctx context.Context, userID int64, orderID string) (*model.CourseOrderDetail, error) {
	order, err := s.orders.FindDetailByID(ctx, orderID)
	if errors.Is(err, repository.ErrOrderNotFound) {
		return nil, repository.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, repository.ErrOrderForbidden
	}
	return order, nil
}

var ErrInvalidRole = errors.New("invalid role")

func (s *RoomService) Join(ctx context.Context, userID, mentorID int64, orderID string, role string) (*model.JoinRoomResponse, error) {
	if role != "user" && role != "mentor" {
		return nil, ErrInvalidRole
	}

	order, err := s.orders.FindDetailByID(ctx, orderID)
	if errors.Is(err, repository.ErrOrderNotFound) {
		return nil, repository.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}

	if role == "user" && order.UserID != userID {
		return nil, repository.ErrOrderForbidden
	}
	if role == "mentor" && (mentorID <= 0 || order.MentorID != mentorID) {
		return nil, repository.ErrOrderForbidden
	}

	if order.Status != model.OrderStatusReserved && order.Status != model.OrderStatusActive {
		return nil, repository.ErrOrderNotJoinable
	}

	respStatus := order.Status
	if role == "mentor" {
		if err := s.orders.Activate(ctx, orderID); err != nil {
			return nil, err
		}
		_ = s.sessions.MarkActive(ctx, orderID)
		respStatus = model.OrderStatusActive
	}

	_ = s.sessions.RecordHeartbeat(ctx, orderID, role)
	status := s.roomStatus(ctx, orderID, order.SlotDate, order.EndTime)

	uid := agora.UIDForRole(order.UserID, order.MentorID, role)
	token, err := s.agora.BuildRTCToken(order.AgoraChannelName, uid)
	if err != nil {
		return nil, err
	}

	return &model.JoinRoomResponse{
		AppID:           s.agora.AppID(),
		Channel:         order.AgoraChannelName,
		Token:           token,
		UID:             uid,
		Role:            role,
		MockMode:        s.agora.IsMockMode(),
		OrderID:         order.ID,
		OrderStatus:     respStatus,
		EndAt:           status.endAt.Format(time.RFC3339),
		DurationMinutes: s.lessonMin,
		LessonStarted:   status.lessonStarted,
		StartedAt:       status.startedAt,
		ElapsedMinutes:  status.elapsed,
		MentorOnline:    status.mentorOnline,
		UserOnline:      status.userOnline,
		MentorName:      order.MentorName,
	}, nil
}

func (s *RoomService) Heartbeat(ctx context.Context, userID, mentorID int64, orderID, role string) (*model.HeartbeatResponse, error) {
	order, err := s.orders.FindDetailByID(ctx, orderID)
	if errors.Is(err, repository.ErrOrderNotFound) {
		return nil, repository.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	if role == "user" && order.UserID != userID {
		return nil, repository.ErrOrderForbidden
	}
	if role == "mentor" && (mentorID <= 0 || order.MentorID != mentorID) {
		return nil, repository.ErrOrderForbidden
	}
	if order.Status != model.OrderStatusActive && order.Status != model.OrderStatusReserved {
		return nil, repository.ErrOrderNotJoinable
	}

	if err := s.sessions.RecordHeartbeat(ctx, orderID, role); err != nil {
		return nil, err
	}

	status := s.roomStatus(ctx, orderID, order.SlotDate, order.EndTime)
	if status.lessonStarted {
		_ = s.orders.UpdateActualMinutes(ctx, orderID, status.elapsed)
	}

	resp := &model.HeartbeatResponse{
		OK:            true,
		ActualMinutes: status.elapsed,
		LessonStarted: status.lessonStarted,
		StartedAt:     status.startedAt,
		MentorOnline:  status.mentorOnline,
		UserOnline:    status.userOnline,
	}
	if status.lessonStarted {
		resp.EndAt = status.endAt.Format(time.RFC3339)
	}
	return resp, nil
}

func (s *RoomService) Complete(ctx context.Context, mentorID int64, orderID string) (*model.CompleteRoomResponse, error) {
	order, err := s.orders.FindDetailByID(ctx, orderID)
	if errors.Is(err, repository.ErrOrderNotFound) {
		return nil, repository.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	if mentorID <= 0 || order.MentorID != mentorID {
		return nil, repository.ErrOrderForbidden
	}

	mins, _ := s.sessions.ElapsedMinutes(ctx, orderID)
	if mins == 0 {
		mins = order.ActualMinutes
	}
	if err := s.orders.Complete(ctx, orderID, mins); err != nil {
		return nil, err
	}

	return &model.CompleteRoomResponse{
		OrderID:       orderID,
		Status:        model.OrderStatusCompleted,
		ActualMinutes: mins,
	}, nil
}

// lessonTiming: 学霸首次进房后开始计时；未开课前返回预约时段结束时刻供展示窗口。
func (s *RoomService) lessonTiming(ctx context.Context, orderID, slotDate, slotEndTime string) (started bool, endAt time.Time) {
	slotEnd, err := repository.SlotEndTime(slotDate, slotEndTime)
	if err != nil {
		slotEnd = time.Now().Add(time.Duration(s.lessonMin) * time.Minute)
	}

	activeAt, ok, err := s.sessions.ActiveAt(ctx, orderID)
	if err != nil || !ok {
		return false, slotEnd
	}

	lessonEnd := activeAt.Add(time.Duration(s.lessonMin) * time.Minute)
	if lessonEnd.Before(slotEnd) {
		return true, lessonEnd
	}
	return true, slotEnd
}

type roomLiveStatus struct {
	lessonStarted bool
	endAt         time.Time
	startedAt     string
	elapsed       int
	mentorOnline  bool
	userOnline    bool
}

func (s *RoomService) roomStatus(ctx context.Context, orderID, slotDate, slotEndTime string) roomLiveStatus {
	lessonStarted, endAt := s.lessonTiming(ctx, orderID, slotDate, slotEndTime)
	st := roomLiveStatus{
		lessonStarted: lessonStarted,
		endAt:         endAt,
	}
	if lessonStarted {
		if activeAt, ok, err := s.sessions.ActiveAt(ctx, orderID); err == nil && ok {
			st.startedAt = activeAt.Format(time.RFC3339)
		}
		st.elapsed, _ = s.sessions.ElapsedMinutes(ctx, orderID)
	}
	st.mentorOnline, _ = s.sessions.IsOnline(ctx, orderID, "mentor")
	st.userOnline, _ = s.sessions.IsOnline(ctx, orderID, "user")
	return st
}
