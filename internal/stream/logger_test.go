package stream

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
)

func TestBroadcastHandlerSerializesErrorAttr(t *testing.T) {
	broadcaster := NewLogBroadcaster()
	logger := slog.New(NewBroadcastHandler(broadcaster))

	wantErr := "securecookie: the value is not valid"
	logger.Error("[sessions] ERROR!", "err", errors.New(wantErr))

	history := broadcaster.GetHistory()
	if len(history) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(history))
	}

	gotErrValue, ok := history[0].Attrs["err"].(string)
	if !ok {
		t.Fatalf("expected err attr to be string, got %T", history[0].Attrs["err"])
	}

	if gotErrValue != wantErr {
		t.Fatalf("expected err=%q, got %q", wantErr, gotErrValue)
	}
}

func TestBroadcastHandlerWithAttrsAndGroups(t *testing.T) {
	broadcaster := NewLogBroadcaster()
	handler := NewBroadcastHandler(broadcaster)

	logger := slog.New(handler).
		With("service", "api").
		WithGroup("http")

	logger.LogAttrs(context.Background(), slog.LevelInfo, "request completed",
		slog.Int("status", 200),
		slog.Duration("latency", 150*time.Millisecond),
		slog.Group("request", slog.String("method", "GET")),
	)

	history := broadcaster.GetHistory()
	if len(history) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(history))
	}

	attrs := history[0].Attrs
	service, ok := attrs["service"].(string)
	if !ok || service != "api" {
		t.Fatalf("expected service attr to be %q, got %#v", "api", attrs["service"])
	}

	httpGroup, ok := attrs["http"].(map[string]any)
	if !ok {
		t.Fatalf("expected http group to be map[string]any, got %T", attrs["http"])
	}

	status, ok := httpGroup["status"].(int64)
	if !ok || status != 200 {
		t.Fatalf("expected http.status=200, got %#v (%T)", httpGroup["status"], httpGroup["status"])
	}

	latency, ok := httpGroup["latency"].(string)
	if !ok || latency != "150ms" {
		t.Fatalf("expected http.latency=%q, got %#v", "150ms", httpGroup["latency"])
	}

	requestGroup, ok := httpGroup["request"].(map[string]any)
	if !ok {
		t.Fatalf(
			"expected http.request group to be map[string]any, got %T",
			httpGroup["request"],
		)
	}

	method, ok := requestGroup["method"].(string)
	if !ok || method != "GET" {
		t.Fatalf("expected http.request.method=%q, got %#v", "GET", requestGroup["method"])
	}
}
