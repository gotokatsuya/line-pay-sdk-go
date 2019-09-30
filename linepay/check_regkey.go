package linepay

import (
	"context"
	"fmt"
	"net/http"
)

// CheckRegKey method
// 継続決済 API を使用する前に、regKey が使用可能な状態であるかどうかを確認します。
func (c *Client) CheckRegKey(ctx context.Context, regKey string, req *CheckRegKeyRequest) (*CheckRegKeyResponse, *http.Response, error) {
	endpoint := fmt.Sprintf("v3/payments/preapprovedPay/%s/check", regKey)
	httpReq, err := c.NewRequest(http.MethodGet, endpoint, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(CheckRegKeyResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// CheckRegKeyRequest type
type CheckRegKeyRequest struct {
	CreditCardAuth bool `url:"creditCardAuth"`
}

// CheckRegKeyResponse type
type CheckRegKeyResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
}
