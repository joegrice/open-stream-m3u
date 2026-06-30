package parser

import (
	"encoding/xml"
	"strings"
	"time"
)

type Programme struct {
	Start       time.Time
	Stop        time.Time
	Title       string
	Description string
}

type xmltvTV struct {
	XMLName    xml.Name       `xml:"tv"`
	Programmes []xmltvProgramme `xml:"programme"`
}

type xmltvProgramme struct {
	Channel string `xml:"channel,attr"`
	Start   string `xml:"start,attr"`
	Stop    string `xml:"stop,attr"`
	Title   []xmltvText `xml:"title"`
	Desc    []xmltvText `xml:"desc"`
}

type xmltvText struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

func ParseXMLTV(content string) (map[string][]Programme, error) {
	var tv xmltvTV
	if err := xml.Unmarshal([]byte(content), &tv); err != nil {
		return nil, err
	}

	epgData := make(map[string][]Programme)
	for _, prog := range tv.Programmes {
		start := parseXMLTVTime(prog.Start)
		stop := parseXMLTVTime(prog.Stop)

		title := ""
		if len(prog.Title) > 0 {
			title = strings.TrimSpace(prog.Title[0].Value)
		}

		desc := ""
		if len(prog.Desc) > 0 {
			desc = strings.TrimSpace(prog.Desc[0].Value)
		}

		p := Programme{
			Start:       start,
			Stop:        stop,
			Title:       title,
			Description: desc,
		}

		epgData[prog.Channel] = append(epgData[prog.Channel], p)
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
	now := time.Now()
	for i := range programmes {
		p := &programmes[i]
		if now.After(p.Start) && now.Before(p.Stop) {
			return p
		}
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
