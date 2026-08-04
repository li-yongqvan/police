# Optimize Discourse Mobile Topic Detail and Composer with CSS-Only Preview Flow

For MUI-03A, the next 8080 Discourse mobile UI batch will focus on the core read-to-write loop: topic detail, reply composer, and new-topic composer.

The implementation will remain limited to the Discourse theme's CSS/SCSS files. We will not add theme JavaScript, template overrides, plugins, or behavior changes in this batch. The observed problems are primarily visual and ergonomic: low touch-target height, weak icon visibility, toolbar overflow, composer height, crowded controls, and unclear action hierarchy.

## Decisions

MUI-03A covers only:

- Topic detail page
- Reply composer opened from a topic detail page
- New-topic composer opened from the latest topic list

MUI-03A does not cover search, category pages, user menus, notification drawers, topic list redesign, or admin pages.

The mobile composer should behave like a quiet immersive writing surface. It may remain full-screen or near full-screen, but the editor body must not consume the whole viewport. The publish area and close/minimize controls should remain visible and tappable.

The composer toolbar will stay on one horizontal row and may scroll horizontally. Common actions should be visible first where Discourse's DOM order allows it. We will not shrink toolbar controls below mobile touch-target expectations.

Topic detail posts will not be strongly cardified. The page should keep Discourse's reading structure while improving title hierarchy, author/time metadata calmness, body rhythm, and post action controls.

Post action controls will stay primarily icon-based on mobile. Like, copy, more, and similar secondary actions should be clear 44px-class icon buttons. Reply should be visually stronger as the primary action. We will not add permanent text labels to every post action.

If a blank-looking button cannot be reliably fixed with CSS/SCSS, implementation stops at the CSS boundary. We will report the remaining gap with screenshots before deciding whether template or JavaScript changes are allowed.

MUI-03A changes should first be uploaded to a separate preview theme. After screenshot and metric validation, the accepted changes may be synchronized into the current default theme.

## Acceptance

Validate at 390x844 first, then at 375x667 after the first pass succeeds.

Topic detail acceptance:

- No horizontal scrolling
- Topic title and cooked post body do not overflow
- Post body line-height and paragraph spacing are readable on mobile
- Post action buttons are at least 44px-class tap targets where practical
- Reply is clearly identifiable as the primary action
- No visually empty action blocks remain unless documented as CSS-only gaps

Reply composer acceptance:

- No horizontal scrolling
- Toolbar controls do not run off-screen; horizontal toolbar scrolling is acceptable
- Body input is visible without hiding the publish area
- Reply button is visible and at least 44px high
- Close/minimize controls are visible and at least 40px-class tap targets

New-topic composer acceptance:

- Title input, category selector, optional tags, toolbar, body input, and publish button read as one coherent writing flow
- Toolbar controls do not run off-screen
- Publish button is visible and at least 44px high
- No confusing blank operation blocks remain unless documented as CSS-only gaps

Human acceptance:

The user can tell at first glance where they are, what they are writing, how to publish, and how to exit.

## Consequences

This keeps the mobile UI work upgrade-tolerant and reversible, but it limits how far we can change composer semantics and toolbar behavior. CSS can improve spacing, hierarchy, visibility, and touch ergonomics; it cannot safely redefine editor behavior or rebuild unstable Discourse controls.

The preview-theme-first flow adds one extra deployment step, but it prevents unfinished composer changes from immediately affecting the live default 8080 forum.
