package payment

import "fmt"
// Configure WX_MCH_ID + WX_MCH_API_KEY + WX_APP_ID to enable in live mode.
type WechatJSAPIProvider struct {
	cfg MerchantConfig
}

func (p *WechatJSAPIProvider) Channel() Channel { return ChannelWechatJSAPI }

func (p *WechatJSAPIProvider) Prepare(req PrepareRequest) (PrepareResult, error) {
	if p.cfg.WxMchID == "" || p.cfg.WxMchAPIKey == "" || p.cfg.WxAppID == "" {
		return PrepareResult{}, ErrNotConfigured
	}
	if req.OpenID == "" {
		return PrepareResult{}, fmt.Errorf("openid required for wechat_jsapi")
	}
	// TODO: call WeChat unified order API and return prepay_id signed params
	return PrepareResult{
		AutoComplete: false,
		JsapiParams: map[string]string{
			"appId":     p.cfg.WxAppID,
			"timeStamp": "PLACEHOLDER",
			"nonceStr":  "PLACEHOLDER",
			"package":   "prepay_id=PLACEHOLDER",
			"signType":  "RSA",
			"paySign":   "PLACEHOLDER",
		},
	}, nil
}
