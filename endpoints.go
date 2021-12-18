package flutterwave

import (
	"errors"
	"fmt"
	"log"
)

func (fg *flutterGo) VerifyTransaction(id int64) (*TransactionResponse, error) {

	log.Printf("\n\n verifying transaction %v", id)
	if id == 0 {
		return nil, errors.New("fluttergo: ID is required")
	}

	var respTarget TransactionResponse
	err := fg.makeRequest("GET", fmt.Sprintf("%v/transactions/%v/verify", fg.apiUrl, id), nil, nil, &respTarget)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n\n verify transaction response: \n %+v", respTarget)

	return &respTarget, nil
}

func (fg *flutterGo) GetTransactionsByTransactionReference(txRef string) ([]TransactionData, error) {

	log.Printf("\n\n Geting transaction %v", txRef)
	if txRef == "" {
		return nil, errors.New("fluttergo: txRef is required")
	}

	var respTarget AllTransactionsResponse
	err := fg.makeRequest("GET", fmt.Sprintf("%v/transactions?tx_ref=%v", fg.apiUrl, txRef), nil, nil, &respTarget)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n\n get transactions response: \n %+v", respTarget)

	return respTarget.Data, nil
}

func (fg *flutterGo) ValidateAccountNumber(account, bank string) (*AccountValidationResponse, error) {

	log.Printf("\n\n verifying account #%v, bank code %v", account, bank)

	var respTarget AccountValidationResponse

	body := map[string]interface{}{
		"account_number": account,
		"account_bank":   bank,
	}

	payload, err := fg.preparePayload(body)
	if err != nil {
		return nil, err
	}

	err = fg.makeRequest(
		"POST",
		fmt.Sprintf("%v/accounts/resolve", fg.apiUrl),
		payload, nil, &respTarget)
	if err != nil {
		return &respTarget, err
	}

	fmt.Printf("\n\n verify acct number response: \n %+v", respTarget)

	return &respTarget, nil
}

func (fg *flutterGo) ChargeTokenizedCard(token, email, narration, txRef string, amount float64) (response *TransactionResponse, err error) {

	if token == "" {
		return nil, errors.New("fluttergo: token is required")
	}
	if email == "" {
		return nil, errors.New("fluttergo: email is required")
	}
	if narration == "" {
		return nil, errors.New("fluttergo: narration is required")
	}
	if txRef == "" {
		return nil, errors.New("fluttergo: txRef is required")
	}
	if amount == 0 {
		return nil, errors.New("fluttergo: amount is required")
	}

	body := map[string]interface{}{
		"token":     token,
		"currency":  "NGN",
		"country":   "NG",
		"amount":    amount,
		"email":     email,
		"narration": narration,
		"tx_ref":    txRef,
	}

	payload, err := fg.preparePayload(body)
	if err != nil {
		return nil, err
	}

	var respTarget TransactionResponse
	err = fg.makeRequest("POST", fmt.Sprintf("%v/tokenized-charges", fg.apiUrl), payload, nil, &respTarget)
	if err != nil {
		return &respTarget, err
	}

	return &respTarget, nil
}

func (fg *flutterGo) CreatePermanentVirtualAccount(email, firstName, lastName, bvn, txRef string) (response *VirtualAccountResponse, err error) {

	if email == "" {
		return nil, errors.New("fluttergo: email is required")
	}
	if firstName == "" {
		return nil, errors.New("fluttergo: firstName is required")
	}
	if lastName == "" {
		return nil, errors.New("fluttergo: lastName is required")
	}
	if bvn == "" {
		return nil, errors.New("fluttergo: bvn is required")
	}
	if txRef == "" {
		return nil, errors.New("fluttergo: txRef is required")
	}

	body := map[string]interface{}{
		"email": email,
		"bvn": bvn,
		"tx_ref": txRef,
		"firstname": firstName,
		"lastname": lastName,
		"is_permanent": true,
		"narration": fmt.Sprintf("%v %v",firstName, lastName),
	}

	payload, err := fg.preparePayload(body)
	if err != nil {
		return nil, err
	}

	var respTarget VirtualAccountResponse
	err = fg.makeRequest("POST", fmt.Sprintf("%v/virtual-account-numbers", fg.apiUrl), payload, nil, &respTarget)
	if err != nil {
		return &respTarget, err
	}

	return &respTarget, nil
}
