package linepay

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Refund method
// 決済完了(売上確定済み)された取引を返金します。
// 返金時は、LINE Payユーザーの決済取引番号を必ず渡す必要があります。一部返金も可能です。
func (c *Client) Refund(ctx context.Context, transactionID int64, req *RefundRequest) (*RefundResponse, *http.Response, error) {
	endpoint := fmt.Sprintf("v3/payments/%d/refund", transactionID)
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
	RefundAmount int `json:"refundAmount,omitempty"`
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
