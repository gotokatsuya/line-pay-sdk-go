package linepay

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Refund method
// LINE Pay で決済された取引の払い戻しをリクエストします。
// 払い戻しの際には、LINE Pay 会員の決済取引番号を必ず指定しなければなりません。
// 一部払い戻しも可能です。払い戻し可能期間は売上(capture)から 30 日間となります。
func (c *Client) Refund(ctx context.Context, transactionID string, req *RefundRequest) (*RefundResponse, *http.Response, error) {
	endpoint := fmt.Sprintf("v2/payments/%s/refund", transactionID)
	httpReq, err := c.NewRequest(http.MethodPost, endpoint, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(RefundResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// RefundRequest type
type RefundRequest struct {
	RefundAmount int `json:"refundAmount"`
}

// RefundResponse type
type RefundResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Info          struct {
		RefundTransactionID   int64     `json:"refundTransactionId"`
		RefundTransactionDate time.Time `json:"refundTransactionDate"`
	} `json:"info"`
}
