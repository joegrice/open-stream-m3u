package parser

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ItemType string

const (
	TypeTV     ItemType = "tv"
	TypeMovie  ItemType = "movie"
	TypeSeries ItemType = "series"
)

type MediaItem struct {
	ID        string
	Name      string
	URL       string
	Type      ItemType
	Logo      string
	Group     string
	EPGID     string
	Season    int
	Episode   int
	Year      int
	Plot      string
	Attrs     map[string]string
}

type Episode struct {
	ID        string
	Title     string
	Season    int
	Episode   int
	URL       string
	Thumbnail string
}

var (
	reExtinf      = regexp.MustCompile(`#EXTINF:(-?\d+)(?:\s+(.*))?,(.*)`)
	reAttr        = regexp.MustCompile(`([\w-]+)="([^"]*)"`)
	reMovieFormat = regexp.MustCompile(`\(\d{4}\)|\d{4}\.|HD$|FHD$|4K$`)
	reSeriesSE    = regexp.MustCompile(`\bS(\d{1,2})E(\d{1,2})\b`)
	reSeasonEp    = regexp.MustCompile(`\bSeason\s?(\d{1,2}).*?\bEpisode\s?(\d{1,3})\b`)
	reSeasonEp2   = regexp.MustCompile(`\bSeason\s?(\d{1,2}).*?\bEp\s?(\d{1,3})\b`)
	reYear        = regexp.MustCompile(`\((\d{4})\)`)
)

func ParseM3U(content string) ([]MediaItem, error) {
	lines := strings.Split(content, "\n")
	var items []MediaItem
	var currentItem *MediaItem

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#EXTINF:") {
			matches := reExtinf.FindStringSubmatch(line)
			if matches == nil {
				continue
			}

			currentItem = &MediaItem{
				Attrs: make(map[string]string),
			}

			attrs := parseAttributes(matches[2])
			currentItem.Attrs = attrs
			currentItem.Name = strings.TrimSpace(matches[3])
			currentItem.Logo = attrs["tvg-logo"]
			currentItem.EPGID = attrs["tvg-id"]
			if currentItem.EPGID == "" {
				currentItem.EPGID = attrs["tvg-name"]
			}
			currentItem.Group = attrs["group-title"]
		} else if !strings.HasPrefix(line, "#") && currentItem != nil {
			currentItem.URL = line
			currentItem.Type = classifyItem(currentItem)
			currentItem.ID = generateID(currentItem.Name + currentItem.URL)

			if m := reYear.FindStringSubmatch(currentItem.Name); m != nil {
				currentItem.Year, _ = strconv.Atoi(m[1])
			}

			if currentItem.Type == TypeSeries {
				currentItem.Season, currentItem.Episode = extractSeasonEpisode(currentItem.Name)
			}

			items = append(items, *currentItem)
			currentItem = nil
		}
	}

	return items, nil
}

func parseAttributes(s string) map[string]string {
	attrs := make(map[string]string)
	matches := reAttr.FindAllStringSubmatch(s, -1)
	for _, m := range matches {
		attrs[m[1]] = m[2]
	}
	return attrs
}

func classifyItem(item *MediaItem) ItemType {
	group := strings.ToLower(item.Group)
	name := strings.ToLower(item.Name)

	if strings.Contains(group, "movie") || strings.Contains(name, "movie") || reMovieFormat.MatchString(item.Name) {
		return TypeMovie
	}

	if strings.Contains(group, "series") || strings.Contains(group, "show") ||
		reSeriesSE.MatchString(item.Name) || reSeasonEp.MatchString(item.Name) || reSeasonEp2.MatchString(item.Name) {
		return TypeSeries
	}

	return TypeTV
}

func extractSeasonEpisode(name string) (int, int) {
	if m := reSeriesSE.FindStringSubmatch(name); m != nil {
		s, _ := strconv.Atoi(m[1])
		e, _ := strconv.Atoi(m[2])
		return s, e
	}
	if m := reSeasonEp.FindStringSubmatch(name); m != nil {
		s, _ := strconv.Atoi(m[1])
		e, _ := strconv.Atoi(m[2])
		return s, e
	}
	if m := reSeasonEp2.FindStringSubmatch(name); m != nil {
		s, _ := strconv.Atoi(m[1])
		e, _ := strconv.Atoi(m[2])
		return s, e
	}
	return 1, 0
}

func generateID(s string) string {
	hash := md5.Sum([]byte(s))
	return fmt.Sprintf("iptv_%x", hash[:8])
}

func GetSeriesBaseName(name string) string {
	name = reSeriesSE.ReplaceAllString(name, "")
	name = reSeasonEp.ReplaceAllString(name, "")
	name = reSeasonEp2.ReplaceAllString(name, "")
	name = strings.TrimRight(name, "-._ ")
	return strings.TrimSpace(name)
}

func GroupSeries(items []MediaItem) (map[string][]Episode, map[string]*MediaItem) {
	seriesMap := make(map[string]*MediaItem)
	episodesMap := make(map[string][]Episode)

	for _, item := range items {
		if item.Type != TypeSeries {
			continue
		}

		baseName := GetSeriesBaseName(item.Name)
		if baseName == "" {
			continue
		}

		seriesID := "iptv_series_" + generateID(baseName)

		if _, exists := seriesMap[seriesID]; !exists {
			seriesMap[seriesID] = &MediaItem{
				ID:    seriesID,
				Name:  baseName,
				Type:  TypeSeries,
				Logo:  item.Logo,
				Group: item.Group,
				Plot:  item.Plot,
				Attrs: item.Attrs,
			}
		}

		epID := "iptv_series_ep_" + generateID(seriesID+item.URL+strconv.Itoa(item.Season)+strconv.Itoa(item.Episode))
		episode := Episode{
			ID:        epID,
			Title:     item.Name,
			Season:    item.Season,
			Episode:   item.Episode,
			URL:       item.URL,
			Thumbnail: item.Logo,
		}
		episodesMap[seriesID] = append(episodesMap[seriesID], episode)
	}

	return episodesMap, seriesMap
}
