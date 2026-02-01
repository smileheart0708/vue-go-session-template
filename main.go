package main

import (
	"embed"
	"log/slog"

	"kiroapi/internal/middleware"
	"kiroapi/internal/stream"
)

//go:embed web/dist
var distFS embed.FS

func main() {
	logBroadcaster := stream.NewLogBroadcaster()
	logger := middleware.InitLogger("info", logBroadcaster)
	slog.SetDefault(logger)

	slog.Info("服务启动", "embed", "enabled")
}
