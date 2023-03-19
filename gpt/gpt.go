package gpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/qingconglaixueit/wechatbot/config"
)

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
// 	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
	Error   struct {
		Message string      `json:"message"`
		Type    string      `json:"type"`
		Param   interface{} `json:"param"`
		Code    interface{} `json:"code"`
	} `json:"error"`
}

type ChoiceItem struct {
	Message      Message `json:"message"`
	Index        int    `json:"index"`
}


type Message struct {
	Role      string `json:"role"`
	Content        string    `json:"content"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string  `json:"model"`
	Password            string  `json:"password"`
	Messages         []Message  `json:"messages"`
// 	MaxTokens        uint    `json:"max_tokens"`
// 	Temperature      float64 `json:"temperature"`
// 	TopP             int     `json:"top_p"`
// 	FrequencyPenalty int     `json:"frequency_penalty"`
// 	PresencePenalty  int     `json:"presence_penalty"`
}


// Completions gtp文本模型回复
//curl https://api.openai.com/v1/completions
//-H "Content-Type: application/json"
//-H "Authorization: Bearer your chatGPT key"
//-d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string) (string, error) {
    var reply string
	var resErr error
	for retry := 1; retry <= 3; retry++ {
		if retry > 1 {
			time.Sleep(time.Duration(retry-1) * 100 * time.Millisecond)
		}
		reply, resErr = httpRequestCompletions(msg, retry)
		if resErr != nil {
			log.Printf("gpt request(%d) error: %v\n", retry, resErr)
			continue
		}
		if reply != "" {
			break
		}
	}
	if resErr != nil {
		return "", resErr
	}
// 	var reply string
// 	if gptResponseBody != nil && len(gptResponseBody.Choices) > 0 {
// 		reply = gptResponseBody.Choices[0].Message.Content
// 	}
	return reply, nil
}

func httpRequestCompletions(msg string, runtimes int) (string, error) {
	cfg := config.LoadConfig()
	if cfg.ApiKey == "" {
		return "", errors.New("api key required")
	}
	message := []Message{
	    {Role: "user",Content:msg},
	}

	requestBody := ChatGPTRequestBody{
		Model:            cfg.Model,
		Messages:         message,
		Password: "1234",
	}
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("json.Marshal requestBody error: %v", err)
	}

	log.Printf("gpt request(%d) json: %s\n", runtimes, string(requestData))

// 	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestData))
	req, err := http.NewRequest(http.MethodPost, "https://chatgpt-vercel-gamma-six.vercel.app/api", bytes.NewBuffer(requestData))
	if err != nil {
		return "", fmt.Errorf("http.NewRequest error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.ApiKey)
	client := &http.Client{Timeout: 15 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("client.Do error: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

    var reply string
    reply = string(body)
	log.Printf("gpt response(%d) json: %s\n", runtimes, reply)

	return reply, nil
}
