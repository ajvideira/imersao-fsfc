package repository

import (
	"fmt"
	"time"

	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/domain/model"
	"github.com/jinzhu/gorm"
)

/*type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction,error)
}*/

// TransactionRepositoryDb represents a repository for transaction operations
type TransactionRepositoryDb struct {
	Db *gorm.DB
}

// Register registers a new transaction
func (repository *TransactionRepositoryDb) Register(transaction *model.Transaction) error {
	err := repository.Db.Create(transaction).Error

	if err != nil {
		return err
	}

	return nil
}

// Save saves a existing transaction
func (repository *TransactionRepositoryDb) Save(transaction *model.Transaction) error {
	transaction.UpdatedAt = time.Now()
	err := repository.Db.Save(transaction).Error

	if err != nil {
		return err
	}

	return nil
}

// Find search for a transaction
func (repository *TransactionRepositoryDb) Find(id string) (*model.Transaction, error) {
	var transaction model.Transaction

	repository.Db.Preload("AccountFrom.Bank").First(&transaction, "id = ?", id)

	if transaction.ID == "" {
		return nil, fmt.Errorf("transaction not found")
	}
	return &transaction, nil
}
