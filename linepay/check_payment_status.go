package linepay

import (
	"context"
	"fmt"
	"net/http"
)

// CheckPaymentStatus method
// LINE Pay でのオーソリ履歴の内訳を照会する API です。オーソリ済み、またはオーソリ無効処理データのみ照会できます。売上が確
// 定されたデータは「決済内訳照会 API」で照会できます。
func (c *Client) CheckPaymentStatus(ctx context.Context, transactionID int64, req *CheckPaymentStatusRequest) (*CheckPaymentStatusResponse, *http.Response, error) {
	endpoint := fmt.Sprintf("v3/payments/requests/%d/check", transactionID)
	httpReq, err := c.NewRequest(http.MethodGet, endpoint, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(CheckPaymentStatusResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// CheckPaymentStatusRequest type
type CheckPaymentStatusRequest struct {
}

// CheckPaymentStatusResponse type
type CheckPaymentStatusResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
}
