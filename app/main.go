package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/nikmy/locky/app/bot"
	"github.com/nikmy/locky/app/db"
	"go.uber.org/zap"
)

func main() {
	zlog, _ := zap.NewDevelopment()
	log := zlog.Sugar()

	api, err := db.NewStorage(db.Config{
		Host: "localhost",
		Port: 5432,
		Credentials: db.Credentials{
			Username: "postresql",
			Password: "postresql",
		},
		SSLMode: false,
	})
	if err != nil {
		log.Fatalf("cannot initialize db storage: %s", err)
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	bot.Run(ctx, log, api, "")
}
