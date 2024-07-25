package usecase

import (
	"context"

	"gitlab.com/back1ng1/messaggio-test/internal/entity"
	"gitlab.com/back1ng1/messaggio-test/internal/kafka"
)

type Usecase interface {
	StoreMessage(msg entity.Message) (entity.Message, error)
	ProcessMessage(msg entity.Message) (entity.Message, error)
	GetStats() (int, error)
}

type repository interface {
	StoreMessage(ctx context.Context, msg entity.Message) (entity.Message, error)
	ProcessMessage(ctx context.Context, msg entity.Message) (entity.Message, error)
	GetStats(ctx context.Context) (int, error)
}

type usecase struct {
	repo   repository
	broker kafka.Kafka
}

func New(repo repository, broker kafka.Kafka) usecase {
	return usecase{
		repo:   repo,
		broker: broker,
	}
}

func (uc usecase) StoreMessage(msg entity.Message) (entity.Message, error) {
	ctx := context.Background()

	msg, err := uc.repo.StoreMessage(ctx, msg)
	if err != nil {
		return msg, err
	}

	err = uc.broker.StoreMessage(msg)
	if err != nil {
		return msg, err
	}

	return msg, nil
}

func (uc usecase) ProcessMessage(msg entity.Message) (entity.Message, error) {
	ctx := context.Background()

	msg, err := uc.repo.ProcessMessage(ctx, msg)
	if err != nil {
		return msg, err
	}

	return msg, nil
}

func (uc usecase) GetStats() (int, error) {
	return uc.repo.GetStats(context.Background())
}
