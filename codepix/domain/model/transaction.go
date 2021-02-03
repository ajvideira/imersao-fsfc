package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	// TransactionPending represents a pending transaction
	TransactionPending string = "pending"
	// TransactionConfirmed represents a confirmed transaction
	TransactionConfirmed string = "confirmed"
	// TransactionCompleted represents a completed transaction
	TransactionCompleted string = "completed"
	// TransactionError represents a error transaction
	TransactionError string = "error";
)

// TransactionRepositoryInterface represents a interface of all operations
type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction,error)
}

// Transactions represents a list of transactions
type Transactions struct {
	Transaction []*Transaction
}

// Transaction represents a model transaction
type Transaction struct {
	Base `valid:"required"`
	AccountFrom *Account `valid:"-"`
	Amount float64 `json:"amount" valid:"notnull"`
	PixKeyTo *PixKey `valid:"-"`
	Status string `json:"status" valid:"notnull"`
	Description string `json:"description" valid:"notnull"`
	CancelDescription string `json:"cancel_description" valid:"-"`
}

//isValid perform validation of a pix transaction
func (transaction *Transaction) isValid() error {
	_,err := govalidator.ValidateStruct(transaction);

	if transaction.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if transaction.Status != TransactionPending && transaction.Status != TransactionConfirmed && transaction.Status != TransactionCompleted && transaction.Status != TransactionError {
		return errors.New("invalid status")
	}

	if transaction.PixKeyTo.AccountID == transaction.AccountFrom.ID {
		return errors.New("source and destination account cannot be the same")
	}

	if err != nil {
		return err;
	}
	return nil;
}

// NewTransaction return a new instance of a Transaction
func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom: accountFrom,
		Amount: amount,
		PixKeyTo: pixKeyTo,
		Status: TransactionPending,
		Description: description,
	}

	err := transaction.isValid()
	if (err != nil) {
		return nil, err;
	}

	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	return &transaction, nil;
}

// Complete completes a transaction
func (transaction *Transaction) Complete() error {
	transaction.Status = TransactionCompleted
	transaction.UpdatedAt = time.Now()
	err := transaction.isValid()
	return err
}
// Confirm confirm a transaction
func (transaction *Transaction) Confirm() error {
	transaction.Status = TransactionConfirmed
	transaction.UpdatedAt = time.Now()
	err := transaction.isValid()
	return err
}

// Cancel cancels a transaction
func (transaction *Transaction) Cancel(description string) error {
	transaction.Status = TransactionError
	transaction.UpdatedAt = time.Now()
	transaction.CancelDescription = description
	err := transaction.isValid()
	return err
}