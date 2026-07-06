package parser

import (
	"strings"
	"testing"
)

// TestM3UStramingSelfCheck verifies the bufio.Scanner-based parser produces
// the same structural results as the old strings.Split version, plus that
// lines longer than the default scanner buffer (64KB) still parse.
func TestM3USelfCheck(t *testing.T) {
	m3u := `#EXTM3U
#EXTINF:-1 tvg-id="bbc1" tvg-logo="http://x/logo.png" group-title="News",UK: BBC 1
http://example.com/bbc1
#EXTINF:-1 group-title="Movies",Some Movie (2024)
http://example.com/movie
#EXTINF:-1 group-title="Series",Show Name S01E02 Pilot
http://example.com/series
#EXTINF:-1 tvg-id="bad" group-title="News",No URL follows
#EXTINF:-1,Good After Bad
http://example.com/recovered
# worthless comment line
http://example.com/orphan-no-extinf
`

	items, err := ParseM3U(strings.NewReader(m3u))
	if err != nil {
		t.Fatalf("ParseM3U error: %v", err)
	}
	if len(items) != 4 {
		t.Fatalf("want 4 items (no-URL entry dropped, orphan dropped), got %d: %+v", len(items), items)
	}
	if items[0].Name != "UK: BBC 1" || items[0].Group != "News" || items[0].EPGID != "bbc1" {
		t.Errorf("item[0] = %+v", items[0])
	}
	if items[0].NameLower != "uk: bbc 1" {
		t.Errorf("NameLower = %q", items[0].NameLower)
	}
	if items[1].Type != TypeMovie {
		t.Errorf("item[1] type = %v, want movie", items[1].Type)
	}
	if items[2].Type != TypeSeries || items[2].Season != 1 || items[2].Episode != 2 {
		t.Errorf("item[2] series/SE = %v %d %d", items[2].Type, items[2].Season, items[2].Episode)
	}
	if items[3].Name != "Good After Bad" {
		t.Errorf("item[3] (recovered after no-URL entry) = %+v", items[3])
	}
}

// TestM3ULongLine verifies a #EXTINF line exceeding the default scanner
// buffer (64KB) still parses with the raised 1MB cap.
func TestM3ULongLine(t *testing.T) {
	longAttrs := strings.Repeat("tvg-id=\"x\" ", 9000) // ~90KB of attrs
	m3u := "#EXTINF:-1 " + longAttrs + "group-title=\"Big\",Huge Line\nhttp://example.com/big\n"

	items, err := ParseM3U(strings.NewReader(m3u))
	if err != nil {
		t.Fatalf("ParseM3U error on long line: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("want 1 item, got %d", len(items))
	}
	if items[0].Name != "Huge Line" {
		t.Errorf("name = %q", items[0].Name)
	}
	if items[0].Group != "Big" {
		t.Errorf("group = %q (long attrs broke parsing)", items[0].Group)
	}
}