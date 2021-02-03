package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

// Bank represents a model bank
type Bank struct {
	Base `valid:"required"`
	Code string `json:"code" valid:"notnull"`
	Name string `json:"name" valid:"notnull"`
	Accounts []*Account `valid:"-"`
}

//isValid perform validation of a bank
func (bank *Bank) isValid() error {
	_,err := govalidator.ValidateStruct(bank);
	if err != nil {
		return err;
	}
	return nil;
}

// NewBank return a new instance of a Bank
func NewBank(code string, name string) (*Bank, error) {
	bank := Bank{
		Code: code,
		Name: name,
	}

	err := bank.isValid()
	if (err != nil) {
		return nil, err;
	}

	bank.ID = uuid.NewV4().String()
	bank.CreatedAt = time.Now()

	return &bank, nil;
}