package server

import (
	"embed"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gorillasessions "github.com/gorilla/sessions"
	sloggin "github.com/samber/slog-gin"

	"main/internal/config"
	"main/internal/handlers"
	"main/internal/middleware"
	"main/internal/stream"
)

const sessionDuration = 7 * 24 * time.Hour

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
	r.Use(sloggin.NewWithConfig(slog.Default().WithGroup("http"), httpLogConfig))
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions(cfg.SessionName, newSessionStore(cfg)))

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
	sessionDir := filepath.Join(cfg.DataDir, "sessions")
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		panic(err)
	}

	var store *gorillasessions.FilesystemStore
	if cfg.SessionEncKey == "" {
		store = gorillasessions.NewFilesystemStore(sessionDir, []byte(cfg.SessionAuthKey))
	} else {
		store = gorillasessions.NewFilesystemStore(
			sessionDir,
			[]byte(cfg.SessionAuthKey),
			[]byte(cfg.SessionEncKey),
		)
	}

	wrappedStore := &filesystemSessionStore{FilesystemStore: store}
	wrappedStore.Options(sessions.Options{
		Path:     "/",
		MaxAge:   int(sessionDuration.Seconds()),
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
