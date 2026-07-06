package parser

import (
	"encoding/xml"
	"io"
	"sort"
	"strings"
	"time"
)

type Programme struct {
	Start       time.Time
	Stop        time.Time
	Title       string
	TitleLower  string
	Description string
	DescLower   string
}

type xmltvProgramme struct {
	Channel string      `xml:"channel,attr"`
	Start   string      `xml:"start,attr"`
	Stop    string      `xml:"stop,attr"`
	Title   []xmltvText `xml:"title"`
	Desc    []xmltvText `xml:"desc"`
}

type xmltvText struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

// ParseXMLTV streams an XMLTV document from r and returns a map of channel
// ID -> sorted programmes. It decodes one <programme> at a time, so peak
// memory is O(largest programme element) rather than O(whole document) —
// the body bytes never materialize as a single buffer.
//
// Malformed individual programmes are skipped, not fatal: a single bad element
// no longer discards the whole EPG.
func ParseXMLTV(r io.Reader) (map[string][]Programme, error) {
	dec := xml.NewDecoder(r)
	epgData := make(map[string][]Programme)

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		se, ok := tok.(xml.StartElement)
		if !ok || se.Name.Local != "programme" {
			continue
		}

		var xp xmltvProgramme
		if err := dec.DecodeElement(&xp, &se); err != nil {
			// ponytail: skip malformed programme, keep the rest of the document.
			continue
		}

		title := ""
		if len(xp.Title) > 0 {
			title = strings.TrimSpace(xp.Title[0].Value)
		}
		desc := ""
		if len(xp.Desc) > 0 {
			desc = strings.TrimSpace(xp.Desc[0].Value)
		}

		p := Programme{
			Start:       parseXMLTVTime(xp.Start),
			Stop:        parseXMLTVTime(xp.Stop),
			Title:       title,
			TitleLower:  strings.ToLower(title),
			Description: desc,
			DescLower:   strings.ToLower(desc),
		}

		epgData[xp.Channel] = append(epgData[xp.Channel], p)
	}

	// Sort each channel's programmes by Start so GetCurrentProgramme can
	// binary-search (see GetCurrentProgramme).
	for _, progs := range epgData {
		sort.Slice(progs, func(i, j int) bool {
			return progs[i].Start.Before(progs[j].Start)
		})
	}

	return epgData, nil
}

func parseXMLTVTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}

	s = strings.TrimSpace(s)

	if len(s) >= 14 {
		base := s[:14]
		tz := ""
		if len(s) > 14 {
			tz = strings.TrimSpace(s[14:])
		}

		year := parseInt(base[0:4])
		month := parseInt(base[4:6])
		day := parseInt(base[6:8])
		hour := parseInt(base[8:10])
		min := parseInt(base[10:12])
		sec := parseInt(base[12:14])

		loc := time.UTC
		if tz != "" {
			if offset, err := parseTimezone(tz); err == nil {
				loc = time.FixedZone(tz, offset)
			}
		}

		return time.Date(year, time.Month(month), day, hour, min, sec, 0, loc)
	}

	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}

	return time.Time{}
}

func parseInt(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}

func parseTimezone(tz string) (int, error) {
	tz = strings.TrimSpace(tz)
	if len(tz) < 5 {
		return 0, nil
	}

	sign := 1
	if tz[0] == '-' {
		sign = -1
	} else if tz[0] != '+' {
		return 0, nil
	}

	hours := parseInt(tz[1:3])
	mins := parseInt(tz[3:5])
	return sign * (hours*3600 + mins*60), nil
}

func GetCurrentProgramme(programmes []Programme) *Programme {
	if len(programmes) == 0 {
		return nil
	}
	now := time.Now()
	// ponytail: assumes non-overlapping programmes per channel (sorted by Start in ParseXMLTV).
	// Find last programme with Start <= now via binary search for first Start > now.
	idx := sort.Search(len(programmes), func(i int) bool {
		return programmes[i].Start.After(now)
	})
	if idx == 0 {
		return nil
	}
	candidate := &programmes[idx-1]
	if now.Before(candidate.Stop) {
		return candidate
	}
	return nil
}

func GetUpcomingProgrammes(programmes []Programme, limit int) []Programme {
	now := time.Now()
	var upcoming []Programme
	for _, p := range programmes {
		if p.Start.After(now) {
			upcoming = append(upcoming, p)
			if len(upcoming) >= limit {
				break
			}
		}
	}
	return upcoming
}