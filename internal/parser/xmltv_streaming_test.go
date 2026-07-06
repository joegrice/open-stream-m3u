package parser

import (
	"strings"
	"testing"
)

// TestStreamingSelfCheck runs the package-level self-check so the
// streaming-decode logic runs in `go test ./...`. If any invariant breaks the
// demoStreamingSelfCheck func panics and the test fails.
func TestStreamingSelfCheck(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("streaming self-check failed: %v", r)
		}
	}()
	demoStreamingSelfCheck()
}

// TestStreamingMalformedSkipped verifies a single malformed <programme>
// element doesn't discard the whole document.
func TestStreamingMalformedSkipped(t *testing.T) {
	xml := `<?xml version="1.0"?>
<tv>
  <programme channel="ok" start="20260706120000 +0000" stop="20260706123000 +0000"><title>Good</title></programme>
  <programme channel="broken"><title>Missing start/stop attrs - gets zero times but still parses</title></programme>
  <programme channel="ok" start="20260706123000 +0000" stop="20260706130000 +0000"><title>Also Good</title></programme>
</tv>`

	epg, err := ParseXMLTV(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("ParseXMLTV error: %v", err)
	}
	if len(epg["ok"]) != 2 {
		t.Errorf("ok channel want 2 programmes, got %d", len(epg["ok"]))
	}
	if len(epg["broken"]) != 1 {
		t.Errorf("broken channel want 1 (zero-time entry), got %d", len(epg["broken"]))
	}
}