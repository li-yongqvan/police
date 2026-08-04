# Keep Discourse Mobile Theme Changes CSS-Only

For MUI-02, the 8080 Discourse mobile UI work will be limited to the theme's CSS/SCSS files. We are choosing style-layer fixes over theme JavaScript, template overrides, or plugin logic because the observed issues are primarily touch target size, visual hierarchy, spacing, and mobile density; expanding into templates or behavior would raise maintenance risk during Discourse upgrades.

**Consequences**

If a requirement cannot be met with CSS/SCSS, implementation stops and the gap is reported before expanding scope.
