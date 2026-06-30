package slotlock

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocker_TryLock(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
	locker := NewLocker(client)
	ctx := context.Background()

	ok, err := locker.TryLock(ctx, 1, "2026-07-01", "19:00:00", 100)
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = locker.TryLock(ctx, 1, "2026-07-01", "19:00:00", 200)
	require.NoError(t, err)
	assert.False(t, ok)

	require.NoError(t, locker.Unlock(ctx, 1, "2026-07-01", "19:00:00"))

	ok, err = locker.TryLock(ctx, 1, "2026-07-01", "19:00:00", 200)
	require.NoError(t, err)
	assert.True(t, ok)
}
