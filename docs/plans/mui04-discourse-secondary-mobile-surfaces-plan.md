# MUI-04 Discourse Secondary Mobile Surfaces Plan

Date: 2026-08-12

## 1. Goal

MUI-04 covers the mobile pages and panels that were intentionally left out of MUI-03A:

- Search form, search panel, search results, empty/error states
- Hamburger menu
- User menu, notifications, and private-message entry points
- Categories page
- Read-only regression for the latest list and composer entry points

The goal is to make these secondary surfaces feel consistent with the MUI-03A reading and writing flow without changing Discourse behavior.

## 2. Confirmed Decisions

- CSS/SCSS only.
- No theme JavaScript.
- No template overrides.
- No plugin changes.
- No posting, replying, liking, deleting, or changing user data during validation.
- First deployment target is a preview theme, not the live default theme.

## 3. Preview Theme

- Preview theme name: `ai-forum-mui04-preview`
- Preview theme ID: `3`
- Preview entry: `http://122.51.233.225:8080/?preview_theme_id=3`
- Default live theme remains: `ai-forum-premium-preview`
- Current default theme ID remains: `1`

## 4. Current Local Changes

The implemented MUI-04 CSS changes are in:

- `discourse-themes/ai-forum-premium-preview/mobile/mobile.scss`

Commits:

- `9a41eab fix: align Discourse mobile secondary panels (#MUI-04)`
- `9e8c698 fix: polish Discourse mobile search and categories (#MUI-04)`

## 5. Validation Pages

Primary viewport:

- `390 x 844`

Secondary viewport:

- `375 x 667`

Desktop smoke viewport:

- `1440 x 900`

Pages and states:

- `/latest?preview_theme_id=3`
- `/categories?preview_theme_id=3`
- `/search?preview_theme_id=3`
- Search panel opened from header
- Hamburger menu opened from header
- Search result/empty/error states
- Logged-in user menu
- Logged-in notification panel
- Logged-in private-message entry

## 6. Acceptance Criteria

- No page-level horizontal overflow at `390px` or `375px`.
- Header search, hamburger, login/user controls are visible and at least 44px-class touch targets.
- Search fields and filters do not overflow.
- Search empty/error states do not dominate the mobile viewport with oversized layout.
- Hamburger menu items have readable text, visible icons, and clear grouping.
- Categories page uses readable category blocks and does not mix oversized headings with cramped topic rows.
- User menu, notification, and private-message entries are readable in logged-in mobile state.
- Desktop has no obvious regression from mobile-only rules.

## 7. Stop Conditions

Stop and ask before continuing if:

- A required MUI-04 fix needs template or JavaScript changes.
- Logged-in validation cannot be performed without accessing cookies, localStorage, passwords, or private messages.
- Preview theme validation shows a regression in MUI-03A topic detail or composer flows.
- The next step would synchronize changes into the default live theme.

## 8. Next Step

Run visual validation against preview theme ID `3`, capture screenshots, and list any remaining CSS-only gaps. Synchronization into the default live theme requires separate user approval.
