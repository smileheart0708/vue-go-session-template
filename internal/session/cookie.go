package session

import (
	"net/http"

	ginsessions "github.com/gin-contrib/sessions"
)

// SetCookieOptions applies a consistent cookie policy for session operations.
func SetCookieOptions(sess ginsessions.Session, secure bool, maxAge int) {
	sess.Options(ginsessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// ExpireCookie clears session values and expires the session cookie immediately.
func ExpireCookie(sess ginsessions.Session, secure bool) {
	sess.Clear()
	SetCookieOptions(sess, secure, -1)
}
