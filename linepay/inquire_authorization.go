package linepay

import (
	"context"
	"net/http"
	"time"
)

// InquireAuthorization method
// LINE Pay でのオーソリ履歴の内訳を照会する API です。オーソリ済み、またはオーソリ無効処理データのみ照会できます。売上が確
// 定されたデータは「決済内訳照会 API」で照会できます。
func (c *Client) InquireAuthorization(ctx context.Context, req *InquireAuthorizationRequest) (*InquireAuthorizationResponse, *http.Response, error) {
	endpoint, err := mergeQuery("v2/payments/authorizations", req)
	if err != nil {
		return nil, nil, err
	}
	httpReq, err := c.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	resp := new(InquireAuthorizationResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// InquireAuthorizationRequest type
type InquireAuthorizationRequest struct {
	TransactionID int64  `url:"transactionId,omitempty"`
	OrderID       string `url:"orderId,omitempty"`
}

// InquireAuthorizationResponse type
type InquireAuthorizationResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Info          []struct {
		TransactionID   int64     `json:"transactionId"`
		TransactionDate time.Time `json:"transactionDate"`
		TransactionType string    `json:"transactionType"`
		PayInfo         []struct {
			Method string `json:"method"`
			Amount int    `json:"amount"`
		} `json:"payInfo"`
		ProductName             string    `json:"productName"`
		Currency                string    `json:"currency"`
		OrderID                 string    `json:"orderId"`
		PayStatus               string    `json:"payStatus"`
		AuthorizationExpireDate time.Time `json:"authorizationExpireDate"`
	} `json:"info"`
}
