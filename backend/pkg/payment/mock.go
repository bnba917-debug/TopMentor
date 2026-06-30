package payment

import "fmt"

type MockProvider struct{}

func (p *MockProvider) Channel() Channel { return ChannelMock }

func (p *MockProvider) Prepare(req PrepareRequest) (PrepareResult, error) {
	return PrepareResult{
		AutoComplete:  true,
		ProviderTxnID: fmt.Sprintf("mock_txn_%s", req.OrderID),
	}, nil
}
