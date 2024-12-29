package config

import (
	"testing"

	"github.com/codejago/polypully-openai-checker/internal/config"
)

func TestLoadConfig(t *testing.T) {
	t.Run("load config", func(t *testing.T) {
		_, err := config.LoadConfig("testdata/application.yaml")
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
	})
}

func TestGetSystemMessage(t *testing.T) {
	t.Run("get system message", func(t *testing.T) {
		c, err := config.LoadConfig("testdata/application.yaml")
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		_, err = c.GetSystemMessage()
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
	})
}
