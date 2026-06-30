package payment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockProvider_AutoComplete(t *testing.T) {
	p := &MockProvider{}
	res, err := p.Prepare(PrepareRequest{OrderID: "R001", AmountCents: 9900})
	require.NoError(t, err)
	assert.True(t, res.AutoComplete)
	assert.Contains(t, res.ProviderTxnID, "mock_txn_")
}

func TestFactory_MockMode(t *testing.T) {
	f := NewFactory(MerchantConfig{Mode: "mock"})
	assert.True(t, f.IsMockMode())

	p, err := f.Provider(ChannelWechatJSAPI)
	require.NoError(t, err)
	_, ok := p.(*MockProvider)
	assert.True(t, ok, "mock mode should return MockProvider for any channel")
}

func TestWechatJSAPI_NotConfigured(t *testing.T) {
	f := NewFactory(MerchantConfig{Mode: "live"})
	p, err := f.Provider(ChannelWechatJSAPI)
	require.NoError(t, err)

	_, err = p.Prepare(PrepareRequest{OrderID: "R1", OpenID: "oid"})
	assert.ErrorIs(t, err, ErrNotConfigured)
}

func TestWechatJSAPI_ConfiguredPlaceholder(t *testing.T) {
	f := NewFactory(MerchantConfig{
		Mode:        "live",
		WxAppID:     "wx123",
		WxMchID:     "mch123",
		WxMchAPIKey: "key123",
	})
	p, err := f.Provider(ChannelWechatJSAPI)
	require.NoError(t, err)

	res, err := p.Prepare(PrepareRequest{OrderID: "R1", OpenID: "oid", AmountCents: 100})
	require.NoError(t, err)
	assert.False(t, res.AutoComplete)
	assert.NotEmpty(t, res.JsapiParams["appId"])
}

func TestParseChannel(t *testing.T) {
	ch, err := ParseChannel("wechat_h5")
	require.NoError(t, err)
	assert.Equal(t, ChannelWechatH5, ch)

	_, err = ParseChannel("invalid")
	assert.Error(t, err)
}
