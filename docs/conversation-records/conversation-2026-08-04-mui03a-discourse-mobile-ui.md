# Conversation Record - 2026-08-04 - MUI-03A Discourse Mobile UI

## What happened

This discussion continued the 8080 Discourse mobile UI work after MUI-02 was deployed.

The user first reported that `http://122.51.233.225:8080/` still looked like the default Discourse theme. Investigation showed that the updated theme existed as a preview theme but was not the default theme. After user confirmation, theme ID `1` (`ai-forum-premium-preview`) was set as the Discourse default theme. The formal root URL was then validated with mobile screenshots and metrics.

After the user expressed satisfaction with the homepage/list result, the conversation moved to optimizing other functional pages. A read-only audit was performed for topic detail, reply composer, new-topic composer, search, category, and user menu surfaces. The highest-impact issues were found in the topic detail and composer flow.

The user then requested a `grill-with-docs` alignment session. A sequence of decisions defined MUI-03A as a CSS/SCSS-only batch for topic detail, reply composer, and new-topic composer. An ADR, execution plan, and handoff were created.

## What was decided

- MUI-03A covers only:
  - Topic detail page
  - Reply composer
  - New-topic composer
- MUI-03A does not cover:
  - Search page or search panel
  - Category pages
  - User menu
  - Notification drawer
  - Admin pages
  - Desktop redesign
- Implementation remains CSS/SCSS-only.
- Do not add theme JavaScript, template overrides, plugins, or behavior changes.
- The composer should become a restrained mobile immersive writing surface.
- The composer toolbar should remain one horizontal row with horizontal scrolling.
- The editor body should not consume the whole viewport; publish and exit/minimize controls must remain visible.
- Topic detail should not be strongly cardified.
- Topic detail should receive lighter reading-layout refinements: title hierarchy, calmer metadata, readable body rhythm, and better post action controls.
- Post actions should remain primarily icon buttons.
- Reply should be the visually stronger primary action.
- If a blank-looking control cannot be reliably fixed with CSS, stop and report the CSS-only gap before expanding scope.
- MUI-03A should first be uploaded to a separate preview theme named `ai-forum-mui03a-preview`.
- Only after preview validation and user approval should accepted CSS/SCSS be synchronized into the live default theme.

## Facts recorded from the environment

- Current branch: `codex/discourse-rebuild`.
- Latest pushed MUI-02 commits:
  - `4b655cf fix: improve Discourse mobile topic list UI (#MUI-02)`
  - `be69bf0 fix: constrain Discourse mobile preview overflow (#MUI-02)`
  - `2dd1c95 fix: contain Discourse mobile staff badge overflow (#MUI-02)`
- Discourse default theme was changed to:
  - `default_theme_id = 1`
  - `theme_name = ai-forum-premium-preview`
- Formal 8080 root URL validation passed after the switch:
  - `previewNotice: false`
  - `cannotPreview: false`
  - `overflowX: false`
  - Mobile new-topic action rendered as the agreed one-character publish affordance.
  - Key topic-list action targets were 44px-class
- MUI-03A audit artifacts are under:
  - `work/screenshots/discourse-8080-functional-audit-2026-08-04/`
- Existing rollback reference from the earlier theme backup:
  - `/shared/tmp/ai-forum-premium-preview-before-mui02-20260804-152140.zip`

## What is still open

- MUI-03A implementation has not started.
- The `ai-forum-mui03a-preview` preview theme has not been created or uploaded.
- A new MUI-03A-specific backup of the current live default theme still needs to be created before implementation.
- Search, category, user menu, and notification drawer issues are known but intentionally deferred.
- Some blank-looking buttons may be CSS-fixable; if not, the user must decide whether to allow template or JavaScript work.

## What was written where

- ADR:
  - `docs/adr/0002-discourse-mobile-topic-detail-composer-css-plan.md`
- Execution plan:
  - `docs/plans/mui03a-discourse-topic-detail-composer-mobile-css-plan.md`
- Handoff:
  - `docs/handoffs/handoff-2026-08-04-mui03a-discourse-mobile-composer.md`
- This conversation record:
  - `docs/conversation-records/conversation-2026-08-04-mui03a-discourse-mobile-ui.md`

## Recommended next actions

1. In the next implementation session, read:
   - `AGENTS.md`
   - `CONTEXT.md`
   - `docs/adr/0001-discourse-mobile-theme-css-only.md`
   - `docs/adr/0002-discourse-mobile-topic-detail-composer-css-plan.md`
   - `docs/plans/mui03a-discourse-topic-detail-composer-mobile-css-plan.md`
   - `docs/handoffs/handoff-2026-08-04-mui03a-discourse-mobile-composer.md`

2. Commit the documentation artifacts separately if the user wants the docs persisted in git:
   - ADR
   - execution plan
   - handoff
   - conversation record

3. Before MUI-03A implementation, export a new backup of the current default theme ID `1`.

4. Implement local CSS/SCSS only, primarily in:
   - `discourse-themes/ai-forum-premium-preview/mobile/mobile.scss`
   - `discourse-themes/ai-forum-premium-preview/common/common.scss` only if needed

5. Package and upload to `ai-forum-mui03a-preview`, not directly to the live default theme.

6. Validate at:
   - `390x844`
   - `375x667`

7. Validate pages:
   - `http://122.51.233.225:8080/t/ai/5`
   - Reply composer opened from that topic
   - New-topic composer opened from `http://122.51.233.225:8080/latest`

8. Report screenshots, metrics, and remaining CSS-only gaps before syncing to the live default theme.

## Recovery notes and pitfalls

- Do not confuse theme ID `1` with a preview-only theme anymore; it is currently the live default theme.
- Do not directly modify the live default theme for MUI-03A first pass.
- Use a separate preview theme named `ai-forum-mui03a-preview`.
- Use the `discourse` user for Rails runner commands inside the Discourse container.
- Avoid `&&` in PowerShell commands in this environment.
- Avoid complex local PowerShell quoting for Rails runner strings containing `|`; send remote scripts via here-string to `ssh ... 'bash -s'`.
- Avoid Chinese regex or raw Chinese snippets in Node scripts piped through PowerShell when reliable parsing matters; use Unicode escapes where needed.
- Playwright package exists, but bundled Chromium may not. Use the local Chrome executable path if needed:
  - `C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`
- Do not commit `work/` screenshots, zip files, or unrelated generated dev-records unless the user explicitly asks.
- Do not write passwords, tokens, cookies, or private local file contents into durable project docs.
