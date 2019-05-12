package linepay

import (
	"context"
	"fmt"
	"net/http"
)

// VoidAuthorization method
// オーソリ状態の決済データを無効にします。
// 決済 confirm API を呼び出してオーソリ状態まで進んだ場合に、オーソリを無効にするAPI です。
// オーソリ状態の決済のみ無効に処理できます。売上が確定された決済は、「払い戻し API」を使って払い戻し処理を行って下さい。
func (c *Client) VoidAuthorization(ctx context.Context, transactionID string, req *VoidAuthorizationRequest) (*VoidAuthorizationResponse, *http.Response, error) {
	endpoint := fmt.Sprintf("v2/payments/authorizations/%s/void", transactionID)
	httpReq, err := c.NewRequest(http.MethodPost, endpoint, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(VoidAuthorizationResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// VoidAuthorizationRequest type
type VoidAuthorizationRequest struct {
}

// VoidAuthorizationResponse type
type VoidAuthorizationResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
}
