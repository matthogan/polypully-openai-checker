package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	OPENAI_CHAT_ENDPOINT = "https://api.openai.com/v1/chat/completions"
)

type OpenAI struct {
	ApiKey       string
	ChatEndpoint string
	Proxy        string
}

type ChatRequest struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	Temperature      float64   `json:"temperature"`
	MaxTokens        int       `json:"max_tokens"`
	Topp             float64   `json:"top_p"`
	FrequencyPenalty float64   `json:"frequency_penalty"`
	PresencePenalty  float64   `json:"presence_penalty"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type ChatResponse struct {
	Id      string `json:"id"`
	Cbject  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"` // our data
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type OpenAiInterface interface {
	Chat(request *ChatRequest) (*ChatResponse, error)
}

func (o *OpenAI) Chat(request *ChatRequest) (*ChatResponse, error) {

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	endpoint := o.ChatEndpoint
	if endpoint == "" {
		endpoint = OPENAI_CHAT_ENDPOINT
	}
	req, err := http.NewRequest("POST", endpoint, io.NopCloser(bytes.NewReader(requestBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.ApiKey)

	transport := http.DefaultTransport
	if o.Proxy != "" {
		proxyUrl, err := url.Parse(o.Proxy)
		if err != nil {
			return nil, err
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}
	client := &http.Client{}
	client.Transport = transport

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API returned status code: %d", resp.StatusCode)
	}
	var res ChatResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}
