package service

import (
	"context"
	"encoding/json"
	"strings"

	pb "github.com/codejago/polypully-openai-checker/api/software"
	openai "github.com/codejago/polypully-openai-checker/internal/openai"
)

type Server struct {
	pb.UnimplementedSoftwareInfoServiceServer
	SystemMessage openai.Message
	Oai           openai.OpenAiInterface
	Config        Config
}

type Config struct {
	Model                    string  `yaml:"model"`
	Temperature              float64 `yaml:"temperature"`
	MaxTokens                int     `yaml:"max_tokens"`
	TopP                     float64 `yaml:"topp"`
	FrequencyPenalty         float64 `yaml:"frequency_penalty"`
	PresencePenalty          float64 `yaml:"presence_penalty"`
	SystemMessageContentFile string  `yaml:"system_message_content_file"`
	SystemMessage            string  `yaml:"system_message"`
	ChatEndpoint             string  `yaml:"openai_chat_endpoint"`
	Proxy                    string  `yaml:"proxy"`
}

// the only call in this service
func (s *Server) GetSoftwareInfo(ctx context.Context, req *pb.InfoRequest) (*pb.InfoResponse, error) {

	userMessage := openai.Message{
		Role: "user",
	}
	userMessage.Content = append(userMessage.Content, openai.Content{
		Text: req.Name,
		Type: "text",
	})
	request := openai.ChatRequest{
		Model:            s.Config.Model,
		Temperature:      s.Config.Temperature,
		MaxTokens:        s.Config.MaxTokens,
		Topp:             s.Config.TopP,
		FrequencyPenalty: s.Config.FrequencyPenalty,
		PresencePenalty:  s.Config.PresencePenalty,
	}
	request.Messages = append(request.Messages, s.SystemMessage)
	request.Messages = append(request.Messages, userMessage)

	openAiResponse, err := s.Oai.Chat(&request)
	var res pb.InfoResponse
	if err != nil {
		res.Error = err.Error()
		return &res, nil
	}

	content := openAiResponse.Choices[0].Message.Content
	if content == "" {
		res.Error = "empty response"
		return &res, nil
	}
	content = strings.ReplaceAll(content, "\\\"", "\"")
	// openai has been instructed to return a response
	// with a specific format
	if err = json.Unmarshal([]byte(content), &res); err != nil {
		res.Error = err.Error()
	}
	return &res, nil
}
