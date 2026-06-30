package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/joe/open-stream-m3u/internal/parser"
)

type DirectProvider struct {
	m3uURL string
	epgURL string
	client *http.Client
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

func (p *DirectProvider) FetchChannels(ctx context.Context) ([]parser.MediaItem, error) {
	items, err := p.fetchAndParse(ctx)
	if err != nil {
		return nil, err
	}

	var channels []parser.MediaItem
	for _, item := range items {
		if item.Type == parser.TypeTV {
			channels = append(channels, item)
		}
	}
	return channels, nil
}

func (p *DirectProvider) FetchMovies(ctx context.Context) ([]parser.MediaItem, error) {
	items, err := p.fetchAndParse(ctx)
	if err != nil {
		return nil, err
	}

	var movies []parser.MediaItem
	for _, item := range items {
		if item.Type == parser.TypeMovie {
			movies = append(movies, item)
		}
	}
	return movies, nil
}

func (p *DirectProvider) FetchSeries(ctx context.Context) ([]parser.MediaItem, error) {
	items, err := p.fetchAndParse(ctx)
	if err != nil {
		return nil, err
	}

	_, seriesMap := parser.GroupSeries(items)
	var series []parser.MediaItem
	for _, s := range seriesMap {
		series = append(series, *s)
	}
	return series, nil
}

func (p *DirectProvider) FetchSeriesInfo(ctx context.Context, seriesID string) ([]parser.Episode, error) {
	items, err := p.fetchAndParse(ctx)
	if err != nil {
		return nil, err
	}

	episodesMap, _ := parser.GroupSeries(items)
	return episodesMap[seriesID], nil
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
