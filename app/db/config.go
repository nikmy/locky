package db

type Credentials struct {
	Username string
	Password string
}

type Config struct {
	Credentials

	Host    string
	Port    uint16
	SSLMode bool
}
