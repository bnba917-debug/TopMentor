package service

import (
	"context"

	"github.com/topmentor/backend/internal/model"
)

type MentorStore interface {
	List(ctx context.Context, q model.MentorListQuery) ([]model.Mentor, int, error)
	FindByID(ctx context.Context, id int64) (*model.Mentor, error)
}

type MentorService struct {
	mentors MentorStore
}

func NewMentorService(mentors MentorStore) *MentorService {
	return &MentorService{mentors: mentors}
}

func (s *MentorService) List(ctx context.Context, q model.MentorListQuery) (*model.MentorListResponse, error) {
	list, total, err := s.mentors.List(ctx, q)
	if err != nil {
		return nil, err
	}
	return &model.MentorListResponse{List: list, Total: total}, nil
}

func (s *MentorService) GetByID(ctx context.Context, id int64) (*model.Mentor, error) {
	return s.mentors.FindByID(ctx, id)
}
