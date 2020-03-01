package linepay

import (
	"context"
	"net/http"
)

// Request method
// LINE Pay決済をリクエストします。このとき、ユーザーの注文情報と決済手段を設定できます。
// リクエストに成功するとLINE Pay取引番号が発行されます。この取引番号を利用して、決済完了・返金を行うことができます。
func (c *Client) Request(ctx context.Context, req *RequestRequest) (*RequestResponse, *http.Response, error) {
	path := "/v3/payments/request"
	httpReq, err := c.NewRequest(http.MethodPost, path, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(RequestResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

// RequestPackageProduct type
type RequestPackageProduct struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name"`
	ImageURL      string `json:"imageUrl,omitempty"`
	Quantity      int    `json:"quantity"`
	Price         int    `json:"price"`
	OriginalPrice int    `json:"originalPrice,omitempty"`
}

// RequestPackage type
type RequestPackage struct {
	ID       string                   `json:"id"`
	Amount   int                      `json:"amount"`
	UserFee  int                      `json:"userFee,omitempty"`
	Name     string                   `json:"name"`
	Products []*RequestPackageProduct `json:"products"`
}

// RequestRedirectURLs type
type RequestRedirectURLs struct {
	AppPackageName string `json:"appPackageName,omitempty"`
	ConfirmURL     string `json:"confirmUrl"`
	ConfirmURLType string `json:"confirmUrlType,omitempty"`
	CancelURL      string `json:"cancelUrl"`
}

// RequestOptionsPayment type
type RequestOptionsPayment struct {
	Capture *bool  `json:"capture,omitempty"`
	PayType string `json:"payType,omitempty"`
}

// RequestOptionsDisplay type
type RequestOptionsDisplay struct {
	Locale                 string `json:"locale,omitempty"`
	CheckConfirmURLBrowser *bool  `json:"checkConfirmUrlBrowser,omitempty"`
}

// RequestOptionsShipping type
type RequestOptionsShipping struct {
	Type           string `json:"type,omitempty"`
	FeeInquiryURL  string `json:"feeInquiryUrl,omitempty"`
	FeeInquiryType string `json:"feeInquiryType,omitempty"`
}

// RequestOptionsExtras type
type RequestOptionsExtras struct {
	FamilyService struct {
		AddFriends []struct {
			Type string   `json:"type,omitempty"`
			IDs  []string `json:"ids,omitempty"`
		} `json:"addFriends,omitempty"`
	} `json:"familyService,omitempty"`
	BranchName string `json:"branchName,omitempty"`
	BranchID   string `json:"branchId,omitempty"`
}

// RequestOptions type
type RequestOptions struct {
	Payment  *RequestOptionsPayment  `json:"payment,omitempty"`
	Display  *RequestOptionsDisplay  `json:"display,omitempty"`
	Shipping *RequestOptionsShipping `json:"shipping,omitempty"`
	Extras   *RequestOptionsExtras   `json:"extras,omitempty"`
}

// RequestRequest type
type RequestRequest struct {
	Amount       int                  `json:"amount"`
	Currency     string               `json:"currency"`
	OrderID      string               `json:"orderId"`
	Packages     []*RequestPackage    `json:"packages"`
	RedirectURLs *RequestRedirectURLs `json:"redirectUrls"`
	Options      *RequestOptions      `json:"options,omitempty"`
}

// RequestResponse type
type RequestResponse struct {
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
