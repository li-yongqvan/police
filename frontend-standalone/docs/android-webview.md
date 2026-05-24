# Android WebView 集成说明

本前端为 **Hash 路由**（`#/community`）+ **相对资源路径**（`base: './'`），便于在 WebView 中从 `file://` 或内网 HTTPS 加载构建产物。

## 构建

```bash
cd frontend-standalone
npm install
npm run build
```

将 `dist/` 目录放入 Android 工程，例如：

- `app/src/main/assets/dist/index.html`

## WebView 最小配置

```kotlin
webView.settings.apply {
    javaScriptEnabled = true
    domStorageEnabled = true   // localStorage 演示态持久化
    allowFileAccess = true
}
webView.loadUrl("file:///android_asset/dist/index.html")
```

若通过 HTTPS 部署，直接 `loadUrl("https://your-host/dist/index.html")` 即可。

## 路由与资源

| 项 | 说明 |
|---|---|
| Hash 路由 | 避免 `file://` 下 History API 路径异常 |
| `base: './'` | `dist/index.html` 中脚本为 `./assets/...`，避免绝对路径 404 |
| 演示数据 | 键名 `ai-forum-standalone-state`（业务态）、`ai-forum-token` / `ai-forum-user`（登录态） |

## 返回键（建议）

抽屉打开时，优先在 WebView 内关闭抽屉，再退出 Activity：

```kotlin
override fun onBackPressed() {
    webView.evaluateJavascript(
        "(function(){var b=document.querySelector('.layout-backdrop.is-visible');if(b){b.click();return true;}return false;})()"
    ) { result ->
        if (result != "true") super.onBackPressed()
    }
}
```

或在页面内通过 `window.__closeDrawer?.()` 与原生桥接（可选扩展）。

## 验证清单

1. 竖屏 360dp：首屏为社区内容，非整块侧栏
2. 顶栏菜单打开/关闭抽屉，点击遮罩关闭
3. 登录 → 社区 → 板块 → 发帖 → 管理员进中台
4. 表单聚焦时不出现异常页面缩放（输入框 16px 字号）
5. 桌面宽度 ≥1024px：侧栏常驻，无顶栏菜单
