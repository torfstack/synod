package main

import (
	"log/slog"
	"os"

	"github.com/torfstack/kayvault/backend/config"
	"github.com/torfstack/kayvault/backend/http"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	cfg, err := config.ParseFile("config.yaml")
	if err != nil {
		slog.Error("Failed to parse config file", "error", err.Error())
		os.Exit(1)
	}
	slog.Debug("Config parsed successfully")
	server := http.NewServer(*cfg)
	server.Start()
}
