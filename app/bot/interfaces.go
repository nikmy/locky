package bot

type storage interface {
	Set(userID int64, service string, login string, password string) error
	Get(userID int64, service string) error
	Del(userID int64, service string) error
}
