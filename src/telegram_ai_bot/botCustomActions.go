package telegram_ai_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"math/rand"
	"time"

	"net/http"
)

func DownloadLatestPhoto(message *tgbotapi.Message) []byte {
	fileId := message.Photo[len(message.Photo)-1].FileID
	return retrieveFile(fileId)
}

func retrieveFile(fileId string) []byte {
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileId})
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	downloadLink := file.Link(telegramConfig.BotToken)

	response, err := http.Get(downloadLink)
	if err != nil {
		fmt.Println("Error downloading file with id :", fileId, err)
		return []byte{}
	}
	defer response.Body.Close()

	body, err2 := io.ReadAll(response.Body)
	if err2 != nil {
		fmt.Println("Error reading body of downloaded file with id :", fileId, err)
		return []byte{}
	}

	return body
}

func sendFridayPoll(chatId int64, messageId int) {
	today := time.Now()

	// Find the next Friday
	daysUntilFriday := (5 - int(today.Weekday()) + 7) % 7
	nextFriday := today.AddDate(0, 0, daysUntilFriday)
	nextFridayStr := nextFriday.Format("Monday 2 January")

	bot.Send(tgbotapi.SendPollConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: chatId,
		},
		Question:              "Who is going to Utrecht on: " + nextFridayStr,
		Options:               []string{"I'll go to the office", "I'll stay home", "I'll be at the client", "I don't work friday", "Other"},
		IsAnonymous:           false,
		AllowsMultipleAnswers: true,
	})

	bot.Request(tgbotapi.NewDeleteMessage(chatId, messageId))
}

func SendWithoutResponse(bot *tgbotapi.BotAPI, c tgbotapi.Chattable) error {
	_, err := bot.Request(c)
	if err != nil {
		return err
	}

	return nil
}

var thinkingStickerIds = []string{
	"CAACAgIAAxkBAAOwZZmZm32z04f3y0XcmrTD0X716-AAAqYXAAJ39xFKJGX2cQOxMaw0BA",
	"CAACAgIAAxkBAAIBy2WqfPzPPc1Tlc3-F0-bukeKCPREAAICAQACVp29Ck7ibIHLQOT_NAQ",
	"CAACAgEAAxkBAAIBzWWqfSFO3ruYzKwYBZxGkmfzB2bOAAIeAQACOA6CEUZYaNdphl79NAQ",
	"CAACAgIAAxkBAAIBz2WqfT-HGIcjJbK95R0e5cAvvMbhAAJvBQACP5XMCsA0UmHcq07INAQ",
	"CAACAgIAAxkBAAIB0WWqfVv41cPtN_BNNJoW4JbsEUQbAAJ_AAOmysgMlaChOYY6gsw0BA",
	"CAACAgIAAxkBAAIB02WqfXHy68ntv1G3xVmEFfwkKVShAAIZAQACUomRI27pC30cRuNINAQ",
	"CAACAgIAAxkBAAIB1WWqfeRYx-3eTVoEpoEdoyWM9RiVAAIzAQAC9wLIDzvK4ZTu2U7NNAQ",
}

func GetRandomThinkStickerId() string {
	return thinkingStickerIds[rand.Intn(len(thinkingStickerIds))]
}
