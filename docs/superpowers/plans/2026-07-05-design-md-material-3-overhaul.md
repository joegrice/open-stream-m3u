# DESIGN.md + Material 3 Token Overhaul — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Adopt the `design.md` format as the visual source of truth and overhaul the Open Stream M3U web UI CSS into a Material 3 role-based token system, with no build step, no frameworks, and no HTML/JS structural changes.

**Architecture:** A new `DESIGN.md` at repo root holds YAML design tokens + rationale prose. A new `web/css/tokens.css` is its flat `--m3-*` CSS reflection (`:root` light + `[data-theme="dark"]` overrides). `web/css/style.css` swaps its ad-hoc vars for the M3 role tokens and gains four one-line behavioral upgrades (pill buttons, elevation cards, outlined-field focus, distinct background). Both HTML pages link `tokens.css` before `style.css`. `AGENTS.md` documents the system and points future agents at it.

**Tech Stack:** Plain CSS custom properties, Go `embed.FS`, standard HTML/JS. No frameworks, no Node, no build step, no web fonts.

## Global Constraints

- **Color identity:** M3 baseline purple seed `#6750A4` (Tone 40); full M3 reference tonal palette derived from it.
- **Typography:** System font stack only — `-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif`. No web fonts.
- **No new dependencies.** Stdlib/embedded-only repo. No npm tooling, no `@google/design.md` linter wired in. The repo stays Go-only.
- **No HTML structural changes.** Both pages keep their existing elements; only `<link>` lines are added.
- **No JS changes.** `web/js/app.js` is untouched. The existing `[data-theme="dark"]` toggle keeps working because dark overrides live in `tokens.css`.
- **1:1 token mapping.** Every `colors.<role>` in `DESIGN.md` front matter becomes `--m3-<role>` in `tokens.css`. Editing one means editing the other in the same commit.
- **Flat layer.** No two-tier source/alias split. Light values in `:root`, dark overrides in `[data-theme="dark"]`.
- **Commit style:** Match the repo's terse `Verb: subject` style (e.g. `Add …`, `Change …`, `Fix …`).

---

### Task 1: Create `DESIGN.md` at repo root

**Files:**
- Create: `DESIGN.md`

**Interfaces:**
- Produces: the YAML front matter `colors` / `typography` / `rounded` / `spacing` / `components` token blocks that Task 2 reflects into CSS. Token names used by Tasks 2 and 4: `primary`, `on-primary`, `primary-container`, `on-primary-container`, `secondary`, `on-secondary`, `secondary-container`, `on-secondary-container`, `tertiary`, `on-tertiary`, `tertiary-container`, `on-tertiary-container`, `error`, `on-error`, `error-container`, `on-error-container`, `background`, `on-background`, `surface`, `on-surface`, `surface-variant`, `on-surface-variant`, `outline`, `outline-variant`, `surface-container`, `surface-container-high`, `inverse-surface`, `inverse-on-surface`. Typography tokens: `body-md`, `title-lg`, `title-md`, `label-lg`, `label-caps`. Rounded scale: `none`, `xs`, `sm`, `md`, `lg`, `xl`, `full`. Spacing scale: `none`, `xs`, `sm`, `md`, `lg`, `xl`.

- [ ] **Step 1: Write `DESIGN.md`**

Create `/home/joe/Code/open-stream-m3u/DESIGN.md` with this exact content (YAML front matter + eight markdown sections in canonical order):

```markdown
---
name: Open Stream M3U
description: Material 3 design system for the Open Stream M3U web UI. Light and dark themes, system fonts only, embedded static assets.
version: "alpha"
colors:
  primary: "#6750A4"
  on-primary: "#FFFFFF"
  primary-container: "#EADDFF"
  on-primary-container: "#21005D"
  secondary: "#625B71"
  on-secondary: "#FFFFFF"
  secondary-container: "#E8DEF8"
  on-secondary-container: "#1D192B"
  tertiary: "#7D5260"
  on-tertiary: "#FFFFFF"
  tertiary-container: "#FFD8E4"
  on-tertiary-container: "#31111D"
  error: "#B3261E"
  on-error: "#FFFFFF"
  error-container: "#F9DEDC"
  on-error-container: "#410E0B"
  background: "#FEF7FF"
  on-background: "#1C1B1F"
  surface: "#FEF7FF"
  on-surface: "#1C1B1F"
  surface-variant: "#E7E0EC"
  on-surface-variant: "#49454F"
  outline: "#79747E"
  outline-variant: "#CAC4D0"
  surface-container: "#F3EDF7"
  surface-container-high: "#E6E0E9"
  inverse-surface: "#313033"
  inverse-on-surface: "#F4EFF4"
typography:
  body-md:
    fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif"
    fontSize: 0.875rem
    fontWeight: 400
    lineHeight: 1.25rem
  title-lg:
    fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif"
    fontSize: 1.375rem
    fontWeight: 400
    lineHeight: 1.75rem
  title-md:
    fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif"
    fontSize: 1rem
    fontWeight: 500
    lineHeight: 1.25rem
  label-lg:
    fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif"
    fontSize: 0.875rem
    fontWeight: 500
    lineHeight: 1.25rem
    letterSpacing: 0.006rem
  label-caps:
    fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif"
    fontSize: 0.6875rem
    fontWeight: 500
    letterSpacing: 0.078rem
rounded:
  none: 0px
  xs: 4px
  sm: 8px
  md: 12px
  lg: 16px
  xl: 28px
  full: 9999px
spacing:
  none: 0
  xs: 4px
  sm: 8px
  md: 16px
  lg: 24px
  xl: 32px
components:
  button-filled:
    backgroundColor: "{colors.primary}"
    textColor: "{colors.on-primary}"
    rounded: "{rounded.full}"
    padding: "10px 24px"
  button-tonal:
    backgroundColor: "{colors.secondary-container}"
    textColor: "{colors.on-secondary-container}"
    rounded: "{rounded.full}"
    padding: "10px 24px"
  button-outlined:
    backgroundColor: "transparent"
    textColor: "{colors.primary}"
    rounded: "{rounded.full}"
    padding: "10px 24px"
  card-elevated:
    backgroundColor: "{colors.surface-container}"
    textColor: "{colors.on-surface}"
    rounded: "{rounded.md}"
    padding: 24px
  text-field-outlined:
    backgroundColor: "transparent"
    textColor: "{colors.on-surface}"
    rounded: "{rounded.xs}"
    padding: "12px 16px"
  checkbox:
    textColor: "{colors.primary}"
    rounded: "{rounded.sm}"
  snackbar:
    backgroundColor: "{colors.inverse-surface}"
    textColor: "{colors.inverse-on-surface}"
    rounded: "{rounded.xs}"
---

## Overview

Open Stream M3U is styled with Material 3 (m3.material.io). The UI is plain
HTML/CSS/JS embedded into a Go binary — no build step, no frameworks, no web
fonts. This document is the source of truth for the visual identity;
`web/css/tokens.css` is its CSS reflection. Edit them together.

The aesthetic is the baseline M3 system: a purple seed (`#6750A4`), full tonal
palette, surface elevation, pill-shaped buttons, outlined text fields.
Forward-compatibility over novelty — when M3 revises the spec, update the
values, not the system.

## Colors

The palette is the Material 3 reference derived from the `#6750A4` seed (Tone
40). Light theme values are the tokens in the front matter above; dark theme
overrides live in `tokens.css` under `[data-theme="dark"]`.

- **primary / on-primary / primary-container / on-primary-container** — the
  accent system. CTAs, focused field borders, active tab indicator, links.
- **secondary** family — neutral-violet; used for tonal buttons and tabs.
- **tertiary** family — rosy accent; reserved for future emphasis.
- **error / on-error / error-container / on-error-container** — validation
  and failure states (only one such state today).
- **background / on-background** — page canvas. Distinct from surface per
  M3; today they share a value but the roles must not be conflated.
- **surface / on-surface / surface-variant / on-surface-variant** — card and
  text-field chrome; `surface-variant` is the muted text ground.
- **outline / outline-variant** — outlines and hairline dividers. Use
  `outline-variant` for borders, `outline` for true outlines.
- **surface-container / surface-container-high** — elevated surfaces
  (raised cards, dialogs). The higher container is one tonal step darker
  than the lower one.
- **inverse-surface / inverse-on-surface** — snackbar / inverse emphasis.

A single non-M3 color survives: `#4caf50` for the green check-marks in the
Features list. It is documented inline in `tokens.css` and must not propagate.

## Typography

System font stack only. No web fonts, no Google Fonts, no bundling. Material 3
permits system fonts and the metrics of `-apple-system`/`Segoe UI`/`Roboto`
are close enough not to break M3 spacing.

Five text styles: `body-md` (default body), `title-lg` (page header),
`title-md` (card heading), `label-lg` (form labels, buttons), `label-caps`
(stat labels, eyebrows). All five share one `fontFamily`; they differ in
size, weight, and letter-spacing.

Line heights are unitless `rem` values (NOT multipliers) so the type ramp
stays rooted in the user's default font size.

## Layout & Spacing

A single 5-step scale: `0 / 4 / 8 / 16 / 24 / 32` px. One container
`max-width: 900px`, centered. Cards live in a `repeat(auto-fit, minmax(280px,
1fr))` grid; on `<600px` screens it collapses to one column. Form padding
`24px` matches `card-elevated.padding`. Mobile breakpoint at `600px`.

Motion durations and easings are M3 reference values, exposed as CSS vars in
`tokens.css` (not in this front matter because `design.md` alpha schema does
not define motion): `duration-short2 100ms`, `duration-short3 150ms`,
`duration-medium2 200ms`, `duration-medium4 250ms`, `duration-long2 500ms`;
`ease-standard` and `ease-emphasized` are both `cubic-bezier(0.2, 0, 0, 1)`.

## Elevation & Depth

Five M3 elevation levels, each a two-shadow composite (ambient + key), exposed
as `--m3-elevation-1` … `--m3-elevation-5` in `tokens.css`. Light theme uses
the M3 reference alphas (0.15–0.30); dark theme bumps them to 0.35–0.60 so
shadows read against the dark canvas. Use elevation, never ad-hoc shadows.

- **elevation-1** — raised cards (default state).
- **elevation-3** — card hover, dialogs.
- Higher levels are reserved and currently unused.

## Shapes

Six corner radii: `0 / 4 / 8 / 12 / 16 / 28 / 9999` px named `none / xs / sm
/ md / lg / xl / full`. Pills (`full`) for buttons; `md` for cards; `xs` for
text fields; `sm` for checkboxes. Do not invent intermediate radii.

## Components

- **button-filled** — primary CTA ( Install Addon , Open in Stremio ).
- **button-tonal** — secondary action ( Configure , Load Categories , Copy ).
- **button-outlined** — reserved for future tertiary actions; not used today.
- **card-elevated** — the two landing cards and all form/panel surfaces.
- **text-field-outlined** — URL / text / password / number inputs.
- **checkbox** — content-type and EPG toggles.
- **snackbar** — reserved for future user feedback; not wired today.

Hover states use the M3 state-layer pattern: switch `background` to
`primary-container` and `color` to `on-primary-container` (filled buttons) so
the surface does not depend on opacity blending. Focus on text fields is
expressed by a thicker `--m3-primary` border, no glow halo — M3's outlined
field drops the box-shadow halo of older Material variants.

## Do's and Don'ts

- **Do** edit `DESIGN.md` and `web/css/tokens.css` in the same commit when
  a token moves.
- **Do** add new components to the `components` block in `DESIGN.md` and
  emit their CSS in `tokens.css` or `style.css` using `--m3-*` tokens.
- **Don't** re-introduce ad-hoc vars (`--bg-*`, `--accent*`, `--text-*`),
  hand-rolled shadows, or off-scale radii. Use the M3 role tokens.
- **Don't** add web fonts, font CDNs, or `@import` rules. The system stack is
  the only stack.
- **Don't** use `!important`. Token cascade is the only override mechanism.
- **Don't** invent a `success` role. The single documented `--success` green
  is the non-M3 exception, kept for the Features check-marks.
```

- [ ] **Step 2: Verify the file parses as expected**

Run: `head -5 DESIGN.md && echo "---" && wc -l DESIGN.md`
Expected: shows `---` front-matter fence on line 1, total ~200 lines.

- [ ] **Step 3: Commit**

```bash
git add DESIGN.md
git commit -m "Add DESIGN.md: Material 3 design system source of truth"
```

---

### Task 2: Create `web/css/tokens.css`

**Files:**
- Create: `web/css/tokens.css`

**Interfaces:**
- Consumes: every color/typography/rounded/spacing token name from Task 1's `DESIGN.md` front matter. The mapping rule is `colors.<role>` → `--m3-<role>`, `typography.<style>` → `--m3-font-<style>` (CSS shorthand), `rounded.<level>` → `--m3-radius-<level>`, `spacing.<level>` → `--m3-space-<level>`.
- Produces: the `--m3-*` custom properties that Task 4 renames the `style.css` rules against. Available names (must match exactly): `--m3-primary`, `--m3-on-primary`, `--m3-primary-container`, `--m3-on-primary-container`, `--m3-secondary`, `--m3-on-secondary`, `--m3-secondary-container`, `--m3-on-secondary-container`, `--m3-tertiary`, `--m3-on-tertiary`, `--m3-tertiary-container`, `--m3-on-tertiary-container`, `--m3-error`, `--m3-on-error`, `--m3-error-container`, `--m3-on-error-container`, `--m3-background`, `--m3-on-background`, `--m3-surface`, `--m3-on-surface`, `--m3-surface-variant`, `--m3-on-surface-variant`, `--m3-outline`, `--m3-outline-variant`, `--m3-surface-container`, `--m3-surface-container-high`, `--m3-inverse-surface`, `--m3-inverse-on-surface`, `--m3-elevation-1` … `--m3-elevation-5`, `--m3-radius-none` … `--m3-radius-full`, `--m3-space-none` … `--m3-space-xl`, `--m3-font-stack`, `--m3-font-body-md`, `--m3-font-title-lg`, `--m3-font-title-md`, `--m3-font-label-lg`, `--m3-font-label-caps`, `--m3-duration-short2`, `--m3-duration-short3`, `--m3-duration-medium2`, `--m3-duration-medium4`, `--m3-duration-long2`, `--m3-ease-standard`, `--m3-ease-emphasized`.

- [ ] **Step 1: Write `web/css/tokens.css`**

Create `/home/joe/Code/open-stream-m3u/web/css/tokens.css` with this exact content:

```css
:root {
    /* ===== Material 3 color roles — light theme =====
       Mirrors DESIGN.md colors.* 1:1. Edit them together. */

    --m3-primary: #6750A4;
    --m3-on-primary: #FFFFFF;
    --m3-primary-container: #EADDFF;
    --m3-on-primary-container: #21005D;

    --m3-secondary: #625B71;
    --m3-on-secondary: #FFFFFF;
    --m3-secondary-container: #E8DEF8;
    --m3-on-secondary-container: #1D192B;

    --m3-tertiary: #7D5260;
    --m3-on-tertiary: #FFFFFF;
    --m3-tertiary-container: #FFD8E4;
    --m3-on-tertiary-container: #31111D;

    --m3-error: #B3261E;
    --m3-on-error: #FFFFFF;
    --m3-error-container: #F9DEDC;
    --m3-on-error-container: #410E0B;

    --m3-background: #FEF7FF;
    --m3-on-background: #1C1B1F;

    --m3-surface: #FEF7FF;
    --m3-on-surface: #1C1B1F;
    --m3-surface-variant: #E7E0EC;
    --m3-on-surface-variant: #49454F;

    --m3-outline: #79747E;
    --m3-outline-variant: #CAC4D0;

    --m3-surface-container: #F3EDF7;
    --m3-surface-container-high: #E6E0E9;

    --m3-inverse-surface: #313033;
    --m3-inverse-on-surface: #F4EFF4;

    /* non-M3: green check-mark color for the Features list.
       Material 3 has no "success" role; keep this exception rare. */
    --m3-success: #4caf50;

    /* ===== Elevation (M3 reference, two-shadow composite) ===== */
    --m3-elevation-1: 0 1px 2px rgba(0, 0, 0, 0.30), 0 1px 3px 1px rgba(0, 0, 0, 0.15);
    --m3-elevation-2: 0 1px 2px rgba(0, 0, 0, 0.30), 0 2px 6px 2px rgba(0, 0, 0, 0.15);
    --m3-elevation-3: 0 1px 3px rgba(0, 0, 0, 0.30), 0 4px 8px 3px rgba(0, 0, 0, 0.15);
    --m3-elevation-4: 0 2px 3px rgba(0, 0, 0, 0.30), 0 6px 10px 4px rgba(0, 0, 0, 0.15);
    --m3-elevation-5: 0 4px 4px rgba(0, 0, 0, 0.30), 0 8px 12px 6px rgba(0, 0, 0, 0.15);

    /* ===== Shape ===== */
    --m3-radius-none: 0;
    --m3-radius-xs: 4px;
    --m3-radius-sm: 8px;
    --m3-radius-md: 12px;
    --m3-radius-lg: 16px;
    --m3-radius-xl: 28px;
    --m3-radius-full: 9999px;

    /* ===== Spacing ===== */
    --m3-space-none: 0;
    --m3-space-xs: 4px;
    --m3-space-sm: 8px;
    --m3-space-md: 16px;
    --m3-space-lg: 24px;
    --m3-space-xl: 32px;

    /* ===== Typography ===== */
    --m3-font-stack: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
    --m3-font-body-md: 400 0.875rem/1.25rem var(--m3-font-stack);
    --m3-font-title-lg: 400 1.375rem/1.75rem var(--m3-font-stack);
    --m3-font-title-md: 500 1rem/1.25rem var(--m3-font-stack);
    --m3-font-label-lg: 500 0.875rem/1.25rem var(--m3-font-stack);
    --m3-font-label-caps: 500 0.6875rem/1rem var(--m3-font-stack);

    /* ===== Motion (M3 reference) ===== */
    --m3-duration-short2: 100ms;
    --m3-duration-short3: 150ms;
    --m3-duration-medium2: 200ms;
    --m3-duration-medium4: 250ms;
    --m3-duration-long2: 500ms;
    --m3-ease-standard: cubic-bezier(0.2, 0, 0, 1);
    --m3-ease-emphasized: cubic-bezier(0.2, 0, 0, 1);
}

[data-theme="dark"] {
    /* ===== Material 3 color roles — dark theme overrides ===== */

    --m3-primary: #D0BCFF;
    --m3-on-primary: #381E72;
    --m3-primary-container: #4F378B;
    --m3-on-primary-container: #EADDFF;

    --m3-secondary: #CCC2DC;
    --m3-on-secondary: #332D41;
    --m3-secondary-container: #4A4458;
    --m3-on-secondary-container: #E8DEF8;

    --m3-tertiary: #EFB8C8;
    --m3-on-tertiary: #492532;
    --m3-tertiary-container: #633B48;
    --m3-on-tertiary-container: #FFD8E4;

    --m3-error: #F2B8B5;
    --m3-on-error: #601410;
    --m3-error-container: #8C1D18;
    --m3-on-error-container: #F9DEDC;

    --m3-background: #1C1B1F;
    --m3-on-background: #E6E1E5;

    --m3-surface: #1C1B1F;
    --m3-on-surface: #E6E1E5;
    --m3-surface-variant: #49454F;
    --m3-on-surface-variant: #CAC4D0;

    --m3-outline: #938F99;
    --m3-outline-variant: #49454F;

    --m3-surface-container: #211F26;
    --m3-surface-container-high: #2B2930;

    --m3-inverse-surface: #E6E1E5;
    --m3-inverse-on-surface: #313033;

    /* non-M3: green check-mark (see :root comment) */
    --m3-success: #81c784;

    /* Dark theme elevations: higher alpha so shadows read on dark canvas */
    --m3-elevation-1: 0 1px 3px rgba(0, 0, 0, 0.50), 0 1px 2px rgba(0, 0, 0, 0.35);
    --m3-elevation-2: 0 1px 3px rgba(0, 0, 0, 0.55), 0 2px 6px 2px rgba(0, 0, 0, 0.40);
    --m3-elevation-3: 0 1px 3px rgba(0, 0, 0, 0.60), 0 4px 8px 3px rgba(0, 0, 0, 0.45);
    --m3-elevation-4: 0 2px 3px rgba(0, 0, 0, 0.60), 0 6px 10px 4px rgba(0, 0, 0, 0.50);
    --m3-elevation-5: 0 4px 4px rgba(0, 0, 0, 0.60), 0 8px 12px 6px rgba(0, 0, 0, 0.55);
}
```

- [ ] **Step 2: Sanity-check the file**

Run: `wc -l web/css/tokens.css`
Expected: ~110 lines, no parse errors.

- [ ] **Step 3: Commit**

```bash
git add web/css/tokens.css
git commit -m "Add Material 3 token layer (web/css/tokens.css)"
```

---

### Task 3: Link `tokens.css` from both HTML pages

**Files:**
- Modify: `web/index.html:7`
- Modify: `web/configure.html:8`

**Interfaces:**
- Consumes: the `/css/tokens.css` route served by the existing `embed.FS` static handler (no Go change needed — `web/` is the embed root).
- Produces: the `--m3-*` cascade for Task 4's `style.css` rework to consume.

- [ ] **Step 1: Insert the link in `web/index.html`**

In `/home/joe/Code/open-stream-m3u/web/index.html`, change:

```html
    <link rel="stylesheet" href="/css/style.css">
```

to:

```html
    <link rel="stylesheet" href="/css/tokens.css">
    <link rel="stylesheet" href="/css/style.css">
```

(tokens must load before style so `style.css` can resolve `--m3-*`.)

- [ ] **Step 2: Insert the link in `web/configure.html`**

Same edit at `web/configure.html:8`.

- [ ] **Step 3: Verify both pages still load**

```bash
go run main.go &
SERVER_PID=$!
sleep 1
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:7001/css/tokens.css
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:7001/
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:7001/configure
kill $SERVER_PID
```
Expected: `200` for all three requests.

- [ ] **Step 4: Commit**

```bash
git add web/index.html web/configure.html
git commit -m "Link tokens.css before style.css on both pages"
```

---

### Task 4: Rework `web/css/style.css` to consume M3 tokens

**Files:**
- Modify: `web/css/style.css` (entire file — var swap + 4 behavioral upgrades)

**Interfaces:**
- Consumes: every `--m3-*` token produced by Task 2.
- Produces: the visible M3-restyled UI. No new selectors, no new class names, no structural CSS changes — only var references and four one-line behavioral upgrades.

**Mapping (from the spec):**

| Old var                     | New M3 role                              |
|-----------------------------|------------------------------------------|
| `--bg-primary`              | `--m3-surface`                            |
| `--bg-secondary`            | `--m3-surface-container`                  |
| `--bg-card`                 | `--m3-surface-container`                  |
| `--text-primary`            | `--m3-on-surface`                         |
| `--text-secondary`          | `--m3-on-surface-variant`                 |
| `--text-muted`              | `--m3-on-surface-variant` w/ `opacity: 0.7` on the element |
| `--border`                  | `--m3-outline-variant`                    |
| `--accent` / `--accent-hover` | `--m3-primary`; hover uses `--m3-primary-container` + `--m3-on-primary-container` (state-layer pattern) |
| `--accent-light`            | `--m3-primary-container`                  |
| `--success`                 | `--m3-success`                            |
| `--error`                   | `--m3-error`                              |
| `--shadow`                  | `--m3-elevation-1`                        |
| `--shadow-lg`               | `--m3-elevation-3`                        |
| `--radius`                  | `--m3-radius-md`                         |
| `--radius-sm`               | `--m3-radius-sm`                         |
| `--transition`              | `var(--m3-duration-short2) var(--m3-ease-standard)` |

Behavioral upgrades (one line each, applied in the relevant rules):
1. **Pill buttons.** `.btn.primary`, `.btn.secondary` → `border-radius: var(--m3-radius-full); padding: 10px 24px;`
2. **Elevation cards.** `.card` → `box-shadow: var(--m3-elevation-1);` and `.card:hover` → `box-shadow: var(--m3-elevation-3);`
3. **Outlined-field focus.** `.form-group input:focus` → `border-color: var(--m3-primary); border-width: 2px; box-shadow: none;`
4. **Body background.** `body` → `background: var(--m3-background); color: var(--m3-on-background);`

- [ ] **Step 1: Replace `web/css/style.css` in full**

Overwrite `/home/joe/Code/open-stream-m3u/web/css/style.css` with this exact content (the `:root` and `[data-theme="dark"]` blocks are deleted because those tokens now live in `tokens.css`; the rest of the file keeps its rules but reads `--m3-*`):

```css
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: var(--m3-font-stack);
    font: var(--m3-font-body-md);
    background: var(--m3-background);
    color: var(--m3-on-background);
    line-height: 1.6;
    min-height: 100vh;
}

.container {
    max-width: 900px;
    margin: 0 auto;
    padding: var(--m3-space-xl) var(--m3-space-lg);
}

header {
    display: flex;
    align-items: center;
    gap: var(--m3-space-md);
    margin-bottom: var(--m3-space-lg);
    flex-wrap: wrap;
}

header h1 {
    font: var(--m3-font-title-lg);
    color: var(--m3-on-surface);
}

header p {
    color: var(--m3-on-surface-variant);
    font-size: 0.9rem;
}

.back-link {
    color: var(--m3-primary);
    text-decoration: none;
    font-size: 0.9rem;
    margin-left: auto;
    transition: var(--m3-duration-short2) var(--m3-ease-standard);
}

.back-link:hover {
    text-decoration: underline;
}

.icon-btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: var(--m3-space-sm);
    border-radius: var(--m3-radius-full);
    color: var(--m3-on-surface-variant);
    transition: var(--m3-duration-short2) var(--m3-ease-standard);
}

.icon-btn:hover {
    background: var(--m3-primary-container);
    color: var(--m3-on-primary-container);
}

.theme-icon {
    width: 24px;
    height: 24px;
    display: none;
}

[data-theme="light"] .light-icon {
    display: block;
}

[data-theme="dark"] .dark-icon {
    display: block;
}

.card-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: var(--m3-space-lg);
    margin-bottom: var(--m3-space-lg);
}

.card {
    background: var(--m3-surface-container);
    border-radius: var(--m3-radius-md);
    padding: var(--m3-space-lg);
    box-shadow: var(--m3-elevation-1);
    transition: box-shadow var(--m3-duration-short2) var(--m3-ease-standard), transform var(--m3-duration-short2) var(--m3-ease-standard);
    border: 1px solid var(--m3-outline-variant);
}

.card:hover {
    box-shadow: var(--m3-elevation-3);
    transform: translateY(-2px);
}

.card-icon {
    width: 48px;
    height: 48px;
    margin-bottom: var(--m3-space-md);
    color: var(--m3-primary);
}

.card-icon svg {
    width: 100%;
    height: 100%;
}

.card h2 {
    font: var(--m3-font-title-md);
    font-size: 1.25rem;
    margin-bottom: var(--m3-space-sm);
    color: var(--m3-on-surface);
}

.card p {
    color: var(--m3-on-surface-variant);
    font-size: 0.9rem;
    margin-bottom: var(--m3-space-md);
}

.btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 10px 24px;
    border-radius: var(--m3-radius-full);
    font: var(--m3-font-label-lg);
    text-decoration: none;
    cursor: pointer;
    border: none;
    transition: background var(--m3-duration-short2) var(--m3-ease-standard), color var(--m3-duration-short2) var(--m3-ease-standard);
}

.btn.primary {
    background: var(--m3-primary);
    color: var(--m3-on-primary);
}

.btn.primary:hover {
    background: var(--m3-primary-container);
    color: var(--m3-on-primary-container);
}

.btn.secondary {
    background: var(--m3-secondary-container);
    color: var(--m3-on-secondary-container);
}

.btn.secondary:hover {
    background: var(--m3-secondary);
    color: var(--m3-on-secondary);
}

.info-section {
    background: var(--m3-surface-container);
    border-radius: var(--m3-radius-md);
    padding: var(--m3-space-md);
    box-shadow: var(--m3-elevation-1);
    border: 1px solid var(--m3-outline-variant);
}

.info-section h3 {
    font: var(--m3-font-title-md);
    font-size: 1rem;
    margin-bottom: var(--m3-space-sm);
    color: var(--m3-on-surface);
}

.info-section ul {
    list-style: none;
    color: var(--m3-on-surface-variant);
    font-size: 0.9rem;
}

.info-section li {
    padding: 0.25rem 0;
    padding-left: 1.5rem;
    position: relative;
}

.info-section li::before {
    content: "\2713";
    position: absolute;
    left: 0;
    color: var(--m3-success);
}

.tabs {
    display: flex;
    gap: var(--m3-space-sm);
    margin-bottom: var(--m3-space-md);
    border-bottom: 1px solid var(--m3-outline-variant);
    padding-bottom: 0;
}

.tab {
    background: none;
    border: none;
    padding: var(--m3-space-sm) var(--m3-space-md);
    font-size: 0.9rem;
    color: var(--m3-on-surface-variant);
    cursor: pointer;
    border-bottom: 2px solid transparent;
    margin-bottom: -1px;
    transition: color var(--m3-duration-short2) var(--m3-ease-standard), border-color var(--m3-duration-short2) var(--m3-ease-standard);
}

.tab:hover {
    color: var(--m3-on-surface);
}

.tab.active {
    color: var(--m3-primary);
    border-bottom-color: var(--m3-primary);
}

.tab-content {
    display: none;
}

.tab-content.active {
    display: block;
}

.config-form {
    background: var(--m3-surface-container);
    border-radius: var(--m3-radius-md);
    padding: var(--m3-space-lg);
    box-shadow: var(--m3-elevation-1);
    border: 1px solid var(--m3-outline-variant);
}

.form-group {
    margin-bottom: 1.25rem;
}

.form-group label {
    display: block;
    font-size: 0.85rem;
    font-weight: 500;
    color: var(--m3-on-surface-variant);
    margin-bottom: var(--m3-space-sm);
}

.form-group input[type="url"],
.form-group input[type="text"],
.form-group input[type="password"],
.form-group input[type="number"] {
    width: 100%;
    padding: 12px 16px;
    border: 1px solid var(--m3-outline-variant);
    border-radius: var(--m3-radius-xs);
    font-size: 0.9rem;
    background: transparent;
    color: var(--m3-on-surface);
    transition: border-color var(--m3-duration-short2) var(--m3-ease-standard);
}

.password-wrapper {
    position: relative;
    display: flex;
    gap: var(--m3-space-sm);
}

.password-wrapper input {
    flex: 1;
}

.password-toggle {
    padding: 12px 16px;
    border: 1px solid var(--m3-outline-variant);
    border-radius: var(--m3-radius-xs);
    background: var(--m3-secondary-container);
    color: var(--m3-on-secondary-container);
    cursor: pointer;
    font-size: 0.85rem;
    transition: background var(--m3-duration-short2) var(--m3-ease-standard);
}

.password-toggle:hover {
    background: var(--m3-secondary);
    color: var(--m3-on-secondary);
}

.form-group input:focus {
    outline: none;
    border-color: var(--m3-primary);
    border-width: 2px;
    box-shadow: none;
}

.checkbox-label {
    display: flex;
    align-items: center;
    gap: var(--m3-space-sm);
    cursor: pointer;
    font-size: 0.9rem;
    color: var(--m3-on-surface);
}

.checkbox-label input[type="checkbox"] {
    width: 18px;
    height: 18px;
    accent-color: var(--m3-primary);
}

.form-actions {
    margin-top: var(--m3-space-md);
    padding-top: var(--m3-space-md);
    border-top: 1px solid var(--m3-outline-variant);
    display: flex;
    gap: var(--m3-space-sm);
}

.form-actions .btn {
    flex: 1;
}

.overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 100;
    backdrop-filter: blur(4px);
}

.overlay.hidden {
    display: none;
}

.overlay-content {
    background: var(--m3-surface-container-high);
    border-radius: var(--m3-radius-lg);
    padding: var(--m3-space-lg);
    text-align: center;
    max-width: 400px;
    width: 90%;
    box-shadow: var(--m3-elevation-3);
}

.overlay-content h2 {
    margin-bottom: var(--m3-space-md);
    color: var(--m3-on-surface);
}

.progress-bar {
    height: 8px;
    background: var(--m3-surface-container);
    border-radius: var(--m3-radius-full);
    overflow: hidden;
    margin-bottom: var(--m3-space-md);
}

.progress-fill {
    height: 100%;
    background: var(--m3-primary);
    border-radius: var(--m3-radius-full);
    width: 0%;
    transition: width var(--m3-duration-medium4) var(--m3-ease-standard);
}

#progressText {
    color: var(--m3-on-surface-variant);
    font-size: 0.9rem;
}

.result-panel {
    background: var(--m3-surface-container);
    border-radius: var(--m3-radius-md);
    padding: var(--m3-space-lg);
    box-shadow: var(--m3-elevation-1);
    border: 1px solid var(--m3-outline-variant);
    margin-top: var(--m3-space-md);
}

.result-panel.hidden {
    display: none;
}

.result-panel h2 {
    color: var(--m3-primary);
    margin-bottom: var(--m3-space-md);
    font: var(--m3-font-title-md);
    font-size: 1.25rem;
}

.input-group {
    display: flex;
    gap: var(--m3-space-sm);
}

.input-group input {
    flex: 1;
    padding: 12px 16px;
    border: 1px solid var(--m3-outline-variant);
    border-radius: var(--m3-radius-xs);
    font-size: 0.85rem;
    font-family: monospace;
    background: var(--m3-surface-container);
    color: var(--m3-on-surface);
}

.result-panel .btn {
    width: 100%;
    margin-top: var(--m3-space-md);
}

.code-block {
    background: var(--m3-surface-container);
    border: 1px solid var(--m3-outline-variant);
    border-radius: var(--m3-radius-xs);
    padding: var(--m3-space-md);
    font-family: 'Courier New', Courier, monospace;
    font-size: 0.85rem;
    word-break: break-all;
    color: var(--m3-on-surface);
    line-height: 1.5;
    user-select: all;
    cursor: text;
}

.code-block:hover {
    border-color: var(--m3-primary);
}

.stats-box {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: var(--m3-space-md);
    margin-bottom: var(--m3-space-md);
    padding: var(--m3-space-md);
    background: var(--m3-surface-container);
    border-radius: var(--m3-radius-xs);
    border: 1px solid var(--m3-outline-variant);
}

.stat-item {
    text-align: center;
}

.stat-value {
    display: block;
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--m3-primary);
    line-height: 1.2;
    min-height: 1.8rem;
}

.stat-value.loading {
    animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.4; }
}

.stat-label {
    display: block;
    font: var(--m3-font-label-caps);
    color: var(--m3-on-surface-variant);
    opacity: 0.7;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-top: 0.25rem;
}

@media (max-width: 600px) {
    .stats-box {
        grid-template-columns: repeat(2, 1fr);
    }
}

footer {
    margin-top: 3rem;
    padding-top: var(--m3-space-md);
    border-top: 1px solid var(--m3-outline-variant);
    text-align: center;
}

footer p {
    color: var(--m3-on-surface-variant);
    opacity: 0.7;
    font-size: 0.8rem;
}

@media (max-width: 600px) {
    .container {
        padding: var(--m3-space-md);
    }

    header h1 {
        font-size: 1.5rem;
    }

    .card-grid {
        grid-template-columns: 1fr;
    }

    .config-form {
        padding: var(--m3-space-md);
    }
}

.group-selector {
    margin-top: var(--m3-space-md);
    background: var(--m3-surface-container);
    border-radius: var(--m3-radius-xs);
    border: 1px solid var(--m3-outline-variant);
    padding: 1.25rem;
}

.group-selector.hidden {
    display: none;
}

.group-selector-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: var(--m3-space-sm);
}

.group-selector-header h3 {
    font: var(--m3-font-title-md);
    font-size: 1rem;
    color: var(--m3-on-surface);
}

.selected-count {
    font-size: 0.8rem;
    color: var(--m3-primary);
    font-weight: 400;
}

.group-actions {
    display: flex;
    gap: var(--m3-space-sm);
}

.btn-link {
    background: none;
    border: none;
    color: var(--m3-primary);
    cursor: pointer;
    font-size: 0.8rem;
    padding: 0;
}

.btn-link:hover {
    text-decoration: underline;
}

.group-search {
    width: 100%;
    padding: 0.6rem 0.85rem;
    border: 1px solid var(--m3-outline-variant);
    border-radius: var(--m3-radius-xs);
    font-size: 0.85rem;
    background: transparent;
    color: var(--m3-on-surface);
    margin-bottom: var(--m3-space-sm);
}

.group-search:focus {
    outline: none;
    border-color: var(--m3-primary);
    border-width: 2px;
    box-shadow: none;
}

.group-list {
    max-height: 300px;
    overflow-y: auto;
    border: 1px solid var(--m3-outline-variant);
    border-radius: var(--m3-radius-xs);
    background: var(--m3-surface);
}

.group-item {
    display: flex;
    align-items: center;
    gap: var(--m3-space-sm);
    padding: 0.6rem 0.85rem;
    cursor: pointer;
    border-bottom: 1px solid var(--m3-outline-variant);
    font-size: 0.85rem;
    transition: background var(--m3-duration-short2) var(--m3-ease-standard);
}

.group-item:last-child {
    border-bottom: none;
}

.group-item:hover {
    background: var(--m3-primary-container);
}

.group-checkbox {
    accent-color: var(--m3-primary);
    width: 16px;
    height: 16px;
    flex-shrink: 0;
}

.group-name {
    flex: 1;
    color: var(--m3-on-surface);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.group-count {
    color: var(--m3-on-surface-variant);
    opacity: 0.7;
    font-size: 0.8rem;
    flex-shrink: 0;
    min-width: 3rem;
    text-align: right;
}

.group-hint {
    margin-top: var(--m3-space-sm);
    font-size: 0.8rem;
    color: var(--m3-on-surface-variant);
    opacity: 0.7;
}
```

- [ ] **Step 2: Verify no stale old-variable references remain**

```bash
! grep -E -- '--(bg-primary|bg-secondary|bg-card|text-primary|text-secondary|text-muted|border|accent|accent-hover|accent-light|success|error|shadow|shadow-lg|radius|radius-sm|transition)\b' web/css/style.css
```
Expected: command exits 0 (no matches — the negated `!` means a grep hit fails the check).

Also verify `tokens.css` is loaded:

```bash
grep -c 'm3-primary' web/css/style.css
```
Expected: line count > 0 (sanity, the new vars are referenced).

- [ ] **Step 3: Compile check**

```bash
go build ./...
```
Expected: no output (clean build; `embed.FS` picks up the new CSS at compile time but built asset roots are unchanged).

- [ ] **Step 4: Visual verification in both themes**

```bash
go run main.go &
SERVER_PID=$!
sleep 1
echo "Open http://localhost:7001/ and http://localhost:7001/configure"
echo "Toggle the theme icon (top-right). Verify:"
echo "  - body background is M3 surface (light purple-tinted / dark near-black"
echo "  - buttons are pill-shaped (full radius) with M3 purple / secondary-violet fills"
echo "  - cards have elevation-1 by default, lift to elevation-3 on hover"
echo "  - text fields have an outline-variant border; focus turns border to primary 2px, no glow"
echo "  - active tab has primary-colored underline"
echo "  - green check-marks in Features list (non-M3 --m3-success)"
echo "  - theme toggle switches light/dark cleanly, no unstyled flash"
read
kill $SERVER_PID
```

- [ ] **Step 5: Commit**

```bash
git add web/css/style.css
git commit -m "Rework style.css to consume Material 3 role tokens"
```

---

### Task 5: Update `AGENTS.md`

**Files:**
- Modify: `AGENTS.md` (add a section, update two existing sections, append Decision Log entries)

**Interfaces:**
- Produces: documentation that points future agents at `DESIGN.md` as the source of truth and `web/css/tokens.css` as its reflection.

- [ ] **Step 1: Add the new "UI / Design System" section**

In `/home/joe/Code/open-stream-m3u/AGENTS.md`, find the "Security & Trust Boundaries" section and insert this new section immediately after it (before "Common Pitfalls"):

```markdown
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
```

- [ ] **Step 2: Update the "Change the UI" workflow**

Replace the existing "### Change the UI" block (currently lines 65-68):

```markdown
### Change the UI
1. Files under `web/` (HTML, CSS, JS).
2. Keep Material Design styling and dark/light mode support.
3. The UI is served via `embed.FS`; rebuild not required for static changes.
```

with:

```markdown
### Change the UI
1. Files under `web/` (HTML, CSS, JS).
2. **Start at `DESIGN.md`** for any visual change; update it and
   `web/css/tokens.css` together when tokens move.
3. Consume `--m3-*` tokens from `tokens.css`; do not re-introduce ad-hoc
   `--bg-*` / `--accent*` / `--text-*` vars.
4. Material 3 styling, dark/light via `[data-theme]`, pill buttons, M3
   elevation, outlined text fields — see "UI / Design System" above.
5. Served via `embed.FS`; rebuild not required for static changes.
```

- [ ] **Step 3: Add a Coding Conventions bullet**

In `/home/joe/Code/open-stream-m3u/AGENTS.md`, in the "Coding Conventions" section, after the "**Naming:** …" bullet, add:

```markdown
- **CSS:** prefer Material 3 role tokens (`--m3-*`); see `DESIGN.md` and "UI / Design System".
```

- [ ] **Step 4: Append Decision Log entries**

At the end of the "Decision Log" section, append:

```markdown
- M3 token system per `DESIGN.md` spec; flat `--m3-*` CSS reflection in
  `tokens.css`, no aliasing layer, no build step.
- System font stack only for Material 3 — no web fonts, preserves the
  zero-external-deps / embedded model.
```

- [ ] **Step 5: Verify AGENTS.md sections**

```bash
grep -n '^## ' AGENTS.md
```
Expected: the new "UI / Design System" section appears between "Security & Trust Boundaries" and "Common Pitfalls". No duplicate `## ` headers.

- [ ] **Step 6: Commit**

```bash
git add AGENTS.md
git commit -m "Document Material 3 token system in AGENTS.md"
```

### Task 6: Final verification

**Files:** (read-only)

- [ ] **Step 1: Run the Go test suite**

```bash
go test ./...
```
Expected: all packages pass (this is a CSS-only change; tests should be unaffected).

- [ ] **Step 2: Run the Go build**

```bash
go build ./...
```
Expected: no output (clean build; `embed.FS` embeds the new `tokens.css` automatically).

- [ ] **Step 3: Verify the static asset routes serve expected content types**

```bash
go run main.go &
SERVER_PID=$!
sleep 1
curl -sI http://localhost:7001/css/tokens.css | grep -i 'content-type'
curl -sI http://localhost:7001/css/style.css | grep -i 'content-type'
kill $SERVER_PID
```
Expected: both return `Content-Type: text/css` (the existing static handler sets this).

- [ ] **Step 4: Manual cross-theme visual sweep**

```bash
go run main.go &
SERVER_PID=$!
sleep 1
echo "Open http://localhost:7001/ and http://localhost:7001/configure in a browser."
echo "Toggle theme. For EACH page in EACH theme, verify:"
echo "  - layout identical to before the overhaul (same components, same structure)"
echo "  - colors are M3 reference (purple primary, violet secondary, near-white surface in light)"
echo "  - pill buttons (full radius) for primary, secondary, copy, install"
echo "  - elevation-1 cards lift to elevation-3 on hover"
echo "  - focus on inputs shows 2px primary border, NO blue glow halo"
echo "  - green check-marks (the only non-M3 color)"
echo "  - no external font requests in the network panel"
read
kill $SERVER_PID
```

- [ ] **Step 5: Confirm no stray old vars**

```bash
! grep -rE -- '--(bg-primary|bg-secondary|bg-card|text-primary|text-secondary|text-muted|accent|accent-hover|accent-light|shadow|shadow-lg|radius|radius-sm|transition)\b' web/
```
Expected: exits 0 (no matches anywhere under `web/`).

(End of plan)
```