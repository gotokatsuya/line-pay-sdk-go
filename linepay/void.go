package linepay

import (
	"context"
	"fmt"
	"net/http"
)

// Void method
// 決済ステータスがオーソリ状態である決済データを無効化するAPIです。
// Confirm APIを呼び出して決済完了したオーソリ状態の取引を取り消すことができます。
// 取り消しできるのはオーソリ状態の取引だけであり、売上確定済みの取引はRefund APIを使用して返金します。
func (c *Client) Void(ctx context.Context, transactionID int64, req *VoidRequest) (*VoidResponse, *http.Response, error) {
	path := fmt.Sprintf("/v3/payments/authorizations/%d/void", transactionID)
	httpReq, err := c.NewRequest(http.MethodPost, path, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(VoidResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// VoidRequest type
type VoidRequest struct {
}

// VoidResponse type
type VoidResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Info          struct {
		RefundTransactionID   int64  `json:"refundTransactionId"`
		RefundTransactionDate string `json:"refundTransactionDate"`
	} `json:"info"`
}
