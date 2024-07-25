package repository

import (
	"context"
	"errors"
	"fmt"
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

	if err := row.Scan(&msg.ID); err != nil {
		return msg, err
	}

	if err := tx.Commit(ctx); err != nil {
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

	var proccessedAt *time.Time
	sql := "SELECT processed_at FROM messages WHERE id = $1"
	row := tx.QueryRow(ctx, sql, msg.ID)
	row.Scan(&proccessedAt)
	if proccessedAt != nil {
		return msg, fmt.Errorf("%d msg already proccessed", msg.ID)
	}

	sql = "UPDATE messages SET processed_at = $1 WHERE id = $2 AND processed_at IS NULL"
	commandTag, err := tx.Exec(ctx, sql, msg.ProcessedAt.Format(time.DateTime), msg.ID)
	if err != nil {
		return msg, err
	}

	if commandTag.RowsAffected() != 1 {
		return msg, errors.New("not expected rows count has affected")
	}

	if err := tx.Commit(ctx); err != nil {
		return msg, err
	}

	return msg, nil
}

// GetStats returns count of processed messages for all time
func (r repo) GetStats(ctx context.Context) (int, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var count int

	sql := "SELECT COUNT(processed_at) FROM messages WHERE processed_at IS NOT NULL"
	row := r.pool.QueryRow(ctx, sql)

	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return count, nil
}
