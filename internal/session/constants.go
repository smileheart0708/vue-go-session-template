package session

import "time"

const (
	SessionCookieName    = "session_id"
	SessionDirectoryName = "sessions"
	SessionFilePrefix    = "session_"
)

const (
	SessionTTL           = 7 * 24 * time.Hour
	SessionMaxAgeSeconds = int(SessionTTL / time.Second)
	CleanupInterval      = 30 * time.Minute
	CleanupGrace         = 10 * time.Minute
)
