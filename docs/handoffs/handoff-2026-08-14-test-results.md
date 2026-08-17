# Handoff · 2026-08-14 全页面功能测试结果

> 接续 `docs/handoffs/handoff-2026-08-14-full-pages-testing-guide.md` 与 `docs/plans/full-pages-testing-plan.md`。
> 本轮在生产环境（122.51.233.225）执行：Vue 轻前端 :8888 + Discourse 论坛 :8080 + SSO。
> 测试工具：Browser‑Use（in-app browser，落盘通信）+ 项目内 Playwright 独立脚本（Browser‑Use 会话后期 CDP 不稳定，后续步骤改用 Playwright）。
> 测试时间：2026-08-14（下午）。

## 1. 执行总结（TL;DR）

- 通过：游客路径（26 项）、学生闭环（发帖/回复/书签/编辑/通知/个人资料四页签）、Vue 管理端登录与用户管理、邀请码生成/查询/作废、SSO 全链路（含 admin=true 映射）、演示数据清理。
- **P0（1 个）**：管理/版主账号登录 Discourse 后全站 500；根因指向服务器内存耗尽（详见 §3）。
- **P1（2 个）**：admin-service 硬编码 127.0.0.1:8002 调用已下线的 forum-service；内容审核闭环未实现（前端无路由 + 后端 TODO）。
- P2/P3 若干（菜单项指向不存在页面、邀请码统计不刷新、图标按钮无 aria 等）。

## 2. 已通过项

### 2.1 Phase 0-1 游客路径（全部通过）
- :8080 的 /latest、话题详情、/categories、/c/ai/5、/top、/search、/badges、/about、/faq、/tos、/signup、/login、/g、/u/... 等 26 项：`docOverflow=0 / bodyOverflow=0`，无空白按钮。
- 游客 /unread、/new 404 属预期（login_required 未强制全站登录）。

### 2.2 Phase 2 学生闭环（demo_student，全部通过）
- :8888 登录 → SSO 跳转 :8080 论坛成功。
- 发话题 #17（移动端测试-2026-08-14-闭环）、#19（并行互动测试-2026-08-14）、#21（插件功能抽查-2026-08-14）成功；回复、书签、编辑成功。
- 通知面板（2 条徽章 + 1 未读）、用户菜单、汉堡菜单、/u/demo_student/messages 可达。
- 个人资料四页签（likes-given/bookmarks/reactions/solved/topics/replies）390 + 375 全绿。
- 内建语法 details/spoiler/checklist 渲染正常。
- 并行回复审核拦截生效：demo02 对 #19 的回复进入待处理审核（未通过，无通知）。
- demo01 无发帖按钮为 Discourse 信任等级限制（预期），非 bug。

### 2.3 Phase 3 demo_admin（通过，含问题）
- Vue 登录 → /admin 成功；用户管理 76 用户 4 页正常。
- demo06 封禁→解封→改等级 5→4→5 恢复闭环成功；plat06 误改已还原 Lv.2；demo06 已恢复 active + Lv.5。
- 管理菜单中 stats/config 重定向回概览为设计行为（handoff 已注明）。

### 2.4 Phase 4 demo_platform_admin（通过，含问题）
- Vue 登录 → /admin 成功（管理员不走 SSO 到论坛，设计如此）。
- 邀请码：生成 `13f536a2-...` → 查询 `unused` → 作废成功（测试数据已清理）。
- 用户管理（76 用户）正常。
- SSO 链路验证：Discourse /session/sso → user-service 校验 cookie → 回传 payload 含 `admin=true` → 建立 Discourse staff 会话（payload 已解码确认）。
- 概览页三个统计接口报错（见 P1-1）。

## 3. 缺陷清单

### P0-1 · 管理/版主登录 Discourse 后全站 500（论坛 HTML 页面不可用）
- 现象：demo_platform_admin（管理员）、demo_admin（版主）经 SSO 登录后，访问 /、/latest、/categories、/about 全部返回 500「Oops」；/admin/dashboard 30s 超时。同一时间学生（demo_student）与游客访问相同页面 200 正常。JSON API（/latest.json、/session/csrf.json、/search.json）正常。
- 证据：
  - 浏览器端：`p4-sso-pages2.json`（管理员 500 列表）、`p4-moderator-now.json`（版主 500）、`p4-student-now.json`（学生 200 对照）。
  - 服务器端（只读排查，未做任何修改）：`docker stats` 显示 app 容器 2.094GiB / 3.636GiB（57.6%）、主机内存 2847/3723MB、Swap 已用 1426MB；`docker logs app` 反复出现 `worker=0 timed out, exiting`；`production.log` 连续 `Completed 500 in 30000ms (ActiveRecord: 0.0ms, 0 queries)`——请求在进入业务逻辑前排队 30s 被掐断。
- 判断：非代码逻辑错误，是**生产服务器内存耗尽导致 unicorn worker 饥饿**；staff 页面不做匿名缓存、渲染更重，故只有 staff 会话可见 500。试运行阶段最需要管理员在线治理，此问题阻断全部论坛侧治理入口，定级 P0。
- 建议（需用户确认后再执行，本轮未动服务器）：
  1. 短期：`cd /var/discourse && ./launcher restart app` 释放泄漏内存；或临时调低 `UNICORN_WORKERS`。
  2. 中期：给服务器加内存/加大 Swap（当前 2G Swap 已用 1.4G）；评估关闭不必要插件（chat/reactions/solved/templates/presence/narrative-bot 等）与 Prometheus 导出。
  3. 排查内存占用随时间的增长曲线（Discourse 2.1GB 偏高，疑有 worker 泄漏或插件放大）。
- 截图：`work/testing-2026-08-14/p4-sso-trace2-final.png`（管理员落地 500）、`p4-admin-*.png`。


### P0-1 修复记录（2026-08-17 · 已修复）

- 处置动作（全部在生产服务器 122.51.233.225 执行）：
  1. **Phase A 重启**：./launcher restart app 因服务器到 Docker Hub 出网超时（launcher 硬编码基础镜像 discourse/base:2.0.20260803-0122，本地只有 2.0.20260726-0220）在镜像预检阶段失败且未停容器；改用等价的 docker restart app 完成重启。
  2. **Phase B**：SiteSetting.version_checks = false（经 docker exec app rails runner 写入）。
  3. **追加修复 A（启动阻塞）**：重启后 mold 启动崩溃循环（Pitchfork::BootFailure）。定位为测试期间站点 port 设置被改为 8080，与 DISCOURSE_HOSTNAME='122.51.233.225:8080'（hostname 已含端口）叠加导致 Discourse.base_url = http://122.51.233.225:8080:8080，URI.parse 抛 InvalidURIError，client_settings_json 静默返回空串，PrettyText.cook 在 __optInput.siteSettings = ; 处 MiniRacer::ParseError。修复：SiteSetting.port = ""，base_url 恢复 http://122.51.233.225:8080。
  4. **追加修复 B（staff 500 直接根因）**：worker 超时堆栈显示 staff 页面渲染时 CurrentUserSerializer#has_unseen_features → DiscourseUpdates.new_features 对 Redis 中缓存的更新条目逐条执行 git merge-base --is-ancestor；该仓库为 partial clone（--filter=blob:none），本地缺失对象时 git 会向 origin 懒拉取，而服务器到 GitHub 出网挂起 → 每个 staff 页面请求同步卡死 → 30s worker timed out → 500，且挂起 git 子进程随 staff 访问不断堆积（与 recon 中观察到的 git fetch/git merge-base 进程吻合，学生/游客因不序列化该 staff 属性幸免）。ersion_checks=false 只阻止更新任务写入新数据，不清理存量缓存，故单独关闭无效。修复：Discourse.redis.del("new_features") + del("latest_new_feature_created_at")，并杀掉容器内遗留挂起 git 进程。
- 复测（Playwright，work/testing-2026-08-14/run-p4-sso-retest-2026-08-17.cjs → p4-sso-retest-2026-08-17.json）：demo_platform_admin 经 :8888 登录 → SSO 落地 :8080 正常；/latest、/admin/users/list/active、/about、/categories 全部 200，管理 UI 正常渲染。（/admin/dashboard 返回 200 但内容为 Discourse 对非规范路由的 404 文案，非回归。）
- 修复后服务器状态：app 容器 1.62GiB / 3.636GiB；主机 Swap 已用 344MB（修复前 1619MB）；容器内 git 挂起进程 0；ersion_checks=false、Redis 
ew_features=nil；production.log 无新增 worker 超时。
- 遗留观察项：内存 1.6GiB 仍偏高（3 worker + sidekiq），计划书 Phase C（UNICORN_WORKERS 3→2）与 Phase D（Swap 扩容）可选，建议观察 24h 内存曲线后再决定；port 设置为何被改成 8080 建议回溯测试操作日志。
### P1-1 · admin-service 硬编码 127.0.0.1:8002 调用已下线的 forum-service
- 现象：Vue 管理端「数据概览」显示 `failed to call forum-service ... dial tcp 127.0.0.1:8002 connection refused`；`/admin-api/api/v1/admin/stats/overview`、`/admin-api/api/v1/admin/boards`、`/admin-api/api/v1/admin/audit/pending` 均 500（浏览器网络面板确认）。
- 原因：当前架构已切换为「Vue 轻前端 + Discourse」，forum-service 容器已移除，但 admin-service 内部 client 仍指向 127.0.0.1:8002。
- 影响：概览统计、板块列表、审核队列数全部不可用。
- 建议：admin-service 改为只读 Discourse 数据（或移除这三个调用），并同步前端概览页降级展示。

### P1-2 · 内容审核闭环未实现
- 现象：侧边栏「内容审核/举报处理/帖子管理/板块管理/敏感词」点击后落入兜底路由 → 跳转 Discourse 首页；前端无对应视图（views 目录仅有 AdminConfig/AdminInvites/AdminLayout/AdminOverview/AdminRoles/AdminUsers）；后端 `ListPendingAudit` 为 TODO，无审核/审批接口。
- 影响：demo02 的待处理回复无人能审批（该记录成为孤儿数据）；敏感词增删无 UI。
- 建议：排期实现审核列表 + 通过/驳回；或明确这些治理能力全部下沉 Discourse（/review、/admin/flags），删除 Vue 侧边栏对应菜单项。

### P2-1 · 菜单项指向不存在页面（用户体验误导）
- 现象：Vue 管理端侧边栏 11 项中，实际可用仅 数据概览/用户管理/邀请码；趋势统计、系统配置、角色权限按设计重定向回概览；其余 5 项（内容审核/举报处理/帖子管理/板块管理/敏感词）直接跳论坛首页。
- 建议：隐藏未实现项或标注「敬请期待」，与 P1-2 一并处理。

### P2-2 · 邀请码作废后统计数字不刷新
- 现象：作废后提示「邀请码已作废」，列表行变为 voided，但「当前可用」仍为 2、「已使用/作废」仍为 0/0（需刷新页面才正确）。另状态文案混排「voided · 未使用」。
- 复现：:8888 管理员登录 → /admin/invites → 生成 1 个 → 作废。

### P2-3 · 论坛字体 CORS 错误
- 现象：页面引用 `http://122.51.233.225/fonts/InterVariable.woff2`（80 端口），从 :8080 页面加载被 CORS 拦截，字体回退。自定义主题中的绝对地址写死 80 端口。
- 建议：主题字体地址改相对路径或统一域名端口。

### P2-4 · 中文搜索无结果 / 报 invalid param
- 现象：搜索「测试」返回 invalid param 或无结果；英文词（Theme preview）有结果。疑似 Discourse zh 分词未配置（search_tokenize_chinese_japanese 等设置）。
- 建议：Discourse 后台开启中文分词设置并重建索引。

### P2-5 · 点赞/reactions/voting/solved 控件在话题内数量为 0 且交互缺失
- 现象：话题内未见可用的 reactions/like/voting/solved 控件（count=0，按钮缺失）。插件已安装但可能未启用或主题未渲染。
- 建议：Discourse 后台确认插件启用状态与分类设置。


### P2-7 · Vue 端（登录页/管理后台）不支持暗色模式
- 现象：浏览器 prefers-color-scheme: dark 下 Vue 页面与亮色完全一致（p5-vue-login-390-dark.png 与 light 截图逐字节相同），仅 Discourse 响应暗色。
- 建议：如产品要求亮/暗双主题，需为 Vue 端补充暗色变量；否则在测试门禁中将 Vue 端标注为「仅亮色」。### P2-6 · Vue 用户管理页无搜索
- 现象：76 用户 4 页，仅有分页按钮，无搜索/筛选；找 demo06 需逐页翻。管理员治理效率低。
- 建议：加用户名/学号搜索框。

### P3-1 · `#create-topic` 为纯图标按钮且无 aria-label
- 现象：新建话题按钮无文字、无 aria-label（`no-text`），依赖图标识别，无障碍/可读性不佳。

### P3-2 · 邀请码行状态文案不一致
- 现象：作废后显示 `voided · 未使用`（作废与未使用同时展示）。见 P2-2 截图。

## 4. 测试数据清理（已完成）
- 话题 #17/#19/#21（demo_student 创建）已用管理员 SSO 会话删除（`/t/{id}.json` DELETE，:8080 确认列表无残留）。
- 邀请码 `13f536a2-d2a0-cb18-0000-000000000000` 已作废；预置 `DEMO2026` 未动。
- 账号状态：demo06 已恢复 active + Lv.5，plat06 已恢复 Lv.2（Phase 3）。
- 遗留 1 项无法自助清理：demo02 对 #19 的回复产生的待处理审核记录（审核闭环未实现，见 P1-2）——话题删除后该记录已失效，但 DB 中可能残留，需后续审核功能上线后或 DBA 手工清理。

## 5. 账号密码校正
- 实测全部演示账号密码为 `123456`：demo_student、demo01..06、demo_admin、admin01..06、demo_platform_admin、plat01..06。
- 旧 handoff 中「demo123456」「并行账号密码与用户名相同」均过期，已在本轮测试中确认。
- 后续 handoff（`handoff-2026-08-14-full-pages-testing-guide.md`）第 3 节需同步更新为 `123456`。

## 6. 待办与建议下一步
1. **优先**：解决 P0-1（服务器内存），管理员才能进入 Discourse 治理（/admin、/review、举报处理）。
2. 修复 P1-1（admin-service 内部调用）与 P1-2（审核闭环）。
3. P2/P3 按排期处理；P2-4（中文搜索）影响核心体验，建议紧随 P0/P1。
4. Phase 5（多视口 390/375/320/1440 × 亮/暗 复扫）结果见 `work/testing-2026-08-14/p5-rescan.json`（staff 页面因 P0-1 无法纳入，学生论坛 + Vue 管理端已覆盖）。

## 7. 证据目录
- 全部过程 JSON/截图：`work/testing-2026-08-14/`（已 gitignore）
  - 关键证据：`p4-sso-pages2.json`、`p4-moderator-now.json`、`p4-student-now.json`、`p4-sso-trace2.json`（SSO 链路 302 轨迹）、`p6-topic-cleanup3.json`、`p4-void-invite.json`、`p3-menu-walk.json`、`results.json`（Phase 0-2 汇总）。
  - 截图：`p4-sso-trace2-final.png`、`p4-discourse-admin-*.png`、`p5-*.png`。
- 服务器只读排查脚本：`work/testing-2026-08-14/remote-diag.sh`（未做任何修改性操作）。

## 8. 风险与红线说明
- 全程未读 token/cookie 值、私信内容；未批量删除文件；未上传隐私数据。
- 生产环境副作用仅限演示账号；测试内容已清理（见 §4）。
- 服务器侧仅执行只读命令（docker logs/stats、grep production.log），未重启任何容器、未改任何配置。