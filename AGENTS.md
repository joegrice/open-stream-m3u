# AGENTS.md

Agent-focused context for working on Open Stream M3U. For human users, see README.md.

## Project Snapshot

- **What:** Self-hosted IPTV addon for M3U playlists and XMLTV EPG.
- **Stack:** Go 1.22+, standard library only (zero external dependencies).
- **Entry point:** `main.go` embeds `web/` and starts `internal/server`.
- **Security-sensitive areas:** encrypted config tokens, CORS bypass proxy, external URL fetching.

## Quick Commands

| Task | Command |
|------|---------|
| Run locally | `go run main.go` |
| Run tests | `go test ./...` |
| Run with debug | `DEBUG=true go run main.go` |
| Docker | `docker compose up -d` |
| Build image | `docker build -t open-stream-m3u .` |

Default port is `7001`. The README mentions `go run ./cmd/server` but the repo uses `main.go`.

## Architecture at a Glance

```
main.go
  └─ server.New()
       ├─ static routes  → web/ (embed.FS)
       ├─ API routes     → encrypt, prefetch, groups, info, debug
       └─ /{token}/{path...}
              └─ getOrBuildInstance()
                     ├─ parse token → decrypt/decode config
                     ├─ createProvider() → direct | xtream
                     ├─ Instance.Initialize() → fetch + parse
                     └─ route to addon handlers (manifest, catalog, stream, meta)
```

- `internal/config` — env-based config.
- `internal/server` — HTTP routing, middleware, API handlers.
- `internal/addon` — addon manifest/catalog/stream/meta handlers and instance cache.
- `internal/provider` — Direct M3U and Xtream Codes providers.
- `internal/parser` — M3U and XMLTV parsing.
- `internal/crypto` — token encoding/decoding and AES-256-GCM encryption.
- `internal/cache` — LRU instance cache.

## Agent Workflows

### Add or modify an endpoint
1. Edit `internal/server/server.go` → `setupRoutes()` and handler.
2. Keep handlers thin: validate, delegate, encode JSON or `http.Error`.
3. Use `context.WithTimeout` for any external or expensive work.
4. Add/update `internal/server/server_test.go` for non-trivial logic.

### Add a provider
1. Implement `provider.Provider` in a new file under `internal/provider/`.
2. Register it in `server.createProvider()` using the `providerType` config key.
3. Keep it stdlib-only; use `http.Client` with timeouts.

### Change parsing
1. `internal/parser/m3u.go` or `internal/parser/xmltv.go`.
2. Avoid regex for large files; prefer streaming/scanners where possible.
3. Add table-driven tests for edge cases.

### Change the UI
1. Files under `web/` (HTML, CSS, JS).
2. **Start at `DESIGN.md`** for any visual change; update it and
   `web/css/tokens.css` together when tokens move.
3. Consume `--m3-*` tokens from `tokens.css`; do not re-introduce ad-hoc
   `--bg-*` / `--accent*` / `--text-*` vars.
4. Material 3 styling, dark/light via `[data-theme]`, pill buttons, M3
   elevation, outlined text fields — see "UI / Design System" below.
5. Served via `embed.FS`; rebuild not required for static changes.

### Fix a bug
1. Reproduce with a test if possible.
2. Fix in the shared function when multiple callers are affected.
3. Run `go test ./...` before considering it done.

### Security change
1. Double-check trust boundaries: prefetch proxy, token parsing, credential handling.
2. Never log secrets or decoded credentials.
3. Fail closed (deny by default).

## Coding Conventions

- **Dependencies:** standard library only. Do not add new modules.
- **Formatting:** `gofmt`.
- **Errors:** log in handlers, propagate in packages. HTTP handlers return 4xx for client errors, 5xx for server errors.
- **Timeouts:** always use `context.WithTimeout` for external fetches.
- **Tests:** table-driven tests for parsers and handlers; `go test ./...` must pass.
- **Naming:** keep package names short and consistent with existing code.
- **CSS:** prefer Material 3 role tokens (`--m3-*`); see `DESIGN.md` and "UI / Design System".

## Security & Trust Boundaries

- `CONFIG_SECRET` enables encrypted tokens. Never log it.
- `POST /api/prefetch` is a CORS bypass proxy; `isBlockedHost()` blocks loopback and RFC1918 addresses.
- Tokens may be URL-encoded; `handleTokenRoute()` decodes before parsing.
- Token decryption fails closed when `CONFIG_SECRET` is missing or invalid.
- Validate external URLs before `http.Get`; respect `PrefetchMaxSize`.

## UI / Design System

- `DESIGN.md` at repo root is the **source of truth** for the visual identity
  — Material 3 design system. Read it before touching anything under `web/`.
- `web/css/tokens.css` is the CSS reflection of `DESIGN.md`. The mapping is 1:1:
  every `colors.<role>` in DESIGN.md front matter becomes `--m3-<role>` in
  `tokens.css`. When you change one, change the other in the same commit.
- Two themes are expressed by overriding the `--m3-*` vars under
  `[data-theme="dark"]` in `tokens.css`. Do **not** introduce new ad-hoc theme
  vars elsewhere; reuse the M3 role tokens.
- Non-M3 additions (the green `--m3-success` check-mark color is the only
  current one) must be documented inline in `tokens.css` with a
  `/* non-M3: … */` comment. Keep them rare.
- Elevation, motion, shape, and spacing scales also live in `tokens.css` and
  follow the M3 reference values documented in the DESIGN.md body. Prefer
  `--m3-elevation-N` and `--m3-radius-N` over hand-rolled shadow/radius.
- Typography uses the system font stack only (no web fonts, no CDN, no
  bundling). Font weights and sizes come from `--m3-font-*` shorthands.

## Common Pitfalls

- `cmd/server` does not exist; run from repo root with `go run main.go`.
- Token routes: `{token}/{path...}` means the token is the first path segment.
- Group catalog IDs are derived from `md5(group)` and must stay stable.
- Cache TTL (`CACHE_TTL`) is separate from per-fetch timeouts.
- `handleGroups` ignores fetch errors intentionally to return partial group lists; document if changing.

## Decision Log

- Stdlib-only to keep the Docker image small and dependency surface zero.
- LRU cache of addon instances to avoid re-fetching playlists on every request.
- AES-256-GCM for encrypted tokens with `CONFIG_SECRET`.
- Web UI is static and embedded; no build step.
- M3 token system per `DESIGN.md` spec; flat `--m3-*` CSS reflection in
  `tokens.css`, no aliasing layer, no build step.
- System font stack only for Material 3 — no web fonts, preserves the
  zero-external-deps / embedded model.
