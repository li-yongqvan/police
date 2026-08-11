# Handoff - 2026-08-10 - MUI-03A Discourse mobile blank affordance fix

> Continue from this note before changing code.

## 1. Current State

- Branch: `codex/discourse-rebuild`
- Default Discourse theme on 8080: `ai-forum-premium-preview`
- MUI-03A CSS is already synced into the live default theme.
- This issue is visual clarity, not backend data or routing.
- Current tracked-code worktree changes are documentation only.

## 2. Problem Statement

Some mobile controls still read as blank blocks or icon-only squares, especially in dark mode. The user cannot tell what they do.

Observed symptoms from the user screenshot:

- A blue square button with no readable label.
- Several square controls in the top area with no obvious name or icon.
- Topic list content is visible, but key controls do not communicate their function.

## 3. What Has Been Checked

- The live theme CSS is active on the 8080 site.
- Public topic-list and topic-detail pages load normally.
- The issue is most likely CSS presentation, not missing data.
- Anonymous `latest` and root pages did not exactly reproduce the screenshot state, so the final fix should be validated with a logged-in mobile session.

## 4. Likely Root Cause

The strongest suspect is the mobile theme's `no-text` and icon-only treatment.

Important files:

- `discourse-themes/ai-forum-premium-preview/mobile/mobile.scss`
- `discourse-themes/ai-forum-premium-preview/common/common.scss`

Important pattern:

- `#create-topic` is forced to `font-size: 0`.
- `#create-topic.no-text .d-icon` is hidden.
- The visible label depends on `#create-topic::before` for readable text.
- Other controls rely on SVG icon color inheritance and `no-text` classes.

Why this is fragile:

- If pseudo-element text has weak contrast in dark mode, the button looks blank.
- If an SVG icon inherits a weak color or is hidden, the control becomes an empty square.
- Controls with only accessibility labels remain screen-reader accessible, but are not self-explanatory visually.

## 5. Recommended Fix Direction

Stay CSS/SCSS-only first.

1. Audit every mobile `no-text` control on the list, home, and topic surfaces.
2. For `#create-topic`, replace the brittle pseudo-label behavior with a stable visible affordance.
3. Ensure button SVGs inherit a readable color in dark mode.
4. Do not hide both text and icon on the same control unless a visible replacement is guaranteed.
5. Keep mobile-first behavior and avoid desktop regressions.

## 6. Files To Inspect First

- `discourse-themes/ai-forum-premium-preview/mobile/mobile.scss`
- `discourse-themes/ai-forum-premium-preview/common/common.scss`

## 7. Validation Plan

Check at minimum:

- `390 x 844`
- `375 x 667`

Pages and states:

- `http://122.51.233.225:8080/`
- `http://122.51.233.225:8080/latest`
- Topic detail page after login
- Logged-in topic-list state where the create-topic action appears
- Dark mode and normal light mode if possible

Pass criteria:

- No blank blue square buttons.
- No icon-only controls without a readable visual clue.
- Search, navigation, list filters, and create-topic controls are understandable.
- No page-level horizontal overflow.
- Desktop remains unchanged.

## 8. Stop Conditions

Stop and ask the user if:

- Fixing the control requires template changes.
- Fixing the control requires JavaScript.
- The screenshot route still cannot be reproduced reliably after logged-in testing.
- A CSS-only change would break the mobile layout.

## 9. Working Tree Notes

- Do not stage or commit `work/` artifacts.
- There is an untracked MUI-03A test manual document from the prior step:
  - `docs/test-manual-mui03a-default-theme.md`
- The new handoff document is:
  - `docs/handoffs/handoff-2026-08-10-mui03a-discourse-mobile-blank-affordance-fix.md`

## 10. Suggested Next Step

Start with a narrow CSS patch that restores visible labels or icons for mobile action controls, then validate with screenshots before syncing the fix into the live default theme.
