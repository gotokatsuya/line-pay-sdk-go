package linepay

import (
	"context"
	"net/http"
)

// Reserve method
// LINE Pay 決済を行う前に、加盟店の状態が正常であるかを判断し、決済のための情報を予約します。
// 決済予約が成功したら、決済完了/払い戻しするまで使用する「取引番号」が発行されます。
func (c *Client) Reserve(ctx context.Context, req *ReserveRequest) (*ReserveResponse, *http.Response, error) {
	endpoint := "v2/payments/request"
	httpReq, err := c.NewRequest(http.MethodPost, endpoint, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(ReserveResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// ReserveRequest type
type ReserveRequest struct {
	ProductName            string `json:"productName"`
	ProductImageURL        string `json:"productImageUrl,omitempty"`
	Amount                 int    `json:"amount"`
	Currency               string `json:"currency"`
	Mid                    string `json:"mid,omitempty"`
	OneTimeKey             string `json:"oneTimeKey,omitempty"`
	ConfirmURL             string `json:"confirmUrl"`
	ConfirmURLType         string `json:"confirmUrlType,omitempty"`
	CheckConfirmURLBrowser *bool  `json:"checkConfirmUrlBrowser,omitempty"`
	CancelURL              string `json:"cancelUrl,omitempty"`
	PackageName            string `json:"packageName,omitempty"`
	OrderID                string `json:"orderId"`
	DeliveryPlacePhone     string `json:"deliveryPlacePhone,omitempty"`
	PayType                string `json:"payType,omitempty"`
	LangCd                 string `json:"langCd,omitempty"`
	Capture                *bool  `json:"capture,omitempty"`
	Extras                 struct {
		AddFriends []struct {
			Type   string   `json:"type"`
			IDList []string `json:"idList"`
		} `json:"addFriends,omitempty"`
		BranchName string `json:"branchName,omitempty"`
	} `json:"extras,omitempty"`
}

// ReserveResponse type
type ReserveResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Info          struct {
		TransactionID int64 `json:"transactionId"`
		PaymentURL    struct {
			Web string `json:"web"`
			App string `json:"app"`
		} `json:"paymentUrl"`
		PaymentAccessToken string `json:"paymentAccessToken"`
	} `json:"info"`
}
