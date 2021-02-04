package usecase

import "github.com/ajvideira/imersao-fullstack-fullcycle/codepix/domain/model"

// TransactionUseCase represents a use case for transactions
type TransactionUseCase struct {
	PixKeyRepository model.PixKeyRepositoryInterface
	TransactionRepository model.TransactionRepositoryInterface
}

// Register registers a new transaction
func (useCase *TransactionUseCase) Register(accountID string, amount float64, pixKeyTo string, pixKeyKindTo string, description string) (*model.Transaction, error) {
	account, err := useCase.PixKeyRepository.FindAccount(accountID)
	if err != nil {
		return nil, err
	}

	pixKey, err := useCase.PixKeyRepository.FindKeyByKind(pixKeyTo, pixKeyKindTo);
	if err != nil {
		return nil, err
	}

	transaction, err := model.NewTransaction(account, amount, pixKey, description)
	if (err != nil) {
		return nil, err
	}
	
	err = useCase.TransactionRepository.Register(transaction)
	if (err != nil) {
		return nil, err
	}
	return transaction, nil
}

// Confirm confirms a transaction
func (useCase *TransactionUseCase) Confirm(transactionID string) (*model.Transaction, error) {
	transaction, err := useCase.TransactionRepository.Find(transactionID)
	if err != nil {
		return nil, err
	}

	transaction.Status = model.TransactionConfirmed
	err = useCase.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// Complete completes a transaction
func (useCase *TransactionUseCase) Complete(transactionID string) (*model.Transaction, error) {
	transaction, err := useCase.TransactionRepository.Find(transactionID)
	if err != nil {
		return nil, err
	}

	transaction.Status = model.TransactionCompleted
	err = useCase.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// Error cancels a transaction
func (useCase *TransactionUseCase) Error(transactionID string, cancelDescription string) (*model.Transaction, error) {
	transaction, err := useCase.TransactionRepository.Find(transactionID)
	if err != nil {
		return nil, err
	}

	transaction.Status = model.TransactionError
	transaction.CancelDescription = cancelDescription
	err = useCase.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}