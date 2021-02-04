package usecase

import "github.com/ajvideira/imersao-fullstack-fullcycle/codepix/domain/model"

// PixKeyUseCase represents a use case for pix keys
type PixKeyUseCase struct {
	PixKeyRepository model.PixKeyRepositoryInterface
}

// RegisterKey registers a new key
func (useCase *PixKeyUseCase) RegisterKey(kind string, key string, accountID string) (*model.PixKey, error) {
	account, err := useCase.PixKeyRepository.FindAccount(accountID)
	if err != nil {
		return nil, err
	}

	pixKey, err := model.NewPixKey(kind, account, key)
	if err != nil {
		return nil, err
	}

	pixKey, err = useCase.PixKeyRepository.RegisterKey(pixKey)
	if (err != nil) {
		return nil, err
	}
	return pixKey, nil
}

// FindKey search a pix key by kind and key
func (useCase *PixKeyUseCase) FindKey(kind string, key string) (*model.PixKey, error) {
	pixKey, err := useCase.PixKeyRepository.FindKeyByKind(key, kind);
	if err != nil {
		return nil, err
	}
	return pixKey, nil;
}