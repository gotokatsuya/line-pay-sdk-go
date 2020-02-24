package linepay

import (
	"context"
	"fmt"
	"net/http"
)

// ExpireRegKey method
// 継続決済で登録された regKey 情報を満了させる API です。
// この API を呼び出した以降は、当該の regKey では継続決済することができなくなります。
func (c *Client) ExpireRegKey(ctx context.Context, regKey string, req *ExpireRegKeyRequest) (*ExpireRegKeyResponse, *http.Response, error) {
	path := fmt.Sprintf("/v3/payments/preapprovedPay/%s/expire", regKey)
	httpReq, err := c.NewRequest(http.MethodPost, path, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(ExpireRegKeyResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// ExpireRegKeyRequest type
type ExpireRegKeyRequest struct {
}

// ExpireRegKeyResponse type
type ExpireRegKeyResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
}
