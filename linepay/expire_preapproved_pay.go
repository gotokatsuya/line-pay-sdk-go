package linepay

import (
	"context"
	"fmt"
	"net/http"
)

// ExpirePreapprovedPay method
// 継続決済で登録された regKey 情報を満了させる API です。
// この API を呼び出した以降は、当該の regKey では継続決済することができなくなります。
func (c *Client) ExpirePreapprovedPay(ctx context.Context, regKey string, req *ExpirePreapprovedPayRequest) (*ExpirePreapprovedPayResponse, *http.Response, error) {
	endpoint := fmt.Sprintf("v2/payments/preapprovedPay/%s/expire", regKey)
	httpReq, err := c.NewRequest(http.MethodPost, endpoint, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(ExpirePreapprovedPayResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// ExpirePreapprovedPayRequest type
type ExpirePreapprovedPayRequest struct {
}

// ExpirePreapprovedPayResponse type
type ExpirePreapprovedPayResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
}
