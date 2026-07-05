# Open Stream M3U

A fast, self-hosted IPTV addon. Written in Go, zero external dependencies.

## Features

- Direct M3U playlist support
- Xtream Codes API (JSON + M3U+ modes)
- XMLTV EPG with "now playing" and upcoming programmes
- Movies, Series, and Live TV catalogs
- Encrypted configuration tokens (AES-256-GCM)
- Fast in-memory LRU cache
- Material 3 web UI with dark/light mode

## Quick Start

### Local

```bash
go run main.go
```

Opens on `http://localhost:7001`

### Docker

```bash
docker compose up -d
```

### Docker (manual)

```bash
docker build -t open-stream-m3u .
docker run -d -p 7001:7001 -e CONFIG_SECRET=my-secret open-stream-m3u
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `7001` | HTTP server port |
| `CONFIG_SECRET` | _(empty)_ | Enables encrypted tokens (min 16 chars) |
| `CACHE_TTL` | `6h` | Cache time-to-live |
| `MAX_CACHE_ENTRIES` | `500` | Max cached addon instances |
| `DEBUG` | `false` | Enable debug logging |
| `PREFETCH_ENABLED` | `true` | Enable CORS bypass proxy |
| `PREFETCH_MAX_SIZE` | `157286400` | Max prefetch response size (150MB) |

## Usage

1. Open `http://localhost:7001` in your browser
2. Choose **Direct M3U** or **Xtream Codes** mode
3. Fill in your playlist/credentials
4. Click **Install Addon**
5. Copy the manifest URL or click **Open in Stremio**

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /` | Landing page |
| `GET /configure` | Configuration UI |
| `GET /health` | Health check |
| `POST /api/prefetch` | CORS bypass proxy |
| `POST /api/encrypt` | Encrypt config token |
| `POST /api/groups` | List playlist groups for selection |
| `GET /api/info` | Addon info for a token |
| `GET /api/debug` | Debug state (requires `DEBUG=true`) |
| `GET /{token}/manifest.json` | Addon manifest |
| `GET /{token}/catalog/{type}/{id}.json` | Catalog |
| `GET /{token}/stream/{type}/{id}.json` | Stream |
| `GET /{token}/meta/{type}/{id}.json` | Metadata |

## License

MIT
