package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/codejago/polypully-openai-checker/internal/openai"

	"gopkg.in/yaml.v3"
)

const DefaultConfiguration = "config/application.yaml"

type Config struct {
	Server  Server  `yaml:"server"`
	Tls     Tls     `yaml:"tls"`
	OpenAI  OpenAI  `yaml:"openai"`
	Proxy   string  `yaml:"proxy"`
	Cache   Cache   `yaml:"cache"`
	Metrics Metrics `yaml:"metrics"`
}

type Metrics struct {
	Enabled     bool   `yaml:"enabled"`
	Localhost   string `yaml:"local_host"`
	ContextRoot string `yaml:"context_root"`
}

type OpenAI struct {
	Model                    string  `yaml:"model"`
	Temperature              float64 `yaml:"temperature"`
	MaxTokens                int     `yaml:"max_tokens"`
	TopP                     float64 `yaml:"topp"`
	FrequencyPenalty         float64 `yaml:"frequency_penalty"`
	PresencePenalty          float64 `yaml:"presence_penalty"`
	SystemMessageContentFile string  `yaml:"system_message_content_file"`
	SystemMessage            string  `yaml:"system_message"`
	ChatEndpoint             string  `yaml:"openai_chat_endpoint"`
}

type Server struct {
	Port int `yaml:"port"`
}

type Cache struct {
	Enabled          bool   `yaml:"enabled"`
	Host             string `yaml:"host"`
	KeyHashAlgorithm string `yaml:"key_hash_algorithm"`
}

type Tls struct {
	Enabled      bool   `yaml:"enabled"`
	SelfSigned   bool   `yaml:"self_signed"`
	ServerCert   string `yaml:"server_cert"`
	ServerKey    string `yaml:"server_key"`
	ServerCaCert string `yaml:"server_ca_cert"`
}

func LoadConfig(configuration string) (*Config, error) {
	if configuration == "" {
		configuration = DefaultConfiguration
	}
	data, err := os.ReadFile(configuration)
	if err != nil {
		return nil, err
	}
	var c Config
	if err = yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) GetSystemMessage() (*openai.Message, error) {
	systemMessage := openai.Message{
		Role: "system",
	}
	systemContent, err := os.ReadFile(c.OpenAI.SystemMessageContentFile)
	if err != nil {
		return nil, err
	}
	text := c.OpenAI.SystemMessage + strings.ReplaceAll(string(systemContent), "\"", "\\\"")
	systemMessage.Content = append(systemMessage.Content, openai.Content{
		Text: text,
		Type: "text",
	})
	return &systemMessage, nil
}

func (c *Config) GetAPIKey() (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("environment variable OPENAI_API_KEY not set")
	}
	return apiKey, nil
}
