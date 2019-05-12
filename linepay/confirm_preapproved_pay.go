package linepay

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ConfirmPreapprovedPay method
// 決済 reserve API で決済タイプ(type)が PREAPPROVED で決済された場合、決済結果の受信時に regKey を受け取ります。
// 継続決済 API は、この regKey を利用し LINE アプリを介さずに直接決済する際に使用します。
func (c *Client) ConfirmPreapprovedPay(ctx context.Context, regKey string, req *ConfirmPreapprovedPayRequest) (*ConfirmPreapprovedPayResponse, *http.Response, error) {
	endpoint := fmt.Sprintf("v2/payments/preapprovedPay/%s/payment", regKey)
	httpReq, err := c.NewRequest(http.MethodPost, endpoint, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(ConfirmPreapprovedPayResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// ConfirmPreapprovedPayRequest type
type ConfirmPreapprovedPayRequest struct {
	ProductName string `json:"productName"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	OrderID     string `json:"orderId"`
	Capture     *bool  `json:"capture,omitempty"`
}

// ConfirmPreapprovedPayResponse type
type ConfirmPreapprovedPayResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Info          struct {
		TransactionID   int64     `json:"transactionId"`
		TransactionDate time.Time `json:"transactionDate"`
	} `json:"info"`
}
