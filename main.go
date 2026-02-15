package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"main/internal/config"
	"main/internal/database"
	"main/internal/middleware"
	"main/internal/server"
	"main/internal/session"
	"main/internal/stream"
)

//go:embed web/dist
var distFS embed.FS

// printBanner 打印启动横幅
func printBanner(cfg *config.Config) {
	separator := "════════════════════════════════════════════════"
	fmt.Println(separator)
	fmt.Printf("服务地址: http://localhost:%d\n", cfg.Port)
	fmt.Printf("日志级别: %s\n", cfg.LogLevel)
	fmt.Printf("数据目录: %s\n", cfg.DataDir)
	fmt.Printf("数据库文件: %s\n", filepath.Join(cfg.DataDir, "data.db"))
	fmt.Println(separator)
}

func main() {
	startTime := time.Now().Unix()

	// 加载配置
	cfg := config.Load()

	// 初始化日志系统
	logBroadcaster := stream.NewLogBroadcaster()
	logger := middleware.InitLogger(cfg.LogLevel, logBroadcaster)
	slog.SetDefault(logger)

	// 打印启动横幅
	printBanner(cfg)

	// 如果是自动生成的 AUTH_KEY，则打印出来
	if cfg.IsAutoAuthKey {
		slog.Info("自动生成的 AUTH_KEY", "auth_key", cfg.AuthKey)
	}

	// 初始化 SQLite 数据库
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()

	dbPath := filepath.Join(cfg.DataDir, "data.db")
	dbContainer, err := database.Open(dbCtx, database.Options{Path: dbPath})
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		return
	}
	defer func() {
		if err := dbContainer.Close(); err != nil {
			slog.Warn("failed to close database", "error", err)
		}
	}()
	slog.Info("database initialized", "path", dbContainer.Path())

	// 初始化 session 管理器
	sessionManager, err := session.NewManager(cfg.DataDir)
	if err != nil {
		slog.Error("failed to initialize session manager", "error", err)
		return
	}

	// 启动定期清理过期 session
	sessionManager.StartCleanup(time.Hour)

	// 创建路由
	r := server.NewRouter(cfg, sessionManager, logBroadcaster, startTime, distFS)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Port)
	slog.Info("启动 HTTP 服务器", "address", addr)
	if err := r.Run(addr); err != nil {
		slog.Error("http server exited with error", "error", err)
	}
}
