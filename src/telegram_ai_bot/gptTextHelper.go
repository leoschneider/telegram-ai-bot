package telegram_ai_bot

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"strings"
)

const gptIoUrl = "https://iogpt-api-management-service.azure-api.net/openai/api/proxy/openai/chat/completions"

func handleGPTMessage(message *tgbotapi.Message, gptModel string) {
	chatProperties, keyExists := telegramConfig.Chats[message.Chat.ID]
	if !keyExists {
		fmt.Println("Key does not exists!")
		return
	}

	payload := strings.NewReader(buildBody(message.Text, gptModel))

	request, err := http.NewRequest("POST", gptIoUrl, payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+chatProperties.Llm.IoGpt.ApiKey)

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

func buildBody(text string, gptModel string) string {
	bodyRequest := BodyRequest{
		Model: gptModel,
		Messages: []BodyRequestMessage{{
			Role:    "user",
			Content: text,
		}},
	}

	jsonBytes, err := json.Marshal(bodyRequest)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(jsonBytes)
}

type BodyRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BodyRequest struct {
	Model    string               `json:"model"`
	Messages []BodyRequestMessage `json:"messages"`
}

type BodyResponseChoiceMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BodyResponseChoice struct {
	Index   int                       `json:"index"`
	Message BodyResponseChoiceMessage `json:"message"`
}

type BodyResponse struct {
	Id      string               `json:"id"`
	Choices []BodyResponseChoice `json:"choices"`
}
