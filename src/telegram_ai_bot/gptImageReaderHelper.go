package telegram_ai_bot

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"strings"
)

const gptOpenAiUrl = "https://api.openai.com/v1/chat/completions"

func explainImage(message *tgbotapi.Message, gptModel string) {
	imageBytes := DownloadLatestPhoto(message)
	imageString := base64.StdEncoding.EncodeToString(imageBytes)

	imageBody := buildImageBody(message.Caption, imageString, gptModel)

	base64FormatImageBody := strings.Replace(imageBody,
		"\"image_url\":{\"url\":\"",
		"\"image_url\":{\"url\":\"data:image/jpeg;base64,",
		1)

	payload := strings.NewReader(base64FormatImageBody)

	chatProperties, keyExists := telegramConfig.Chats[message.Chat.ID]
	if !keyExists {
		fmt.Println("Key does not exists!")
		return
	}
	request, err := http.NewRequest("POST", gptOpenAiUrl, payload)
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

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//fmt.Printf("response body: %s", body)

	var deserializedBody BodyResponse
	err = json.Unmarshal(body, &deserializedBody)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	responseMessage := deserializedBody.Choices[len(deserializedBody.Choices)-1].Message.Content
	telegramMessage := tgbotapi.NewMessage(message.Chat.ID, responseMessage)
	telegramMessage.ParseMode = tgbotapi.ModeMarkdown

	_, err = bot.Send(telegramMessage)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func buildImageBody(caption string, image string, gptModel string) string {
	base64Url := image
	//	f"data:image/jpeg;base64,{base64_image}"
	bodyRequest := BodyRequestImage{
		Model: gptModel,
		Messages: []BodyRequestImageMessage{{
			Role: "user",
			Content: []BodyRequestImageMessageContent{
				BodyRequestImageMessageContentText{
					Type: "text",
					Text: caption,
				},
				BodyRequestImageMessageContentImage{
					Type: "image_url",
					ImageUrl: BodyRequestImageMessageContentImageUrl{
						Url: base64Url,
					},
				},
			},
		}},
		MaxTokens: 300,
	}

	jsonBytes, err := json.Marshal(bodyRequest)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(jsonBytes)
}

type BodyRequestImageMessageContentImageUrl struct {
	Url string `json:"url"`
}

type BodyRequestImageMessageContentText struct {
	Type string `json:"type"`
	Text string `json:"text"`
	//ImageUrl BodyRequestImageMessageContentImageUrl `json:"image_url"`
}

type BodyRequestImageMessageContentImage struct {
	Type string `json:"type"`
	//Text     string                                 `json:"text"`
	ImageUrl BodyRequestImageMessageContentImageUrl `json:"image_url"`
}

type BodyRequestImageMessageContent interface {
}

type BodyRequestImageMessage struct {
	Role    string                           `json:"role"`
	Content []BodyRequestImageMessageContent `json:"content"`
}

type BodyRequestImage struct {
	Model     string                    `json:"model"`
	Messages  []BodyRequestImageMessage `json:"messages"`
	MaxTokens int                       `json:"max_tokens"`
}
