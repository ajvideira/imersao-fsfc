package model

import (
	"encoding/json"
	"fmt"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// Transaction represents a external entity transaction
type Transaction struct {
	ID           string  `json:"id" valid:"notnull,uuid"`
	AccountID    string  `json:"accountId" valid:"notnull,uuid"`
	Amount       float64 `json:"amount" valid:"notnull,numeric"`
	PixKeyTo     string  `json:"pixKeyTo" valid:"notnull"`
	PixKeyKindTo string  `json:"pixKeyKindTo" valid:"notnull"`
	Description  string  `json:"description" valid:"-"`
	Status       string  `json:"status" valid:"notnull"`
	Error        string  `json:"error" valid:"-"`
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
