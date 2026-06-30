package agora

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenService_MockMode(t *testing.T) {
	svc := NewTokenService(Config{MockMode: true})

	token, err := svc.BuildRTCToken("tm_Cabc", 42)
	require.NoError(t, err)
	assert.Contains(t, token, "mock_rtc_")
	assert.Equal(t, "mock_app_id", svc.AppID())
	assert.True(t, svc.IsMockMode())
}

func TestUIDForRole(t *testing.T) {
	assert.Equal(t, uint32(10), UIDForRole(10, 5, "user"))
	assert.Equal(t, uint32(100005), UIDForRole(10, 5, "mentor"))
}
