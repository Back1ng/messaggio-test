package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(ctx context.Context, opt SetupOptions) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(opt.String())
	if err != nil {
		return nil, err
	}

	config.MinConns = 10

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return conn, nil

}

type SetupOptions struct {
	Username string
	Password string
	Host     string
	Port     uint16
	Database string
}

func (opt SetupOptions) String() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		opt.Username,
		opt.Password,
		opt.Host,
		opt.Port,
		opt.Database,
	)
}
