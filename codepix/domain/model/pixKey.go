package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

// PixKeyRepositoryInterface represents a interface of all operations
type PixKeyRepositoryInterface interface {
	RegisterKey(pixKey *PixKey) (*PixKey, error)
	FindKeyByKind(key string, kind string) (*PixKey, error)
	AddBank(bank *Bank) error
	AddAccount(account *Account) error
	FindAccount(id string) (*Account, error)
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// PixKey represents a model pix key
type PixKey struct {
	Base `valid:"required"`
	Kind string `json:"kind" gorm:"column:kind;type:varchar(20)" valid:"notnull"`
	Key string `json:"key" gorm:"column:key;type:varchar(255)" valid:"notnull"`
	AccountID string `json:"account_id" gorm:"column:account_id;type:uuid;not null" valid:"notnull"`
	Account *Account `valid:"-"`
	Status string `json:"status" gorm:"type:varchar(20)" valid:"notnull"`
}

//isValid perform validation of a pix key
func (pixKey *PixKey) isValid() error {
	_,err := govalidator.ValidateStruct(pixKey);

	if pixKey.Kind != "email" && pixKey.Kind != "cpf" {
		return errors.New("invalid type of key")
	}

	if pixKey.Status != "active" && pixKey.Status != "inactive" {
		return errors.New("invalid status")
	}

	if err != nil {
		return err;
	}
	return nil;
}

// NewPixKey return a new instance of a PixKey
func NewPixKey(kind string, account *Account, key string) (*PixKey, error) {
	pixKey := PixKey{
		Kind: kind,
		Account: account,
		AccountID: account.ID,
		Key: key,
		Status: "active",
	}

	pixKey.ID = uuid.NewV4().String()
	pixKey.CreatedAt = time.Now()
	pixKey.UpdatedAt = time.Now()

	err := pixKey.isValid()
	if (err != nil) {
		return nil, err
	}

	return &pixKey, nil
}