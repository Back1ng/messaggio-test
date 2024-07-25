package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/back1ng1/messaggio-test/internal/entity"
)

type repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) repo {
	return repo{
		pool: pool,
	}
}

func (r repo) StoreMessage(ctx context.Context, msg entity.Message) (entity.Message, error) {
	if msg.Message == "" {
		return msg, errors.New("empty message given")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return msg, err
	}
	defer tx.Rollback(ctx)

	msg.CreatedAt = time.Now().UTC()

	sql := "INSERT INTO messages(message, created_at) VALUES($1, $2) RETURNING id"
	row := tx.QueryRow(ctx, sql, msg.Message, msg.CreatedAt.Format(time.DateTime))

	err = row.Scan(&msg.ID)
	if err != nil {
		return msg, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return msg, err
	}

	return msg, nil
}

func (r repo) ProcessMessage(ctx context.Context, msg entity.Message) (entity.Message, error) {
	if msg.ID == 0 {
		return msg, errors.New("incorrect id of message given")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return msg, err
	}
	defer tx.Rollback(ctx)

	t := time.Now().UTC()
	msg.ProcessedAt = &t

	sql := "UPDATE messages SET processed_at = $1 WHERE id = $2"
	commandTag, err := tx.Exec(ctx, sql, msg.ProcessedAt.Format(time.DateTime), msg.ID)
	if err != nil {
		return msg, err
	}

	if commandTag.RowsAffected() != 1 {
		return msg, errors.New("not expected rows count has affected")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return msg, err
	}

	return msg, nil
}
