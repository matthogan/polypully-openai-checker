package main

import (
	"log"
	"net"
	"os"
	"strings"

	pb "github.com/codejago/polypully-openai-checker/api/software"
	openai "github.com/codejago/polypully-openai-checker/internal/openai"
	service "github.com/codejago/polypully-openai-checker/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

const (
	CONFIG_FILE = "config/application.yaml"
	TICK        = "\033[32m✔\033[0m"
)

var (
	oai           openai.OpenAI
	c             service.Config
	systemMessage openai.Message
)

// init function to read config file and set up openai client
func init() {
	// api key should be set in env var
	// preferably in a secure way eg using a secret manager like keyvault
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// critical error if api key is not set
		log.Fatalf("env var not set - OPENAI_API_KEY")
	} else {
		log.Printf("OPENAI_API_KEY %s", TICK)
	}
	data, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		log.Fatalf("error reading config file %s: %v", CONFIG_FILE, err)
	}
	c = service.Config{}
	if err = yaml.Unmarshal(data, &c); err != nil {
		log.Fatalf("error parsing config file %s: %v", CONFIG_FILE, err)
	} else {
		log.Printf("Config file %s %s", CONFIG_FILE, TICK)
	}
	system_content, err := os.ReadFile(c.SystemMessageContentFile)
	if err != nil {
		log.Fatalf("error reading system message content file: %v", err)
	} else {
		log.Printf("System message content file %s %s", c.SystemMessageContentFile, TICK)
	}
	// defaults for openai requests
	oai = openai.OpenAI{
		ApiKey:       apiKey,
		ChatEndpoint: c.ChatEndpoint,
		Proxy:        c.Proxy,
	}
	systemMessage = openai.Message{
		Role: "system",
	}
	text := c.SystemMessage + strings.ReplaceAll(string(system_content), "\"", "\\\"")
	systemMessage.Content = append(systemMessage.Content, openai.Content{
		Text: text,
		Type: "text",
	})
}

// main function to start the grpc server
func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSoftwareInfoServiceServer(grpcServer, &service.Server{
		SystemMessage: systemMessage,
		Oai:           &oai,
		Config:        c,
	})

	reflection.Register(grpcServer)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
