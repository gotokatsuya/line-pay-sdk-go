package linepay

import (
	"context"
	"fmt"
	"net/http"
)

// CheckPreapprovedPay method
// 継続決済 API を使用する前に、regKey が使用可能な状態であるかどうかを確認します。
func (c *Client) CheckPreapprovedPay(ctx context.Context, regKey string, req *CheckPreapprovedPayRequest) (*CheckPreapprovedPayResponse, *http.Response, error) {
	endpoint := fmt.Sprintf("v2/payments/preapprovedPay/%s/check", regKey)
	httpReq, err := c.NewRequest(http.MethodGet, endpoint, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(CheckPreapprovedPayResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// CheckPreapprovedPayRequest type
type CheckPreapprovedPayRequest struct {
	CreditCardAuth bool `url:"creditCardAuth"`
}

// CheckPreapprovedPayResponse type
type CheckPreapprovedPayResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
}
