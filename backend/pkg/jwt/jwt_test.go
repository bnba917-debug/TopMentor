package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_IssueAndParse(t *testing.T) {
	mgr := NewManager("test-secret", 24)

	token, err := mgr.Issue(42, "openid-abc", 0)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := mgr.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, int64(42), claims.UserID)
	assert.Equal(t, "openid-abc", claims.OpenID)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

func TestManager_Parse_InvalidToken(t *testing.T) {
	mgr := NewManager("test-secret", 24)

	_, err := mgr.Parse("not-a-valid-token")
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestManager_IssueAdmin(t *testing.T) {
	mgr := NewManager("test-secret", 24)

	token, err := mgr.IssueAdmin("admin")
	require.NoError(t, err)

	claims, err := mgr.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, RoleAdmin, claims.Role)
	assert.Equal(t, "admin", claims.Subject)
}

func TestManager_Parse_WrongSecret(t *testing.T) {
	issuer := NewManager("secret-a", 24)
	parser := NewManager("secret-b", 24)

	token, err := issuer.Issue(1, "oid", 0)
	require.NoError(t, err)

	_, err = parser.Parse(token)
	assert.ErrorIs(t, err, ErrInvalidToken)
}
