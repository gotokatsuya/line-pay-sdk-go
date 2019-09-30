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
	http.HandleFunc("/pay/request", func(w http.ResponseWriter, r *http.Request) {
		reserveReq := &linepay.RequestRequest{
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
							ID:       "PRIME-M-001",
							Name:     "Prime MemberShip",
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
			Options: &linepay.RequestOptions{
				Payment: &linepay.RequestOptionsPayment{
					PayType: "PREAPPROVED",
				},
			},
		}
		requestResp, _, err := pay.Request(context.Background(), reserveReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if requestResp.ReturnCode != "0000" {
			http.Error(w, requestResp.ReturnMessage, http.StatusInternalServerError)
			return
		}
		transactionID := requestResp.Info.TransactionID
		paymentSession, err := store.Get(r, "payment-transaction")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		paymentSession.Values[transactionID] = &PaymentTransactionSession{
			TransactionID: transactionID,
			OrderID:       reserveReq.OrderID,
			Amount:        reserveReq.Amount,
			Currency:      reserveReq.Currency,
		}
		if err := paymentSession.Save(r, w); err != nil {
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
		if confirmResp.ReturnCode != "0000" {
			http.Error(w, confirmResp.ReturnMessage, http.StatusInternalServerError)
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
		payPreapprovedResp, _, err := pay.PayPreapproved(
			context.Background(),
			user.RegKey,
			&linepay.PayPreapprovedRequest{
				ProductName: "Prime MemberShip",
				Amount:      250,
				Currency:    "JPY",
				OrderID:     uuid.New().String(),
			})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if payPreapprovedResp.ReturnCode != "0000" {
			http.Error(w, payPreapprovedResp.ReturnMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payPreapprovedResp)
	})
	fmt.Println("open http://localhost:8080/pay/request")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
