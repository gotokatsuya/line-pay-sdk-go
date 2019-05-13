package linepay

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestClient_InquireAuthorization(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	req := &InquireAuthorizationRequest{
		TransactionID: 100000,
	}

	mux.HandleFunc("/v2/payments/authorizations", func(w http.ResponseWriter, r *http.Request) {
		v := &InquireAuthorizationRequest{
			TransactionID: MustParseInt64(r.URL.Query().Get("transactionId")),
			OrderID:       r.URL.Query().Get("orderId"),
		}
		if got := r.Method; got != http.MethodGet {
			t.Errorf("Request method: %v, want %v", got, http.MethodGet)
		}
		if !reflect.DeepEqual(v, req) {
			t.Errorf("Request url params = %+v, want %+v", v, req)
		}
		fmt.Fprint(w, `{"returnCode":"0000"}`)
	})

	resp, _, err := client.InquireAuthorization(context.Background(), req)
	if err != nil {
		t.Errorf("InquireAuthorization returned error: %v", err)
	}

	want := &InquireAuthorizationResponse{ReturnCode: "0000"}
	if !reflect.DeepEqual(resp, want) {
		t.Errorf("InquireAuthorization returned %+v, want %+v", resp, want)
	}
}
