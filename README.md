# Open Stream M3U

A fast, self-hosted IPTV addon for Stremio. Written in Go, zero external dependencies.

## Features

- Direct M3U playlist support
- Xtream Codes API (JSON + M3U+ modes)
- XMLTV EPG with "now playing" and upcoming programmes
- Movies, Series, and Live TV catalogs
- Encrypted configuration tokens (AES-256-GCM)
- Fast in-memory LRU cache
- Material Design web UI with dark/light mode
- ~15MB Docker image (distroless)
- Cloudflare Tunnel ready

## Quick Start

### Local

```bash
go run ./cmd/server
```

Opens on `http://localhost:7001`

### Docker

```bash
docker compose up -d
```

### Docker (manual)

```bash
docker build -t open-stream-m3u .
docker run -d -p 7000:7000 -e CONFIG_SECRET=my-secret open-stream-m3u
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

1. Open `http://localhost:7000` in your browser
2. Choose **Direct M3U** or **Xtream Codes** mode
3. Fill in your playlist/credentials
4. Click **Install Addon**
5. Copy the manifest URL or click **Open in Stremio**

## Cloudflare Tunnel

```bash
cloudflared tunnel --url http://localhost:7000
```

Or with a config file:

```yaml
tunnel: <your-tunnel-id>
ingress:
  - hostname: iptv.yourdomain.com
    service: http://localhost:7000
  - service: http_status:404
```

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /` | Landing page |
| `GET /configure` | Configuration UI |
| `GET /health` | Health check |
| `POST /api/prefetch` | CORS bypass proxy |
| `POST /api/encrypt` | Encrypt config token |
| `GET /{token}/manifest.json` | Stremio manifest |
| `GET /{token}/catalog/{type}/{id}.json` | Catalog |
| `GET /{token}/stream/{type}/{id}.json` | Stream |
| `GET /{token}/meta/{type}/{id}.json` | Metadata |

## License

MIT
