package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

// Account represents a model account
type Account struct {
	Base `valid:"required"`
	OwnerName string `json:"owner_name" valid:"notnull"`
	Bank *Bank `valid:"-"`
	Number string `json:"number" valid:"notnull"`
	PixKeys []*PixKey `valid:"-"`
}

//isValid perform validation of a account
func (account *Account) isValid() error {
	_,err := govalidator.ValidateStruct(account);
	if err != nil {
		return err;
	}
	return nil;
}

// NewAccount return a new instance of a Account
func NewAccount(bank *Bank, number string, ownerName string) (*Account, error) {
	account := Account{
		OwnerName: ownerName,
		Bank: bank,
		Number: number,
	}

	err := account.isValid()
	if (err != nil) {
		return nil, err;
	}

	account.ID = uuid.NewV4().String()
	account.CreatedAt = time.Now()

	return &account, nil;
}