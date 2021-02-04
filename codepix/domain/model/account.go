package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// Account represents a model account
type Account struct {
	Base `valid:"required"`
	OwnerName string `json:"owner_name" gorm:"column:owner_name;type:varchar(255);not null" valid:"notnull"`
	BankID string `gorm:"column:bank_id;type:uuid;not null" valid:"-"`
	Bank *Bank `valid:"-"`
	Number string `json:"number" gorm:"column:number;type:varchar(20);not null" valid:"notnull"`
	PixKeys []*PixKey `gorm:"ForeignKey:AccountID" valid:"-"`
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
		BankID: bank.ID,
		Number: number,
	}

	account.ID = uuid.NewV4().String()
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	err := account.isValid()
	if (err != nil) {
		return nil, err
	}

	return &account, nil
}