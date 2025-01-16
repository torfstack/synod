package main

import (
	"log/slog"
	"main/internal/http"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	server := http.NewServer()
	server.Start()
}
