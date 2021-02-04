package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// Bank represents a model bank
type Bank struct {
	Base `valid:"required"`
	Code string `json:"code" gorm:"column:code;type:varchar2(20);not null" valid:"notnull"`
	Name string `json:"name" gorm:"column:name;type:varchar2(255);not null" valid:"notnull"`
	Accounts []*Account `gorm:"ForeignKey:BankID" valid:"-"`
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

	bank.ID = uuid.NewV4().String()
	bank.CreatedAt = time.Now()
	bank.UpdatedAt = time.Now()

	err := bank.isValid()
	if (err != nil) {
		return nil, err;
	}
	
	return &bank, nil;
}