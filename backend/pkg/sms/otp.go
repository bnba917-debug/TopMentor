package sms

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

const (
	keyPrefix     = "sms:otp:"
	ratePrefix    = "sms:rate:"
	defaultTTL    = 5 * time.Minute
	rateLimitTTL  = 60 * time.Second
)

type OTPStore struct {
	client      *goredis.Client
	ttl         time.Duration
	mockMode    bool
	mockCode    string
}

func NewOTPStore(client *goredis.Client, mockMode bool, mockCode string) *OTPStore {
	if mockCode == "" {
		mockCode = "123456"
	}
	return &OTPStore{
		client:   client,
		ttl:      defaultTTL,
		mockMode: mockMode,
		mockCode: mockCode,
	}
}

func (s *OTPStore) Send(ctx context.Context, phone string) (string, error) {
	rateKey := ratePrefix + phone
	set, err := s.client.SetNX(ctx, rateKey, "1", rateLimitTTL).Result()
	if err != nil {
		return "", err
	}
	if !set {
		return "", ErrSendTooFrequent
	}

	code := s.mockCode
	if !s.mockMode {
		generated, err := randomDigits(6)
		if err != nil {
			return "", err
		}
		code = generated
		// TODO: integrate Aliyun/Tencent SMS provider
	}

	if err := s.client.Set(ctx, keyPrefix+phone, code, s.ttl).Err(); err != nil {
		return "", err
	}
	return code, nil
}

func (s *OTPStore) Verify(ctx context.Context, phone, code string) error {
	stored, err := s.client.Get(ctx, keyPrefix+phone).Result()
	if err == goredis.Nil {
		return ErrInvalidCode
	}
	if err != nil {
		return err
	}
	if stored != code {
		return ErrInvalidCode
	}
	_ = s.client.Del(ctx, keyPrefix+phone).Err()
	return nil
}

func randomDigits(n int) (string, error) {
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		v, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		out[i] = byte('0' + v.Int64())
	}
	return string(out), nil
}

var ErrInvalidCode = fmt.Errorf("invalid sms code")
var ErrSendTooFrequent = fmt.Errorf("sms send too frequent")
