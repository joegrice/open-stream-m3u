package addon

type Manifest struct {
	ID            string      `json:"id"`
	Version       string      `json:"version"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	Resources     []string    `json:"resources"`
	Types         []string    `json:"types"`
	IDPrefixes    []string    `json:"idPrefixes"`
	Catalogs      []Catalog   `json:"catalogs"`
	BehaviorHints map[string]any `json:"behaviorHints,omitempty"`
}

type Catalog struct {
	Type   string         `json:"type"`
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Extra  []CatalogExtra `json:"extra,omitempty"`
	Genres []string       `json:"genres,omitempty"`
}

type CatalogExtra struct {
	Name    string   `json:"name"`
	IsRequired bool  `json:"isRequired,omitempty"`
	Options []string `json:"options,omitempty"`
}

type MetaPreview struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Poster      string   `json:"poster,omitempty"`
	PosterShape string   `json:"posterShape,omitempty"`
	Logo        string   `json:"logo,omitempty"`
	Description string   `json:"description,omitempty"`
	ReleaseInfo string   `json:"releaseInfo,omitempty"`
	Year        int      `json:"year,omitempty"`
	Runtime     string   `json:"runtime,omitempty"`
	Genres      []string `json:"genres,omitempty"`
}

type Meta struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Poster      string   `json:"poster,omitempty"`
	Logo        string   `json:"logo,omitempty"`
	Description string   `json:"description,omitempty"`
	ReleaseInfo string   `json:"releaseInfo,omitempty"`
	Year        int      `json:"year,omitempty"`
	Runtime     string   `json:"runtime,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Videos      []Video  `json:"videos,omitempty"`
}

type Video struct {
	ID        string `json:"id"`
	Title     string `json:"title,omitempty"`
	Season    int    `json:"season,omitempty"`
	Episode   int    `json:"episode,omitempty"`
	Released  string `json:"released,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Stream    *Stream `json:"stream,omitempty"`
}

type Stream struct {
	URL            string            `json:"url,omitempty"`
	ExternalURL    string            `json:"externalUrl,omitempty"`
	YouTube        string            `json:"ytId,omitempty"`
	InfoHash       string            `json:"infoHash,omitempty"`
	Title          string            `json:"title,omitempty"`
	BehaviorHints  map[string]any    `json:"behaviorHints,omitempty"`
}

type CatalogResponse struct {
	Metas []MetaPreview `json:"metas"`
}

type MetaResponse struct {
	Meta *Meta `json:"meta"`
}

type StreamResponse struct {
	Streams []Stream `json:"streams"`
}

func BuildManifest() *Manifest {
	return &Manifest{
		ID:          "org.openstream.m3u",
		Version:     "1.0.0",
		Name:        "Open Stream M3U",
		Description: "IPTV addon for M3U playlists and EPG data",
		Resources:   []string{"catalog", "stream", "meta"},
		Types:       []string{"tv", "movie", "series"},
		IDPrefixes:  []string{"iptv_"},
		Catalogs: []Catalog{
			{
				Type: "tv",
				ID:   "iptv_channels",
				Name: "IPTV Channels",
				Extra: []CatalogExtra{
					{Name: "genre"},
					{Name: "search"},
					{Name: "skip"},
				},
			},
			{
				Type: "movie",
				ID:   "iptv_movies",
				Name: "IPTV Movies",
				Extra: []CatalogExtra{
					{Name: "genre"},
					{Name: "search"},
					{Name: "skip"},
				},
			},
			{
				Type: "series",
				ID:   "iptv_series",
				Name: "IPTV Series",
				Extra: []CatalogExtra{
					{Name: "genre"},
					{Name: "search"},
					{Name: "skip"},
				},
			},
		},
		BehaviorHints: map[string]any{
			"configurable":        true,
			"configurationRequired": false,
		},
	}
}
