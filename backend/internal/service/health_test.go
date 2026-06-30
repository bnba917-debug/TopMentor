package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPinger struct {
	err error
}

func (m mockPinger) Ping(_ context.Context) error {
	return m.err
}

func TestHealthService_AllUp(t *testing.T) {
	svc := NewHealthService(mockPinger{}, mockPinger{})
	result := svc.Check(context.Background())

	assert.Equal(t, StatusUp, result.DB)
	assert.Equal(t, StatusUp, result.Redis)
	assert.Equal(t, "ok", result.Status)
	assert.True(t, result.IsHealthy())
}

func TestHealthService_DBDown(t *testing.T) {
	svc := NewHealthService(mockPinger{err: assert.AnError}, mockPinger{})
	result := svc.Check(context.Background())

	assert.Equal(t, StatusDown, result.DB)
	assert.Equal(t, StatusUp, result.Redis)
	assert.Equal(t, "degraded", result.Status)
	assert.False(t, result.IsHealthy())
}

func TestHealthService_RedisDown(t *testing.T) {
	svc := NewHealthService(mockPinger{}, mockPinger{err: assert.AnError})
	result := svc.Check(context.Background())

	assert.Equal(t, StatusUp, result.DB)
	assert.Equal(t, StatusDown, result.Redis)
	assert.False(t, result.IsHealthy())
}
