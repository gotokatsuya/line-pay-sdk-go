package linepay

import (
	"context"
	"fmt"
	"net/http"
)

// Confirm method
// confirmUrlまたはCheck Payment Status APIによってユーザーが決済要求を承認した後、加盟店側で決済を完了させるためのAPIです。
// Request APIの"options.payment.capture"をfalseに設定するとオーソリと売上確定が分離された決済になり、決済を完了させても決済ステータスは売上確定待ち(オーソリ)状態のままとなります。
// 売上を確定するには、Capture APIを呼び出して売上確定を行う必要があります。
func (c *Client) Confirm(ctx context.Context, transactionID int64, req *ConfirmRequest) (*ConfirmResponse, *http.Response, error) {
	endpoint := fmt.Sprintf("v3/payments/%d/confirm", transactionID)
	httpReq, err := c.NewRequest(http.MethodPost, endpoint, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(ConfirmResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// ConfirmRequest type
type ConfirmRequest struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

// ConfirmResponse type
type ConfirmResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Info          struct {
		OrderID                 string `json:"orderId"`
		TransactionID           int64  `json:"transactionId"`
		AuthorizationExpireDate string `json:"authorizationExpireDate,omitempty"`
		RegKey                  string `json:"regKey,omitempty"`
		PayInfo                 []struct {
			Method string `json:"method"`
			Amount int    `json:"amount"`
		} `json:"payInfo"`
	} `json:"info"`
}
