package wx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Session struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
}

type Client struct {
	appID     string
	appSecret string
	mockMode  bool
	httpClient *http.Client
}

func NewClient(appID, appSecret string, mockMode bool) *Client {
	return &Client{
		appID:     appID,
		appSecret: appSecret,
		mockMode:  mockMode,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) Code2Session(ctx context.Context, code string) (*Session, error) {
	if c.mockMode {
		return &Session{
			OpenID:     "mock_" + code,
			SessionKey: "mock_session_key",
		}, nil
	}

	if c.appID == "" || c.appSecret == "" {
		return nil, fmt.Errorf("wechat app credentials not configured")
	}

	endpoint := "https://api.weixin.qq.com/sns/jscode2session"
	q := url.Values{}
	q.Set("appid", c.appID)
	q.Set("secret", c.appSecret)
	q.Set("js_code", code)
	q.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.ErrCode != 0 {
		return nil, fmt.Errorf("wechat api error: %d %s", payload.ErrCode, payload.ErrMsg)
	}
	if payload.OpenID == "" {
		return nil, fmt.Errorf("wechat api returned empty openid")
	}

	return &Session{
		OpenID:     payload.OpenID,
		SessionKey: payload.SessionKey,
		UnionID:    payload.UnionID,
	}, nil
}
