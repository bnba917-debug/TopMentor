package sms

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestStore(t *testing.T, mock bool) *OTPStore {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)

	client := goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
	return NewOTPStore(client, mock, "123456")
}

func TestOTPStore_MockSendAndVerify(t *testing.T) {
	store := newTestStore(t, true)

	code, err := store.Send(context.Background(), "13800138000")
	require.NoError(t, err)
	assert.Equal(t, "123456", code)

	err = store.Verify(context.Background(), "13800138000", "123456")
	require.NoError(t, err)

	err = store.Verify(context.Background(), "13800138000", "123456")
	assert.ErrorIs(t, err, ErrInvalidCode)
}

func TestOTPStore_SendTooFrequent(t *testing.T) {
	store := newTestStore(t, true)

	_, err := store.Send(context.Background(), "13800138000")
	require.NoError(t, err)

	_, err = store.Send(context.Background(), "13800138000")
	assert.ErrorIs(t, err, ErrSendTooFrequent)
}

func TestOTPStore_WrongCode(t *testing.T) {
	store := newTestStore(t, true)

	_, err := store.Send(context.Background(), "13800138000")
	require.NoError(t, err)

	err = store.Verify(context.Background(), "13800138000", "000000")
	assert.ErrorIs(t, err, ErrInvalidCode)
}
