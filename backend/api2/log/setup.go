package log

import (
	"log/slog"
	"os"
)

func SetupDebugLevelLogger() {
	var debugLevel = new(slog.LevelVar)
	debugLevel.Set(slog.LevelDebug)
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: debugLevel}))
}
