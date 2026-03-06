package server

import (
	"embed"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gorillasessions "github.com/gorilla/sessions"
	sloggin "github.com/samber/slog-gin"

	"main/internal/config"
	"main/internal/handlers"
	"main/internal/middleware"
	"main/internal/session"
	"main/internal/stream"
)

func NewRouter(
	cfg *config.Config,
	logBroadcaster *stream.LogBroadcaster,
	startTime int64,
	distFS embed.FS,
) *gin.Engine {
	authHandler := handlers.NewAuthHandler(cfg.AuthKey, cfg.CookieSecure)
	logsHandler := handlers.NewLogsHandler(logBroadcaster)
	systemHandler := handlers.NewSystemHandler(startTime)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	httpLogConfig := sloggin.DefaultConfig()
	httpLogConfig.WithRequestID = false
	if cfg.DisableStaticAssetLogs {
		httpLogConfig.Filters = append(httpLogConfig.Filters, func(c *gin.Context) bool {
			return !shouldSkipStaticAssetAccessLog(c.Request.URL.Path)
		})
	}
	r.Use(sloggin.NewWithConfig(slog.Default().WithGroup("http"), httpLogConfig))
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions(session.SessionCookieName, newSessionStore(cfg)))

	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.GET("/session", authHandler.Session)
		api.POST("/logout", authHandler.Logout)

		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware(cfg.CookieSecure))
		{
			authenticated.GET("/dashboard/stats", systemHandler.GetStats)
			authenticated.GET("/logs/stream", logsHandler.StreamLogs)
			authenticated.GET("/logs/history", logsHandler.GetHistory)
		}
	}

	r.NoRoute(spaHandler(distFS))
	return r
}

func newSessionStore(cfg *config.Config) sessions.Store {
	sessionDir := filepath.Join(cfg.DataDir, session.SessionDirectoryName)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		panic(err)
	}

	store := gorillasessions.NewFilesystemStore(sessionDir, []byte(cfg.AuthKey))

	wrappedStore := &filesystemSessionStore{FilesystemStore: store}
	wrappedStore.Options(sessions.Options{
		Path:     "/",
		MaxAge:   session.SessionMaxAgeSeconds,
		HttpOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})

	return wrappedStore
}

type filesystemSessionStore struct {
	*gorillasessions.FilesystemStore
}

func (s *filesystemSessionStore) Options(options sessions.Options) {
	s.FilesystemStore.Options = options.ToGorillaOptions()
	s.FilesystemStore.MaxAge(options.MaxAge)
}

var staticAssetLogExtensions = map[string]struct{}{
	".js":    {},
	".css":   {},
	".map":   {},
	".ico":   {},
	".svg":   {},
	".png":   {},
	".jpg":   {},
	".jpeg":  {},
	".gif":   {},
	".webp":  {},
	".avif":  {},
	".woff":  {},
	".woff2": {},
	".ttf":   {},
	".eot":   {},
}

func shouldSkipStaticAssetAccessLog(path string) bool {
	if path == "/" || path == "/index.html" {
		return true
	}

	if strings.HasPrefix(path, "/assets/") {
		return true
	}

	ext := strings.ToLower(filepath.Ext(path))
	if _, ok := staticAssetLogExtensions[ext]; ok {
		return true
	}

	return false
}
