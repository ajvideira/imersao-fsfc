package factory

import (
	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/application/usecase"
	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/infrastructure/repository"
	"github.com/jinzhu/gorm"
)

// TransactionUseCaseFactory returns a transaction use case
func TransactionUseCaseFactory(database *gorm.DB) usecase.TransactionUseCase {
	pixRepository := repository.PixKeyRepositoryDb{Db: database}
	transactionRepository := repository.TransactionRepositoryDb{Db: database}

	transactionUseCase := usecase.TransactionUseCase{
		PixKeyRepository: &pixRepository,
		TransactionRepository: &transactionRepository,
	}
	return transactionUseCase
}