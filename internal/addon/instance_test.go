package addon

import (
	"testing"
	"time"

	"github.com/joe/open-stream-m3u/internal/parser"
)

func TestGetCatalogSearchEPG(t *testing.T) {
	now := time.Now()

	channels := []parser.MediaItem{
		{
			ID:    "iptv_ch_bbc1",
			Type:  parser.TypeTV,
			Name:  "UK: BBC 1",
			EPGID: "bbc1",
			Group: "UK",
		},
		{
			ID:    "iptv_ch_bbc2",
			Type:  parser.TypeTV,
			Name:  "UK: BBC 2",
			EPGID: "bbc2",
			Group: "UK",
		},
		{
			ID:    "iptv_ch_noepg",
			Type:  parser.TypeTV,
			Name:  "No EPG Channel",
			EPGID: "",
			Group: "Other",
		},
	}

	epgData := map[string][]parser.Programme{
		"bbc1": {
			{
				Start:       now.Add(-2 * time.Hour),
				Stop:        now.Add(2 * time.Hour),
				Title:       "MOTD FIFA World Cup 2026",
				Description: "Match of the Day coverage of the World Cup.",
			},
		},
		"bbc2": {
			{
				Start:       now.Add(-2 * time.Hour),
				Stop:        now.Add(2 * time.Hour),
				Title:       "Documentary",
				Description: "A film about penguins in Antarctica.",
			},
		},
	}

	inst := &Instance{
		channels:      channels,
		epgData:       epgData,
		groupCatalogs: map[string]string{},
		enabledTypes:  map[string]bool{"tv": true, "movie": true, "series": true},
		manifest:      BuildManifest(nil, map[string]bool{"tv": true, "movie": true, "series": true}),
	}

	tests := []struct {
		name    string
		query   string
		wantIDs []string
	}{
		{
			name:    "match by channel name",
			query:   "bbc 1",
			wantIDs: []string{"iptv_ch_bbc1"},
		},
		{
			name:    "match by current programme title",
			query:   "MOTD",
			wantIDs: []string{"iptv_ch_bbc1"},
		},
		{
			name:    "match by current programme description",
			query:   "penguins",
			wantIDs: []string{"iptv_ch_bbc2"},
		},
		{
			name:    "case insensitive match",
			query:   "motd",
			wantIDs: []string{"iptv_ch_bbc1"},
		},
		{
			name:    "channel name still works without epg",
			query:   "No EPG",
			wantIDs: []string{"iptv_ch_noepg"},
		},
		{
			name:    "no match returns empty",
			query:   "news",
			wantIDs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extra := map[string]string{"search": tt.query}
			got := inst.GetCatalog("channel", "iptv_channels", extra)

			if len(got) != len(tt.wantIDs) {
				t.Fatalf("got %d results, want %d: %+v", len(got), len(tt.wantIDs), got)
			}

			for i, want := range tt.wantIDs {
				if got[i].ID != want {
					t.Errorf("result[%d].ID = %q, want %q", i, got[i].ID, want)
				}
			}
		})
	}
}

func TestDisabledTypesHideContent(t *testing.T) {
	inst := &Instance{
		channels: []parser.MediaItem{
			{ID: "iptv_tv_1", Type: parser.TypeTV, Name: "TV Channel", URL: "http://example.com/tv"},
		},
		movies: []parser.MediaItem{
			{ID: "iptv_movie_1", Type: parser.TypeMovie, Name: "A Movie", URL: "http://example.com/movie"},
		},
		series: []parser.MediaItem{
			{ID: "iptv_series_1", Type: parser.TypeSeries, Name: "A Series"},
		},
		episodes: map[string][]parser.Episode{
			"iptv_series_1": {{ID: "iptv_series_ep_1", Title: "Pilot", URL: "http://example.com/ep1"}},
		},
		groupCatalogs: map[string]string{},
		enabledTypes:  map[string]bool{"tv": true, "movie": false, "series": false},
	}

	if got := inst.GetCatalog("movie", "iptv_movies", nil); len(got) != 0 {
		t.Errorf("GetCatalog for disabled movie type returned %d results, want 0", len(got))
	}
	if got := inst.GetCatalog("series", "iptv_series", nil); len(got) != 0 {
		t.Errorf("GetCatalog for disabled series type returned %d results, want 0", len(got))
	}
	if inst.GetStream("movie", "iptv_movie_1") != nil {
		t.Error("GetStream returned a stream for a disabled movie type")
	}
	if inst.GetMeta("series", "iptv_series_1") != nil {
		t.Error("GetMeta returned a meta for a disabled series type")
	}

	// Live TV should still work.
	if got := inst.GetCatalog("channel", "iptv_channels", nil); len(got) != 1 {
		t.Errorf("GetCatalog for enabled tv type returned %d results, want 1", len(got))
	}
}
