package payment

import "fmt"

// AlipayH5Provider handles Alipay mobile website payment.
type AlipayH5Provider struct {
	cfg MerchantConfig
}

func (p *AlipayH5Provider) Channel() Channel { return ChannelAlipayH5 }

func (p *AlipayH5Provider) Prepare(req PrepareRequest) (PrepareResult, error) {
	if p.cfg.AlipayAppID == "" || p.cfg.AlipayPrivKey == "" {
		return PrepareResult{}, ErrNotConfigured
	}
	// TODO: call alipay.trade.wap.pay
	return PrepareResult{
		AutoComplete: false,
		PayURL:       fmt.Sprintf("%s/pay/alipay/pending?order_id=%s", p.cfg.SiteURL, req.OrderID),
	}, nil
}
