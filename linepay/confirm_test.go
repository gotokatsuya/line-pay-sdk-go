package linepay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestClient_Confirm(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	req := &ConfirmRequest{
		Amount:   100,
		Currency: "JPY",
	}

	mux.HandleFunc("/v3/payments/20190513000000/confirm", func(w http.ResponseWriter, r *http.Request) {
		v := new(ConfirmRequest)
		json.NewDecoder(r.Body).Decode(v)
		if got := r.Method; got != http.MethodPost {
			t.Errorf("Request method: %v, want %v", got, http.MethodPost)
		}
		if !reflect.DeepEqual(v, req) {
			t.Errorf("Request body = %+v, want %+v", v, req)
		}
		fmt.Fprint(w, `{"returnCode":"0000"}`)
	})

	resp, _, err := client.Confirm(context.Background(), 20190513000000, req)
	if err != nil {
		t.Errorf("Confirm returned error: %v", err)
	}

	want := &ConfirmResponse{ReturnCode: "0000"}
	if !reflect.DeepEqual(resp, want) {
		t.Errorf("Confirm returned %+v, want %+v", resp, want)
	}
}
