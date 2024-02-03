package telegram_ai_bot

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"strings"
)

const gptOpenAiImageGenerationUrl = "https://api.openai.com/v1/images/generations"

type ResponseBodyImageGenerationData struct {
	Url string `json:"url"`
}

type ResponseBodyImageGeneration struct {
	Created int                               `json:"created"`
	Data    []ResponseBodyImageGenerationData `json:"data"`
}

func createDallE3Image(message *tgbotapi.Message, prompt string) {
	chatProperties, keyExists := telegramConfig.Chats[message.Chat.ID]
	if !keyExists {
		fmt.Println("Key does not exists!")
		return
	}

	generateImageBody := buildGenerateImageBody(prompt, dallE3)

	payload := strings.NewReader(generateImageBody)

	request, err := http.NewRequest("POST", gptOpenAiImageGenerationUrl, payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+chatProperties.Llm.OpenAiGpt.ApiKey)

	response, err := (&http.Client{}).Do(request)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//fmt.Printf("response status: %s\n", response.Status)

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//fmt.Printf("response body: %s", body)

	var deserializedBody ResponseBodyImageGeneration
	err = json.Unmarshal(body, &deserializedBody)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = SetChatStatus(message.Chat.ID, tgbotapi.ChatUploadPhoto)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, err = bot.Send(tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileURL(deserializedBody.Data[0].Url)))
	if err != nil {
		fmt.Println("Telegram bot cannot upload photo :", err)
		return
	}
}

type BodyRequestGenerateImage struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

func buildGenerateImageBody(prompt string, gptModel string) string {
	bodyRequest := BodyRequestGenerateImage{
		Model:  gptModel,
		Prompt: prompt,
		N:      1,
		Size:   "1024x1024",
	}

	jsonBytes, err := json.Marshal(bodyRequest)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(jsonBytes)
}
