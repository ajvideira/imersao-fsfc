package repository

import (
	"fmt"

	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/domain/model"
	"github.com/jinzhu/gorm"
)

/*type PixKeyRepositoryInterface interface {
	RegisterKey(pixKey *PixKey) (*PixKey, error)
	FindKeyByKind(kind string) (*PixKey, error)
	AddBank(bank *Bank) error
	AddAccount(account *Account) error
	FindAccount(id string) (*Account, error)
}*/

// PixKeyRepositoryDb represents a repository for pixKey operations
type PixKeyRepositoryDb struct {
	Db *gorm.DB
}

// AddBank adds a new bank to the database
func (repository *PixKeyRepositoryDb) AddBank(bank *model.Bank) error {
	err := repository.Db.Create(bank).Error

	if err != nil {
		return err
	}
	return nil
}

// AddAccount adds a new account to the database
func (repository *PixKeyRepositoryDb) AddAccount(account *model.Account) error {
	err := repository.Db.Create(account).Error

	if err != nil {
		return err
	}
	return nil 
}

// RegisterKey registers a new pixKey to the database
func (repository *PixKeyRepositoryDb) RegisterKey(pixKey *model.PixKey) (*model.PixKey, error) {
	err := repository.Db.Create(pixKey).Error

	if err != nil {
		return nil, err
	}
	return pixKey, nil 
}

//FindKeyByKind search for a key and kind
func (repository *PixKeyRepositoryDb) FindKeyByKind(key string, kind string) (*model.PixKey, error) {
	var pixKey model.PixKey;

	repository.Db.Preload("Account.Bank").First(&pixKey, "kind = ? AND key = ?", kind, key)

	if pixKey.ID == "" {
		return nil, fmt.Errorf("no key was found")
	}

	return &pixKey, nil
}

//FindAccount search for a account
func (repository *PixKeyRepositoryDb) FindAccount(id string) (*model.Account, error) {
	var account model.Account;

	repository.Db.Preload("Bank").First(&account, "id ?", id)

	if account.ID == "" {
		return nil, fmt.Errorf("no account was found")
	}

	return &account, nil
}

//FindBank search for a bank
func (repository *PixKeyRepositoryDb) FindBank(id string) (*model.Bank, error) {
	var bank model.Bank;

	repository.Db.First(&bank, "id ?", id)

	if bank.ID == "" {
		return nil, fmt.Errorf("no bank was found")
	}

	return &bank, nil
}