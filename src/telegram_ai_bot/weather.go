package telegram_ai_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	weather_api "telegram-ai-bot/src/telegram_ai_bot/weather-api"
)

const WeatherCurrentRequestLocationMessage = "Current Weather: which location?"
const WeatherForecastRequestLocationMessage = "Forecast Weather: which location?"

func replyWeatherRequest(message *tgbotapi.Message, weatherType string) {
	locationQuery := message.Text
	text, markup := handleWeatherRequest(message, weatherType, locationQuery)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = markup
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("NewEditMessageTextAndMarkup: %s", err.Error())
	}
}

func handleWeatherRequest(message *tgbotapi.Message, weatherType string, queryLocation string) (string, tgbotapi.InlineKeyboardMarkup) {

	telegramConfigChat, telegramConfigChatFound := telegramConfig.Chats[message.Chat.ID]

	if !telegramConfigChatFound || telegramConfigChat.Weather.ApiKey == "" {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "I could not find your key 😕"))
		return "I could not find your key 😕", tgbotapi.NewInlineKeyboardMarkup()
	}

	text := weather_api.GetWeatherResponse(weatherType, telegramConfigChat.Weather.ApiKey, queryLocation)
	markup := buildCurrentWeatherMenu(weather_api.WeatherUiEndpoint + queryLocation)

	return text, markup
}
