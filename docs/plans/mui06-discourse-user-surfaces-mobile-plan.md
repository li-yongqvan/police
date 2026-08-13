# MUI-06 Discourse User Surfaces Mobile Polish Plan

Date: 2026-08-13

## 1. Goal

MUI-06 covers the logged-in user-facing surfaces around identity and post-level actions, after MUI-03A/04/05 have been synced to the default Discourse theme:

- User profile page (own profile, logged-in state)
- User card popover (triggered from a topic post author)
- Post action row and post bookmark menu on topic detail
- Notifications and account panels from the header user menu
- Tags page
- Login entry (logged-out state)
- Read-only regression check for `/latest` and all MUI-03A/04/05 surfaces

The goal is to make identity and per-post interactions feel consistent with the MUI-04/05 mobile polish while preserving native Discourse behavior. The pre-implementation audit found exactly one substantive defect — the user-profile activity navigation is clipped on mobile — so the CSS work is narrow and defensive rather than a redesign.

## 2. Current Screenshots

Captured before implementation in `work/mui06-audit/`:

- Latest logged-out 390: `01-latest-loggedout-390x844.png`
- Tags logged-out 390: `02-tags-loggedout-390x844.png`
- Login logged-out 390: `03-login-loggedout-390x844.png`
- Latest logged-out desktop: `04-latest-loggedout-desktop-1440x900.png`
- Latest logged-in 390: `05-latest-loggedin-390x844.png`
- Latest logged-in 375: `05b-latest-loggedin-375x667.png`
- Latest logged-in desktop: `05c-latest-loggedin-desktop-1440x900.png`
- Tags logged-in 390: `06-tags-loggedin-390x844.png`
- Topic detail logged-in 390: `07-topic-detail-loggedin-390x844.png`
- Post bookmark menu logged-in 390: `07b-post-bookmark-menu-loggedin-390x844.png`
- Topic detail logged-in 375: `08-topic-detail-loggedin-375x667.png`
- Topic detail logged-in desktop: `09-topic-detail-loggedin-desktop-1440x900.png`
- User card logged-in 390: `10-user-card-loggedin-390x844.png`
- Notifications logged-in 390: `11-notifications-loggedin-390x844.png`
- User menu notifications logged-in 390: `11b-user-menu-notifications-loggedin-390x844.png`
- User profile logged-in 390: `12-user-profile-loggedin-390x844.png`
- User profile logged-in 375: `13-user-profile-loggedin-375x667.png`
- User profile logged-in desktop: `14-user-profile-loggedin-desktop-1440x900.png`
- Topic detail logged-out 390: `15-topic-detail-loggedout-390x844.png`
- Topic detail logged-out 375: `16-topic-detail-loggedout-375x667.png`
- User profile logged-out 390: `17-user-profile-loggedout-390x844.png`

Measured baseline (all pages, logged-in and logged-out):

- Document and body horizontal overflow: `0` everywhere.
- Only flagged elements: the four activity tabs on the user profile top navigation extend beyond the 390px and 375px viewport width. The outer nav container scrolls internally, so there is no page-level overflow, but the visual tail is cut off.

## 3. Confirmed Constraints

- CSS/SCSS only.
- No theme JavaScript.
- No template overrides.
- No plugin changes.
- Do not post, reply, like, bookmark, delete, upload, or change user data during validation.
- Use preview theme first. Sync to default theme only after explicit user approval.
- Do not read browser cookies, localStorage, passwords, tokens, or private message contents.
- Before logged-in Discourse validation, first visit the platform community route `/community` to let the frontend restore the SSO cookie and complete the Discourse session.

## 4. Observed Problems

### User Profile — Activity Navigation (substantive)

- `.user-nav` activity tabs (获赞/书签/回应/已解决列表, i.e. `.user-nav__activity-likes`, `__activity-bookmarks`, the reactions outlet, and `solved-list`) do not wrap or shrink on narrow screens.
- Measured at 390px viewport: likes `left=375 right=427`, bookmarks `left=434 right=499`, reactions `left=506 right=572`, solved-list `left=579 right=659`. The tab row is a scrollable container (`docOverflow=0`), but the trailing tabs are visually clipped with no affordance hint.
- At 375px the geometry is identical, so the same fix must cover both viewports.
- The profile primary nav (activity/summary/notifications tabs) is intact and must not be disturbed by the fix.

### User Card, Post Actions, Panels, Tags, Login (audited, no substantive defect)

- User card popover fits at 390px; only defensive wrapping is planned for long display names, bio text, and badges.
- Post action row (reply/copy-link/bookmark, like, reaction) and the bookmark popup menu are visible and functional; the planned work is touch-target normalization to match MUI-05 (40–44px class targets) and icon fill consistency in light/dark mode.
- Notifications and account panels from the header user menu fit at 390px; the planned work is containment hardening (panel width, scroll area, close control) only.
- Tags page fits at 390px in both login states; defensive wrapping is planned for very long tag names.
- Login entry is unchanged from MUI-03A and passes checks. The decorative orb on the login page slightly exceeds the viewport during animation; this is animation noise with no user impact and is explicitly out of scope.

## 5. CSS Work Items

Primary file:

- `discourse-themes/ai-forum-premium-preview/mobile/mobile.scss`

Optional only if a shared token is needed:

- `discourse-themes/ai-forum-premium-preview/common/common.scss`

Planned CSS areas:

1. User profile activity navigation (the only substantive fix)
   - Make `.user-navigation .user-nav` (activity tabs) wrap onto a second row or shrink on mobile so all four activity tabs are fully visible at 390px and 375px.
   - Preserve the existing horizontal-scroll container only as a fallback; never hide a tab without a stable visible alternative.
   - Keep tab labels and icons visible in light/dark mode.

2. User profile shell
   - Keep primary nav (activity/summary/notifications) untouched unless the activity fix requires scoped alignment.
   - Defensive wrapping for long usernames and profile fields inside the profile header.

3. User card popover
   - Contain the card inside the mobile viewport.
   - Defensive wrapping for long display names, bios, and badge rows.
   - Keep the dismiss control visible at 375px.

4. Post action row and bookmark menu
   - Normalize post action buttons to clear 40–44px touch targets on mobile.
   - Keep icon color and SVG fill visible in light/dark states.
   - Keep the bookmark popup contained and its buttons reachable.

5. Notifications and account panels
   - Keep panel width bounded to the mobile viewport.
   - Ensure the scroll area and close control remain visible and tappable.
   - No behavior changes; containment and spacing only.

6. Tags page
   - Defensive wrapping for very long tag names inside tag tiles.
   - Keep tag tiles tappable at 40px+ where practical.

7. Login entry
   - No visual changes. Only a regression check against MUI-03A styling.

8. Read-only regression guard
   - `/latest` stays at the MUI-04/05 final state; any preview validation that shows a search/category/menu regression stops the work.

## 6. Validation Plan

Primary mobile viewport:

- `390 x 844`

Secondary mobile viewport:

- `375 x 667`

Desktop smoke:

- `1440 x 900`

Pages and states (logged-in unless noted):

- `/latest` (logged-in and logged-out, read-only regression)
- `/tags` (logged-in and logged-out)
- `/login` (logged-out)
- Topic detail, including the post action row and bookmark popup menu (logged-in and logged-out)
- User card popover triggered from a topic post author
- Header user menu → notifications panel and account panel
- User profile `/u/<current-user>` (logged-in and logged-out for a public profile)

Validation checks:

- No document or body horizontal overflow.
- All four user-profile activity tabs (获赞/书签/回应/已解决列表) fully visible at 390px and 375px.
- No login or permission page leakage.
- All controls keep visible text or icons in light and dark mode; no blank buttons.
- Post action and panel controls remain at least 40px-class touch targets where practical.
- MUI-03A login visuals and MUI-04 search/category/menu and MUI-05 topic detail/composer states show no regression.
- Desktop has no obvious regression from mobile-only CSS.

## 7. Stop Conditions

Stop and ask before continuing if:

- The activity-nav fix needs Discourse template changes, theme JavaScript, or plugin code.
- Validation would require reading private messages, cookies, localStorage, passwords, or tokens.
- An issue can only be solved by changing Discourse behavior.
- Preview validation shows regression in MUI-03A/04/05 surfaces.
- The next step would sync MUI-06 into the default live theme.

## 8. Delivery Flow

1. Patch CSS locally in the `// MUI-06:` block of `mobile.scss`.
2. Build/import as a preview theme, for example `ai-forum-mui06-preview` (expected theme ID 5).
3. Validate the preview theme with logged-in and logged-out mobile screenshots per section 6.
4. Fix any CSS-only issues revealed by preview validation.
5. Commit and push MUI-06 changes (`fix: polish Discourse mobile user surfaces (#MUI-06)`).
6. Ask for explicit approval before syncing MUI-06 into default theme ID `1`.