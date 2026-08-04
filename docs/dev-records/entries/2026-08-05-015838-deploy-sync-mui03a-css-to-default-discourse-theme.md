# deploy: sync MUI-03A CSS to default Discourse theme (#MUI-03A)

- Time: 2026-08-05 01:58:38 +08:00
- Category: deploy
- Branch: codex/discourse-rebuild
- Commit: remote-theme-sync
- Scope: Discourse 8080 default theme

## Context

The MUI-03A mobile topic detail and composer CSS had already been implemented, uploaded to the separate preview theme `ai-forum-mui03a-preview`, and accepted by the user after screenshot review. The next planned step was to synchronize the accepted preview CSS into the live default Discourse theme.

## Outcome

- Confirmed remote Discourse default theme:
  - Default theme ID: `1`
  - Default theme name: `ai-forum-premium-preview`
- Confirmed preview theme:
  - Preview theme ID: `2`
  - Preview theme name: `ai-forum-mui03a-preview`
- Confirmed `common.scss` and `desktop.scss` were already identical between the default and preview themes.
- Copied only the mobile SCSS theme field from preview theme ID `2` to default theme ID `1`.
- Wrote a remote rollback backup before the sync:
  - `/shared/tmp/ai-forum-premium-preview-default-mobile-before-mui03a-20260804-175532.scss`
- Cleared Discourse theme cache and precompiled theme CSS.

## Verification

- Verified default theme ID `1` and preview theme ID `2` now have identical mobile SCSS:
  - mobile SCSS length: `10154`
  - SHA-256 prefix: `00a23feb9f19f805`
- Verified public routes returned HTTP 200:
  - `http://122.51.233.225:8080/`
  - `http://122.51.233.225:8080/t/ai/5`
- Ran mobile viewport checks on the formal non-preview URLs:
  - `390x844`: `/`, `/latest`, `/t/ai/5` all had `overflowX = 0`
  - `375x667`: `/`, `/latest`, `/t/ai/5` all had `overflowX = 0`
- A text probe initially matched the word `preview`, but this was from an existing topic title (`Theme preview acceptance...`), not from a Discourse theme preview notice.

## Follow-ups

- Composer interaction was not re-tested with a fresh authenticated browser session during this sync pass. The synced mobile SCSS is byte-for-byte identical to the user-accepted preview theme.
