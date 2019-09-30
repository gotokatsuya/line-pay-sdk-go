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
	http.HandleFunc("/pay/request", func(w http.ResponseWriter, r *http.Request) {
		requestReq := &linepay.RequestRequest{
			Amount:   250,
			Currency: "JPY",
			OrderID:  uuid.New().String(),
			Packages: []*linepay.RequestPackage{
				&linepay.RequestPackage{
					ID:     "1",
					Amount: 250,
					Name:   "PACKAGE_SHOP_1",
					Products: []*linepay.RequestPackageProduct{
						&linepay.RequestPackageProduct{
							ID:       "PEN-B-001",
							Name:     "Pen Brown",
							Quantity: 1,
							Price:    250,
						},
					},
				},
			},
			RedirectURLs: &linepay.RequestRedirectURLs{
				ConfirmURL: os.Getenv("LINE_PAY_CONFIRM_URL"),
				CancelURL:  os.Getenv("LINE_PAY_CANCEL_URL"),
			},
		}
		requestResp, _, err := pay.Request(context.Background(), requestReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if requestResp.ReturnCode != "0000" {
			http.Error(w, requestResp.ReturnMessage, http.StatusInternalServerError)
			return
		}
		transactionID := requestResp.Info.TransactionID
		session, err := store.Get(r, "payment-transaction")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values[transactionID] = &PaymentTransactionSession{
			TransactionID: transactionID,
			OrderID:       requestReq.OrderID,
			Amount:        requestReq.Amount,
			Currency:      requestReq.Currency,
		}
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, requestResp.Info.PaymentURL.Web, http.StatusFound)
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
		if confirmResp.ReturnCode != "0000" {
			http.Error(w, confirmResp.ReturnMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(confirmResp)
	})
	fmt.Println("open http://localhost:8080/pay/request")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
