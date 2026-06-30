package provider

import (
	"context"

	"github.com/joe/open-stream-m3u/internal/parser"
)

type Provider interface {
	FetchChannels(ctx context.Context) ([]parser.MediaItem, error)
	FetchMovies(ctx context.Context) ([]parser.MediaItem, error)
	FetchSeries(ctx context.Context) ([]parser.MediaItem, error)
	FetchSeriesInfo(ctx context.Context, seriesID string) ([]parser.Episode, error)
	FetchEPG(ctx context.Context) (map[string][]parser.Programme, error)
}
