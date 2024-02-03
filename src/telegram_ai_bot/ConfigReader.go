package telegram_ai_bot

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ChatPropertiesLLMGPT struct {
	ApiKey string `json:"apiKey"`
}

type ChatPropertiesLLM struct {
	OpenAiGpt ChatPropertiesLLMGPT `json:"openai-gpt"`
	IoGpt     ChatPropertiesLLMGPT `json:"io-gpt"`
}

type WeatherProperties struct {
	ApiKey string `json:"apiKey"`
}

type ChatProperties struct {
	//Id  int64             `json:"id"`
	Weather WeatherProperties `json:"weather"`
	Llm     ChatPropertiesLLM `json:"llm"`
}

type TelegramConfig struct {
	Chats    map[int64]ChatProperties `json:"chats"`
	BotToken string                   `json:"botToken"`
}

func readConfig() TelegramConfig {
	fileName := "src/resources/TelegramConfig.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		fmt.Printf("chatsConfig File does not exist, skipping config\n")
		return TelegramConfig{}
	}

	return readFile(fileName)
}

func readFile(fileName string) TelegramConfig {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("err: ", err)
		return TelegramConfig{}
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return TelegramConfig{}
	}

	var config TelegramConfig
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return TelegramConfig{}
	}

	return config
}
