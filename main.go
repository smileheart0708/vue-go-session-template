package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"main/internal/middleware"
	"main/internal/stream"
)

//go:embed web/dist
var distFS embed.FS

// spaHandler 实现 SPA 路由回退逻辑
func spaHandler(distFS embed.FS) gin.HandlerFunc {
	// 从 embed.FS 中提取 web/dist 子目录
	subFS, err := fs.Sub(distFS, "web/dist")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(subFS))

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 尝试打开请求的文件
		f, err := subFS.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			// 文件存在,检查是否为目录
			stat, err := f.Stat()
			f.Close()
			if err == nil && !stat.IsDir() {
				// 是文件,直接提供服务
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		// 文件不存在或是目录,检查是否为静态资源请求
		ext := filepath.Ext(path)
		if ext != "" && ext != ".html" {
			// 有扩展名但不是 HTML,可能是缺失的静态资源,返回 404
			c.Status(http.StatusNotFound)
			return
		}

		// 其他情况返回 index.html,让前端路由处理
		c.Request.URL.Path = "/"
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	logBroadcaster := stream.NewLogBroadcaster()
	logger := middleware.InitLogger("info", logBroadcaster)
	slog.SetDefault(logger)

	// 创建 gin 路由
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// 使用自定义的 SPA handler
	r.NoRoute(spaHandler(distFS))

	// 启动服务器
	r.Run(":8080")
}
