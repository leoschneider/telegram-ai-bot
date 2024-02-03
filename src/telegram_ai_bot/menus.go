package telegram_ai_bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var menus = map[string]tgbotapi.InlineKeyboardMarkup{
	"baseMenu":    baseMenu,
	"weatherMenu": weatherMenu,
	"imageMenu":   imageMenu,
	"soundMenu":   soundMenu,
}

var baseMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("weather", "weatherMenu"),
		tgbotapi.NewInlineKeyboardButtonData("Images", "imageMenu"),
		tgbotapi.NewInlineKeyboardButtonData("Sound", "soundMenu"),
	),
)

var weatherMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Back", "baseMenu")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("current weather", "weather-current")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Weather forecast", "weather-forecast"),
	),
)

func buildCurrentWeatherMenu(url string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "weatherMenu")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Details", url),
		))
}

var imageMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Back", "baseMenu")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("generate Image", "image-generate")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("edit image", "image-edit")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("describe image", "image-describe")),
)

var soundMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Back", "baseMenu")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("transcript voice", "voice-transcript")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Text to speech", "speech-creation")),
)
