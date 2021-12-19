package flutterwave

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

type (
	FlutterGo interface {
		GetTransactionsByTransactionReference(txRef string) ([]TransactionData, error)
		VerifyTransaction(id int64) (*TransactionResponse, error)
		ChargeTokenizedCard(token, email, narration, txRef string, amount float64) (response *TransactionResponse, err error)
		ValidateAccountNumber(account, bank string) (*AccountValidationResponse, error)
		CreatePermanentVirtualAccount(email, firstName, lastName, bvn, txRef string) (response *VirtualAccountResponse, err error)
	}

	flutterGo struct {
		authKey string
		client  *http.Client
		apiUrl  string
	}

	Error struct {
		Code     int
		Body     string
		Endpoint string
	}

	Config struct {
		AuthKey string
		Client  *http.Client
		ApiUrl  string
	}

	header struct {
		Key   string
		Value string
	}
)

func New(cfg Config) (FlutterGo, error) {
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	g := &flutterGo{
		authKey: cfg.AuthKey,
		client:  cfg.Client,
		apiUrl:  cfg.ApiUrl,
	}

	return g, nil
}

func NewDefaultConfig(secretKey string) Config {
	return Config{
		AuthKey: secretKey,
		Client: &http.Client{
			Timeout: 20 * time.Second,
		},
		ApiUrl: os.Getenv("FW_API_URL"),
	}
}

func validateConfig(cfg *Config) error {
	if cfg.AuthKey == "" {
		return errors.New("fluttergo: Missing Auth Key")
	}

	if cfg.Client == nil {
		return errors.New("fluttergo: HTTP Client Cannot Be Nil")
	}

	if cfg.ApiUrl == "" {
		return errors.New("fluttergo: Missing API Url")
	}

	return nil
}

func (fg *flutterGo) preparePayload(body interface{}) (io.Reader, error) {
	b, err := json.Marshal(body)
	log.Printf("\nPayload Bodhy string: %v", string(b)) //debug
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (fg *flutterGo) makeRequest(method, url string, body io.Reader, headers []header, responseTarget interface{}) error {
	if reflect.TypeOf(responseTarget).Kind() != reflect.Ptr {
		return errors.New("fluttergo: responseTarget must be a pointer to a struct for JSON unmarshalling")
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	for _, h := range headers {
		req.Header.Set(h.Key, h.Value)
	}
	req.Header.Set("Authorization", "Bearer "+fg.authKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := fg.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, responseTarget)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		return nil
	}

	err = Error{
		Code:     resp.StatusCode,
		Body:     string(b),
		Endpoint: req.URL.String(),
	}
	return err
}

func (e Error) Error() string {
	return fmt.Sprintf("Request To %v Endpoint Failed With Status Code %v | Body: %v", e.Endpoint, e.Code, e.Body)
}
