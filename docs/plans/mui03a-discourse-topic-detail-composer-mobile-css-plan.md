# MUI-03A Discourse Topic Detail and Composer Mobile CSS Plan

Date: 2026-08-04

## 1. Background

The 8080 Discourse forum now uses `ai-forum-premium-preview` as the default theme. MUI-02 improved the mobile topic list, but the functional pages still show default Discourse ergonomics.

The MUI-03A audit found the highest-impact issues in the read-to-write loop:

- Topic detail page
- Reply composer
- New-topic composer

The main problems are small or blank-looking action buttons, composer toolbar overflow, weak icon visibility, oversized editor body, and publish/exit controls that are not clear enough on mobile.

## 2. Goal

Make the mobile topic detail and composer experience feel like a coherent, restrained mobile writing and reading flow while preserving Discourse behavior.

The user should be able to tell at first glance:

- Where they are
- What they are reading or writing
- How to publish
- How to exit or minimize the composer

## 3. Confirmed Decisions

- Scope is MUI-03A only.
- Modify CSS/SCSS only.
- Do not add theme JavaScript.
- Do not override Discourse templates.
- Do not add plugins.
- Do not change Discourse data, permissions, or behavior.
- Use a separate preview theme named `ai-forum-mui03a-preview` first.
- After validation and user approval, synchronize the accepted CSS/SCSS changes into the current default theme.

## 4. In Scope

### Topic Detail Page

- Topic title hierarchy
- Author, time, and floor metadata spacing
- Cooked post body line-height and paragraph rhythm
- Quote, code, and callout spacing
- Post action row
- Reply primary action styling
- 44px-class mobile touch targets where practical

### Reply Composer

- Full-screen or near-full-screen mobile composer shell
- Visible close and minimize controls
- Single-row horizontal scrolling toolbar
- Body input height control
- Visible reply button
- Bottom safe-area spacing
- Icon visibility and button affordance

### New-Topic Composer

- Title input
- Category selector
- Optional tag selector
- Toolbar
- Body input
- Publish button
- Close and minimize controls

## 5. Out of Scope

MUI-03A will not optimize:

- Search page or search panel
- Category pages
- User menu
- Notification drawer
- Admin pages
- Desktop redesign
- Topic list redesign beyond avoiding regressions
- Server-side Discourse settings except creating/updating the preview theme

## 6. CSS/SCSS Work Items

All implementation should stay inside the Discourse theme files, primarily:

- `discourse-themes/ai-forum-premium-preview/mobile/mobile.scss`
- `discourse-themes/ai-forum-premium-preview/common/common.scss` only when a shared token or safe base style is needed

Expected CSS areas:

1. Topic detail reading layout
   - Adjust topic title spacing and font size for mobile.
   - Make topic metadata quieter.
   - Normalize cooked body line-height, paragraph margins, image width, quote spacing, and code block overflow.

2. Post action controls
   - Normalize `.post-controls` and post action buttons to 44px-class tap targets.
   - Ensure SVG icons inherit `currentColor` and remain visible.
   - Make reply action visually stronger than secondary actions.

3. Composer shell
   - Keep `#reply-control` full-height or near-full-height on mobile.
   - Ensure the composer content uses viewport-aware layout.
   - Keep publish and exit/minimize controls visible.

4. Composer toolbar
   - Use one horizontal row.
   - Enable horizontal scrolling without page-level overflow.
   - Keep toolbar buttons 40px to 44px class targets.
   - Avoid wrapping into a visually noisy two-row toolbar.

5. Composer inputs
   - Give title, category, tags, and body input consistent spacing.
   - Cap the editor body height so it does not consume the full viewport.
   - Use internal scrolling for long body content.

6. Blank-looking buttons
   - First fix icon color, SVG fill, size, and alignment with CSS.
   - Use `::before` text only for stable, high-value actions where the selector is reliable.
   - If a button cannot be reliably fixed with CSS, stop and report it as a CSS-only boundary issue.

## 7. Preview Theme Deployment Flow

Do not directly update the live default theme during first implementation.

1. Export a backup of the current default theme.
2. Create or update preview theme:
   - Name: `ai-forum-mui03a-preview`
3. Upload the packaged theme zip to the server.
4. Import or update the preview theme in the Discourse container.
5. Validate using `preview_theme_id` for the preview theme.
6. Show screenshots and metrics to the user.
7. After user approval, synchronize the same accepted CSS/SCSS changes into the default theme.
8. Clear Discourse theme cache and precompile theme CSS.
9. Validate the formal `http://122.51.233.225:8080/` entry after synchronization.

## 8. Validation Pages

Primary viewport:

- `390x844`

Secondary viewport after primary passes:

- `375x667`

Pages:

- Topic detail: `http://122.51.233.225:8080/t/ai/5`
- Reply composer: open from the topic detail page
- New-topic composer: open from `http://122.51.233.225:8080/latest`

Use an existing staff/admin session for validation because composer controls require an authenticated session. Do not record account passwords in this plan.

## 9. Acceptance Criteria

### Topic Detail

- No horizontal scrolling.
- Topic title and cooked body do not overflow.
- Body line-height and paragraph spacing are readable on mobile.
- Quotes, code blocks, and images do not break the viewport.
- Post action buttons are 44px-class tap targets where practical.
- Reply is clearly identifiable as the primary action.
- No confusing blank action blocks remain unless explicitly documented as CSS-only gaps.

### Reply Composer

- No horizontal scrolling.
- Toolbar controls do not run off-screen.
- Toolbar horizontal scrolling is available when needed.
- Body input is visible without hiding the reply button.
- Reply button is visible and at least 44px high.
- Close and minimize controls are visible and at least 40px-class tap targets.
- Bottom safe-area spacing prevents controls from feeling clipped.

### New-Topic Composer

- Title input, category selector, optional tags, toolbar, body input, and publish button read as one coherent flow.
- Toolbar controls do not run off-screen.
- Publish button is visible and at least 44px high.
- Close and minimize controls are visible and usable.
- No confusing blank operation blocks remain unless explicitly documented as CSS-only gaps.

## 10. Regression Checks

- `http://122.51.233.225:8080/`
- `http://122.51.233.225:8080/latest`
- Topic list reply count remains a right-side heat anchor.
- The mobile "new topic" action still renders as the representative single character entry.
- No page-level horizontal scrolling at 390px and 375px widths.
- Desktop must not receive obvious layout breakage from mobile-only rules.

## 11. Rollback

Rollback options:

1. If still in preview:
   - Do not synchronize to the default theme.
   - Keep or delete `ai-forum-mui03a-preview` later.

2. If synchronized to default theme:
   - Re-upload the pre-change backup zip to the default theme.
   - Or restore the previous Git commit and re-upload the theme.
   - Clear Discourse theme cache and precompile CSS again.

The previous MUI-02 backup path remains relevant for older rollback:

- `/shared/tmp/ai-forum-premium-preview-before-mui02-20260804-152140.zip`

Before MUI-03A implementation, create a new MUI-03A-specific backup.

## 12. Stop Conditions

Stop and ask the user before continuing if:

- A required fix cannot be done reliably with CSS/SCSS.
- Fixing a blank-looking control requires template or JavaScript changes.
- The preview theme cannot be created or updated safely.
- Validation shows a desktop regression caused by shared styles.
- A server command would delete broad directories, upload unrelated local files, or expose private data.
- The user challenges the visual direction or asks to change the implementation scope.

## 13. Implementation Sequence

1. Create a new MUI-03A-specific backup of the current default theme.
2. Inspect current Discourse selectors for topic detail and composer controls.
3. Implement local CSS/SCSS changes in small groups:
   - Topic detail reading and post actions
   - Reply composer shell and toolbar
   - New-topic composer inputs and publish area
4. Package the theme zip.
5. Upload to `ai-forum-mui03a-preview`.
6. Validate preview at 390x844.
7. Fix preview issues within CSS-only boundaries.
8. Validate preview at 375x667.
9. Report screenshots, metrics, and remaining CSS-only gaps.
10. Wait for user approval before synchronizing to the live default theme.
