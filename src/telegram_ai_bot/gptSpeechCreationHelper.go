package telegram_ai_bot

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"strings"
	"time"
)

const gptOpenAiAudioCreationUrl = "https://api.openai.com/v1/audio/speech"

const SpeechCreationRequestMessage = "Please provide a text to synthesize as speech:"

func handleSpeechCreation(message *tgbotapi.Message, gptModel string) {
	chatProperties, keyExists := telegramConfig.Chats[message.Chat.ID]
	if !keyExists {
		fmt.Println("Key does not exists!")
		return
	}

	body := createBodyRequest(message, gptModel)
	bodyReader := strings.NewReader(body)

	req, _ := http.NewRequest("POST", gptOpenAiAudioCreationUrl, bodyReader)
	req.Header.Set("Authorization", "Bearer "+chatProperties.Llm.OpenAiGpt.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, _ := client.Do(req)
	defer response.Body.Close()

	if response.StatusCode > 299 {
		fmt.Println("Speech creation helper failed with status: " + response.Status)
		return
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	audioMessage := tgbotapi.NewAudio(message.Chat.ID, tgbotapi.FileBytes{
		Name:  "created-speech-" + time.Now().UTC().Format(time.DateOnly) + ".mp3",
		Bytes: responseBody,
	})

	_, err = bot.Send(audioMessage)
	if err != nil {
		fmt.Println("Error sending audio response message:", err)
		return
	}
}

type BodyRequestCreateSpeech struct {
	Model string `json:"model"`
	Input string `json:"input"`
	Voice string `json:"voice"`
}

func createBodyRequest(message *tgbotapi.Message, gptModel string) string {
	bodyRequest := BodyRequestCreateSpeech{
		Model: gptModel,
		Input: message.Text,
		Voice: "alloy",
	}

	jsonBytes, err := json.Marshal(bodyRequest)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(jsonBytes)
}
