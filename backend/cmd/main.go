package main

import (
	"github.com/torfstack/kayvault/internal/http"
	"log/slog"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	server := http.NewServer()
	server.Start()
}
