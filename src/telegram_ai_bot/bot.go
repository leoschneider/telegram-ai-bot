package telegram_ai_bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const gpt35Turbo = "gpt-3.5-turbo-1106"
const gpt4Vision = "gpt-4-vision-preview"
const dallE3 = "dall-e-3"
const gptTranscriptionModel = "whisper-1"
const gptTextToSpeechModel = "tts-1"

var bot *tgbotapi.BotAPI
var telegramConfig *TelegramConfig

func Main() {
	var err error

	readConfiguration := readConfig()
	telegramConfig = &readConfiguration

	bot, err = tgbotapi.NewBotAPI(telegramConfig.BotToken)
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}
	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, updates)

	log.Println("Telegram bot ready! Start listening for updates.")

	// Create a channel to handle OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Block the process until a signal is received
	<-sigCh

	// Handle additional cleanup or stop logic if required
	fmt.Println("Received signal from OS. Stopping process...")
}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	if update.CallbackQuery == nil && update.Message != nil {
		thinkStickerId := GetRandomThinkStickerId()
		loadingSticker := tgbotapi.NewSticker(update.Message.Chat.ID, tgbotapi.FileID(thinkStickerId))
		loadingSticker.DisableNotification = true
		loadingMessage, err := bot.Send(loadingSticker)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer bot.Request(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, loadingMessage.MessageID))
	}

	switch {
	// Handle callbacks and buttons (menu operations)
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery)
		break

	// Handle replies and text queries for bot
	case update.Message != nil && update.Message.ReplyToMessage != nil:
		isHandled := handleReplyResponse(update.Message)
		if isHandled {
			break
		}

	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message)
		break
	}
}

func handleReplyResponse(message *tgbotapi.Message) bool {
	switch message.ReplyToMessage.Text {
	case WeatherCurrentRequestLocationMessage:
		replyWeatherRequest(message, "current")
		return true
	case WeatherForecastRequestLocationMessage:
		replyWeatherRequest(message, "forecast")
		return true
	case ImageDescribeRequestMessage:
		explainImage(message, gpt4Vision)
		return true
	case ImageGenerateRequestMessage:
		createDallE3Image(message, message.Text)
		return true
	case ImageEditRequestMessage:
		editImage(message)
		return true
	case VoiceTranscriptRequestMessage:
		if message.Voice != nil && message.Voice.FileID != "" {
			handleGPTVoice(message, gptTranscriptionModel)
			return true
		}
	case SpeechCreationRequestMessage:
		handleSpeechCreation(message, gptTextToSpeechModel)
		return true
	}
	return false
}

func handleMessage(message *tgbotapi.Message) {
	err := setChatTyping(message)
	if err != nil {
		return
	}

	user := message.From
	if user == nil {
		return
	}

	log.Printf("update from: %s in chat: %d", user.FirstName, message.Chat.ID)

	if strings.HasPrefix(message.Text, "/") {
		err = handleCommand(message)
	} else if message.Photo != nil && len(message.Photo) > 0 && message.Caption != "" {
		explainImage(message, gpt4Vision)
	} else if message.Voice != nil && message.Voice.FileID != "" {
		handleGPTVoice(message, gptTranscriptionModel)
	} else if message.Sticker != nil && message.Sticker.FileID != "" {
		println(message.Sticker.FileID)
	} else {
		handleGPTMessage(message, gpt35Turbo)
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

func setChatTyping(message *tgbotapi.Message) error {
	err := SetChatStatus(message.Chat.ID, tgbotapi.ChatTyping)
	if err != nil {
		fmt.Println("Error setting the typing status:", err)
		return err
	}
	return nil
}

func SetChatStatus(chatId int64, status string) error {
	_, err := bot.Request(tgbotapi.NewChatAction(chatId, status))
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return err
}

func handleCommand(message *tgbotapi.Message) error {
	var err error
	command := message.Text
	chatId := message.Chat.ID

	if command == "/menu" {
		err = sendMenu(message)
	} else if strings.HasPrefix(strings.ToLower(command), "/setkey ") {
		telegramConfig.Chats[chatId] = ChatProperties{Llm: ChatPropertiesLLM{IoGpt: ChatPropertiesLLMGPT{ApiKey: command[len("/setkey "):]}}}
		//(*config).Chats[chatId].Llm.OpenAiGpt.ApiKey = command[len("/setkey "):]
		//(*keys)[chatId] = command[len("/setkey "):]
		_, err = bot.Send(tgbotapi.NewMessage(chatId, "I Registered the new key 👍"))
	} else if strings.HasPrefix(strings.ToLower(command), "/getkey") {
		telegramConfigChat, found := telegramConfig.Chats[chatId]
		if found {
			_, err = bot.Send(tgbotapi.NewMessage(chatId, telegramConfigChat.Llm.IoGpt.ApiKey))
		} else {
			_, err = bot.Send(tgbotapi.NewMessage(chatId, "I could not find your key 😕"))
		}
	} else if strings.HasPrefix(strings.ToLower(command), "/image ") {
		prompt := command[len("/image "):]
		createDallE3Image(message, prompt)
	} else if strings.ToLower(command) == "/polloffice" {
		sendFridayPoll(chatId, message.MessageID)
	} else if command == "/help" {
		_, err = bot.Send(tgbotapi.NewMessage(chatId, "the available commands are: \n"+
			"  - Going to the main menu: /menu\n"+
			"  - Show this page: /help\n"+
			"You can also send a text message to chat with ai.\n"))
	} else {
		_, err = bot.Send(tgbotapi.NewMessage(chatId, "I don't know this command 😕\nThe available commands are: \n  - /setkey {key}\n  - /getkey\n  - /help"))
	}

	return err
}

func handleButton(query *tgbotapi.CallbackQuery) {
	message := query.Message
	menu, foundMenu := menus[query.Data]
	if foundMenu {
		text := "Menu"
		markup := menu

		err := SendWithoutResponse(bot, tgbotapi.NewCallback(query.ID, ""))
		if err != nil {
			log.Printf("tgbotapi.NewCallback: %s", err.Error())
			return
		}

		// Replace menu text and keyboard
		msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
		msg.ParseMode = tgbotapi.ModeHTML
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("NewEditMessageTextAndMarkup: %s", err.Error())
		}
	} else if query.Data == "weather-current" {
		sendReplyRequestMessage(message, WeatherCurrentRequestLocationMessage)
	} else if query.Data == "weather-forecast" {
		sendReplyRequestMessage(message, WeatherForecastRequestLocationMessage)
	} else if query.Data == "image-generate" {
		sendReplyRequestMessage(message, ImageGenerateRequestMessage)
	} else if query.Data == "image-edit" {
		sendReplyRequestMessage(message, ImageEditRequestMessage)
	} else if query.Data == "image-describe" {
		sendReplyRequestMessage(message, ImageDescribeRequestMessage)
	} else if query.Data == "voice-transcript" {
		sendReplyRequestMessage(message, VoiceTranscriptRequestMessage)
	} else if query.Data == "speech-creation" {
		sendReplyRequestMessage(message, SpeechCreationRequestMessage)
	}
}

func sendReplyRequestMessage(message *tgbotapi.Message, requestLocationMessage string) {
	messageToSend := tgbotapi.NewMessage(message.Chat.ID, requestLocationMessage)
	messageToSend.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "",
		Selective:             false,
	}
	bot.Send(messageToSend)
}

func sendMenu(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Choose your option")
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = baseMenu
	msg.DisableNotification = true
	_, err := bot.Send(msg)

	_, err = bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
	return err
}
