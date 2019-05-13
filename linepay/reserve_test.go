package linepay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestClient_Reserve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	req := &ReserveRequest{
		ProductName: "100yen",
		Amount:      100,
		Currency:    "JPY",
		ConfirmURL:  "http://localhost:5000/pay/confirm",
		OrderID:     "test-order",
	}

	mux.HandleFunc("/v2/payments/request", func(w http.ResponseWriter, r *http.Request) {
		v := new(ReserveRequest)
		json.NewDecoder(r.Body).Decode(v)
		if got := r.Method; got != http.MethodPost {
			t.Errorf("Request method: %v, want %v", got, http.MethodPost)
		}
		if !reflect.DeepEqual(v, req) {
			t.Errorf("Request body = %+v, want %+v", v, req)
		}
		fmt.Fprint(w, `{"returnCode":"0000"}`)
	})

	resp, _, err := client.Reserve(context.Background(), req)
	if err != nil {
		t.Errorf("Reserve returned error: %v", err)
	}

	want := &ReserveResponse{ReturnCode: "0000"}
	if !reflect.DeepEqual(resp, want) {
		t.Errorf("Reserve returned %+v, want %+v", resp, want)
	}
}
