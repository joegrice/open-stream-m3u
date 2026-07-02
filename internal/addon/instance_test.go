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
		manifest:      BuildManifest(nil),
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
