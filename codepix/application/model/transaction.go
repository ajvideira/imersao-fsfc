package model

import (
	"encoding/json"
	"fmt"

	"github.com/asaskevich/govalidator"
)

// Transaction represents a external entity transaction
type Transaction struct {
	ID           string  `json:"id" validate:"required,uuid"`
	AccountID    string  `json:"accountId" validate:"required,uuid"`
	Amount       float64 `json:"amount" validate:"required,numeric"`
	PixKeyTo     string  `json:"pixKeyTo" validate:"required"`
	PixKeyKindTo string  `json:"pixKeyKindTo" validate:"required"`
	Description  string  `json:"description" validate:"required"`
	Status       string  `json:"status" validate:"required"`
	Error        string  `json:"error"`
}

//isValid perform validation of a pix transaction
func (transaction *Transaction) isValid() error {
	_,err := govalidator.ValidateStruct(transaction);
	if err != nil {
		fmt.Println("error during  transaction validation", err)
		return err;
	}
	return nil;
}
// ParseJSON parses json into transaction struct
func (transaction *Transaction) ParseJSON(data []byte) error {
	err := json.Unmarshal(data, transaction)
	if err != nil {
		return err
	}

	err = transaction.isValid()
	if err != nil {
		return err
	}

	return nil
}

// ToJSON converts a transaction struct to json
func (transaction *Transaction) ToJSON() ([]byte, error) {
	err := transaction.isValid()
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}
	return result, err
}
