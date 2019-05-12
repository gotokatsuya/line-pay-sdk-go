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
	mux.HandleFunc("/v2/payments/authorizations", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"returnCode":"0000"}`)
	})

	req := &InquireAuthorizationRequest{}
	resp, _, err := client.InquireAuthorization(context.Background(), req)
	if err != nil {
		t.Errorf("InquireAuthorization returned error: %v", err)
	}

	want := &InquireAuthorizationResponse{ReturnCode: "0000"}
	if !reflect.DeepEqual(resp, want) {
		t.Errorf("InquireAuthorization returned %+v, want %+v", resp, want)
	}
}
