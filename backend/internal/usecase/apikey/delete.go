package apikey

import (
	"context"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	apikeyvo "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/model/valueobject"
	apikeyrepo "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/repository"
	"github.com/jimi024rion/todo-go/backend/internal/domain/tx"
)

type DeleteUseCase struct {
	apiKeyRepo apikeyrepo.APIKeyRepository
	txManager  tx.TxManager
}

func NewDeleteUseCase(apiKeyRepo apikeyrepo.APIKeyRepository, txManager tx.TxManager) *DeleteUseCase {
	return &DeleteUseCase{
		apiKeyRepo: apiKeyRepo,
		txManager:  txManager,
	}
}

type DeleteInput struct {
	ID string
}

type DeleteOutput struct{}

func (uc *DeleteUseCase) Execute(ctx context.Context, input *DeleteInput) (*DeleteOutput, error) {
	apiKeyID, err := apikeyvo.APIKeyIDFromString(input.ID)
	if err != nil {
		return nil, errs.NewErr(errs.BadRequest, err)
	}

	_, err = uc.txManager.Do(ctx, func(txCtx context.Context) (any, error) {
		if err := uc.apiKeyRepo.Delete(txCtx, apiKeyID); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return &DeleteOutput{}, nil
}
