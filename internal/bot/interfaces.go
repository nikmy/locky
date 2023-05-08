package bot

import "context"

type storage interface {
	Set(ctx context.Context, userID int64, service string, login string, password string) error
	Get(ctx context.Context, userID int64, service string) (string, string, error)
	Del(ctx context.Context, userID int64, service string) error
}
