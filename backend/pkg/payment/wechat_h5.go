package payment

import "fmt"

// WechatH5Provider handles mobile browser redirect to WeChat app.
type WechatH5Provider struct {
	cfg MerchantConfig
}

func (p *WechatH5Provider) Channel() Channel { return ChannelWechatH5 }

func (p *WechatH5Provider) Prepare(req PrepareRequest) (PrepareResult, error) {
	if p.cfg.WxMchID == "" || p.cfg.WxMchAPIKey == "" || p.cfg.WxAppID == "" {
		return PrepareResult{}, ErrNotConfigured
	}
	// TODO: call WeChat H5 unified order API (trade_type=MWEB)
	return PrepareResult{
		AutoComplete: false,
		PayURL:       fmt.Sprintf("%s/pay/wechat-h5/pending?order_id=%s", p.cfg.SiteURL, req.OrderID),
	}, nil
}
