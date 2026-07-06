package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/joe/open-stream-m3u/internal/parser"
)

type XtreamProvider struct {
	baseURL    string
	username   string
	password   string
	useM3U     bool
	client     *http.Client
	catNames   map[string]string
	catLoaded  bool

	mu             sync.Mutex
	cachedEpisodes map[string][]parser.Episode
}

func NewXtreamProvider(baseURL, username, password string, useM3U bool) *XtreamProvider {
	return &XtreamProvider{
		baseURL:  baseURL,
		username: username,
		password: password,
		useM3U:   useM3U,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *XtreamProvider) apiURL(action string) string {
	return fmt.Sprintf("%s/player_api.php?username=%s&password=%s&action=%s",
		p.baseURL,
		url.QueryEscape(p.username),
		url.QueryEscape(p.password),
		action,
	)
}

func (p *XtreamProvider) loadCategories(ctx context.Context) {
	if p.catLoaded {
		return
	}
	p.catLoaded = true
	p.catNames = make(map[string]string)

	for _, action := range []string{"get_live_categories", "get_vod_categories", "get_series_categories"} {
		req, err := http.NewRequestWithContext(ctx, "GET", p.apiURL(action), nil)
		if err != nil {
			continue
		}
		resp, err := p.client.Do(req)
		if err != nil {
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}
		var cats []struct {
			CategoryID   string `json:"category_id"`
			CategoryName string `json:"category_name"`
		}
		if err := json.Unmarshal(body, &cats); err != nil {
			continue
		}
		for _, c := range cats {
			if c.CategoryID != "" {
				p.catNames[c.CategoryID] = c.CategoryName
			}
		}
	}
}

func (p *XtreamProvider) resolveCategory(catID string) string {
	if name, ok := p.catNames[catID]; ok {
		return name
	}
	return catID
}

func (p *XtreamProvider) FetchAll(ctx context.Context) (channels, movies, series []parser.MediaItem, err error) {
	if p.useM3U {
		items, err := p.fetchM3U(ctx)
		if err != nil {
			return nil, nil, nil, err
		}
		eps, seriesMap := parser.GroupSeries(items)
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
		p.cachedEpisodes = eps
		p.mu.Unlock()
		return channels, movies, series, nil
	}

	channels, err = p.fetchChannelsFromAPI(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	movies, err = p.fetchMoviesFromAPI(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	series, err = p.fetchSeriesFromAPI(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	return channels, movies, series, nil
}

func (p *XtreamProvider) FetchSeriesInfo(ctx context.Context, seriesID string) ([]parser.Episode, error) {
	if p.useM3U {
		p.mu.Lock()
		eps := p.cachedEpisodes
		p.mu.Unlock()
		if eps != nil {
			return eps[seriesID], nil
		}
		return nil, nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", p.apiURL("get_series_info")+"&series_id="+url.QueryEscape(seriesID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info struct {
		Episodes map[string][]struct {
			ID                int    `json:"id"`
			Title             string `json:"title"`
			Season            int    `json:"season"`
			EpisodeNum        int    `json:"episode_num"`
			ContainerExtension string `json:"container_extension"`
			Info              struct {
				MovieImage string `json:"movie_image"`
			} `json:"info"`
		} `json:"episodes"`
	}

	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	var episodes []parser.Episode
	for seasonKey, eps := range info.Episodes {
		for _, ep := range eps {
			epURL := fmt.Sprintf("%s/series/%s/%s/%d.%s",
				p.baseURL,
				url.QueryEscape(p.username),
				url.QueryEscape(p.password),
				ep.ID,
				ep.ContainerExtension,
			)

			episodes = append(episodes, parser.Episode{
				ID:        fmt.Sprintf("iptv_series_ep_%d", ep.ID),
				Title:     ep.Title,
				Season:    ep.Season,
				Episode:   ep.EpisodeNum,
				URL:       epURL,
				Thumbnail: ep.Info.MovieImage,
			})
		}
		_ = seasonKey
	}

	return episodes, nil
}

func (p *XtreamProvider) FetchEPG(ctx context.Context) (map[string][]parser.Programme, error) {
	epgURL := fmt.Sprintf("%s/xmltv.php?username=%s&password=%s",
		p.baseURL,
		url.QueryEscape(p.username),
		url.QueryEscape(p.password),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", epgURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parser.ParseXMLTV(resp.Body)
}

func (p *XtreamProvider) fetchChannelsFromAPI(ctx context.Context) ([]parser.MediaItem, error) {
	p.loadCategories(ctx)

	req, err := http.NewRequestWithContext(ctx, "GET", p.apiURL("get_live_streams"), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var streams []struct {
		StreamID     int    `json:"stream_id"`
		Name         string `json:"name"`
		StreamIcon   string `json:"stream_icon"`
		EPGChannelID string `json:"epg_channel_id"`
		CategoryID   string `json:"category_id"`
	}

	if err := json.Unmarshal(body, &streams); err != nil {
		return nil, err
	}

	var channels []parser.MediaItem
	for _, s := range streams {
		streamURL := fmt.Sprintf("%s/live/%s/%s/%d.m3u8",
			p.baseURL,
			url.QueryEscape(p.username),
			url.QueryEscape(p.password),
			s.StreamID,
		)

		channels = append(channels, parser.MediaItem{
			ID:    fmt.Sprintf("iptv_live_%d", s.StreamID),
			Name:  s.Name,
			URL:   streamURL,
			Type:  parser.TypeTV,
			Logo:  s.StreamIcon,
			EPGID: s.EPGChannelID,
			Group: p.resolveCategory(s.CategoryID),
			Attrs: map[string]string{
				"tvg-logo":   s.StreamIcon,
				"tvg-id":     s.EPGChannelID,
				"group-title": p.resolveCategory(s.CategoryID),
			},
		})
	}

	return channels, nil
}

func (p *XtreamProvider) fetchMoviesFromAPI(ctx context.Context) ([]parser.MediaItem, error) {
	p.loadCategories(ctx)

	req, err := http.NewRequestWithContext(ctx, "GET", p.apiURL("get_vod_streams"), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if response is an error message
	var errorResp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
		return nil, fmt.Errorf("API error: %s", errorResp.Error)
	}

	// Try to parse as array
	var streams []struct {
		StreamID          int    `json:"stream_id"`
		Name              string `json:"name"`
		StreamIcon        string `json:"stream_icon"`
		Plot              string `json:"plot"`
		CategoryID        string `json:"category_id"`
		ContainerExtension string `json:"container_extension"`
	}

	if err := json.Unmarshal(body, &streams); err != nil {
		return nil, fmt.Errorf("failed to parse VOD response: %w", err)
	}

	var movies []parser.MediaItem
	for _, s := range streams {
		streamURL := fmt.Sprintf("%s/movie/%s/%s/%d.%s",
			p.baseURL,
			url.QueryEscape(p.username),
			url.QueryEscape(p.password),
			s.StreamID,
			s.ContainerExtension,
		)

		movies = append(movies, parser.MediaItem{
			ID:    fmt.Sprintf("iptv_vod_%d", s.StreamID),
			Name:  s.Name,
			URL:   streamURL,
			Type:  parser.TypeMovie,
			Logo:  s.StreamIcon,
			Group: p.resolveCategory(s.CategoryID),
			Plot:  s.Plot,
			Attrs: map[string]string{
				"tvg-logo":    s.StreamIcon,
				"group-title": p.resolveCategory(s.CategoryID),
				"plot":        s.Plot,
			},
		})
	}

	return movies, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (p *XtreamProvider) fetchSeriesFromAPI(ctx context.Context) ([]parser.MediaItem, error) {
	p.loadCategories(ctx)

	req, err := http.NewRequestWithContext(ctx, "GET", p.apiURL("get_series"), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var series []struct {
		SeriesID   int    `json:"series_id"`
		Name       string `json:"name"`
		Cover      string `json:"cover"`
		Plot       string `json:"plot"`
		CategoryID string `json:"category_id"`
	}

	if err := json.Unmarshal(body, &series); err != nil {
		return nil, err
	}

	var result []parser.MediaItem
	for _, s := range series {
		result = append(result, parser.MediaItem{
			ID:    fmt.Sprintf("iptv_series_%d", s.SeriesID),
			Name:  s.Name,
			Type:  parser.TypeSeries,
			Logo:  s.Cover,
			Group: p.resolveCategory(s.CategoryID),
			Plot:  s.Plot,
			Attrs: map[string]string{
				"tvg-logo":    s.Cover,
				"group-title": p.resolveCategory(s.CategoryID),
				"plot":        s.Plot,
			},
		})
	}

	return result, nil
}

func (p *XtreamProvider) fetchM3U(ctx context.Context) ([]parser.MediaItem, error) {
	m3uURL := fmt.Sprintf("%s/get.php?username=%s&password=%s&type=m3u_plus",
		p.baseURL,
		url.QueryEscape(p.username),
		url.QueryEscape(p.password),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", m3uURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parser.ParseM3U(string(body))
}
