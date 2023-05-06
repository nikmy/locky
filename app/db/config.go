package db

type Config struct {
	Host        string
	Port        uint16
	Credentials struct {
		Username string
		Password string
	}
	SSLMode bool
}
