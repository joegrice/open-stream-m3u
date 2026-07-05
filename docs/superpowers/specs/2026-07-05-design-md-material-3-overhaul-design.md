# DESIGN.md + Material 3 Token Overhaul — Design Spec

Date: 2026-07-05
Status: Approved (brainstormed 2026-07-05)

## Goal

Adopt the [`design.md`](https://github.com/google-labs-code/design.md) format
(alpha) as the persistent, structured source of truth for the Open Stream M3U
web UI's visual identity, and overhaul the existing ad-hoc CSS variable system
into a Material 3 (m3.material.io) role-based token system — without adding a
build step, npm tooling, web fonts, or any non-Go dependency.

## Constraints (decided in brainstorming)

1. **Color identity:** M3 baseline purple seed `#6750A4` (Tone 40). Full M3
   tonal system; no custom brand palette.
2. **Typography:** System font stack only
   (`-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif`).
   No web fonts, no Google Fonts CDN, no bundled woff2. Preserves the
   stdlib-only / embedded / offline model.
3. **Scope:** Token system only. Replace ad-hoc CSS vars with a full M3 role
   token set (surface, surface-variant, on-surface, primary/container, outline,
   elevation, shape, motion). Restyle existing components to consume them.
   Pages keep their current structure. No new M3 components, no layout
   restructure, no FAB / Snackbar / NavigationBar additions.
4. **Integration:** `DESIGN.md` at repo root. `AGENTS.md` gets pointers telling
   agents to read it before touching `web/` and to keep
   `web/css/tokens.css` aligned. `@google/design.md` npm linter is **not**
   wired in — keeps the repo Go-only with zero node tooling.
5. **Implementation approach:** Approach A — one flat `--m3-*` token layer.
   Light values in `:root`, dark values overridden in `[data-theme="dark"]`.
   No two-tier source/alias split. Shortest diff that delivers a real M3 system,
   and an agent can diff DESIGN.md against `tokens.css` by eye.

## Deliverables

1. **`DESIGN.md`** (new, repo root) — conforms to the design.md alpha spec.
   - YAML front matter: `name`, `description`, `colors` (full M3 role set),
     `typography` (body-md, title-lg, title-md, label-lg, label-caps),
     `rounded` (none/xs/sm/md/lg/xl/full), `spacing` (none/xs/sm/md/lg/xl),
     `components` (button-filled, button-tonal, button-outlined, card-elevated,
     text-field-outlined, checkbox, snackbar).
   - Markdown body in canonical order: Overview, Colors, Typography,
     Layout & Spacing, Elevation & Depth, Shapes, Components, Do's and Don'ts.
   - Elevation and motion token values are documented in the body (the
     design.md front-matter schema does not define elevation/motion keys; the
     linter treats unknown keys as warnings, so we keep them in prose).

2. **`web/css/tokens.css`** (new, ~250 lines) — flat `--m3-*` custom properties.
   - `:root` — light theme: every M3 color role, 5-level elevation
     (box-shadow strings), shape scale, spacing scale, type shorthands
     (`--m3-font-body-md` etc.), motion durations + easings.
   - `[data-theme="dark"]` — overrides for every color role and every
     elevation (darker shadow strings). Shape / spacing / type / motion are
     theme-invariant and so not repeated.
   - One documented non-M3 raw value: `--success: #4caf50` (used only for the
     `.info-section li::before` check-marks). M3 defines no "success" role;
     the comment `/* non-M3: … */` marks it inline.

3. **`web/css/style.css`** — pure variable swap, structure unchanged.
   The component CSS (`.card`, `.btn`, `.form-group`, `.tabs`, …) keeps its
   selectors and rules; only the var references change.

   | Old var                     | New M3 role                              |
   |-----------------------------|------------------------------------------|
   | `--bg-primary`              | `--m3-surface`                            |
   | `--bg-secondary`            | `--m3-surface-container`                  |
   | `--bg-card`                 | `--m3-surface-container`                  |
   | `--text-primary`            | `--m3-on-surface`                         |
   | `--text-secondary`          | `--m3-on-surface-variant`                 |
   | `--text-muted`              | `--m3-on-surface-variant` with `opacity: 0.7` on the element (keeps the muted-vs-secondary visual distinction; M3 has no tertiary text role) |
   | `--border`                  | `--m3-outline-variant`                    |
   | `--accent` / `--accent-hover` | `--m3-primary`; hover via a `:hover` rule that switches `background` to `--m3-primary-container` + `color: var(--m3-on-primary-container)` (M3 state-layer pattern, no custom brightness hacks) |
   | `--accent-light`            | `--m3-primary-container`                  |
   | `--error`                   | `--m3-error`                              |
   | `--shadow` / `--shadow-lg`  | `--m3-elevation-1` / `--m3-elevation-3`   |
   | `--radius` / `--radius-sm`  | `--m3-radius-md` / `--m3-radius-sm`       |
   | `--transition`              | `var(--m3-duration-short2) var(--m3-ease-standard)` |

   Four one-line behavioral upgrades come free with the rework (no JS):

   1. **Pill buttons.** `.btn.primary`, `.btn.secondary` →
      `border-radius: var(--m3-radius-full); padding: 10px 24px;` per
      DESIGN.md `button-filled` / `button-tonal`.
   2. **Elevation cards.** `.card` → `--m3-elevation-1`; `:hover` →
      `--m3-elevation-3` (today's `--shadow-lg` rename).
   3. **Outlined text-field focus.** `:focus` border → `--m3-primary` 2px;
      no box-shadow (M3 outlined-field pattern drops the glow halo in favor
      of a thicker primary border).
   4. **Body background** = `--m3-background` (M3 distinguishes background
      from surface; current code blurs them).

4. **`web/index.html`, `web/configure.html`** — add
   `<link rel="stylesheet" href="/css/tokens.css">` immediately **before**
   the existing `<link rel="stylesheet" href="/css/style.css">` in both files.
   No other HTML changes.

5. **`AGENTS.md`** updates:
   - New "UI / Design System" section (after "Security & Trust Boundaries",
     before "Common Pitfalls") explaining that `DESIGN.md` is the source of
     truth, `tokens.css` is its 1:1 reflection, dark theme lives under
     `[data-theme="dark"]`, non-M3 additions are documented inline,
     elevation/radius/spacing come from `--m3-*` scales, and typography uses
     the system font stack only.
   - The "Change the UI" workflow (lines 65-68) gets rewritten to start at
     `DESIGN.md`, keep it and `tokens.css` in lockstep, and consume `--m3-*`
     tokens rather than re-introducing ad-hoc vars.
   - One Coding Conventions bullet about CSS naming.
   - Two Decision Log entries (M3 token system adopted; system font stack
     only).

## M3 Token Values

All token values are defined in `DESIGN.md` front matter (color, typography,
rounded, spacing) and `tokens.css` (elevation, motion) per the rules above.
The implementation plan will list every value exactly; the values are the
Material 3 reference palette derived from the `#6750A4` seed, the M3 elevation
box-shadow reference set, the M3 shape/spacing/motion scales, and no custom
inventions beyond the documented `--success` non-M3 green.

## Out of Scope

- No `cmd/server` or Go-side changes.
- No JS changes (`web/js/app.js` unchanged — `[data-theme]` toggle keeps
  working because the new dark overrides live in `tokens.css`).
- No new M3 components (FAB, Snackbar, NavigationBar, Dialog) and no layout
  reflow.
- No shadcn, no Tailwind, no build step, no npm tooling, no linter wiring,
  no web fonts.
- No tests; CSS-only visual change. `go test ./...` is unaffected.

## Verification

- `go run main.go` + manual load of `/` and `/configure` in light and dark
  modes; confirm pill buttons, elevated cards, focused outlined text fields,
  and check-marks.
- Visual diff against the existing UI should show: same layout, same
  components, same dark/light toggle, M3-toned colors, pill buttons, M3
  elevations.
- No link to any external font/CDN; assets panel stays on
  `/css/tokens.css` and `/css/style.css` only.
- `go test ./...` still passes (sanity, should be a no-op).

## Risk

Low. Structural HTML/JS untouched. Var rename is mechanical. The one
behavioral change is a visual upgrade with no logic change. The
`--success` green is the only non-M3 leave-behind and is documented inline.