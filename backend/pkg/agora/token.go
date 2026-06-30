package agora

import (
	"fmt"
	"time"

	rtctokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
)

type Config struct {
	AppID          string
	AppCertificate string
	MockMode       bool
	TokenExpireSec uint32
}

type TokenService struct {
	cfg Config
}

func NewTokenService(cfg Config) *TokenService {
	if cfg.TokenExpireSec == 0 {
		cfg.TokenExpireSec = 3600
	}
	return &TokenService{cfg: cfg}
}

func (s *TokenService) IsMockMode() bool {
	return s.cfg.MockMode || s.cfg.AppID == "" || s.cfg.AppCertificate == ""
}

func (s *TokenService) BuildRTCToken(channel string, uid uint32) (string, error) {
	if s.IsMockMode() {
		return fmt.Sprintf("mock_rtc_%s_%d", channel, uid), nil
	}

	expire := uint32(time.Now().Unix()) + s.cfg.TokenExpireSec
	token, err := rtctokenbuilder.BuildTokenWithUid(
		s.cfg.AppID,
		s.cfg.AppCertificate,
		channel,
		uid,
		rtctokenbuilder.RolePublisher,
		expire,
	)
	if err != nil {
		return "", fmt.Errorf("build agora token: %w", err)
	}
	return token, nil
}

func (s *TokenService) AppID() string {
	if s.IsMockMode() {
		return "mock_app_id"
	}
	return s.cfg.AppID
}

// UIDForRole returns a stable Agora uid for parent or mentor in the same channel.
func UIDForRole(orderUserID, orderMentorID int64, role string) uint32 {
	if role == "mentor" {
		return uint32(100000 + orderMentorID)
	}
	return uint32(orderUserID)
}
