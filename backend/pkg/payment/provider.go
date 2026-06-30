package payment

import "fmt"

type Channel string

const (
	ChannelMock        Channel = "mock"
	ChannelWechatJSAPI Channel = "wechat_jsapi"
	ChannelWechatH5    Channel = "wechat_h5"
	ChannelAlipayH5    Channel = "alipay_h5"
)

func ParseChannel(s string) (Channel, error) {
	switch Channel(s) {
	case ChannelMock, ChannelWechatJSAPI, ChannelWechatH5, ChannelAlipayH5:
		return Channel(s), nil
	default:
		return "", fmt.Errorf("unsupported payment channel: %s", s)
	}
}

type PrepareRequest struct {
	OrderID     string
	UserID      int64
	AmountCents int
	Description string
	ClientIP    string
	OpenID      string // required for wechat_jsapi
}

type PrepareResult struct {
	AutoComplete  bool
	PayURL        string
	JsapiParams   map[string]string
	ProviderTxnID string
}

type Provider interface {
	Channel() Channel
	Prepare(req PrepareRequest) (PrepareResult, error)
}

var ErrNotConfigured = fmt.Errorf("payment provider not configured")

type MerchantConfig struct {
	Mode          string // mock | live
	WxAppID       string
	WxMchID       string
	WxMchAPIKey   string
	WxNotifyURL   string
	AlipayAppID   string
	AlipayPrivKey string
	AlipayPubKey  string
	AlipayNotifyURL string
	SiteURL       string
}

type Factory struct {
	cfg MerchantConfig
}

func NewFactory(cfg MerchantConfig) *Factory {
	if cfg.Mode == "" {
		cfg.Mode = "mock"
	}
	return &Factory{cfg: cfg}
}

func (f *Factory) IsMockMode() bool {
	return f.cfg.Mode == "mock"
}

func (f *Factory) Provider(ch Channel) (Provider, error) {
	if f.IsMockMode() || ch == ChannelMock {
		return &MockProvider{}, nil
	}

	switch ch {
	case ChannelWechatJSAPI:
		return &WechatJSAPIProvider{cfg: f.cfg}, nil
	case ChannelWechatH5:
		return &WechatH5Provider{cfg: f.cfg}, nil
	case ChannelAlipayH5:
		return &AlipayH5Provider{cfg: f.cfg}, nil
	default:
		return nil, fmt.Errorf("unsupported channel")
	}
}

func (f *Factory) ConfiguredChannels() []Channel {
	if f.IsMockMode() {
		return []Channel{ChannelMock, ChannelWechatJSAPI, ChannelWechatH5, ChannelAlipayH5}
	}
	var out []Channel
	out = append(out, ChannelMock)
	if f.cfg.WxMchID != "" && f.cfg.WxMchAPIKey != "" {
		out = append(out, ChannelWechatJSAPI, ChannelWechatH5)
	}
	if f.cfg.AlipayAppID != "" && f.cfg.AlipayPrivKey != "" {
		out = append(out, ChannelAlipayH5)
	}
	return out
}
