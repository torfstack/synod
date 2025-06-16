package main

import (
	"fmt"
	"log/slog"

	"github.com/torfstack/kayvault/backend"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	fmt.Println(backend.NewApplication().Run())
}
