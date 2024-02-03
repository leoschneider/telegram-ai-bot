package telegram_ai_bot

import (
	"bytes"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const gptOpenAiAudioTranscriptionUrl = "https://api.openai.com/v1/audio/transcriptions"

const VoiceTranscriptRequestMessage = "Please provide a voice message to transcript:"

func handleGPTVoice(message *tgbotapi.Message, gptModel string) {
	tmpFile, err := os.CreateTemp("src/generated", "audio-*.mp3")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	audioBytes := retrieveVoice(message)
	if _, err = tmpFile.Write(audioBytes); err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return
	}

	tmpFilePath := tmpFile.Name()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add the file to the form data
	file, _ := os.Open(tmpFilePath)
	defer file.Close()
	part, _ := writer.CreateFormFile("file", tmpFilePath)
	io.Copy(part, file)

	// Add other form fields to the form data
	writer.WriteField("model", gptModel)
	writer.WriteField("response_format", "text")
	writer.Close()

	chatProperties, keyExists := telegramConfig.Chats[message.Chat.ID]
	if !keyExists {
		fmt.Println("Key does not exists!")
		return
	}

	req, _ := http.NewRequest("POST", gptOpenAiAudioTranscriptionUrl, body)
	req.Header.Set("Authorization", "Bearer "+chatProperties.Llm.OpenAiGpt.ApiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	result, _ := io.ReadAll(resp.Body)
	fmt.Println(string(result))

	telegramMessage := tgbotapi.NewMessage(message.Chat.ID, string(result))
	telegramMessage.ParseMode = tgbotapi.ModeMarkdown
	telegramMessage.ReplyToMessageID = message.MessageID

	_, err = bot.Send(telegramMessage)
	if err != nil {
		fmt.Println("Error sending response message:", err)
		return
	}
}

func retrieveVoice(message *tgbotapi.Message) []byte {
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: message.Voice.FileID})
	if err != nil {
		fmt.Println("Error retrieving voice:", err)
		return nil
	}

	downloadLink := file.Link(telegramConfig.BotToken)

	response, err := http.Get(downloadLink)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err2 := io.ReadAll(response.Body)
	if err2 != nil {
		panic(err2)
	}

	return body
}
