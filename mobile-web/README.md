# mobile-web — 手机浏览器适配专用目录

学院 AI 智联论坛的**移动端 CSS / 小 composable** 集中放在此目录，与 `frontend/` 主应用解耦，便于单独评审、测试和迭代。

## 目录结构

```
mobile-web/
  composables/     # 抽屉滚动锁等（由 frontend 引用）
  styles/          # 按 Phase 拆分的样式（仅窄屏生效）
  testing/         # 真机测试清单与设备矩阵
```

## 接入方式

`frontend/src/main.js` 在 `style.css` 之后引入：

```js
import '../../mobile-web/styles/index.css'
```

**Docker 构建**：`frontend/Dockerfile` 的 context 为**仓库根目录**（见 `docker-compose.yml`），会同时 `COPY mobile-web/`。

`frontend/src/composables/useDrawerNav.js` 使用 `useBodyScrollLock` 锁定背景滚动。

## 断点策略

- 默认样式写在 `max-width: 767.98px` 内，**不覆盖**桌面 `@media (min-width: 768px)` 布局。
- 与主样式 `1024px` 常驻侧栏策略兼容。

## 相关文档

- [docs/mobile-web-adaptation-plan-chris-coyier.md](../docs/mobile-web-adaptation-plan-chris-coyier.md)
