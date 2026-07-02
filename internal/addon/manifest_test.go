package addon

import (
	"slices"
	"testing"
)

func TestBuildManifest_AllEnabled(t *testing.T) {
	m := BuildManifest(nil, map[string]bool{"tv": true, "movie": true, "series": true})
	if len(m.Catalogs) != 3 {
		t.Fatalf("expected 3 catalogs, got %d", len(m.Catalogs))
	}
	wantTypes := []string{"tv", "channel", "movie", "series"}
	for _, want := range wantTypes {
		if !slices.Contains(m.Types, want) {
			t.Errorf("expected type %q in %v", want, m.Types)
		}
	}
}

func TestBuildManifest_MoviesDisabled(t *testing.T) {
	m := BuildManifest(nil, map[string]bool{"tv": true, "movie": false, "series": true})
	if slices.Contains(m.Types, "movie") {
		t.Error("expected movie type to be removed")
	}
	for _, c := range m.Catalogs {
		if c.Type == "movie" {
			t.Error("expected no movie catalog")
		}
	}
}

func TestBuildManifest_GroupCatalogsKeepChannelType(t *testing.T) {
	m := BuildManifest([]string{"Sports"}, map[string]bool{"tv": false, "movie": true, "series": false})
	if len(m.Catalogs) != 1 {
		t.Fatalf("expected 1 group catalog, got %d", len(m.Catalogs))
	}
	if !slices.Contains(m.Types, "channel") {
		t.Error("expected channel type for group catalogs")
	}
}
