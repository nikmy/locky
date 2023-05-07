package bot

import (
	"context"
	"database/sql"
	"errors"
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
		u.Timeout = updateTimeout
		updates = bot.GetUpdatesChan(u)
	}

	log.Debug("listen for updates...")

	workerpool.New[tgbotapi.Update](nWorkers).
		WithContext(ctx).
		WithHandler(runner(ctx, log, bot, api)).
		Range(updates)
}

func runner(ctx context.Context, log *zap.SugaredLogger, bot *tgbotapi.BotAPI, storage storage) func(update tgbotapi.Update) {
	badCommand := func(update tgbotapi.Update) {
		telegramAPIRequest(log, bot, tgbotapi.NewMessage(update.Message.Chat.ID, "I don't understand :("))
	}
	return func(update tgbotapi.Update) {
		log.Debug("update has been received")
		if update.Message == nil || len(update.Message.Text) == 0 {
			badCommand(update)
			return
		}

		split := strings.Split(update.Message.Text, " ")
		cmd, args := split[0], split[1:]

		switch cmd {
		case "/start":
			telegramAPIRequest(log, bot, tgbotapi.NewMessage(update.Message.Chat.ID, helloMessage))
		case "/get":
			if len(args) != 1 {
				badCommand(update)
				return
			}
			login, password, err := storage.Get(ctx, update.Message.From.ID, args[0])
			if err != nil {
				if errors.As(err, &sql.ErrNoRows) {
					telegramAPIRequest(log, bot, tgbotapi.NewMessage(update.Message.Chat.ID, "no matches found"))
					return
				}
				log.Error(storageAPIError("Get", err))
				return
			}

			m := tgbotapi.NewMessage(update.Message.Chat.ID, fmtCreds(login, password))
			m.ParseMode = "MarkdownV2"
			msg, err := bot.Send(m)
			if err != nil {
				log.Errorf("cannot send message: %s", err)
			}
			go func() {
				select {
				case <-ctx.Done():
					return
				case <-time.After(20 * time.Second): // TODO: timer pool
				}
				telegramAPIRequest(log, bot, tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID))
			}()
		case "/set":
			if len(args) != 3 {
				badCommand(update)
				return
			}
			err := storage.Set(ctx, update.Message.From.ID, args[0], args[1], args[2])
			if err != nil {
				log.Error(storageAPIError("Set", err))
				return
			}
			telegramAPIRequest(log, bot, tgbotapi.NewMessage(update.Message.Chat.ID, "saved."))
		case "/del":
			if len(args) != 1 {
				badCommand(update)
				return
			}
			err := storage.Del(ctx, update.Message.From.ID, args[0])
			if err != nil {
				if errors.As(err, &sql.ErrNoRows) {
					telegramAPIRequest(log, bot, tgbotapi.NewMessage(update.Message.Chat.ID, "no matches found"))
					return
				}
				log.Error(storageAPIError("Del", err))
				return
			}
			telegramAPIRequest(log, bot, tgbotapi.NewMessage(update.Message.Chat.ID, "ok"))
		default:
			badCommand(update)
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

func telegramAPIRequest(log *zap.SugaredLogger, bot *tgbotapi.BotAPI, c tgbotapi.Chattable) {
	_, err := bot.Request(c)
	if err != nil {
		log.Errorf("telegram API error: %s", err)
	}
}

func fmtCreds(login, password string) string {
	return fmt.Sprintf("login: `%s`\npassword: `%s`", login, password)
}

func storageAPIError(method string, err error) error {
	return fmt.Errorf("storage API error: %s: %w", method, err)
}
