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
	_, err := s.db.ExecContext(ctx, "SELECT set_password($1, $2, $3, $4);", userID, service, login, password)
	return err
}

func (s *pgStorage) Get(ctx context.Context, userID int64, service string) (login string, password string, err error) {
	err = s.db.QueryRowxContext(ctx, "SELECT * FROM get_credentials($1, $2);", userID, service).Scan(&login, &password)
	return
}

func (s *pgStorage) Del(ctx context.Context, userID int64, service string) error {
	_, err := s.db.ExecContext(ctx, "SELECT delete_credentials($1, $2);", userID, service)
	return err
}
