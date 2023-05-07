package bot

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/nikmy/locky/app/workerpool"
)

func Run(ctx context.Context, log *zap.SugaredLogger, api storage, token string, webhookEnabled bool) {
	log.Debug("launching bot...")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("cannot create bot API: %s", err)
	}
	log.Debug("successfully connected to telegram API")

	var updates tgbotapi.UpdatesChannel
	if webhookEnabled {
		updates = fromWebhook(log, bot)
	} else {
		_, _ = bot.Request(tgbotapi.WebhookConfig{})
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates = bot.GetUpdatesChan(u)
	}

	log.Debug("listen for updates...")

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
		log.Debug("update has been received!")
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
			login, password, err := api.Get(ctx, update.Message.From.ID, args[0])
			if err != nil {
				log.Error(storageAPIError("Get", err))
				return
			}

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
			err := api.Set(ctx, update.Message.From.ID, args[0], args[1], args[2])
			if err != nil {
				log.Error(storageAPIError("Set", err))
				return
			}
		case "/del":
			if len(args) != 1 {
				wrongMsg(update)
				return
			}
			err := api.Del(ctx, update.Message.From.ID, args[0])
			if err != nil {
				log.Error(storageAPIError("Del", err))
				return
			}
		default:
			wrongMsg(update)
		}
	}
}

func fromWebhook(log *zap.SugaredLogger, bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	log.Debug("setting up webhook...")
	wh, err := tgbotapi.NewWebhook(os.Getenv("WEBHOOK") + bot.Token)
	if err != nil {
		log.Fatalf("cannot init webhook: %s", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("cannot init webhook: %s", err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatalf("cannot get webhook info: %s", err)
	}
	if info.LastErrorMessage != "" {
		log.Errorf("webhook error: %s", info.LastErrorMessage)
	}

	return bot.ListenForWebhook("/" + bot.Token)
}

func fmtCreds(login, password string) string {
	return fmt.Sprintf("login: %s\npassword: %s", login, password)
}

func storageAPIError(method string, err error) error {
	return fmt.Errorf("storage API error: %s: %w", method, err)
}
