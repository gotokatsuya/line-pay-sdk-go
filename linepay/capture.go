package linepay

import (
	"context"
	"fmt"
	"net/http"
)

// Capture method
// Request APIを使って決済をリクエストする際に"options.payment.capture"をfalseに設定した場合、Confirm APIで決済を完了させると決済ステータスは売上確定待ち状態になります。
// 決済を完全に確定するためには、Capture APIを呼び出して売上確定を行う必要があります。
func (c *Client) Capture(ctx context.Context, transactionID int64, req *CaptureRequest) (*CaptureResponse, *http.Response, error) {
	path := fmt.Sprintf("/v3/payments/authorizations/%d/capture", transactionID)
	httpReq, err := c.NewRequest(http.MethodPost, path, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(CaptureResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// CaptureRequest type
type CaptureRequest struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

// CaptureResponse type
type CaptureResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Info          struct {
		TransactionID int64  `json:"transactionId"`
		OrderID       string `json:"orderId"`
		PayInfo       []struct {
			Method string `json:"method"`
			Amount int    `json:"amount"`
		} `json:"payInfo"`
	} `json:"info"`
}
