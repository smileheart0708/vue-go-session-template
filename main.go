package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"main/configs"
	"main/internal/handlers"
	"main/internal/middleware"
	"main/internal/session"
	"main/internal/stream"
)

//go:embed web/dist
var distFS embed.FS

// spaHandler 实现支持预压缩文件的 SPA 路由逻辑
func spaHandler(distFS embed.FS) gin.HandlerFunc {
	subFS, err := fs.Sub(distFS, "web/dist")
	if err != nil {
		panic(err)
	}

	// 压缩格式优先级：br > gzip
	compressionFormats := []struct {
		encoding string
		ext      string
	}{
		{"br", ".br"},
		{"gzip", ".gz"},
	}

	return func(c *gin.Context) {
		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		ae := c.GetHeader("Accept-Encoding")

		// 按优先级尝试压缩文件
		for _, cf := range compressionFormats {
			if strings.Contains(ae, cf.encoding) && tryServeFile(c, subFS, path+cf.ext, cf.encoding) {
				return
			}
		}

		// 如果没有匹配的压缩格式，返回 gzip 版本作为兜底
		if tryServeFile(c, subFS, path+".gz", "gzip") {
			return
		}

		// SPA Fallback
		if filepath.Ext(path) == "" || filepath.Ext(path) == ".html" {
			// 对 index.html 尝试压缩版本
			for _, cf := range compressionFormats {
				if strings.Contains(ae, cf.encoding) && tryServeFile(c, subFS, "index.html"+cf.ext, cf.encoding) {
					return
				}
			}
			// 最终回退到 gzip 版本的 index.html
			if tryServeFile(c, subFS, "index.html.gz", "gzip") {
				return
			}
		}

		c.Status(http.StatusNotFound)
	}
}

// tryServeFile 辅助函数：尝试发送压缩文件
func tryServeFile(c *gin.Context, subFS fs.FS, filePath string, encoding string) bool {
	f, err := subFS.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return false
	}

	// 优化：直接根据原始文件名判断 MIME
	// 根据编码格式移除对应的后缀
	extMap := map[string]string{
		"br":   ".br",
		"gzip": ".gz",
	}
	if suffix, ok := extMap[encoding]; ok {
		filePath = strings.TrimSuffix(filePath, suffix)
	}

	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 设置基础响应头
	c.Writer.Header().Set("Content-Type", contentType)
	c.Writer.Header().Set("Content-Encoding", encoding)
	c.Writer.Header().Set("Vary", "Accept-Encoding")

	// 针对带 Hash 的静态资源开启强缓存 (Vite 默认生成带 hash 的文件名)
	if strings.Contains(filePath, "/assets/") {
		c.Writer.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		c.Writer.Header().Set("Cache-Control", "no-cache") // index.html 等不带 hash 的建议不缓存
	}

	// 使用 http.ServeContent，它会自动处理 ETag (基于 ModTime 和 Size)
	// 这样可以避免手动 io.ReadAll 导致的内存拷贝
	http.ServeContent(c.Writer, c.Request, filePath, stat.ModTime(), f.(io.ReadSeeker))
	return true
}

func main() {
	// 加载配置
	cfg := configs.Load()

	// 初始化日志系统
	logBroadcaster := stream.NewLogBroadcaster()
	logger := middleware.InitLogger(cfg.LogLevel, logBroadcaster)
	slog.SetDefault(logger)

	// 输出配置信息
	slog.Info("应用配置已加载",
		"port", cfg.Port,
		"data_dir", cfg.DataDir,
		"log_level", cfg.LogLevel,
		"auth_key", cfg.AuthKey,
	)

	// 初始化 session 管理器
	sessionManager, err := session.NewManager(cfg.DataDir)
	if err != nil {
		slog.Error("failed to initialize session manager", "error", err)
		return
	}

	// 定期清理过期 session
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			sessionManager.CleanExpiredSessions()
		}
	}()

	// 创建认证处理器
	authHandler := handlers.NewAuthHandler(cfg.AuthKey, sessionManager)

	// 创建 gin 路由
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// API 路由
	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/validate-session", authHandler.ValidateSession)
		api.POST("/logout", authHandler.Logout)
	}

	// 使用自定义的 SPA handler
	r.NoRoute(spaHandler(distFS))

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Port)
	slog.Info("启动 HTTP 服务器", "address", addr)
	r.Run(addr)
}
