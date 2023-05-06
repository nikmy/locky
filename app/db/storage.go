package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const dbName = "locky_user_data"

func NewStorage(cfg Config) (*pgStorage, error) {
	source := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port,
		cfg.Credentials.Username, cfg.Credentials.Password,
		dbName,
		func() string {
			if cfg.SSLMode {
				return "enable"
			}
			return "disable"
		}(),
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

func (s *pgStorage) Set(userID int64, service string, login string, password string) error {
	// TODO implement me
	panic("implement me")
}

func (s *pgStorage) Get(userID int64, service string) error {
	// TODO implement me
	panic("implement me")
}

func (s *pgStorage) Del(userID int64, service string) error {
	// TODO implement me
	panic("implement me")
}
