package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"

	"github.com/gotokatsuya/line-pay-sdk-go/linepay"
)

var (
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
)

// PaymentTransactionSession type
type PaymentTransactionSession struct {
	TransactionID int64  `json:"transactionId"`
	OrderID       string `json:"orderId"`
	ProductName   string `json:"productName"`
	Amount        int    `json:"amount"`
	Currency      string `json:"currency"`
}

func init() {
	gob.Register(&PaymentTransactionSession{})
}

func main() {
	pay, err := linepay.New(
		os.Getenv("LINE_PAY_CHANNEL_ID"),
		os.Getenv("LINE_PAY_CHANNEL_SECRET"),
		linepay.WithSandbox(),
	)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/pay/reserve", func(w http.ResponseWriter, r *http.Request) {
		reserveReq := &linepay.ReserveRequest{
			ProductName: "basic demo product",
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
		session, err := store.Get(r, "payment-transaction")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values[transactionID] = &PaymentTransactionSession{
			TransactionID: transactionID,
			OrderID:       reserveReq.OrderID,
			ProductName:   reserveReq.ProductName,
			Amount:        reserveReq.Amount,
			Currency:      reserveReq.Currency,
		}
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, reserveResp.Info.PaymentURL.Web, http.StatusFound)
	})
	http.HandleFunc("/pay/confirm", func(w http.ResponseWriter, r *http.Request) {
		transactionID, err := linepay.ParseInt64(r.URL.Query().Get("transactionId"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session, err := store.Get(r, "payment-transaction")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		transactionVal, ok := session.Values[transactionID]
		if !ok {
			http.Error(w, fmt.Sprintf("PaymentTransaction is not found. id:%d", transactionID), http.StatusInternalServerError)
			return
		}
		transaction := transactionVal.(*PaymentTransactionSession)
		confirmReq := &linepay.ConfirmRequest{
			Amount:   transaction.Amount,
			Currency: transaction.Currency,
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
	fmt.Println("open http://localhost:8080/pay/reserve")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
