# Frontend Standalone

This copy is a self-contained frontend demo. It does not require the Go services under `services/`.

## Run

```bash
npm install
npm run dev
```

The Vite dev server runs on `http://127.0.0.1:4174`. Routes use Hash mode (`#/community`).

## Build

```bash
npm run build
```

Output in `dist/` uses relative asset paths (`base: './'`) for Android WebView.

## Android APK

完整 Android Studio 工程在 [`android/`](android/) 目录。

```powershell
# 构建前端并打包 Debug APK
.\scripts\build-apk.ps1
```

APK 路径：`android/app/build/outputs/apk/debug/app-debug.apk`

详细步骤见 [android/README.md](android/README.md)。WebView 技术说明见 [docs/android-webview.md](docs/android-webview.md)。

### Quick verification

1. `npm run build` — confirm `dist/index.html` references `./assets/...`
2. Chrome DevTools → Pixel 5 / Galaxy S20 (360×800): main content visible first; menu opens drawer
3. Login → community → board → new post → admin role → admin console
4. Width ≥1024px: persistent sidebar, no top menu bar
5. Form focus: no unwanted page zoom (inputs use 16px font size)

## Notes

- Demo data is bundled under `src/mock-data/`.
- Runtime edits such as posting, commenting, moderation, and user status changes are persisted in browser `localStorage`.
- To reset the demo back to seed data, clear the `ai-forum-standalone-state` entry from browser storage.

## Mobile layout

- **&lt;1024px**: top bar + off-canvas navigation drawer (content first)
- **≥1024px**: classic sidebar + main area (same as desktop demo)
