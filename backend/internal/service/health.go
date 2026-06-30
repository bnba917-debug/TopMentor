package service

import "context"

type ComponentStatus string

const (
	StatusUp   ComponentStatus = "up"
	StatusDown ComponentStatus = "down"
)

type HealthResult struct {
	Status string          `json:"status"`
	DB     ComponentStatus `json:"db"`
	Redis  ComponentStatus `json:"redis"`
}

func (r HealthResult) IsHealthy() bool {
	return r.DB == StatusUp && r.Redis == StatusUp
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type HealthService struct {
	db    Pinger
	redis Pinger
}

func NewHealthService(db, redis Pinger) *HealthService {
	return &HealthService{db: db, redis: redis}
}

func (s *HealthService) Check(ctx context.Context) HealthResult {
	result := HealthResult{
		Status: "ok",
		DB:     StatusDown,
		Redis:  StatusDown,
	}

	if s.db != nil {
		if err := s.db.Ping(ctx); err == nil {
			result.DB = StatusUp
		}
	}

	if s.redis != nil {
		if err := s.redis.Ping(ctx); err == nil {
			result.Redis = StatusUp
		}
	}

	if !result.IsHealthy() {
		result.Status = "degraded"
	}

	return result
}
