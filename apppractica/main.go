package main

import (
	tgbotapi "github.com/crocone/tg-bot"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var bot *tgbotapi.BotAPI

type button struct {
	name string
	data string
}

func startMenu() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Привет", data: "hi"},
		{name: "Пока", data: "buy"},
		{name: "О чем бот?", data: "help"},
	}

	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, state := range states {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(state.name, state.data),
		))
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func helpMenu() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Назад", data: "exit"},
		{name: "Дополнительная помощь", data: "helpe"},
	}

	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, state := range states {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(state.name, state.data),
		))
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env not loaded")
	}

	bot, err = tgbotapi.NewBotAPI(os.Getenv("TG_API_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot API: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			callbacks(update)
		} else if update.Message != nil && update.Message.IsCommand() {
			commands(update)
		}
	}
}

func callbacks(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	switch data {
	case "hi":
		deleteMess := tgbotapi.NewDeleteMessage(chatID, messageID)
		if _, err := bot.Send(deleteMess); err != nil {
			log.Printf("ERROR")
		}

		text := "Привет!"
		sendText(chatID, text)
	case "buy":
		deleteMess := tgbotapi.NewDeleteMessage(chatID, messageID)
		if _, err := bot.Send(deleteMess); err != nil {
			log.Printf("ERROR")
		}

		text := "Пока!"
		sendText(chatID, text)
	case "help":

		// Переход в меню помощи
		newMenu := helpMenu()
		edit := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, newMenu)
		if _, err := bot.Send(edit); err != nil {
			log.Printf("Failed to edit message: %v", err)
		}
	case "exit":

		// Возврат в стартовое меню
		newMenu := startMenu()
		edit := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, newMenu)
		if _, err := bot.Send(edit); err != nil {
			log.Printf("Failed to edit message: %v", err)
		}
	case "helpe":
		deleteMess := tgbotapi.NewDeleteMessage(chatID, messageID)
		if _, err := bot.Send(deleteMess); err != nil {
			log.Printf("ERROR")
		}
		text := "Дополнительная помощь недоступна в данный момент."
		sendText(chatID, text)
	default:
		sendText(chatID, "Неизвестная команда")
	}
}

func commands(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	switch update.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(chatID, "Выберите действие")
		msg.ReplyMarkup = startMenu()
		sendMessage(msg)
	case "help":
		msg := tgbotapi.NewMessage(chatID, "Выберите действие:")
		msg.ReplyMarkup = helpMenu()
		sendMessage(msg)
	}
}

func sendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	sendMessage(msg)
}

func sendMessage(msg tgbotapi.Chattable) {
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
