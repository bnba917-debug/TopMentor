package roomsession

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)
	return NewStore(goredis.NewClient(&goredis.Options{Addr: mr.Addr()}))
}

func TestStore_MarkActiveAndElapsed(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	require.NoError(t, store.MarkActive(ctx, "C001"))
	require.NoError(t, store.MarkActive(ctx, "C001"))

	mins, err := store.ElapsedMinutes(ctx, "C001")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, mins, 0)
}

func TestStore_Heartbeat(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	require.NoError(t, store.RecordHeartbeat(ctx, "C001", "user"))
	ts, ok, err := store.LastHeartbeat(ctx, "C001", "user")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.WithinDuration(t, time.Now(), ts, 2*time.Second)

	online, err := store.IsOnline(ctx, "C001", "user")
	require.NoError(t, err)
	assert.True(t, online)

	online, err = store.IsOnline(ctx, "C001", "mentor")
	require.NoError(t, err)
	assert.False(t, online)
}
