package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/joe/open-stream-m3u/internal/parser"
)

type DirectProvider struct {
	m3uURL string
	epgURL string
	client *http.Client

	mu             sync.Mutex
	cachedItems    []parser.MediaItem
	cachedEpisodes map[string][]parser.Episode
}

func NewDirectProvider(m3uURL, epgURL string) *DirectProvider {
	return &DirectProvider{
		m3uURL: m3uURL,
		epgURL: epgURL,
		client: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

func (p *DirectProvider) FetchAll(ctx context.Context) (channels, movies, series []parser.MediaItem, err error) {
	items, err := p.fetchAndParse(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	var eps map[string][]parser.Episode
	var seriesMap map[string]*parser.MediaItem
	eps, seriesMap = parser.GroupSeries(items)
	for _, s := range seriesMap {
		series = append(series, *s)
	}
	for _, item := range items {
		switch item.Type {
		case parser.TypeTV:
			channels = append(channels, item)
		case parser.TypeMovie:
			movies = append(movies, item)
		}
	}

	// ponytail: struct cache valid for one refresh cycle, overwritten by next FetchAll.
	p.mu.Lock()
	p.cachedItems = items
	p.cachedEpisodes = eps
	p.mu.Unlock()

	return channels, movies, series, nil
}

func (p *DirectProvider) FetchSeriesInfo(ctx context.Context, seriesID string) ([]parser.Episode, error) {
	p.mu.Lock()
	eps := p.cachedEpisodes
	p.mu.Unlock()
	if eps != nil {
		return eps[seriesID], nil
	}

	// Cold path: cache miss (e.g. /api/groups probe path that never called FetchAll).
	// Populate lazily so future calls are O(1).
	items, err := p.fetchAndParse(ctx)
	if err != nil {
		return nil, err
	}
	eMap, _ := parser.GroupSeries(items)
	p.mu.Lock()
	p.cachedItems = items
	p.cachedEpisodes = eMap
	p.mu.Unlock()
	return eMap[seriesID], nil
}

func (p *DirectProvider) FetchEPG(ctx context.Context) (map[string][]parser.Programme, error) {
	if p.epgURL == "" {
		return nil, nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", p.epgURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "open-stream-m3u/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("EPG fetch failed: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parser.ParseXMLTV(string(body))
}

func (p *DirectProvider) fetchAndParse(ctx context.Context) ([]parser.MediaItem, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.m3uURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "open-stream-m3u/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("M3U fetch failed: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parser.ParseM3U(string(body))
}