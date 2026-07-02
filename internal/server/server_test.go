package server

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"

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
	cfg := &config.Config{Port: 7001, CacheTTL: time.Hour, MaxCacheEntries: 10}
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
	if _, ok := s.cache.Get(token); !ok {
		t.Fatalf("instance not in cache before request")
	}

	tests := []struct {
		suffix    string
		want      int
		wantMetas int
	}{
		{"/manifest.json", http.StatusOK, 0},
		{"/catalog/tv/iptv_channels.json", http.StatusOK, 1},
		{"/stream/tv/iptv_test.json", http.StatusOK, 0},
		{"/meta/tv/iptv_test.json", http.StatusOK, 0},
	}

	for _, tt := range tests {
		t.Run(tt.suffix, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/"+token+tt.suffix, nil)
			rec := httptest.NewRecorder()
			s.mux.ServeHTTP(rec, req)
			if rec.Code != tt.want {
				t.Errorf("%s: got status %d, want %d", tt.suffix, rec.Code, tt.want)
			}
			if tt.wantMetas > 0 {
				var resp struct {
					Metas []json.RawMessage `json:"metas"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Fatalf("%s: decode response: %v", tt.suffix, err)
				}
				if len(resp.Metas) != tt.wantMetas {
					t.Errorf("%s: got %d metas, want %d", tt.suffix, len(resp.Metas), tt.wantMetas)
				}
			}
		})
	}
}
