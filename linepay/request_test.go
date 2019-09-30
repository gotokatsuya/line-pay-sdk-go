package linepay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestClient_Request(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	req := &RequestRequest{
		Amount:   100,
		Currency: "JPY",
		OrderID:  "MKSI_S_20180904_1000001",
		Packages: []*RequestPackage{
			&RequestPackage{
				ID:     "1",
				Amount: 100,
				Name:   "PACKAGE_1",
				Products: []*RequestPackageProduct{
					&RequestPackageProduct{
						ID:       "PEN-B-001",
						Name:     "Pen Brown",
						ImageURL: "https://pay-store.line.com/images/pen_brown.jpg",
						Quantity: 1,
						Price:    100,
					},
				},
			},
		},
		RedirectURLs: &RequestRedirectURLs{
			ConfirmURL: "https://pay-store.line.com/order/payment/authorize",
			CancelURL:  "https://pay-store.line.com/order/payment/cancel",
		},
	}

	mux.HandleFunc("/v3/payments/request", func(w http.ResponseWriter, r *http.Request) {
		v := new(RequestRequest)
		json.NewDecoder(r.Body).Decode(v)
		if got := r.Method; got != http.MethodPost {
			t.Errorf("Request method: %v, want %v", got, http.MethodPost)
		}
		if !reflect.DeepEqual(v, req) {
			t.Errorf("Request body = %+v, want %+v", v, req)
		}
		fmt.Fprint(w, `{"returnCode":"0000"}`)
	})

	resp, _, err := client.Request(context.Background(), req)
	if err != nil {
		t.Errorf("Request returned error: %v", err)
	}

	want := &RequestResponse{ReturnCode: "0000"}
	if !reflect.DeepEqual(resp, want) {
		t.Errorf("Request returned %+v, want %+v", resp, want)
	}
}
