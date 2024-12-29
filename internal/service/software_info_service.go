package service

import (
	"context"
	"encoding/json"
	pb "github.com/codejago/polypully-openai-checker/api/software"
	"github.com/codejago/polypully-openai-checker/internal/cache"
	"github.com/codejago/polypully-openai-checker/internal/config"
	"github.com/codejago/polypully-openai-checker/internal/openai"
	"log/slog"
	"strings"
)

type Service struct {
	pb.UnimplementedSoftwareInfoServiceServer
	SystemMessage openai.Message
	Oai           openai.OpenAiInterface
	Config        config.Config
	Cache         *cache.Cache
}

type SoftwareInfoService interface {
	GetSoftwareInfo(ctx context.Context, req *pb.InfoRequest) (*pb.InfoResponse, error)
}

// GetSoftwareInfo is the only call in this service
func (s *Service) GetSoftwareInfo(ctx context.Context, req *pb.InfoRequest) (*pb.InfoResponse, error) {

	slog.Debug("received request", "name", req.GetName())
	userMessage := openai.Message{
		Role: "user",
	}
	userMessage.Content = append(userMessage.Content, openai.Content{
		Text: req.Name,
		Type: "text",
	})
	request := openai.ChatRequest{
		Model:            s.Config.OpenAI.Model,
		Temperature:      s.Config.OpenAI.Temperature,
		MaxTokens:        s.Config.OpenAI.MaxTokens,
		Topp:             s.Config.OpenAI.TopP,
		FrequencyPenalty: s.Config.OpenAI.FrequencyPenalty,
		PresencePenalty:  s.Config.OpenAI.PresencePenalty,
	}
	request.Messages = append(request.Messages, s.SystemMessage)
	request.Messages = append(request.Messages, userMessage)

	// caching based on the request
	key, err := s.Cache.KeyInit(request)
	if err != nil {
		slog.Warn("unable to marshall a request for the cache key", "name", req.GetName(), "error", err)
	} else if s.Config.Cache.Enabled {
		if cachedResponse, err := s.Cache.Get(key); err == nil {
			var res pb.InfoResponse
			if err = json.Unmarshal(cachedResponse, &res); err != nil {
				slog.Warn("failed to unmarshal cached response", "name", req.GetName(), "error", err)
			} else {
				slog.Debug("cache returned a hit for request", "name", req.GetName())
				return &res, nil
			}
		}
	}

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
	// openai has been instructed to return a response with a specific format
	if err = json.Unmarshal([]byte(content), &res); err != nil {
		res.Error = err.Error()
	} else if key != "" && s.Config.Cache.Enabled {
		if err = s.Cache.Set(key, content); err != nil {
			slog.Warn("failed to cache response", "name", req.GetName(), "error", err)
		}
	}
	return &res, nil
}
