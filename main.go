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

	"github.com/gin-gonic/gin"

	"main/configs"
	"main/internal/middleware"
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
	fileServer := http.FileServer(http.FS(subFS))

	return func(c *gin.Context) {
		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		ae := c.GetHeader("Accept-Encoding")

		// 2026 优先级：zstd > br > gzip
		served := false
		if strings.Contains(ae, "zstd") && tryServeFile(c, subFS, path+".zst", "zstd") {
			served = true
		} else if strings.Contains(ae, "br") && tryServeFile(c, subFS, path+".br", "br") {
			served = true
		} else if strings.Contains(ae, "gzip") && tryServeFile(c, subFS, path+".gz", "gzip") {
			served = true
		}

		if served {
			return
		}

		// 原始文件兜底
		f, err := subFS.Open(path)
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// SPA Fallback
		if filepath.Ext(path) == "" || filepath.Ext(path) == ".html" {
			// 同样对 index.html 尝试压缩版本
			if strings.Contains(ae, "zstd") && tryServeFile(c, subFS, "index.html.zst", "zstd") {
				return
			}
			if strings.Contains(ae, "br") && tryServeFile(c, subFS, "index.html.br", "br") {
				return
			}
			// 最终回退到原始 index.html
			c.Request.URL.Path = "/"
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
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
	// 需要根据不同的编码格式移除对应的后缀
	var originalPath string
	switch encoding {
	case "zstd":
		originalPath = strings.TrimSuffix(filePath, ".zst")
	case "br":
		originalPath = strings.TrimSuffix(filePath, ".br")
	case "gzip":
		originalPath = strings.TrimSuffix(filePath, ".gz")
	default:
		originalPath = filePath
	}

	contentType := mime.TypeByExtension(filepath.Ext(originalPath))
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

	// 创建 gin 路由
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// 使用自定义的 SPA handler
	r.NoRoute(spaHandler(distFS))

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Port)
	slog.Info("启动 HTTP 服务器", "address", addr)
	r.Run(addr)
}
