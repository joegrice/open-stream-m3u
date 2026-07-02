package server

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/joe/open-stream-m3u/internal/addon"
	"github.com/joe/open-stream-m3u/internal/config"
	"github.com/joe/open-stream-m3u/internal/crypto"
	"github.com/joe/open-stream-m3u/internal/parser"
)

type staticProvider struct{}

func (staticProvider) FetchChannels(context.Context) ([]parser.MediaItem, error) {
	return []parser.MediaItem{
		{ID: "iptv_test", Name: "Test Channel", Type: parser.TypeTV, URL: "http://example.com/stream"},
	}, nil
}

func (staticProvider) FetchMovies(context.Context) ([]parser.MediaItem, error) { return nil, nil }
func (staticProvider) FetchSeries(context.Context) ([]parser.MediaItem, error) { return nil, nil }
func (staticProvider) FetchSeriesInfo(context.Context, string) ([]parser.Episode, error) {
	return nil, nil
}
func (staticProvider) FetchEPG(context.Context) (map[string][]parser.Programme, error) {
	return nil, nil
}

func TestTokenRoutes(t *testing.T) {
	cfg := &config.Config{Port: 7001}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	s := New(cfg, logger, fstest.MapFS{})

	addonCfg := map[string]any{"provider": "direct"}
	token, err := crypto.EncodeToken(addonCfg)
	if err != nil {
		t.Fatalf("encode token: %v", err)
	}

	instance := addon.NewInstance(addonCfg, staticProvider{}, logger)
	if err := instance.Initialize(context.Background()); err != nil {
		t.Fatalf("initialize instance: %v", err)
	}
	s.cache.Set(token, instance)

	tests := []struct {
		suffix string
		want   int
	}{
		{"/manifest.json", http.StatusOK},
		{"/catalog/tv/iptv_channels.json", http.StatusOK},
		{"/stream/tv/iptv_test.json", http.StatusOK},
		{"/meta/tv/iptv_test.json", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.suffix, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/"+token+tt.suffix, nil)
			rec := httptest.NewRecorder()
			s.mux.ServeHTTP(rec, req)
			if rec.Code != tt.want {
				t.Errorf("%s: got status %d, want %d", tt.suffix, rec.Code, tt.want)
			}
		})
	}
}
