package linepay

import (
	"context"
	"net/http"
	"time"
)

// InquirePayment method
// LINE Pay の決済データを照会します。売上が確定されたデータのみ照会できます。
func (c *Client) InquirePayment(ctx context.Context, req *InquirePaymentRequest) (*InquirePaymentResponse, *http.Response, error) {
	endpoint, err := mergeQuery("v2/payments/payments", req)
	if err != nil {
		return nil, nil, err
	}
	httpReq, err := c.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	resp := new(InquirePaymentResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// InquirePaymentRequest type
type InquirePaymentRequest struct {
	TransactionID string `url:"transactionId,omitempty"`
	OrderID       string `url:"orderId,omitempty"`
}

// InquirePaymentResponse type
type InquirePaymentResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Info          []struct {
		TransactionID      int64     `json:"transactionId"`
		TransactionDate    time.Time `json:"transactionDate"`
		TransactionType    string    `json:"transactionType"`
		ProductName        string    `json:"productName"`
		MerchantName       string    `json:"merchantName"`
		Currency           string    `json:"currency"`
		AuthorizationExpir string    `json:"authorizationExpir"`
		PayInfo            []struct {
			Method string `json:"method"`
			Amount int    `json:"amount"`
		} `json:"payInfo"`

		// 原決済取引照会、および払い戻し取引がある場合
		RefundList []struct {
			RefundTransactionID   string    `json:"refundTransactionId"`
			TransactionType       string    `json:"transactionType"`
			RefundAmount          int       `json:"refundAmount"`
			RefundTransactionDate time.Time `json:"refundTransactionDate"`
		} `json:"refundList,omitempty"`

		// 払い戻し取引の照会の場合
		OriginalTransactionID int64 `json:"originalTransactionId,omitempty"`
	} `json:"info"`
}
