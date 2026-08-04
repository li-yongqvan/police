# Handoff - 2026-08-04 - MUI-03A Discourse 移动端详情页与编辑器优化

> 新会话建议先读本文件，然后说：继续 MUI-03A，先按计划做 CSS/SCSS 预览主题实现。

## 1. 现在做了什么

- 当前任务：继续优化 `http://122.51.233.225:8080/` Discourse 论坛的移动端功能页面。
- 当前分支：`codex/discourse-rebuild`。
- 当前 8080 正式主题状态：
  - Discourse 默认主题已切到主题 ID `1`。
  - 主题名：`ai-forum-premium-preview`。
  - 正式入口 `http://122.51.233.225:8080/` 已加载 MUI-02 改版主题。
- 最新对齐的下一批任务：`MUI-03A`。
- MUI-03A 范围：
  - 帖子详情页
  - 回复编辑器
  - 新建主题编辑器
- 用户已确认的关键决策：
  - 继续只改 CSS/SCSS。
  - 不改模板、不加 JS、不加插件。
  - 编辑器走低调的移动端沉浸式写作体验。
  - 工具栏保持单行横向滚动。
  - 帖子详情页不强卡片化，只做轻分隔和阅读层级优化。
  - 帖子操作以图标按钮为主，回复作为更明显主操作。
  - 如果 CSS 无法可靠修复空白按钮，停下给用户看截图再决定是否扩展到模板/JS。
  - MUI-03A 先上传到独立预览主题 `ai-forum-mui03a-preview`，用户确认后再同步到正式默认主题。

## 2. 已经完成了什么

- 已完成 MUI-02：
  - 首页/最新话题列表移动端优化。
  - 发帖入口使用代表性单字“发”。
  - 回复数保持右侧热度锚点，并扩大到 44px 级触控区。
  - 隐藏移动端欢迎 banner。
  - 修复管理员审核徽标横向溢出。
- 已上传并启用正式默认主题：
  - `ai-forum-premium-preview` 已成为 8080 默认主题。
  - 验证过正式入口无预览提示、无横向溢出。
- 已推送到 GitHub 的最近提交：
  - `4b655cf fix: improve Discourse mobile topic list UI (#MUI-02)`
  - `be69bf0 fix: constrain Discourse mobile preview overflow (#MUI-02)`
  - `2dd1c95 fix: contain Discourse mobile staff badge overflow (#MUI-02)`
- 已创建但尚未提交的 MUI-03A 文档：
  - `docs/adr/0002-discourse-mobile-topic-detail-composer-css-plan.md`
  - `docs/plans/mui03a-discourse-topic-detail-composer-mobile-css-plan.md`
- 已完成 MUI-03A 只读审计，截图和指标在：
  - `work/screenshots/discourse-8080-functional-audit-2026-08-04/`
- MUI-03A 审计结论：
  - 帖子详情页可读，但底部操作按钮有空白方块感，按钮高度约 39px。
  - 回复/发帖编辑器问题最重：工具栏图标弱、按钮跑出屏幕、正文框过高、底部操作区不够清晰。
  - 搜索、分类、用户菜单也有问题，但已明确排除在 MUI-03A 之外。

## 3. 卡在哪里

- 暂无明确阻塞。
- 尚未开始 MUI-03A CSS/SCSS 实现。
- 尚未创建或上传 `ai-forum-mui03a-preview` 预览主题。
- 尚未提交新增的 ADR 和计划书。
- 当前工作区不是干净状态：
  - `docs/dev-records/timeline.md` 有 hook 生成的修改。
  - `docs/dev-records/entries/...` 有多个未跟踪开发记录。
  - `docs/handoffs/handoff-2026-08-04-8080-frontend-ui-plan-prep.md` 未跟踪。
  - `docs/adr/0002-discourse-mobile-topic-detail-composer-css-plan.md` 未跟踪。
  - `docs/plans/` 未跟踪。
  - `work/` 下有截图、zip、审计工件。
- 不要误把 `work/` 截图/zip 工件或无关 dev-records 全部提交，除非用户明确要求。

## 4. 下一步做什么

1. 先让新会话读取：
   - `AGENTS.md`
   - `CONTEXT.md`
   - `docs/adr/0001-discourse-mobile-theme-css-only.md`
   - `docs/adr/0002-discourse-mobile-topic-detail-composer-css-plan.md`
   - `docs/plans/mui03a-discourse-topic-detail-composer-mobile-css-plan.md`

2. 检查当前工作区：
   ```powershell
   git status --short --branch
   git log --oneline -6
   ```

3. 如果用户同意继续实现，先提交文档类变更：
   - 建议单独提交 ADR + plan + 本 handoff。
   - 不要把 `work/` 截图/zip 混入。

4. 开始 MUI-03A 前，先备份当前正式主题 ID 1：
   - 目标是创建新的 MUI-03A 专用备份 zip。
   - 旧备份仍在：`/shared/tmp/ai-forum-premium-preview-before-mui02-20260804-152140.zip`

5. 只改本地主题 CSS/SCSS：
   - 主文件：`discourse-themes/ai-forum-premium-preview/mobile/mobile.scss`
   - 必要时才改：`discourse-themes/ai-forum-premium-preview/common/common.scss`

6. 按计划实现 MUI-03A：
   - 帖子详情页阅读排版和操作按钮。
   - 回复编辑器外壳、工具栏、输入区、发布区。
   - 新建主题编辑器标题/分类/标签/正文/发布区。

7. 打包并上传到独立预览主题：
   - 预览主题名：`ai-forum-mui03a-preview`
   - 不要先覆盖正式默认主题 ID 1。

8. 验收页面：
   - 帖子详情页：`http://122.51.233.225:8080/t/ai/5`
   - 回复编辑器：从上述详情页点击回复打开
   - 新建主题编辑器：从 `http://122.51.233.225:8080/latest` 点击“发”打开
   - 先测 `390x844`，通过后测 `375x667`。

9. 验证通过后向用户报告：
   - 预览 URL
   - 截图
   - 指标：横向溢出、关键按钮尺寸、工具栏是否溢出、发布按钮是否可见
   - CSS-only 仍无法修的残留问题

10. 用户确认后，再把同一批 CSS/SCSS 同步到正式默认主题 ID 1。

## 5. 哪些坑不要再踩

- 不要把 `http://122.51.233.225:8080/` 和 `preview_theme_id=1` 混淆：
  - 现在 ID 1 已经是正式默认主题。
  - MUI-03A 必须新建独立预览主题 `ai-forum-mui03a-preview`，避免直接影响正式站点。
- 不要使用 root 用户跑 Discourse Rails runner：
  - root 方式可能触发 PostgreSQL peer auth 问题。
  - 使用容器用户 `discourse`：
    ```bash
    docker exec -u discourse app bash -lc 'cd /var/www/discourse && RAILS_ENV=production bundle exec rails runner "..."'
    ```
- PowerShell 里不要用 `&&` 串联命令：
  - 当前环境曾报错：`The token '&&' is not a valid statement separator`。
  - 分开执行，或用 PowerShell 原生命令结构。
- PowerShell 本地解析会误处理 Rails runner 里的 `|`：
  - 复杂远端脚本建议用 here-string 通过 stdin 传给 `ssh ... 'bash -s'`。
- PowerShell 管道可能污染包含中文的 Node 脚本/JSON 正则：
  - 脚本里需要匹配中文时优先使用 Unicode escape，或避免中文正则。
- Playwright 包存在，但 bundled Chromium 可能不存在：
  - 本机可用 Chrome 路径曾验证为：`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`
  - Playwright 启动时可指定 `executablePath`。
- 不要对用户质疑直接反驳或直接改：
  - 先判断质疑是否合理。
  - 合理时先给计划，再执行。
- 严守用户确认过的 CSS-only 边界：
  - 如果空白按钮需要模板/JS 才能可靠解决，必须停下汇报。
- 不要批量删除工作区文件：
  - 之前只清理过明确指定的旧审计截图目录。
  - `work/` 里可能还有有用截图、zip 和对比工件。
- 不要把密码、token、cookie 写入 handoff 或提交。
