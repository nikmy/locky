package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nikmy/locky/app/workerpool"
	"go.uber.org/zap"
)

const (
	webhookURL = "TODO"
)

func Run(ctx context.Context, log *zap.SugaredLogger, api storage, token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("cannot create bot API: %s", err)
	}

	wh, err := tgbotapi.NewWebhook(webhookURL + bot.Token)
	if err != nil {
		log.Fatalf("cannot init webhook: %s", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("cannot init webhook: %s", err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)

	workerpool.New[tgbotapi.Update](8).
		WithContext(ctx).
		WithHandler(runner(ctx, log, bot, api)).
		Range(updates)
}

func runner(ctx context.Context, log *zap.SugaredLogger, bot *tgbotapi.BotAPI, api storage) func(update tgbotapi.Update) {
	wrongMsg := func(update tgbotapi.Update) {
		_, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "I don't understand :("))
		if err != nil {
			log.Errorf("cannot send msg: %s", err)
		}
	}
	return func(update tgbotapi.Update) {
		if update.Message == nil || len(update.Message.Text) == 0 {
			wrongMsg(update)
			return
		}

		split := strings.Split(update.Message.Text, " ")
		if len(split) < 2 {
			wrongMsg(update)
			return
		}
		cmd, args := split[0], split[1:]

		switch cmd {
		case "/get":
			if len(args) != 1 {
				wrongMsg(update)
				return
			}
			login, password, _ := api.Get(ctx, update.Message.From.ID, args[0]) // TODO: error
			msg, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmtCreds(login, password)))
			if err != nil {
				log.Errorf("cannot send msg: %s", err)
			}
			go func() {
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Minute): // TODO: timer pool
				}
				_, err := bot.Request(tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID))
				if err != nil {
					log.Errorf("cannot delete sent message: %s", err)
				}
			}()
		case "/set":
			if len(args) != 3 {
				wrongMsg(update)
				return
			}
			_ = api.Set(ctx, update.Message.From.ID, args[0], args[1], args[2]) // TODO: error
		case "/del":
			if len(args) != 1 {
				wrongMsg(update)
				return
			}
			_ = api.Del(ctx, update.Message.From.ID, args[0]) // TODO: error
		default:
			wrongMsg(update)
		}
	}
}

func fmtCreds(login, password string) string {
	return fmt.Sprintf("login: %s\npassword: %s", login, password)
}
