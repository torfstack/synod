package main

import (
	"github.com/torfstack/kayvault/backend"
	"log/slog"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	_ = backend.NewApplication().Run()
}
