package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/topmentor/backend/pkg/payment"
)

type Config struct {
	AppEnv     string
	ServerPort int

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret      string
	JWTExpireHours int

	WxAppID     string
	WxAppSecret string
	WxMockMode  bool
	SMSMockMode bool
	SMSMockCode string

	AdminUsername  string
	AdminPassword  string

	CORSOrigins []string

	PaymentMode       string
	SiteURL           string
	WxMchID           string
	WxMchAPIKey       string
	WxNotifyURL       string
	AlipayAppID       string
	AlipayPrivateKey  string
	AlipayPublicKey   string
	AlipayNotifyURL   string

	AgoraAppID          string
	AgoraAppCertificate string
	AgoraMockMode       bool
	LessonDurationMin   int
	LessonEarnYuan      float64
	WithdrawMockMode    bool
	UploadDir           string
}

func (c *Config) PaymentMerchantConfig() payment.MerchantConfig {
	return payment.MerchantConfig{
		Mode:            c.PaymentMode,
		WxAppID:         c.WxAppID,
		WxMchID:         c.WxMchID,
		WxMchAPIKey:     c.WxMchAPIKey,
		WxNotifyURL:     c.WxNotifyURL,
		AlipayAppID:     c.AlipayAppID,
		AlipayPrivKey:   c.AlipayPrivateKey,
		AlipayPubKey:    c.AlipayPublicKey,
		AlipayNotifyURL: c.AlipayNotifyURL,
		SiteURL:         c.SiteURL,
	}
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	jwtExpire, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "168"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRE_HOURS: %w", err)
	}

	wxAppID := getEnv("WX_APP_ID", "")
	wxMock := getEnv("WX_MOCK_MODE", "") == "true" || wxAppID == ""
	smsMock := getEnv("SMS_MOCK_MODE", "") != "false" // default true for dev

	corsRaw := getEnv("CORS_ORIGINS", "http://localhost:5173,http://127.0.0.1:5173,http://localhost:5174,http://127.0.0.1:5174")
	var corsOrigins []string
	for _, o := range splitComma(corsRaw) {
		if o != "" {
			corsOrigins = append(corsOrigins, o)
		}
	}

	lessonMin, err := strconv.Atoi(getEnv("LESSON_DURATION_MINUTES", "45"))
	if err != nil {
		return nil, fmt.Errorf("invalid LESSON_DURATION_MINUTES: %w", err)
	}
	agoraAppID := getEnv("AGORA_APP_ID", "")
	agoraCert := getEnv("AGORA_APP_CERTIFICATE", "")
	agoraMock := resolveAgoraMockMode(getEnv("AGORA_MOCK_MODE", ""), agoraAppID, agoraCert)

	earnYuan, err := strconv.ParseFloat(getEnv("LESSON_EARN_YUAN", "80"), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid LESSON_EARN_YUAN: %w", err)
	}
	withdrawMock := getEnv("WITHDRAW_MOCK_MODE", "") != "false"

	return &Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		ServerPort:     port,
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         dbPort,
		DBUser:         getEnv("DB_USER", "topmentor"),
		DBPassword:     getEnv("DB_PASSWORD", "topmentor_dev"),
		DBName:         getEnv("DB_NAME", "topmentor"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		RedisDB:        redisDB,
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpireHours: jwtExpire,
		WxAppID:        wxAppID,
		WxAppSecret:    getEnv("WX_APP_SECRET", ""),
		WxMockMode:     wxMock,
		SMSMockMode:    smsMock,
		SMSMockCode:    getEnv("SMS_MOCK_CODE", "123456"),
		CORSOrigins:    corsOrigins,
		PaymentMode:    getEnv("PAYMENT_MODE", "mock"),
		SiteURL:        getEnv("SITE_URL", "http://localhost:5173"),
		WxMchID:        getEnv("WX_MCH_ID", ""),
		WxMchAPIKey:    getEnv("WX_MCH_API_KEY", ""),
		WxNotifyURL:    getEnv("WX_NOTIFY_URL", ""),
		AlipayAppID:    getEnv("ALIPAY_APP_ID", ""),
		AlipayPrivateKey: getEnv("ALIPAY_PRIVATE_KEY", ""),
		AlipayPublicKey:  getEnv("ALIPAY_PUBLIC_KEY", ""),
		AlipayNotifyURL:  getEnv("ALIPAY_NOTIFY_URL", ""),
		AgoraAppID:          agoraAppID,
		AgoraAppCertificate: agoraCert,
		AgoraMockMode:       agoraMock,
		LessonDurationMin:   lessonMin,
		LessonEarnYuan:      earnYuan,
		WithdrawMockMode:    withdrawMock,
		AdminUsername:       getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:       getEnv("ADMIN_PASSWORD", "admin123"),
		UploadDir:           getEnv("UPLOAD_DIR", "uploads"),
	}, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// resolveAgoraMockMode: explicit true → mock; explicit false → live if credentials present.
func resolveAgoraMockMode(envVal, appID, cert string) bool {
	switch envVal {
	case "true":
		return true
	case "false":
		return appID == "" || cert == ""
	default:
		return appID == "" || cert == ""
	}
}

func splitComma(s string) []string {
	var parts []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	return parts
}
