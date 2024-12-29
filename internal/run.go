package internal

import (
	"fmt"
	"github.com/codejago/polypully-openai-checker/internal/config"
	"github.com/codejago/polypully-openai-checker/internal/metrics"
	"github.com/codejago/polypully-openai-checker/internal/openai"
	"github.com/codejago/polypully-openai-checker/internal/server"
	"log/slog"
)

var (
	oai openai.OpenAI
)

func Run() error {
	// cleanup on exit
	defer cleanup()

	c, err := config.LoadConfig("")
	if err != nil {
		return fmt.Errorf("config load 〤")
	}
	slog.Info("config load ✔")

	apiKey, err := c.GetAPIKey()
	if err != nil {
		return fmt.Errorf("OPENAI_API_KEY 〤: %v", err)
	}
	slog.Info("OPENAI_API_KEY ✔")

	systemMessage, err := c.GetSystemMessage()
	if err != nil {
		return fmt.Errorf("system message load 〤")
	}
	slog.Info("system message load ✔")

	if c.Tls.Enabled {
		slog.Info("tls enabled ✔")
	} else {
		slog.Info("tls enabled 〤")
	}

	// Serve Prometheus metrics on /metrics
	metrics.StartMetrics(&c.Metrics)

	oai = openai.OpenAI{
		ApiKey:       apiKey,
		ChatEndpoint: c.OpenAI.ChatEndpoint,
		Proxy:        c.Proxy,
	}

	grpcServer := server.NewServer(c)
	if err = grpcServer.StartServer(systemMessage, &oai); err != nil {
		return fmt.Errorf("server start 〤: %v", err)
	}

	return nil
}

func cleanup() {
	slog.Info("cleaning up")
	metrics.StopMetrics()
}
