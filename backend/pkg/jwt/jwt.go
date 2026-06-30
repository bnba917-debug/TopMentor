package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v4"
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	UserID   int64  `json:"user_id,omitempty"`
	OpenID   string `json:"openid,omitempty"`
	MentorID int64  `json:"mentor_id,omitempty"`
	Role     string `json:"role,omitempty"`
	jwtlib.RegisteredClaims
}

type Manager struct {
	secret     []byte
	expireHour int
}

func NewManager(secret string, expireHours int) *Manager {
	return &Manager{
		secret:     []byte(secret),
		expireHour: expireHours,
	}
}

const RoleAdmin = "admin"
const RoleUser = "user"

func (m *Manager) Issue(userID int64, openID string, mentorID int64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		OpenID:   openID,
		MentorID: mentorID,
		Role:     RoleUser,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(now.Add(time.Duration(m.expireHour) * time.Hour)),
			IssuedAt:  jwtlib.NewNumericDate(now),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) IssueAdmin(username string) (string, error) {
	now := time.Now()
	claims := Claims{
		Role: RoleAdmin,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(now.Add(time.Duration(m.expireHour) * time.Hour)),
			IssuedAt:  jwtlib.NewNumericDate(now),
			Subject:   username,
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) Parse(tokenString string) (*Claims, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &Claims{}, func(t *jwtlib.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
