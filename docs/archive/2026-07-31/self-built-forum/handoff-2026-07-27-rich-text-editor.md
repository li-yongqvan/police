# Handoff · 富文本编辑器（Tiptap）集成

> 日期：2026-07-27 | 分支：`experiment/bg-animation` | 状态：✅ 已完成并部署

---

## 1. 功能概述

将论坛原有的纯文本 `<Textarea>` / `<Input>` 替换为基于 Tiptap 的富文本编辑器，支持加粗、斜体、颜色、高亮、链接、列表、代码块、表情等格式，并提供只读渲染器用于展示。

---

## 2. Ticket 完成清单

### Ticket 1 — Tiptap 组件脚手架 ✅
- **依赖**：`@tiptap/vue-3@3.28.0`、starter-kit、underline、color、highlight、link、placeholder、`emoji-picker-element`（`npm install --legacy-peer-deps`）
- **新建文件**：
  - `frontend/src/composables/useRichTextEditor.js` — `getEditorExtensions({ placeholder })`
  - `frontend/src/components/editor/EmojiPicker.vue` — web component 表情选择器
  - `frontend/src/components/editor/EditorToolbar.vue` — 工具栏（full=14按钮 / light=4按钮），移动端折叠
  - `frontend/src/components/editor/RichTextEditor.vue` — 完整编辑器（v-model JSON）
  - `frontend/src/components/editor/RichTextEditorLight.vue` — 精简版（评论用）
  - `frontend/src/components/editor/RichTextRenderer.vue` — 只读渲染器（自动识别 JSON/纯文本）

### Ticket 2 — NewPost.vue 集成 ✅
- `<Textarea>` → `<RichTextEditor>`
- 新增 `hasActualContent()` 辅助函数（检测 JSON 是否含真实文本）
- `canSubmit` 逻辑更新

### Ticket 3 — EditPost.vue + PostDetail.vue 集成 ✅
- `EditPost.vue`：`<Textarea>` → `<RichTextEditor>`，disabled 逻辑更新
- `PostDetail.vue`：移除 `parsePostContent()` 和 `postContentBlocks`，帖子正文改用 `<RichTextRenderer>`

### Ticket 4 — 评论区集成 ✅
- `PostDetail.vue`：两处 `<Input>` → `<RichTextEditorLight>`，新增 `hasActualText()`，disabled 绑定更新
- `GxCommentNode.vue`：`<p>{{ node.content }}</p>` → `<RichTextRenderer :content="node.content" />`

---

## 3. 关键代码模式

### hasActualContent / hasActualText
```js
function hasActualContent(jsonStr) {
  if (!jsonStr) return false
  try {
    const doc = JSON.parse(jsonStr)
    function hasText(node) {
      if (node.type === "text" && node.text?.trim()) return true
      if (node.content?.length) return node.content.some(hasText)
      return false
    }
    return hasText(doc)
  } catch { return false }
}
```
定义位置：`NewPost.vue`（hasActualContent）、`EditPost.vue`（hasActualContent）、`PostDetail.vue`（hasActualText）

### RichTextRenderer 降级逻辑
```
content 以 { 开头且可 JSON.parse 为 ProseMirror doc → Tiptap readonly 渲染
否则 → <p class="whitespace-pre-wrap">{{ content }}</p>
```

---

## 4. Git 提交

| Commit | 描述 |
|---|---|
| `6c02c81` | feat: rich-text editor (Tiptap) integration + post recommendations |
| `c8244b6` | fix: deployment fixes for domestic server |

---

## 5. 部署记录

**目标**：云服务器 `122.51.233.225:8888`（`experiment/bg-animation` 分支）

### 部署中遇到的问题及修复

| # | 问题 | 原因 | 修复 |
|---|---|---|---|
| 1 | `pipefail: invalid option name` | 服务器 shell 不支持 `pipefail` | `deploy-domestic.sh`：`set -euo pipefail` → `set -eu` |
| 2 | `#!/usr/bin/env: No such file` + `$'\r'` | 脚本含 BOM + CRLF 行尾 | UTF-8 无 BOM + LF 行尾 |
| 3 | `UnicodeEncodeError: 'gbk'` | Windows 控制台 GBK 无法编码输出中的 BOM 字符 | 设置 `PYTHONIOENCODING=utf-8` |
| 4 | Docker Hub `i/o timeout` | 国内服务器无法直连 Docker Hub | 为已有 `docker.m.daocloud.io` 镜像打 Docker Hub 标签 |
| 5 | `npm install` peer 依赖冲突 | `@tiptap/vue-3@3.28.0` vs `@tiptap/core@3.29.0` | `frontend/Dockerfile`：`npm install --legacy-peer-deps` |
| 6 | 端口 80 `address already in use` | `docker-compose.yml` + `docker-compose.server.yml` 端口合并而非覆盖 | 从 `docker-compose.yml` 移除 `80:80`，由 override 文件各自提供 |

### 服务器信息
- IP：`122.51.233.225`
- 用户：`liyongquan`
- 项目路径：`/home/liyongquan/projects/ai-forum`
- 部署脚本：`scripts/remote-deploy-domestic.py`（需 `DEPLOY_PASSWORD` 环境变量）
- 访问：http://122.51.233.225:8888/

---

## 6. 已知注意事项

- **Tiptap 版本**：`@tiptap/vue-3@3.28.0`，其余扩展 `@3.29.0`，npm install 必须加 `--legacy-peer-deps`
- **Tailwind v4**：scoped `<style>` 不能使用 `@apply` 配合自定义主题类，工具栏样式用原始 CSS + `var(--color-*)` 实现
- **移动端工具栏**：`max-sm:` 下默认只显示 4 个按钮 + 「更多」下拉菜单
- **RichTextRenderer**：自动兼容 JSON 和纯文本两种格式，无需单独处理遗留数据
- **docker-compose.server.yml**：端口为追加而非覆盖，base compose 不应包含会被冲突的端口映射

---

## 7. 未改动范围

- `services/` 下除 `forum-service` 的 post_handler 和 forum_service 外的所有 Go 代码
- `frontend/src/api.js`、`frontend/src/api/http.js`、`frontend/src/stores/`
- `frontend/src/components/ui/Textarea.vue`、`frontend/src/components/ui/Input.vue`
- `frontend/src/components/gx/GxCommentTree.vue`
- 数据库 schema、Nginx 配置
- 不涉及任何新的第三方运行时依赖（仅 Tiptap + emoji-picker-element）