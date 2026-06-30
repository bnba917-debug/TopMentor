package slotlock

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

const lockTTL = 10 * time.Second

type Locker struct {
	client *goredis.Client
}

func NewLocker(client *goredis.Client) *Locker {
	return &Locker{client: client}
}

func Key(mentorID int64, slotDate, startTime string) string {
	return fmt.Sprintf("mentor:slot:%d:%s:%s", mentorID, slotDate, startTime)
}

func (l *Locker) TryLock(ctx context.Context, mentorID int64, slotDate, startTime string, userID int64) (bool, error) {
	return l.client.SetNX(ctx, Key(mentorID, slotDate, startTime), userID, lockTTL).Result()
}

func (l *Locker) Unlock(ctx context.Context, mentorID int64, slotDate, startTime string) error {
	return l.client.Del(ctx, Key(mentorID, slotDate, startTime)).Err()
}
