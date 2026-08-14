# 服务器减压计划书 · P0 修复（Discourse staff 全站 500）

Date: 2026-08-14
状态：待用户确认后执行（默认先执行 Phase A，Phase B-D 逐项确认）
分支：codex/discourse-rebuild

## 1. 背景与根因（已实测确认）

现象：管理员/版主经 SSO 登录 Discourse 后访问 /、/latest、/categories、/about 全部 500（30s 超时），/admin/dashboard 超时；学生、游客正常。详见 `docs/handoffs/handoff-2026-08-14-test-results.md` P0-1。

根因链条（全部为只读排查证据，未做任何修改）：

1. **直接原因**：容器内自动更新检查的 `git fetch origin`（`--filter=blob:none` 部分克隆）已挂起 18 分钟以上（`work/testing-2026-08-14/remote-mem-recon3.sh` 输出），另有 6 个 `git merge-base --is-ancestor` 子进程挂起 23 分钟+。
2. 挂起的 git 进程每个 RSS 约 230MB，6 个合计约 1.4GB；且更新检查周期性触发，进程反复堆积（docker logs 中 `reaped unknown subprocess exit 128` 与其吻合）。
3. 结果：app 容器内存 2.16GB / 3.6GB（59%），主机 3.7GB 内存 + 1.9GB Swap（已用 1.4GB）几乎耗尽 → unicorn（3 worker）饥饿，`worker=0 timed out`，请求 30s 被掐断。
4. staff 页面不做匿名缓存、渲染更重，所以只有 staff 会话可见 500；学生/游客走缓存幸免。
5. 推测挂起原因：容器到 GitHub 的出网受限/超慢（`git fetch` 无超时保护）。

## 2. 方案总览

| 阶段 | 动作 | 风险/中断 | 是否可逆 |
|------|------|-----------|----------|
| Phase A | `./launcher restart app` 立即重启 Discourse | 论坛中断 1–3 分钟 | 是（无状态变更） |
| Phase B | 关闭 Discourse 自动更新检查（`version_checks=false`） | 无中断 | 是 |
| Phase C | UNICORN_WORKERS 3→2（app.yml + rebuild） | 中断 10–20 分钟 | 是 |
| Phase D | Swap 扩容 1.9G→4G | 无中断 | 是 |

## 3. Phase A · 立即重启（止血）

```bash
# 服务器上执行（app 由 /home/liyongquan/discourse 的 launcher 管理）
cd /home/liyongquan/discourse && ./launcher restart app
```

- 作用：杀掉挂起的 git 进程与旧 worker，释放约 1.5GB 内存。
- 验证：
  1. `docker stats --no-stream`：app 内存回落到 ~0.6–0.9GB；
  2. 复用 `work/testing-2026-08-14/run-p4-sso-pages2.cjs`（demo_platform_admin SSO 后扫 /latest、/admin/dashboard）确认 500 消失；
  3. `docker exec app ps aux --sort=-rss | head`：无 `git merge-base`/`git fetch` 挂起进程。
- 回滚：无。数据在 PG volume，重启不丢数据。

## 4. Phase B · 阻止复发（关键）

挂起的 `git fetch` 来自 Discourse 自带的版本检查任务。两个途径（任选其一，推荐 2）：

1. 后台 UI：管理 → 设置 → 软件 → 「检查更新」（version_checks）关掉。但 staff 页面当前 500，UI 可能打不开，所以作为 Phase A 后的首选尝试。
2. 兜底（不依赖 Web UI）：
   ```bash
   cd /home/liyongquan/discourse && ./launcher enter app
   # 进入容器后：
   rails runner "SiteSetting.version_checks = false"
   exit
   ```
   如 rails runner 报环境警告可加 `DISCOURSE_DEV_HOSTS=` 前缀或改用 `RAILS_ENV=production`。

- 验证：连续观察 30 分钟，确认不再出现 `git fetch`/`git merge-base` 挂起进程；`docker stats` 内存曲线平稳。
- 备注：若确认是服务器出网到 GitHub 被限，此步是根治；若将来想恢复自动检查，把设置改回 true 即可。

## 5. Phase C · 降低单机负载（可选，需选时间窗口）

- 编辑 `/home/liyongquan/discourse/containers/app.yml`：`UNICORN_WORKERS: 2`，然后 `./launcher rebuild app`（重建镜像，约 10–20 分钟，期间论坛不可用）。
- 收益：每 worker 稳态 ~400MB，3→2 可省 ~400MB，给 3.7GB 小机器留出缓冲。
- 回滚：改回 3 再 rebuild。
- 注意：rebuild 会跑 bundle，对内存本就紧张的机器压力大；建议先完成 Phase A/B，观察内存水位再决定。

## 6. Phase D · Swap 扩容（可选，需 root）

- 现状：`/swap.img` 1.9G，已用 1.4G；主机 3.7G 内存。
- 步骤（需 sudo/root，执行前我会再给出精确命令并逐条确认）：
  1. 新增 2G swapfile（`fallocate -l 2G /swap2.img` 等）；
  2. `mkswap`/`swapon`，加入 `/etc/fstab`；
  3. 可调 `vm.swappiness`。
- 回滚：`swapoff /swap2.img` 并删除 fstab 行。

## 7. Phase E · 验证与观察

- [ ] staff SSO 复测通过（P0 关闭）
- [ ] 学生/游客基线不回归（/latest、/categories、/about、/top 抽查 200）
- [ ] 观察 24h：`docker stats app` 内存不再爬升、无 git 挂起进程、无 `worker timed out`
- [ ] 复测记录写回 `docs/handoffs/handoff-2026-08-14-test-results.md` 并 commit + push

## 8. 中长期建议（不在本轮执行）

- 服务器加内存至 8G（Discourse 单容器建议 2G+ 余量）；或将 Discourse 迁到独立服务器。
- 评估裁剪插件（chat/reactions/solved/templates/presence/narrative-bot/lazy-videos/poll/spoiler/footnote/voting 等）以降内存。
- 若服务器无稳定出网到 GitHub，将更新检查固定为手动（Phase B 已覆盖）。

## 9. 风险清单

- Phase A 重启会造成 1–3 分钟论坛不可用：请确认执行时间窗口（建议避开使用高峰）。
- Phase C rebuild 时间较长且内存压力大：非必要不执行，等 A/B 见效后观察再定。
- 全部操作仅动 Discourse 侧，不碰 ai-forum 栈（frontend/user-service/admin-service/postgres/redis）。
- 每步可逆，且执行前逐条向用户确认。