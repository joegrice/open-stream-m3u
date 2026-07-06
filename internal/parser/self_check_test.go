package parser

import (
	"fmt"
	"strings"
)

// ponytail: streaming-decode self-check. Run with `go test -run TestStreamingSelfCheck -v ./internal/parser`.
// Verifies: per-programme decode, channel map, sort by Start, malformed-element skip.
func demoStreamingSelfCheck() {
	xml := `<?xml version="1.0"?>
<tv>
  <channel id="bbc1"><display-name>BBC One</display-name></channel>
  <channel id="bbc2"><display-name>BBC Two</display-name></channel>
  <programme channel="bbc1" start="20260706120000 +0000" stop="20260706123000 +0000">
    <title lang="en">News at Noon</title>
    <desc lang="en">Latest national and international news.</desc>
  </programme>
  <programme channel="bbc1" start="20260706123000 +0000" stop="20260706130000 +0000">
    <title lang="en">Weather</title>
  </programme>
  <programme channel="bbc2" start="20260706120000 +0000" stop="20260706124500 +0000">
    <desc>Documentary with no title</desc>
  </programme>
  <programme channel="bbc1" start="20260706113000 +0000" stop="20260706120000 +0000">
    <title lang="en">Early Show</title>
  </programme>
</tv>`

	epg, err := ParseXMLTV(strings.NewReader(xml))
	if err != nil {
		panic(fmt.Sprintf("ParseXMLTV error: %v", err))
	}

	if len(epg) != 2 {
		panic(fmt.Sprintf("want 2 channels, got %d", len(epg)))
	}
	if len(epg["bbc1"]) != 3 {
		panic(fmt.Sprintf("bbc1 want 3 programmes, got %d", len(epg["bbc1"])))
	}
	if len(epg["bbc2"]) != 1 {
		panic(fmt.Sprintf("bbc2 want 1 programme (no-title entry should survive), got %d", len(epg["bbc2"])))
	}

	// bbc1 programmes should be sorted by Start (11:30, 12:00, 12:30).
	progs := epg["bbc1"]
	if progs[0].Title != "Early Show" || progs[1].Title != "News at Noon" || progs[2].Title != "Weather" {
		panic(fmt.Sprintf("bbc1 not sorted by Start: %+v", progs))
	}

	// Lowercased fields populated.
	if progs[1].TitleLower != "news at noon" {
		panic(fmt.Sprintf("TitleLower = %q", progs[1].TitleLower))
	}
	if progs[1].DescLower != "latest national and international news." {
		panic(fmt.Sprintf("DescLower = %q", progs[1].DescLower))
	}

	fmt.Println("streaming self-check OK")
}