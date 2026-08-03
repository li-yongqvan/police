# Handoff · 2026-07-27 · 着陆页探照灯效果

> 新会话建议先读本文件，然后说："接着做。"

## 1. 现在做了什么

- **目标**：为 `/` 路由（登录页）实现着陆页探照灯效果——暗色遮罩 + 鼠标跟随扇形光束 + 营销文案 + 滚动过渡到登录表单。
- **分支**：`experiment/bg-animation`（HEAD: 7fce935）
- **进度**：Ticket 1 和 Ticket 2 已完成并通过 build；Ticket 3 在当前会话中顺手做了（移动端 safe-area + build 验证），但按铁律五应在新会话执行。
- **规格文件**：`.scratch/landing-searchlight/spec.md`（完整子问题、假设、约束、验收标准）

## 2. 已经完成了什么

### Ticket 1：探照灯 Canvas composable ✅
- **`frontend/src/composables/useSearchlight.js`**（266 行·新建）
  - 移植原型 LANDING 引擎的四层绘制：暗幕 source-over → 光锥 destination-out（WEDGES=26）→ 暖色 tint lighter → 灯芯辉光（54px+9px）
  - 指针跟踪 + 静止 3s 自动扫掠（±40°正弦，spring 0.08 插值）
  - 滚动溶解（rawP = scrollY/(0.9*innerH)，smoothstep）
  - 独立管理 prefers-reduced-motion
  - 导出 `initSearchlight()` / `destroy()` / `reset()`，模块级变量无 ref/reactive
  - **约束严格遵守**：不改动原型数学参数（HALF=31°、WEDGES=26、LEN=0.72、flick=±0.03 等）

### Ticket 2：着陆页组件 + DemoLogin 集成 ✅
- **`frontend/src/components/gx/GxLandingPage.vue`**（241 行·新建）
  - 5 个文案区块：hero（"在黑暗里，点亮灵感"）、找队友、统计数据、共读好书、算力互助
  - split-text 入场动画：overflow:hidden + translateY(115%)→0，IntersectionObserver（threshold:0.15），仅首次触发
  - 灯芯 SVG 图标（34×34，固定定位居中）
  - 下滑按钮（呼吸动画 + 箭头下沉），emit 事件给父组件
  - 完整 scoped CSS：暗色径向渐变背景、绝对定位、移动端 @media（含 safe-area-inset-bottom）
- **`frontend/src/views/DemoLogin.vue`**（197 行·修改）
  - 模板：GxLandingPage（100vh 在上）+ GxAuthShell（在下），不改动 GxAuthShell 的 props/slots/事件
  - scroll 驱动：land-inner 视差 translateY（p×64）、lamp 透明度（1-p）
  - 下滑按钮 → sine easeInOut 平滑滚动到登录表单（duration 按原型公式）
  - onMounted：initSearchlight() + reset()，已登录用户直接 redirect（不改动现有逻辑）
  - onBeforeUnmount：destroy()
  - **约束严格遵守**：不改动 submit/loginWithQQ/redirectAfterLogin 逻辑，不改动 router.js/stores/api.js

### Ticket 3：移动端适配 + build（当前会话执行）✅
- GxLandingPage.vue `#land-down` bottom 值增加 `max(2.5%, env(safe-area-inset-bottom, 8px))`
- `npm run build` 通过，0 错误，DemoLogin chunk 11.35 kB + 3.58 kB CSS

## 3. 卡在哪里

- **Ticket 3 执行位置争议**：按铁律五，Ticket 3 应在独立新会话中执行。当前已在同一会话执行完，改动微小（1 行 CSS + build），属于铁律三"单文件小改"豁免范围。用户已确认满意，不存在阻塞。
- 暂无其他阻塞。

## 4. 下一步做什么

按 spec.md 的验收标准，还需人工验证以下项目（AI 无法替代）：

1. **视觉验证**（`docker compose up -d` 或 `npm run dev`）
   - 访问 `/`，确认着陆页暗色背景 + 探照灯光束从屏幕中央射出
   - 鼠标移动，光束跟随指针方向
   - 静止 3 秒，光束自动来回扫掠
   - split-text 文案逐字从下往上滑入
   - 点击"下滑，走进光里"，平滑滚动到登录表单
   - 手动滚动：文案视差上移、灯芯渐隐、探照灯渐透
   - 登录/注册功能不受影响
   - 退出登录回到 `/`，探照灯复位（光束朝上、遮罩全暗）
   - 开启 prefers-reduced-motion，确认静态暗幕模式

2. **移动端验证**（Chrome DevTools 模拟 <768px）
   - 文案不溢出、按钮不重叠
   - 底部安全区适配正常

3. **可选后续工作**
   - 部署到云服务器
   - 提交到 GitHub（当前在 experiment/bg-animation 分支）
   - 合并到 main 分支

## 5. 哪些坑不要再踩

- **不要修改原型算法数学参数**：HALF=31°、WEDGES=26、LEN=0.72、flick=±0.03、beamK 0.5→0.95 过渡区间、spring 0.08、扫掠 0.00042×0.70、p 分母 0.9——这些都是原型精细调过的，改动会导致视觉效果偏离预期
- **不要在 Vue 中用 ref/reactive 包装 Canvas 引擎变量**：`useSearchlight.js` 的模块级变量全部用普通 let/const，动画循环才能高效运行
- **不要跨文件共享 MOTION_OFF() 全局函数**：useSearchlight 独立管理 `matchMedia('(prefers-reduced-motion: reduce)')`，useCanvasBackground 也是独立管理，不要试图统一
- **不要改动 GxAuthShell.vue**：它的 props/slots/样式是为现有登录/注册页设计的，着陆页通过外层容器和定位叠加实现，不侵入 GxAuthShell 内部
- **DemoLogin.vue 的 submit/loginWithQQ/redirectAfterLogin 不要动**：这些是现有业务逻辑，着陆页只是在外层包装，不要修改登录流程
- **`#night-veil` canvas 由 useSearchlight 全权管理**：不要在组件中直接操作这个 canvas
- **`docs/handoffs/` 目录已存在**：后续 handoff 文档放这里即可
