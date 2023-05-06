package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewStorage(cfg Config) (*pgStorage, error) {
	source := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	db, err := sqlx.Connect("postgres", source)
	if err != nil {
		return nil, fmt.Errorf("cannot connect database: %w", err)
	}
	return &pgStorage{
		db: db,
	}, nil
}

type pgStorage struct {
	db *sqlx.DB
}

func (s *pgStorage) Set(ctx context.Context, userID int64, service string, login string, password string) error {
	// TODO implement me
	panic("implement me")
}

func (s *pgStorage) Get(ctx context.Context, userID int64, service string) error {
	// TODO implement me
	panic("implement me")
}

func (s *pgStorage) Del(ctx context.Context, userID int64, service string) error {
	// TODO implement me
	panic("implement me")
}
