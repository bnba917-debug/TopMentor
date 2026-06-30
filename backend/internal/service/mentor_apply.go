package service

import (
	"context"
	"errors"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
)

type MentorApplyStore interface {
	GetApplyStatus(ctx context.Context, phone string) (*model.MentorApplyStatus, error)
	SubmitApplication(ctx context.Context, userID int64, phone string, req model.SubmitMentorApplyRequest) error
}

type MentorApplyUserStore interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)
}

type MentorApplyService struct {
	mentors MentorApplyStore
	users   MentorApplyUserStore
}

func NewMentorApplyService(mentors MentorApplyStore, users MentorApplyUserStore) *MentorApplyService {
	return &MentorApplyService{mentors: mentors, users: users}
}

func (s *MentorApplyService) GetStatus(ctx context.Context, userID int64) (*model.MentorApplyStatus, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.mentors.GetApplyStatus(ctx, user.Phone)
}

func (s *MentorApplyService) Submit(ctx context.Context, userID int64, req model.SubmitMentorApplyRequest) (*model.MentorApplyStatus, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.mentors.SubmitApplication(ctx, userID, user.Phone, req); err != nil {
		if errors.Is(err, repository.ErrAlreadyVerifiedMentor) {
			return nil, err
		}
		if errors.Is(err, repository.ErrApplicationPending) {
			return nil, err
		}
		return nil, err
	}

	return s.mentors.GetApplyStatus(ctx, user.Phone)
}
