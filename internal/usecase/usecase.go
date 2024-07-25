package usecase

import (
	"context"

	"gitlab.com/back1ng1/messaggio-test/internal/entity"
)

type Usecase interface {
	StoreMessage(msg entity.Message) (entity.Message, error)
	ProcessMessage(msg entity.Message) (entity.Message, error)
}

type repository interface {
	StoreMessage(ctx context.Context, msg entity.Message) (entity.Message, error)
	ProcessMessage(ctx context.Context, msg entity.Message) (entity.Message, error)
}

type usecase struct {
	repo repository
}

func New(repo repository) usecase {
	return usecase{
		repo: repo,
	}
}

func (uc usecase) StoreMessage(msg entity.Message) (entity.Message, error) {
	ctx := context.Background()

	return uc.repo.StoreMessage(ctx, msg)
}

func (uc usecase) ProcessMessage(msg entity.Message) (entity.Message, error) {
	ctx := context.Background()

	return uc.repo.ProcessMessage(ctx, msg)
}
