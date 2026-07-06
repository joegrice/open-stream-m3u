package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/joe/open-stream-m3u/internal/addon"
	"github.com/joe/open-stream-m3u/internal/config"
	"github.com/joe/open-stream-m3u/internal/crypto"
	"github.com/joe/open-stream-m3u/internal/parser"
	"github.com/joe/open-stream-m3u/internal/provider"
)

type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	cache  *addon.InstanceCache
	mux    *http.ServeMux
	webFS  fs.FS
}

func New(cfg *config.Config, logger *slog.Logger, webFS fs.FS) *Server {
	s := &Server{
		cfg:    cfg,
		logger: logger,
		cache:  addon.NewInstanceCache(cfg.MaxCacheEntries, cfg.CacheTTL),
		mux:    http.NewServeMux(),
		webFS:  webFS,
	}
	s.setupRoutes()

	// Background sweep: evict expired instance cache entries without a re-touch.
	// ponytail: fixed TTL/2 cadence; tighten if quiet tokens linger in memory.
	if cfg.CacheTTL > 0 {
		go func() {
			interval := cfg.CacheTTL / 2
			if interval < time.Minute {
				interval = time.Minute
			}
			t := time.NewTicker(interval)
			defer t.Stop()
			for range t.C {
				s.cache.Sweep()
			}
		}()
	}

	return s
}

func (s *Server) setupRoutes() {
	fileServer := http.FileServer(http.FS(s.webFS))

	s.mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, s.webFS, "index.html")
	})
	s.mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	s.mux.HandleFunc("GET /configure", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, s.webFS, "configure.html")
	})
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("POST /api/prefetch", s.handlePrefetch)
	s.mux.HandleFunc("POST /api/groups", s.handleGroups)
	s.mux.HandleFunc("POST /api/encrypt", s.handleEncrypt)
	s.mux.HandleFunc("GET /api/info", s.handleInfo)
	s.mux.HandleFunc("GET /api/debug", s.handleDebug)

	s.mux.HandleFunc("GET /css/", fileServer.ServeHTTP)
	s.mux.HandleFunc("GET /js/", fileServer.ServeHTTP)

	s.mux.HandleFunc("GET /{token}/{path...}", s.handleTokenRoute)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
	})
}

func (s *Server) handlePrefetch(w http.ResponseWriter, r *http.Request) {
	if !s.cfg.PrefetchEnabled {
		http.Error(w, "Prefetch disabled", http.StatusForbidden)
		return
	}

	var req struct {
		URL     string `json:"url"`
		Purpose string `json:"purpose"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.URL == "" || !strings.HasPrefix(req.URL, "http") {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	if isBlockedHost(req.URL) {
		http.Error(w, "Blocked host", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "GET", req.URL, nil)
	if err != nil {
		http.Error(w, "Request error", http.StatusInternalServerError)
		return
	}
	httpReq.Header.Set("User-Agent", "open-stream-m3u/1.0")

	client := &http.Client{Timeout: 45 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		http.Error(w, "Fetch failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Fetch failed: %d", resp.StatusCode), http.StatusBadGateway)
		return
	}

	limited := io.LimitReader(resp.Body, s.cfg.PrefetchMaxSize)
	body, err := io.ReadAll(limited)
	if err != nil {
		http.Error(w, "Read error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"ok":      true,
		"bytes":   len(body),
		"content": string(body),
	})
}

func (s *Server) handleGroups(w http.ResponseWriter, r *http.Request) {
	var cfg map[string]any
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	prov, err := s.createProvider(cfg)
	if err != nil {
		http.Error(w, "Invalid provider config", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	enabledTypes := addon.EnabledTypesFromConfig(cfg)

	channels, movies, series, _ := prov.FetchAll(ctx)
	// ponytail: handleGroups ignores fetch errors intentionally to return
	// partial group lists — keep that behavior with the single FetchAll call.
	if !enabledTypes["tv"] {
		channels = nil
	}
	if !enabledTypes["movie"] {
		movies = nil
	}
	if !enabledTypes["series"] {
		series = nil
	}

	all := make([]parser.MediaItem, 0, len(channels)+len(movies)+len(series))
	all = append(all, channels...)
	all = append(all, movies...)
	all = append(all, series...)

	type groupInfo struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
		Type  string `json:"type"`
	}

	groupMap := make(map[string]*groupInfo)
	for _, item := range all {
		if item.Group == "" {
			continue
		}
		g, ok := groupMap[item.Group]
		if !ok {
			g = &groupInfo{Name: item.Group, Type: string(item.Type)}
			groupMap[item.Group] = g
		}
		g.Count++
	}

	groups := make([]groupInfo, 0, len(groupMap))
	for _, g := range groupMap {
		groups = append(groups, *g)
	}
	sort.Slice(groups, func(i, j int) bool {
		if groups[i].Count != groups[j].Count {
			return groups[i].Count > groups[j].Count
		}
		return groups[i].Name < groups[j].Name
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func (s *Server) handleEncrypt(w http.ResponseWriter, r *http.Request) {
	if s.cfg.ConfigSecret == "" {
		http.Error(w, "Encryption not configured", http.StatusBadRequest)
		return
	}

	var cfg map[string]any
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "Invalid config", http.StatusBadRequest)
		return
	}

	token, err := crypto.EncryptConfig(cfg, s.cfg.ConfigSecret)
	if err != nil {
		http.Error(w, "Encryption failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	instance, err := s.getOrBuildInstance(token)
	if err != nil {
		s.logger.Error("Failed to build instance for info", "error", err)
		http.Error(w, "Invalid configuration", http.StatusBadRequest)
		return
	}

	stats := instance.GetStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) handleDebug(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	instance, err := s.getOrBuildInstance(token)
	if err != nil {
		s.logger.Error("Failed to build instance for debug", "error", err)
		http.Error(w, "Invalid configuration", http.StatusBadRequest)
		return
	}

	stats := instance.GetStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"stats":  stats,
		"config": instance.GetConfig(),
	})
}

func (s *Server) handleTokenRoute(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	// URL decode the token (in case it contains encoded characters)
	decodedToken, err := url.PathUnescape(token)
	if err != nil {
		s.logger.Error("Failed to decode token", "error", err)
		http.Error(w, "Invalid token encoding", http.StatusBadRequest)
		return
	}

	instance, err := s.getOrBuildInstance(decodedToken)
	if err != nil {
		s.logger.Error("Failed to build instance", "error", err)
		http.Error(w, "Invalid configuration", http.StatusBadRequest)
		return
	}

	path := r.URL.Path
	afterToken := strings.TrimPrefix(path, "/"+token)

	switch {
	case afterToken == "/manifest.json":
		addon.HandleManifest(instance)(w, r)
	case strings.HasPrefix(afterToken, "/catalog/"):
		addon.HandleCatalog(instance)(w, r)
	case strings.HasPrefix(afterToken, "/stream/"):
		addon.HandleStream(instance)(w, r)
	case strings.HasPrefix(afterToken, "/meta/"):
		addon.HandleMeta(instance)(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) getOrBuildInstance(token string) (*addon.Instance, error) {
	if instance, ok := s.cache.Get(token); ok {
		return instance, nil
	}

	cfg, err := s.parseToken(token)
	if err != nil {
		return nil, err
	}

	prov, err := s.createProvider(cfg)
	if err != nil {
		return nil, err
	}

	instance := addon.NewInstance(cfg, prov, s.logger)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := instance.Initialize(ctx); err != nil {
		return nil, err
	}

	s.cache.Set(token, instance)
	return instance, nil
}

func (s *Server) parseToken(token string) (map[string]any, error) {
	if strings.HasPrefix(token, "enc:") {
		if s.cfg.ConfigSecret == "" {
			return nil, fmt.Errorf("encryption not configured")
		}
		return crypto.DecryptConfig(token, s.cfg.ConfigSecret)
	}
	return crypto.DecodeToken(token)
}

func (s *Server) createProvider(cfg map[string]any) (provider.Provider, error) {
	providerType, _ := cfg["provider"].(string)

	switch providerType {
	case "xtream":
		baseURL, _ := cfg["xtreamUrl"].(string)
		username, _ := cfg["xtreamUsername"].(string)
		password, _ := cfg["xtreamPassword"].(string)
		useM3U, _ := cfg["xtreamUseM3U"].(bool)
		return provider.NewXtreamProvider(baseURL, username, password, useM3U), nil

	case "direct":
		m3uURL, _ := cfg["m3uUrl"].(string)
		epgURL, _ := cfg["epgUrl"].(string)
		return provider.NewDirectProvider(m3uURL, epgURL), nil

	default:
		return nil, fmt.Errorf("unknown provider type: %s", providerType)
	}
}

func isBlockedHost(rawURL string) bool {
	host := extractHost(rawURL)
	if host == "" {
		return true
	}

	if host == "localhost" || host == "0.0.0.0" || host == "::1" ||
		strings.HasPrefix(host, "127.") || strings.HasPrefix(host, "10.") ||
		strings.HasPrefix(host, "192.168.") || strings.HasPrefix(host, "172.16.") {
		return true
	}

	return false
}

func extractHost(rawURL string) string {
	if !strings.Contains(rawURL, "://") {
		return ""
	}
	parts := strings.SplitN(rawURL, "://", 2)
	if len(parts) < 2 {
		return ""
	}
	host := strings.Split(parts[1], "/")[0]
	host = strings.Split(host, ":")[0]
	return host
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.cfg.Port)
	s.logger.Info("Starting server", "addr", addr, "debug", s.cfg.Debug)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Handler:      s.loggingMiddleware(s.mux),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return httpServer.Serve(listener)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		if s.cfg.Debug {
			s.logger.Info("Request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
			)
		}
	})
}
