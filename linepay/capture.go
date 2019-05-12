package linepay

import (
	"context"
	"fmt"
	"net/http"
)

// Capture method
// 決済 reserve API を呼び出す時に“capture”パラメータを“false”で指定した場合は、売上処理を行ったうえで決済を完了させることができます。
// 決済 confirm API から 30 日以内に売上を確定させてください。期限を過ぎるとキャンセルとなります。
func (c *Client) Capture(ctx context.Context, transactionID string, req *CaptureRequest) (*CaptureResponse, *http.Response, error) {
	endpoint := fmt.Sprintf("v2/payments/authorizations/%s/capture", transactionID)
	httpReq, err := c.NewRequest(http.MethodPost, endpoint, req)
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
