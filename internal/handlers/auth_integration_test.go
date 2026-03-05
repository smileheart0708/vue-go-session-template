package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	"main/internal/handlers"
	"main/internal/middleware"
)

const testSessionMaxAgeSeconds = 7 * 24 * 60 * 60

type loginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type sessionStatusResponse struct {
	Authenticated bool   `json:"authenticated"`
	Message       string `json:"message,omitempty"`
}

func newAuthTestRouter(authKey string) *gin.Engine {
	gin.SetMode(gin.TestMode)

	store := memstore.NewStore([]byte("0123456789abcdef0123456789abcdef"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   testSessionMaxAgeSeconds,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	authHandler := handlers.NewAuthHandler(authKey, false)

	router := gin.New()
	router.Use(sessions.Sessions("session_id", store))

	api := router.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.GET("/session", authHandler.Session)
		api.POST("/logout", authHandler.Logout)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(false))
		protected.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})
	}

	return router
}

func performRequest(
	router *gin.Engine,
	method string,
	path string,
	body []byte,
	cookies ...*http.Cookie,
) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	for _, cookie := range cookies {
		if cookie != nil {
			req.AddCookie(cookie)
		}
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func findCookieByName(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie != nil && cookie.Name == name {
			return cookie
		}
	}
	return nil
}

func TestLoginSetsSessionCookie(t *testing.T) {
	router := newAuthTestRouter("top-secret-auth-key")
	recorder := performRequest(router, http.MethodPost, "/api/login", []byte(`{"auth_key":"top-secret-auth-key"}`))
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response loginResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	if !response.Success {
		t.Fatalf("expected login success response, got %+v", response)
	}

	setCookie := recorder.Header().Get("Set-Cookie")
	if setCookie == "" {
		t.Fatal("expected Set-Cookie header")
	}

	expectedParts := []string{"session_id=", "HttpOnly", "SameSite=Lax"}
	for _, part := range expectedParts {
		if !strings.Contains(setCookie, part) {
			t.Fatalf("expected Set-Cookie to contain %q, got %q", part, setCookie)
		}
	}
}

func TestProtectedEndpointRequiresSession(t *testing.T) {
	router := newAuthTestRouter("top-secret-auth-key")
	recorder := performRequest(router, http.MethodGet, "/api/protected", nil)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestLoginAllowsProtectedAndSessionEndpoints(t *testing.T) {
	router := newAuthTestRouter("top-secret-auth-key")

	loginRecorder := performRequest(
		router,
		http.MethodPost,
		"/api/login",
		[]byte(`{"auth_key":"top-secret-auth-key"}`),
	)
	if loginRecorder.Code != http.StatusOK {
		t.Fatalf("expected login status 200, got %d", loginRecorder.Code)
	}

	sessionCookie := findCookieByName(loginRecorder.Result().Cookies(), "session_id")
	if sessionCookie == nil {
		t.Fatal("expected session_id cookie in login response")
	}

	protectedRecorder := performRequest(
		router,
		http.MethodGet,
		"/api/protected",
		nil,
		sessionCookie,
	)
	if protectedRecorder.Code != http.StatusOK {
		t.Fatalf("expected protected endpoint status 200, got %d", protectedRecorder.Code)
	}

	sessionRecorder := performRequest(
		router,
		http.MethodGet,
		"/api/session",
		nil,
		sessionCookie,
	)
	if sessionRecorder.Code != http.StatusOK {
		t.Fatalf("expected session endpoint status 200, got %d", sessionRecorder.Code)
	}

	var sessionResponse sessionStatusResponse
	if err := json.Unmarshal(sessionRecorder.Body.Bytes(), &sessionResponse); err != nil {
		t.Fatalf("failed to parse session response: %v", err)
	}
	if !sessionResponse.Authenticated {
		t.Fatalf("expected authenticated session response, got %+v", sessionResponse)
	}
}

func TestSessionEndpointReturnsUnauthorizedWithoutCookie(t *testing.T) {
	router := newAuthTestRouter("top-secret-auth-key")
	recorder := performRequest(router, http.MethodGet, "/api/session", nil)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}

	setCookie := recorder.Header().Get("Set-Cookie")
	if !strings.Contains(setCookie, "session_id=") {
		t.Fatalf("expected unauthorized session response to clear cookie, got %q", setCookie)
	}
}

func TestProtectedEndpointClearsCookieWhenUnauthorized(t *testing.T) {
	router := newAuthTestRouter("top-secret-auth-key")
	recorder := performRequest(router, http.MethodGet, "/api/protected", nil)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}

	setCookie := recorder.Header().Get("Set-Cookie")
	if !strings.Contains(setCookie, "session_id=") {
		t.Fatalf("expected unauthorized protected response to clear cookie, got %q", setCookie)
	}
}

func TestLogoutInvalidatesSession(t *testing.T) {
	router := newAuthTestRouter("top-secret-auth-key")

	loginRecorder := performRequest(
		router,
		http.MethodPost,
		"/api/login",
		[]byte(`{"auth_key":"top-secret-auth-key"}`),
	)
	if loginRecorder.Code != http.StatusOK {
		t.Fatalf("expected login status 200, got %d", loginRecorder.Code)
	}

	sessionCookie := findCookieByName(loginRecorder.Result().Cookies(), "session_id")
	if sessionCookie == nil {
		t.Fatal("expected session_id cookie in login response")
	}

	logoutRecorder := performRequest(
		router,
		http.MethodPost,
		"/api/logout",
		nil,
		sessionCookie,
	)
	if logoutRecorder.Code != http.StatusOK {
		t.Fatalf("expected logout status 200, got %d", logoutRecorder.Code)
	}

	protectedRecorder := performRequest(
		router,
		http.MethodGet,
		"/api/protected",
		nil,
		sessionCookie,
	)
	if protectedRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected protected endpoint status 401 after logout, got %d", protectedRecorder.Code)
	}
}
