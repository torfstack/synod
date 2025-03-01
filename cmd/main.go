package main

import (
	"log/slog"

	"github.com/torfstack/kayvault/backend"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	_ = backend.NewApplication().Run()
}
