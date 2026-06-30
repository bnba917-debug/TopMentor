package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/topmentor/backend/internal/model"
)

type mockMentorStore struct {
	list  []model.Mentor
	total int
	byID  map[int64]*model.Mentor
}

func (m *mockMentorStore) List(_ context.Context, q model.MentorListQuery) ([]model.Mentor, int, error) {
	return m.list, m.total, nil
}

func (m *mockMentorStore) FindByID(_ context.Context, id int64) (*model.Mentor, error) {
	if mentor, ok := m.byID[id]; ok {
		return mentor, nil
	}
	return nil, nil
}

func TestMentorService_List(t *testing.T) {
	store := &mockMentorStore{
		list: []model.Mentor{
			{ID: 1, RealName: "张同学", SchoolName: "清华大学", Tags: []string{"阳光幽默"}},
		},
		total: 1,
	}
	svc := NewMentorService(store)

	result, err := svc.List(context.Background(), model.MentorListQuery{Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.List, 1)
}

func TestMentorService_GetByID_Found(t *testing.T) {
	store := &mockMentorStore{
		byID: map[int64]*model.Mentor{
			2: {ID: 2, RealName: "李同学"},
		},
	}
	svc := NewMentorService(store)

	mentor, err := svc.GetByID(context.Background(), 2)
	require.NoError(t, err)
	assert.Equal(t, "李同学", mentor.RealName)
}

func TestMentorService_GetByID_NotFound(t *testing.T) {
	svc := NewMentorService(&mockMentorStore{byID: map[int64]*model.Mentor{}})

	mentor, err := svc.GetByID(context.Background(), 999)
	require.NoError(t, err)
	assert.Nil(t, mentor)
}
