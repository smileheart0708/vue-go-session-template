package server

import (
	"embed"

	"github.com/gin-gonic/gin"

	"main/internal/config"
	"main/internal/handlers"
	"main/internal/middleware"
	"main/internal/session"
	"main/internal/stream"
)

func NewRouter(
	cfg *config.Config,
	sessionManager *session.Manager,
	logBroadcaster *stream.LogBroadcaster,
	startTime int64,
	distFS embed.FS,
) *gin.Engine {
	authHandler := handlers.NewAuthHandler(cfg.AuthKey, sessionManager)
	logsHandler := handlers.NewLogsHandler(logBroadcaster)
	systemHandler := handlers.NewSystemHandler(startTime)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/validate-session", authHandler.ValidateSession)
		api.POST("/logout", authHandler.Logout)

		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware(sessionManager))
		{
			authenticated.GET("/dashboard/stats", systemHandler.GetStats)
			authenticated.GET("/logs/stream", logsHandler.StreamLogs)
			authenticated.GET("/logs/history", logsHandler.GetHistory)
		}
	}

	r.NoRoute(spaHandler(distFS))
	return r
}
