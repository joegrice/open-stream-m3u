package provider

import (
	"context"

	"github.com/joe/open-stream-m3u/internal/parser"
)

type Provider interface {
	FetchAll(ctx context.Context) (channels, movies, series []parser.MediaItem, err error)
	FetchSeriesInfo(ctx context.Context, seriesID string) ([]parser.Episode, error)
	FetchEPG(ctx context.Context) (map[string][]parser.Programme, error)
}