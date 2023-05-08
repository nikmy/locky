package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/nikmy/locky/internal/bot"
	"github.com/nikmy/locky/internal/db"
)

func main() {
	zlog, _ := zap.NewDevelopment()
	log := zlog.Sugar()

	log.Debug("connecting to storage...")
	api, err := db.NewStorage(loadConfigFromEnv())
	if err != nil {
		log.Fatalf("cannot initialize db storage: %s", err)
	}
	log.Debug("storage has been successfully connected")

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	<-bot.Run(ctx, log, api, os.Getenv("TOKEN"))
	log.Info("Graceful shutdown")
}

func loadConfigFromEnv() db.Config {
	var cfg db.Config

	cfg.Host = os.Getenv("HOST")
	if p, err := strconv.ParseUint(os.Getenv("PORT"), 10, 16); err != nil {
		panic(fmt.Errorf("cannot parse port number: %s", err))
	} else {
		cfg.Port = uint16(p)
	}

	cfg.Username = os.Getenv("USER")
	cfg.Password = os.Getenv("PASSWORD")
	cfg.DBName = os.Getenv("DBNAME")
	cfg.SSLMode = os.Getenv("SSLMODE")

	return cfg
}
