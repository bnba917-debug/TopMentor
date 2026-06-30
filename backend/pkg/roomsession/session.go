package roomsession

import (
	"context"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

const heartbeatTTL = 2 * time.Minute

type Store struct {
	client *goredis.Client
}

func NewStore(client *goredis.Client) *Store {
	return &Store{client: client}
}

func activeKey(orderID string) string {
	return fmt.Sprintf("room:mentor_active_at:%s", orderID)
}

func heartbeatKey(orderID, role string) string {
	return fmt.Sprintf("room:hb:%s:%s", orderID, role)
}

func (s *Store) MarkActive(ctx context.Context, orderID string) error {
	now := time.Now().Unix()
	ok, err := s.client.SetNX(ctx, activeKey(orderID), now, 24*time.Hour).Result()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return nil
}

func (s *Store) ActiveAt(ctx context.Context, orderID string) (time.Time, bool, error) {
	val, err := s.client.Get(ctx, activeKey(orderID)).Result()
	if err == goredis.Nil {
		return time.Time{}, false, nil
	}
	if err != nil {
		return time.Time{}, false, err
	}
	sec, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return time.Time{}, false, err
	}
	return time.Unix(sec, 0), true, nil
}

func (s *Store) RecordHeartbeat(ctx context.Context, orderID, role string) error {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	return s.client.Set(ctx, heartbeatKey(orderID, role), now, heartbeatTTL).Err()
}

func (s *Store) LastHeartbeat(ctx context.Context, orderID, role string) (time.Time, bool, error) {
	val, err := s.client.Get(ctx, heartbeatKey(orderID, role)).Result()
	if err == goredis.Nil {
		return time.Time{}, false, nil
	}
	if err != nil {
		return time.Time{}, false, err
	}
	sec, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return time.Time{}, false, err
	}
	return time.Unix(sec, 0), true, nil
}

func (s *Store) IsOnline(ctx context.Context, orderID, role string) (bool, error) {
	ts, ok, err := s.LastHeartbeat(ctx, orderID, role)
	if err != nil || !ok {
		return false, err
	}
	return time.Since(ts) <= heartbeatTTL, nil
}

func (s *Store) ElapsedMinutes(ctx context.Context, orderID string) (int, error) {
	activeAt, ok, err := s.ActiveAt(ctx, orderID)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, nil
	}
	mins := int(time.Since(activeAt).Minutes())
	if mins < 0 {
		return 0, nil
	}
	return mins, nil
}
