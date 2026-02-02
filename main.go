package main

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"kiroapi/internal/middleware"
	"kiroapi/internal/stream"
)

//go:embed web/dist
var distFS embed.FS

func main() {
	logBroadcaster := stream.NewLogBroadcaster()
	logger := middleware.InitLogger("info", logBroadcaster)
	slog.SetDefault(logger)
	// 创建 gin 路由
	r := gin.New()
	// 静态文件服务
	r.StaticFS("/", http.FS(distFS))
	// 启动服务器
	r.Run(":8080")
}
