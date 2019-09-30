package linepay

import (
	"context"
	"net/http"
	"time"
)

// PaymentDetails method
// LINE Payの取引履歴を照会するAPIです。オーソリと売上確定の取引を照会できます。
// "fields"を設定することで、取引情報または注文情報を選択的に照会することができます。
func (c *Client) PaymentDetails(ctx context.Context, req *PaymentDetailsRequest) (*PaymentDetailsResponse, *http.Response, error) {
	endpoint := "v3/payments"
	httpReq, err := c.NewRequest(http.MethodGet, endpoint, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(PaymentDetailsResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// PaymentDetailsRequest type
type PaymentDetailsRequest struct {
	TransactionID []int64  `url:"transactionId,omitempty"`
	OrderID       []string `url:"orderId,omitempty"`
	Fields        string   `url:"fields,omitempty"`
}

// PaymentDetailsResponse type
type PaymentDetailsResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Info          []struct {
		TransactionID           int64     `json:"transactionId"`
		TransactionDate         time.Time `json:"transactionDate"`
		TransactionType         string    `json:"transactionType"`
		PayStatus               string    `json:"payStatus"`
		ProductName             string    `json:"productName"`
		MerchantName            string    `json:"merchantName"`
		Currency                string    `json:"currency"`
		AuthorizationExpireDate string    `json:"authorizationExpireDate"`
		PayInfo                 []struct {
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
