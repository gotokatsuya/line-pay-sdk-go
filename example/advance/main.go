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

// UserSession type
type UserSession struct {
	RegKey string `json:"regKey"`
}

func init() {
	gob.Register(&PaymentTransactionSession{})
	gob.Register(&UserSession{})
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
			ProductName: "advance demo product",
			Amount:      1,
			Currency:    "JPY",
			ConfirmURL:  os.Getenv("LINE_PAY_CONFIRM_URL"),
			OrderID:     uuid.New().String(),
			PayType:     "PREAPPROVED",
		}
		reserveResp, _, err := pay.Reserve(context.Background(), reserveReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		transactionID := reserveResp.Info.TransactionID
		paymentSession, err := store.Get(r, "payment-transaction")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		paymentSession.Values[transactionID] = &PaymentTransactionSession{
			TransactionID: transactionID,
			OrderID:       reserveReq.OrderID,
			ProductName:   reserveReq.ProductName,
			Amount:        reserveReq.Amount,
			Currency:      reserveReq.Currency,
		}
		if err := paymentSession.Save(r, w); err != nil {
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
		paymentSession, err := store.Get(r, "payment-transaction")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		transactionVal, ok := paymentSession.Values[transactionID]
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
		paymentSession.Values["user"] = &UserSession{
			RegKey: confirmResp.Info.RegKey,
		}
		if err := paymentSession.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(confirmResp)
	})
	http.HandleFunc("/pay/regKey", func(w http.ResponseWriter, r *http.Request) {
		paymentSession, err := store.Get(r, "payment-transaction")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userVal, ok := paymentSession.Values["user"]
		if !ok {
			http.Error(w, "UserTransaction is not found", http.StatusInternalServerError)
			return
		}
		user := userVal.(*UserSession)
		confirmResp, _, err := pay.ConfirmPreapprovedPay(
			context.Background(),
			user.RegKey,
			&linepay.ConfirmPreapprovedPayRequest{
				ProductName: "advance demo product",
				Amount:      1,
				Currency:    "JPY",
				OrderID:     uuid.New().String(),
			})
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
