package linepay

import (
	"context"
	"fmt"
	"net/http"
)

// Confirm method
// 加盟店が決済を最終的に完了させるための API です。加盟店で決済 confirm API を呼び出すことによって、
// 実際の決済が完了し ます。決済 reserve 時に“capture”パラメータが“false”の場合、confirm API 実行時はオーソリ状態になるため、
// 「capture API」実行時に決済完了となります。
func (c *Client) Confirm(ctx context.Context, transactionID string, req *ConfirmRequest) (*ConfirmResponse, *http.Response, error) {
	endpoint := fmt.Sprintf("v2/payments/%s/confirm", transactionID)
	httpReq, err := c.NewRequest(http.MethodPost, endpoint, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(ConfirmResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// ConfirmRequest type
type ConfirmRequest struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

// ConfirmResponse type
type ConfirmResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Info          struct {
		OrderID       string `json:"orderId"`
		TransactionID int64  `json:"transactionId"`
		PayInfo       []struct {
			Method string `json:"method"`
			Amount int    `json:"amount"`
		} `json:"payInfo"`
	} `json:"info"`
}
