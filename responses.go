package flutterwave

import (
	"errors"
	"fmt"

	"log"
	"time"
)

type (
	AccountValidationResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			AccountNumber string `json:"account_number"`
			AccountName   string `json:"account_name"`
		} `json:"data"`
	}

	TransactionResponse struct {
		Status  string          `json:"status"`
		Message string          `json:"message"`
		Data    TransactionData `json:"data"`
	}

	AllTransactionsResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Meta    struct {
			PageInfo struct {
				Total       int `json:"total"`
				CurrentPage int `json:"current_page"`
				TotalPages  int `json:"total_pages"`
			} `json:"page_info"`
		} `json:"meta"`
		Data []TransactionData `json:"data"`
	}

	TransactionData struct {
		Id                int64   `json:"id"`
		TxRef             string  `json:"tx_ref"`
		FlwRef            string  `json:"flw_ref"`
		DeviceFingerprint string  `json:"device_fingerprint"`
		Amount            float64 `json:"amount"`
		Currency          string  `json:"currency"`
		ChargedAmount     float64 `json:"charged_amount"`
		AppFee            float64 `json:"app_fee"`
		MerchantFee       int64   `json:"merchant_fee"`
		ProcessorResponse string  `json:"processor_response"`
		AuthModel         string  `json:"auth_model"`
		Ip                string  `json:"ip"`
		Narration         string  `json:"narration"`
		Status            string  `json:"status"`
		PaymentType       string  `json:"payment_type"`
		CreatedAt         string  `json:"created_at"`
		AccountId         int64   `json:"account_id"`
		Card              struct {
			First6Digits string `json:"first_6digits"`
			Last4Digits  string `json:"last_4digits"`
			Issuer       string `json:"issuer"`
			Country      string `json:"country"`
			Type         string `json:"type"`
			Token        string `json:"token"`
			Expiry       string `json:"expiry"`
		} `json:"card"`
		Meta          interface{} `json:"meta"`
		AmountSettled float64     `json:"amount_settled"`
		Customer      struct {
			Id          int64  `json:"id"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
			Email       string `json:"email"`
			CreatedAt   string `json:"created_at"`
		} `json:"customer"`
	}

	VirtualAccountResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			ResponseCode    string `json:"response_code"`
			ResponseMessage string `json:"response_message"`
			OrderRef        string `json:"order_ref"`
			AccountNumber   string `json:"account_number"`
			AccountStatus   string `json:"account_status"`
			Frequency       string `json:"frequency"`
			BankName        string `json:"bank_name"`
			CreatedAt       string `json:"created_at"`
			ExpiryDate      string `json:"expiry_date"`
			Amount          string `json:"amount"`
		} `json:"data"`
	}

	WebhookMessage struct {
		Event string `json:"event"`
		Data  struct {
			Id                int64   `json:"id"`
			TxRef             string  `json:"tx_ref"`
			FlwRef            string  `json:"flw_ref"`
			DeviceFingerprint string  `json:"device_fingerprint"`
			Amount            float64 `json:"amount"`
			Currency          string  `json:"currency"`
			ChargedAmount     float64 `json:"charged_amount"`
			AppFee            float64 `json:"app_fee"`
			MerchantFee       int64   `json:"merchant_fee"`
			ProcessorResponse string  `json:"processor_response"`
			AuthModel         string  `json:"auth_model"`
			Ip                string  `json:"ip"`
			Narration         string  `json:"narration"`
			Status            string  `json:"status"`
			PaymentType       string  `json:"payment_type"`
			CreatedAt         string  `json:"created_at"`
			AccountId         int64   `json:"account_id"`
			Card              struct {
				First6Digits string `json:"first_6digits"`
				Last4Digits  string `json:"last_4digits"`
				Issuer       string `json:"issuer"`
				Country      string `json:"country"`
				Type         string `json:"type"`
				Token        string `json:"token"`
				Expiry       string `json:"expiry"`
			} `json:"card"`
			Meta          interface{} `json:"meta"`
			AmountSettled float64     `json:"amount_settled"`
			Customer      struct {
				Id          int64  `json:"id"`
				Name        string `json:"name"`
				PhoneNumber string `json:"phone_number"`
				Email       string `json:"email"`
				CreatedAt   string `json:"created_at"`
			} `json:"customer"`
		} `json:"data"`
	}

)

func (wh WebhookMessage) Verify(fwkey, fwUrl string) (*TransactionResponse, error) {
	cfg := NewDefaultConfig(fwkey, fwUrl)
	fg, err := New(cfg)

	if err != nil {
		return nil, err
	}

	tr, err := fg.VerifyTransaction(wh.Data.Id)
	if err != nil {
		return nil, err
	}

	if tr.Data.Status == "successful" {
		return tr, nil
	}

	return nil, errors.New(fmt.Sprintf("transaction status %v", tr.Data.Status))
}

func (tr TransactionResponse) IsTokenised() bool {
	return tr.Data.Card.Token != ""
}

func (tr TransactionResponse) CardWillExpireBeforeDate(date time.Time) (bool, string) {
	cardExpiryDate, err := time.Parse("01/06", tr.Data.Card.Expiry)
	if err != nil {
		e := fmt.Sprintf("couldn't parse card expiry date: '%v'", tr.Data.Card.Expiry)
		log.Print(e)
		return true, e
	}

	return cardExpiryDate.Before(date), tr.Data.Card.Expiry
}

func (tr TransactionData) IsTokenised() bool {
	return tr.Card.Token != ""
}


