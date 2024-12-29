package main

import (
	"log/slog"
	"os"

	"github.com/codejago/polypully-openai-checker/internal"
)

const (
	ERROR = 1
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("application panicked", "error", r)
			os.Exit(ERROR)
		}
	}()
	if err := internal.Run(); err != nil {
		slog.Error("startup crashed", "error", err)
		os.Exit(ERROR)
	}
}
