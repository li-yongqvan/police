# MUI-07 Discourse Top Lists and Final Mobile Sweep Plan

Date: 2026-08-14

## 1. Goal

MUI-07 is the closing sweep of the Discourse mobile polish series (MUI-03A/04/05/06 already live on the default theme). It covers:

- The `/top` topic list, which is the last audited surface with a real mobile overflow defect.
- A final full-surface regression sweep at 390/375/320/1440 to prove MUI-04/05/06 do not regress.
- Migrating the deployment flow to the new backup/sync toolchain now that preview themes have been removed.

## 2. Current Screenshots and Audit Data

Pre-implementation audit, captured 2026-08-14:

- `work/mui07-audit/01-latest-390.png` (hamburger menu open)
- `work/mui07-audit/02-latest-search-390.png` (search dropdown open)
- `work/mui07-audit/03-categories-390.png` / `03b-375` / `03c-1440`
- `work/mui07-audit/04-category-390.png` / `04b-375` (`/c/ai/5`)
- `work/mui07-audit/05-new-390.png` / `06-unread-390.png` / `07-top-390.png`
- `work/mui07-audit/08-badges-390.png` / `09-about-390.png`
- `work/mui07-audit/10-preferences-390.png` / `10b-375` / `10c-1440`
- `work/mui07-audit/11-notifications-390.png`
- `work/mui07-audit/12-search-results-390.png` / `12b-375`
- `work/mui07-audit/13-topic-recheck-390.png`
- `work/mui07-audit/14-signup-390.png` / `15-faq-390.png` / `16-tos-390.png`
- 320px sweep results logged in the run output; metrics in `work/mui07-audit/metrics-audit.json`

Measured baseline across 17 scenes at 390/375/1440 plus 12 scenes at 320px:

- Only defect: `/top` period chooser (`period-chooser`) overflows horizontally at 390px and 320px (`bodyOverflow=36` / `106`; chooser right edge at 426px while viewport is 390px).
- All other surfaces: `docOverflow=0`, `bodyOverflow=0`.
- Known cosmetic item (not a layout bug): the decorative `gx-auth-spotify__orb` on login/signup extends past the right edge by design; page does not scroll because `html/body` already clip horizontal overflow.

## 3. Confirmed Constraints

- CSS/SCSS only, no theme JavaScript, no template overrides, no plugin changes.
- Do not post, reply, like, delete, upload, or change user data during validation.
- Do not read cookies, localStorage, passwords, tokens, or private message contents.
- No preview theme anymore: all changes go to the default theme (ID 1) via the new toolchain; a full theme backup is taken before every sync (mandatory).
- Keep MUI-04 search/category/menu and MUI-05/06 behavior unchanged; new rules live in a scoped `// MUI-07:` block at the end of `mobile.scss`.

## 4. Observed Problems

### `/top` Topic List

- The time-range selector (`.top-lists .period-chooser`, a select-kit dropdown) has a fixed wide width and pushes past the right edge at 390px and below.
- At 320px the overflow doubles, confirming there is no existing responsive rule for this control.
- The control must remain discoverable and tappable: fix by allowing it to shrink/ellipsize, not by hiding it.

### Optional (user decision, default skip)

- Login/signup decorative orb partially offscreen: intended visual design, no scroll impact. Default recommendation is to leave it unchanged.

## 5. CSS Work Items

Primary file: `discourse-themes/ai-forum-premium-preview/mobile/mobile.scss`, new scoped block `// MUI-07: top lists + final sweep`.

1. Contain `.top-lists .period-chooser` and its `.select-kit-header` to available width (`width: auto; max-width: 100%`) and allow `.selected-name` to truncate with `text-overflow: ellipsis` on narrow viewports.
2. Preserve the existing MUI-04 no-text control sizing inside `.top-lists .list-controls`; add only width containment, not new icon/text rules.
3. Add a 320px guard mirroring the existing narrow-viewport guard pattern so the chooser never exceeds the viewport.
4. Do not touch MUI-03A/04/05/06 selectors. If the DOM layout of the chooser requires a structure-specific selector, verify it live with a scoped Playwright DOM dump first.

## 6. Validation

Reuse the MUI-06 harness pattern with logged-in and logged-out persistent profiles:

1. Before/after screenshots of `/top` at 390px and 320px for the user.
2. Full sweep at 390/375/320/1440 across: latest + menus, categories, `/c/ai/5`, `/new`, `/unread`, `/top`, badges, about, preferences/account, notifications, search results, topic detail, signup, faq, tos.
3. Hard gates: `docOverflow=0` and `bodyOverflow=0` on every scene; no new `bad` elements outside the approved decorative orb.
4. MUI-04/05/06 regression: re-run the MUI-06 default-theme validation suite.

## 7. Rollout

1. Commit plan and CSS locally, push `codex/discourse-rebuild` after validation.
2. `discourse-themes/tools/backup-theme1.ps1` (mandatory backup before sync).
3. `discourse-themes/tools/sync-field.ps1 -FieldFile ...\mobile\mobile.scss` to sync the mobile field into default theme ID 1.
4. Wait for stylesheet digest propagation, then re-run validation against the default theme and show before/after screenshots.

## 8. Out of Scope

- Login/signup decorative orb (pending user decision; default unchanged).
- Admin dashboard and desktop-only surfaces.
- Any change to community data, plugins, or Discourse core behavior.
