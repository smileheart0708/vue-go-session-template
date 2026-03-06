package server

import "testing"

func TestShouldSkipStaticAssetAccessLog(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "root page", path: "/", want: true},
		{name: "index html", path: "/index.html", want: true},
		{name: "vite asset", path: "/assets/index-N8bieKVI.js", want: true},
		{name: "favicon", path: "/favicon.ico", want: true},
		{name: "font", path: "/fonts/inter.woff2", want: true},
		{name: "api request", path: "/api/session", want: false},
		{name: "spa route", path: "/settings/security", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSkipStaticAssetAccessLog(tt.path)
			if got != tt.want {
				t.Fatalf("path %q: got %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
