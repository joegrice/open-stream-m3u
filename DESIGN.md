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
- **error / on-error / error-container / on-error-container** — reserved for
  validation and failure states; not yet consumed by current components.
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
Features list (`#81c784` in dark theme). It is documented inline in
`tokens.css` and must not propagate.

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
