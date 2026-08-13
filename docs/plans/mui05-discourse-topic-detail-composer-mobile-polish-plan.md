# MUI-05 Discourse Topic Detail and Composer Mobile Polish Plan

Date: 2026-08-12

## 1. Goal

MUI-05 focuses on the logged-in reading-to-writing loop after MUI-03A and MUI-04 have been synced to the default Discourse theme:

- Topic detail page
- Reply composer
- New-topic composer
- Narrow mobile regression checks for list, category, search, hamburger menu, and user menu

The goal is to make topic reading and writing feel calmer and more predictable on phones while preserving native Discourse behavior.

## 2. Current Screenshots

Captured before implementation:

- Latest entry: `work/mui05-plan-screenshots/01-latest-entry-390x844.png`
- Topic detail: `work/mui05-plan-screenshots/02-topic-detail-390x844.png`
- Reply composer: `work/mui05-plan-screenshots/03-reply-composer-390x844.png`
- New-topic composer: `work/mui05-plan-screenshots/04-new-topic-composer-390x844.png`
- Topic detail 375px: `work/mui05-plan-screenshots/05-topic-detail-375x667.png`
- Topic detail desktop smoke: `work/mui05-plan-screenshots/06-topic-detail-desktop-1440x900.png`

Measured baseline:

- Max horizontal overflow: `0`
- Login page leakage: `0`

## 3. Confirmed Constraints

- CSS/SCSS only.
- No theme JavaScript.
- No template overrides.
- No plugin changes.
- Do not post, reply, like, delete, upload, or change user data during validation.
- Use preview theme first. Sync to default theme only after explicit user approval.
- Do not read browser cookies, localStorage, passwords, tokens, or private message contents.
- Before logged-in Discourse validation, first visit the platform community route `/community` to let the frontend restore the SSO cookie and complete the Discourse session.

## 4. Observed Problems

### Topic Detail

- Topic detail is functionally usable, but the lower action area is visually busy on mobile.
- The topic timeline and footer controls compete with post action buttons.
- Post action icons are visible but feel weak compared with the primary reply affordance.
- Jump/back tooltips can visually cover content during automated or manual interaction.
- Long unbroken content still needs defensive wrapping even though current samples do not overflow.

### Reply Composer

- The reply composer takes a very large vertical area and reads as a raw editor block rather than a compact mobile writing surface.
- The toolbar is usable but visually dense and needs clearer horizontal-scroll boundaries.
- The bottom action row is functional, but spacing between reply, delete, upload, and image buttons can feel disconnected.
- The textarea area is taller than necessary for short replies, which pushes surrounding context out of view.

### New-Topic Composer

- Title, category, tag, toolbar, editor body, and action buttons are all visible.
- The editor body is still too tall for an initial empty draft state.
- Category and tag controls need consistent touch sizing and stronger containment.
- The publish button is clear, but the overall bottom action row can be tightened.

## 5. CSS Work Items

Primary file:

- `discourse-themes/ai-forum-premium-preview/mobile/mobile.scss`

Optional only if a shared token is needed:

- `discourse-themes/ai-forum-premium-preview/common/common.scss`

Planned CSS areas:

1. Topic detail reading shell
   - Strengthen mobile title spacing and wrapping.
   - Improve post stream spacing without cardifying every post.
   - Keep cooked body text readable and defensive against long strings, code, links, and images.

2. Post action controls
   - Normalize post action buttons to clear 40px to 44px touch targets.
   - Keep icon color and SVG fill visible in light/dark states.
   - Make reply action visibly primary without inventing new behavior.

3. Topic footer and timeline controls
   - Reduce visual competition between footer tools and reply action.
   - Keep timeline widgets contained inside the mobile viewport.
   - Avoid page-level overflow from tooltips and progress controls.

4. Reply composer
   - Make the mobile composer shell more compact and viewport-aware.
   - Cap editor body height and use internal scrolling.
   - Keep toolbar as a single horizontal row with clear overflow handling.
   - Align bottom actions into a stable, reachable row.

5. New-topic composer
   - Normalize title/category/tag controls.
   - Reduce empty editor body height.
   - Preserve clear publish, discard, and media controls.
   - Keep the composer usable with the mobile keyboard open where CSS can help.

## 6. Validation Plan

Primary mobile viewport:

- `390 x 844`

Secondary mobile viewport:

- `375 x 667`

Desktop smoke:

- `1440 x 900`

Pages and states:

- `/latest`
- Topic detail page from the latest list
- Reply composer opened from topic detail
- New-topic composer opened from `/latest`
- `/categories`
- `/search`
- Hamburger menu
- User menu

Validation checks:

- No horizontal overflow.
- No login or permission page leakage.
- Composer open states have visible close/minimize/publish controls.
- Toolbar buttons remain visible and horizontally scrollable when needed.
- Reply/new-topic buttons remain at least 44px-class touch targets where practical.
- Desktop has no obvious regression from mobile-only CSS.

## 7. Stop Conditions

Stop and ask before continuing if:

- A required fix needs Discourse template changes, theme JavaScript, or plugin code.
- Validation would require reading private messages, cookies, localStorage, passwords, or tokens.
- A composer issue can only be solved by changing Discourse behavior.
- Preview validation shows regression in MUI-04 search/category/menu surfaces.
- The next step would sync MUI-05 into the default live theme.

## 8. Delivery Flow

1. Patch CSS locally.
2. Build/import as a preview theme, for example `ai-forum-mui05-preview`.
3. Validate preview theme with logged-in mobile screenshots.
4. Fix any CSS-only issues revealed by preview validation.
5. Commit and push MUI-05 changes.
6. Ask for explicit approval before syncing MUI-05 into default theme ID `1`.
