package notifiers

import (
	"log"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/viktorminko/monitor/helper"
	"path"
)

type TelegramSenderConfig struct {
	Token string
	ChatID int64
}

type TelegramSender struct {
	WorkDir string
	chatID int64
	Bot *tgbotapi.BotAPI
}

func (s *TelegramSender) SendMessage(mID string, mBody map[string]interface{}) error {

	//Load message
	message, err := BuildMessage(s.WorkDir, mID, mBody)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(s.chatID, "*"+message.Subject+"*\n\n"+message.Body)
	msg.ParseMode = tgbotapi.ModeMarkdown

	log.Printf("sending message to telegram, mess: "+msg.Text)

	s.Bot.Send(msg)

	return nil
}

func InitTelegramSender(workDir string) (Sender, error) {

	dir := path.Join(
		path.Dir(workDir),
		path.Dir("notifiers/telegram/"),
	)

	config := &TelegramSenderConfig{}
	err := helper.InitObjectFromJsonFile(path.Join(dir, "config.json"), &config)
	if err != nil {
		return nil, err
	}

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}

	bot.Debug = false

	return &TelegramSender{dir,
		config.ChatID,
		bot,
	}, nil
}

