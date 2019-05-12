package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"

	"github.com/gotokatsuya/line-pay-sdk-go/linepay"
)

var (
	cache = sync.Map{}
)

func main() {
	pay, err := linepay.New(
		os.Getenv("LINE_PAY_CHANNEL_ID"),
		os.Getenv("LINE_PAY_CHANNEL_SECRET"),
		linepay.WithSandbox(),
	)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reserveReq := &linepay.ReserveRequest{
			ProductName: "demo product",
			Amount:      1,
			Currency:    "JPY",
			ConfirmURL:  os.Getenv("LINE_PAY_CONFIRM_URL"),
			OrderID:     uuid.New().String(),
		}
		reserveResp, _, err := pay.Reserve(context.Background(), reserveReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		transactionID := reserveResp.Info.TransactionID
		cache.Store(transactionID, reserveReq)
		http.Redirect(w, r, reserveResp.Info.PaymentURL.Web, http.StatusFound)
	})
	http.HandleFunc("/confirm", func(w http.ResponseWriter, r *http.Request) {
		transactionID := r.URL.Query().Get("transactionId")
		reservationVal, ok := cache.Load(transactionID)
		if !ok {
			http.Error(w, "Reservation not found", http.StatusInternalServerError)
			return
		}
		reservation := reservationVal.(*linepay.ReserveRequest)
		confirmReq := &linepay.ConfirmRequest{
			Amount:   reservation.Amount,
			Currency: reservation.Currency,
		}
		confirmResp, _, err := pay.Confirm(context.Background(), transactionID, confirmReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(confirmResp)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
