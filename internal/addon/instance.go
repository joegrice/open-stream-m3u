package addon

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/joe/open-stream-m3u/internal/cache"
	"github.com/joe/open-stream-m3u/internal/parser"
	"github.com/joe/open-stream-m3u/internal/provider"
)

type Instance struct {
	mu            sync.RWMutex
	config        map[string]any
	provider      provider.Provider
	channels      []parser.MediaItem
	movies        []parser.MediaItem
	series        []parser.MediaItem
	episodes      map[string][]parser.Episode
	itemIndex     map[string]*parser.MediaItem
	episodeIndex  map[string]*parser.Episode
	epgData       map[string][]parser.Programme
	lastUpdate    time.Time
	manifest      *Manifest
	groupCatalogs map[string]string // catalog ID -> group name
	enabledTypes  map[string]bool
	logger        *slog.Logger
}

func NewInstance(config map[string]any, prov provider.Provider, logger *slog.Logger) *Instance {
	selectedGroups := extractSelectedGroups(config)
	enabledTypes := EnabledTypesFromConfig(config)
	manifest := BuildManifest(selectedGroups, enabledTypes)

	groupCatalogs := make(map[string]string)
	for _, g := range selectedGroups {
		groupCatalogs[groupCatalogID(g)] = g
	}

	return &Instance{
		config:        config,
		provider:      prov,
		enabledTypes:  enabledTypes,
		episodes:      make(map[string][]parser.Episode),
		episodeIndex:  make(map[string]*parser.Episode),
		epgData:       make(map[string][]parser.Programme),
		manifest:      manifest,
		groupCatalogs: groupCatalogs,
		logger:        logger,
	}
}

func extractSelectedGroups(config map[string]any) []string {
	raw, ok := config["selectedGroups"]
	if !ok {
		return nil
	}
	switch v := raw.(type) {
	case []string:
		return v
	case []any:
		groups := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				groups = append(groups, s)
			}
		}
		return groups
	}
	return nil
}

func EnabledTypesFromConfig(config map[string]any) map[string]bool {
	return map[string]bool{
		"tv":     configBool(config, "enableTv", true),
		"movie":  configBool(config, "enableMovies", true),
		"series": configBool(config, "enableSeries", true),
	}
}

func configBool(config map[string]any, key string, defaultValue bool) bool {
	if v, ok := config[key].(bool); ok {
		return v
	}
	return defaultValue
}

func (i *Instance) Initialize(ctx context.Context) error {
	return i.refresh(ctx)
}

func (i *Instance) refresh(ctx context.Context) error {
	i.logger.Info("Refreshing addon data")

	var channels, movies, series []parser.MediaItem
	var epgData map[string][]parser.Programme

	channels, movies, series, err := i.provider.FetchAll(ctx)
	if err != nil {
		i.logger.Error("Provider FetchAll failed", "error", err)
	}

	if i.enabledTypes["tv"] && configBool(i.config, "enableEpg", true) {
		epgData, err = i.provider.FetchEPG(ctx)
		if err != nil {
			i.logger.Warn("Failed to fetch EPG", "error", err)
			epgData = nil
		}
	}

	i.mu.Lock()
	i.channels = channels
	i.movies = movies
	i.series = series
	i.epgData = epgData
	i.buildItemIndexLocked()
	i.episodeIndex = make(map[string]*parser.Episode)
	i.lastUpdate = time.Now()
	i.mu.Unlock()

	if len(i.groupCatalogs) == 0 {
		i.updateManifestGenres()
	}
	i.logger.Info("Addon data refreshed",
		"channels", len(channels),
		"movies", len(movies),
		"series", len(series),
		"epg_channels", len(epgData),
	)

	return nil
}

// buildItemIndexLocked builds itemIndex from channels+movies+series. Caller
// must hold i.mu in write mode.
func (i *Instance) buildItemIndexLocked() {
	i.itemIndex = make(map[string]*parser.MediaItem, len(i.channels)+len(i.movies)+len(i.series))
	for idx := range i.channels {
		i.itemIndex[i.channels[idx].ID] = &i.channels[idx]
	}
	for idx := range i.movies {
		i.itemIndex[i.movies[idx].ID] = &i.movies[idx]
	}
	for idx := range i.series {
		i.itemIndex[i.series[idx].ID] = &i.series[idx]
	}
}

func (i *Instance) updateManifestGenres() {
	i.mu.RLock()
	defer i.mu.RUnlock()

	for idx := range i.manifest.Catalogs {
		catalog := &i.manifest.Catalogs[idx]
		var items []parser.MediaItem

		switch catalog.ID {
		case "iptv_channels":
			items = i.channels
		case "iptv_movies":
			items = i.movies
		case "iptv_series":
			items = i.series
		default:
			continue
		}

		genres := extractGenres(items)
		catalog.Genres = genres

		for eidx := range catalog.Extra {
			if catalog.Extra[eidx].Name == "genre" {
				catalog.Extra[eidx].Options = genres
			}
		}
	}
}

func extractGenres(items []parser.MediaItem) []string {
	genreSet := make(map[string]struct{})
	for _, item := range items {
		if item.Group != "" {
			genreSet[item.Group] = struct{}{}
		}
	}

	genres := make([]string, 0, len(genreSet))
	for g := range genreSet {
		genres = append(genres, g)
	}
	sort.Strings(genres)
	return genres
}

func (i *Instance) GetCatalog(catalogType, catalogID string, extra map[string]string) []MetaPreview {
	i.mu.RLock()
	defer i.mu.RUnlock()

	var items []parser.MediaItem

	if groupName, ok := i.groupCatalogs[catalogID]; ok {
		all := make([]parser.MediaItem, 0, len(i.channels)+len(i.movies)+len(i.series))
		all = append(all, i.channels...)
		all = append(all, i.movies...)
		all = append(all, i.series...)
		items = filterByGenre(all, groupName)
	} else {
		switch catalogID {
		case "iptv_channels":
			if !i.enabledTypes["tv"] {
				return nil
			}
			items = i.channels
		case "iptv_movies":
			if !i.enabledTypes["movie"] {
				return nil
			}
			items = i.movies
		case "iptv_series":
			if !i.enabledTypes["series"] {
				return nil
			}
			items = i.series
		default:
			return nil
		}
	}

	filtered := i.filterItems(items, extra)

	skip := 0
	if s, ok := extra["skip"]; ok {
		fmt.Sscanf(s, "%d", &skip)
	}

	limit := 100
	if skip+limit > len(filtered) {
		limit = len(filtered) - skip
	}
	if limit < 0 {
		limit = 0
	}

	var metas []MetaPreview
	for _, item := range filtered[skip : skip+limit] {
		metas = append(metas, i.itemToMetaPreview(item))
	}

	return metas
}

func (i *Instance) filterItems(items []parser.MediaItem, extra map[string]string) []parser.MediaItem {
	filtered := items

	if genre, ok := extra["genre"]; ok && genre != "" {
		filtered = filterByGenre(filtered, genre)
	}

	if search, ok := extra["search"]; ok && search != "" {
		filtered = i.filterBySearch(filtered, search)
	}

	return filtered
}

func filterByGenre(items []parser.MediaItem, genre string) []parser.MediaItem {
	var result []parser.MediaItem
	for _, item := range items {
		if item.Group == genre {
			result = append(result, item)
		}
	}
	return result
}

func (i *Instance) filterBySearch(items []parser.MediaItem, query string) []parser.MediaItem {
	query = strings.ToLower(query)
	var result []parser.MediaItem
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Name), query) {
			result = append(result, item)
			continue
		}
		if item.Type == parser.TypeTV {
			if prog := i.getCurrentProgramme(item.EPGID); prog != nil {
				if strings.Contains(strings.ToLower(prog.Title), query) ||
					strings.Contains(strings.ToLower(prog.Description), query) {
					result = append(result, item)
				}
			}
		}
	}
	return result
}

func (i *Instance) itemToMetaPreview(item parser.MediaItem) MetaPreview {
	meta := MetaPreview{
		ID:   item.ID,
		Type: string(item.Type),
		Name: item.Name,
	}

	if item.Logo != "" {
		meta.Poster = item.Logo
	}

	switch item.Type {
	case parser.TypeTV:
		if prog := i.getCurrentProgramme(item.EPGID); prog != nil {
			meta.Description = fmt.Sprintf("Now: %s", prog.Title)
		} else {
			meta.Description = "Live Channel"
		}
		meta.Runtime = "Live"
		if item.Group != "" {
			meta.Genres = []string{item.Group}
		}

	case parser.TypeMovie:
		if item.Plot != "" {
			meta.Description = item.Plot
		} else {
			meta.Description = fmt.Sprintf("Movie: %s", item.Name)
		}
		if item.Year > 0 {
			meta.Year = item.Year
		}
		if item.Group != "" {
			meta.Genres = []string{item.Group}
		}

	case parser.TypeSeries:
		if item.Plot != "" {
			meta.Description = item.Plot
		} else {
			meta.Description = "Series"
		}
		if item.Group != "" {
			meta.Genres = []string{item.Group}
		}
	}

	return meta
}

func (i *Instance) getCurrentProgramme(epgID string) *parser.Programme {
	if epgID == "" || i.epgData == nil {
		return nil
	}
	programmes, ok := i.epgData[epgID]
	if !ok {
		return nil
	}
	return parser.GetCurrentProgramme(programmes)
}

func (i *Instance) GetStream(itemType, itemID string) *Stream {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if !i.enabledTypes[itemType] {
		return nil
	}

	item, ok := i.itemIndex[itemID]
	if !ok {
		return nil
	}
	return &Stream{
		URL:   item.URL,
		Title: item.Name,
		BehaviorHints: map[string]any{
			"notWebReady": true,
		},
	}
}

func (i *Instance) GetEpisodeStream(episodeID string) *Stream {
	i.mu.RLock()
	defer i.mu.RUnlock()

	ep, ok := i.episodeIndex[episodeID]
	if !ok {
		return nil
	}
	return &Stream{
		URL:   ep.URL,
		Title: ep.Title,
		BehaviorHints: map[string]any{
			"notWebReady": true,
		},
	}
}

func (i *Instance) GetMeta(itemType, itemID string) *Meta {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if itemType == "series" || strings.HasPrefix(itemID, "iptv_series_") {
		if !i.enabledTypes["series"] {
			return nil
		}
		return i.getSeriesMeta(itemID)
	}

	if !i.enabledTypes[itemType] {
		return nil
	}

	item, ok := i.itemIndex[itemID]
	if !ok {
		return nil
	}
	return i.itemToMeta(*item)
}

func (i *Instance) getSeriesMeta(seriesID string) *Meta {
	s, ok := i.itemIndex[seriesID]
	if !ok || s.Type != parser.TypeSeries {
		return nil
	}
	meta := &Meta{
		ID:          s.ID,
		Type:        string(s.Type),
		Name:        s.Name,
		Description: s.Plot,
	}
	if s.Logo != "" {
		meta.Poster = s.Logo
	}
	if s.Group != "" {
		meta.Genres = []string{s.Group}
	}

	if episodes, ok := i.episodes[seriesID]; ok {
		for _, ep := range episodes {
			meta.Videos = append(meta.Videos, Video{
				ID:        ep.ID,
				Title:     ep.Title,
				Season:    ep.Season,
				Episode:   ep.Episode,
				Thumbnail: ep.Thumbnail,
			})
		}
	}

	return meta
}

func (i *Instance) itemToMeta(item parser.MediaItem) *Meta {
	meta := &Meta{
		ID:   item.ID,
		Type: string(item.Type),
		Name: item.Name,
	}

	if item.Logo != "" {
		meta.Poster = item.Logo
	}

	switch item.Type {
	case parser.TypeTV:
		var desc strings.Builder
		desc.WriteString("CHANNEL: " + item.Name)

		if prog := i.getCurrentProgramme(item.EPGID); prog != nil {
			startStr := prog.Start.Format("15:04")
			stopStr := prog.Stop.Format("15:04")
			desc.WriteString(fmt.Sprintf("\n\nNOW: %s (%s-%s)", prog.Title, startStr, stopStr))
			if prog.Description != "" {
				desc.WriteString("\n\n" + prog.Description)
			}
		}

		if programmes, ok := i.epgData[item.EPGID]; ok {
			upcoming := parser.GetUpcomingProgrammes(programmes, 3)
			if len(upcoming) > 0 {
				desc.WriteString("\n\nUPCOMING:\n")
				for _, p := range upcoming {
					desc.WriteString(fmt.Sprintf("%s - %s\n", p.Start.Format("15:04"), p.Title))
				}
			}
		}

		meta.Description = desc.String()
		meta.Runtime = "Live"
		if item.Group != "" {
			meta.Genres = []string{item.Group}
		}

	case parser.TypeMovie:
		if item.Plot != "" {
			meta.Description = item.Plot
		} else {
			meta.Description = fmt.Sprintf("Movie: %s", item.Name)
		}
		if item.Year > 0 {
			meta.Year = item.Year
		}
		if item.Group != "" {
			meta.Genres = []string{item.Group}
		}
	}

	return meta
}

func (i *Instance) LoadSeriesEpisodes(ctx context.Context, seriesID string) error {
	i.mu.Lock()
	if _, ok := i.episodes[seriesID]; ok {
		i.mu.Unlock()
		return nil
	}
	i.mu.Unlock()

	episodes, err := i.provider.FetchSeriesInfo(ctx, seriesID)
	if err != nil {
		return err
	}

	if episodes == nil {
		episodes = []parser.Episode{}
	}

	i.mu.Lock()
	i.episodes[seriesID] = episodes
	for idx := range episodes {
		i.episodeIndex[episodes[idx].ID] = &episodes[idx]
	}
	i.mu.Unlock()

	return nil
}

func (i *Instance) GetManifest() *Manifest {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.manifest
}

type Stats struct {
	Channels   int    `json:"channels"`
	Movies     int    `json:"movies"`
	Series     int    `json:"series"`
	EPG        int    `json:"epgChannels"`
	LastUpdate string `json:"lastUpdate"`
}

func (i *Instance) GetStats() Stats {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return Stats{
		Channels:   len(i.channels),
		Movies:     len(i.movies),
		Series:     len(i.series),
		EPG:        len(i.epgData),
		LastUpdate: i.lastUpdate.Format(time.RFC3339),
	}
}

func (i *Instance) GetConfig() map[string]any {
	return i.config
}

func (i *Instance) LastUpdate() time.Time {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.lastUpdate
}

type InstanceCache struct {
	cache *cache.Cache[string, *Instance]
}

func NewInstanceCache(maxSize int, ttl time.Duration) *InstanceCache {
	return &InstanceCache{
		cache: cache.New[string, *Instance](maxSize, ttl),
	}
}

func (c *InstanceCache) Get(key string) (*Instance, bool) {
	return c.cache.Get(key)
}

func (c *InstanceCache) Set(key string, instance *Instance) {
	c.cache.Set(key, instance)
}

func (c *InstanceCache) Sweep() {
	c.cache.Sweep()
}
