package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/topmentor/backend/internal/config"
	"github.com/topmentor/backend/pkg/agora"
)

func main() {
	loaded := detectEnvFile()

	cfg, err := config.Load()
	if err != nil {
		fmt.Println("config error:", err)
		os.Exit(1)
	}

	svc := agora.NewTokenService(agora.Config{
		AppID:          cfg.AgoraAppID,
		AppCertificate: cfg.AgoraAppCertificate,
		MockMode:       cfg.AgoraMockMode,
	})

	fmt.Println("=== Agora 配置检测 ===")
	fmt.Println("加载的环境文件:", loaded)
	fmt.Println("AGORA_APP_ID:", mask(cfg.AgoraAppID))
	fmt.Println("AGORA_APP_CERTIFICATE:", mask(cfg.AgoraAppCertificate))
	fmt.Println("AGORA_MOCK_MODE (env):", os.Getenv("AGORA_MOCK_MODE"))
	fmt.Println("config.AgoraMockMode:", cfg.AgoraMockMode)
	fmt.Println("runtime IsMockMode():", svc.IsMockMode())
	fmt.Println("runtime AppID():", svc.AppID())

	token, err := svc.BuildRTCToken("test_channel", 1)
	if err != nil {
		fmt.Println("Token 签发失败:", err)
		os.Exit(1)
	}
	if strings.HasPrefix(token, "mock_rtc_") {
		fmt.Println("Token 类型: MOCK（不会连真实声网）")
	} else {
		fmt.Println("Token 类型: LIVE（真实声网 Token）")
		fmt.Println("Token 长度:", len(token))
	}

	fmt.Println()
	printVerdict(cfg, svc)
}

func detectEnvFile() string {
	if _, err := os.Stat(".env"); err == nil {
		return "backend/.env 或 当前目录 .env"
	}
	if _, err := os.Stat("../.env"); err == nil {
		return "../.env"
	}
	return "(未找到 .env)"
}

func mask(s string) string {
	if s == "" {
		return "(未设置)"
	}
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

func printVerdict(cfg *config.Config, svc *agora.TokenService) {
	var issues []string

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		if _, err := os.Stat("../.env"); os.IsNotExist(err) {
			issues = append(issues, "没有 .env 文件；请执行 copy .env.example .env")
		}
	}
	if cfg.AgoraAppID == "" {
		issues = append(issues, "AGORA_APP_ID 为空")
	}
	if cfg.AgoraAppCertificate == "" {
		issues = append(issues, "AGORA_APP_CERTIFICATE 为空（声网控制台需启用 App 证书）")
	}
	if os.Getenv("AGORA_MOCK_MODE") == "false" && cfg.AgoraAppCertificate == "" {
		issues = append(issues, "AGORA_MOCK_MODE=false 但证书为空，无法走 Live")
	}
	if os.Getenv("AGORA_MOCK_MODE") == "false" && cfg.AgoraAppID == "" {
		issues = append(issues, "AGORA_MOCK_MODE=false 但 App ID 为空")
	}
	if svc.IsMockMode() {
		fmt.Println("结论: 当前为 Mock 模式，可测进房/心跳/下课，无真实音视频。")
	} else {
		fmt.Println("结论: 当前为 Live 模式，H5 将请求摄像头/麦克风并连接声网。")
	}
	if len(issues) > 0 {
		fmt.Println("\n待修复项:")
		for _, i := range issues {
			fmt.Println(" -", i)
		}
	} else {
		fmt.Println("\n配置项检查通过。")
	}
}
