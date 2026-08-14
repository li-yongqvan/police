# Discourse 默认主题部署 Runbook

> 适用范围：AI 智联论坛 Discourse 实例的主题同步、备份与回滚。
> 最后更新：2026-08-14

## 1. 环境

| 项 | 值 |
|----|----|
| 服务器 | `122.51.233.225`，SSH 用户 `liyongquan`（密钥直连，有 Docker 权限，无需 sudo） |
| Discourse 容器 | `app`，工作目录 `/var/www/discourse` |
| 访问入口 | 前端平台 `http://122.51.233.225:8888`；Discourse 直连 `http://122.51.233.225:8080` |
| 默认主题 | ID 1 `ai-forum-premium-preview`（default=true） |
| 本地主题源码 | `discourse-themes/ai-forum-premium-preview/`（common/desktop/mobile/about，已随主仓库 git 版本化） |

当前线上主题仅剩：ID 1（默认）+ 内置 `Horizon`（-2）、`Foundation`（-1，勿删）。
历史预览主题 2/3/4/6 已于 2026-08-14 清理，删除前存档见 `work/mui07-prep/preview-theme-backups/`。

## 2. 工具（`discourse-themes/tools/`，已提交 git）

| 脚本 | 用途 |
|------|------|
| `backup-theme1.ps1` | 导出默认主题全部字段/设置/配色方案 → 本地 `work/theme-backup/archives/theme-1-<时间戳>.json` |
| `sync-field.ps1` | 同步单个本地字段文件到默认主题（默认 `scss` / target 2 = mobile） |
| `restore-theme1.ps1 <归档.json>` | 从归档回滚默认主题（含 SCSS 重烘焙） |

用法示例（PowerShell 5.1）：

```powershell
# 每次同步前先备份
.\discourse-themes\tools\backup-theme1.ps1

# 同步 mobile SCSS（默认参数）
.\discourse-themes\tools\sync-field.ps1 -FieldFile .\discourse-themes\ai-forum-premium-preview\mobile\mobile.scss

# 同步 desktop / common / about 时显式指定
.\discourse-themes\tools\sync-field.ps1 -FieldFile .\...\desktop\desktop.scss -TargetId 1
.\discourse-themes\tools\sync-field.ps1 -FieldFile .\...\common\common.scss -TargetId 0

# 回滚（紧急情况）
.\discourse-themes\tools\restore-theme1.ps1 .\work\theme-backup\archives\theme-1-20260814-100109.json
```

target_id 对照：0=common，1=desktop，2=mobile，8=about。

## 3. 执行链与注意事项

1. 所有脚本走同一执行链：`scp → docker cp 进容器 → docker exec -u discourse -w /var/www/discourse app bash -lc 'RAILS_ENV=production bundle exec rails runner /tmp/x.rb'`。
2. **PowerShell 5.1 引号坑**：远程命令里只能用单引号，PowerShell 中用 `''` 转义；嵌入双引号会被 PS 5.1 剥掉导致命令静默变错（如 `bash -lc "echo C"` 会变成 `bash -lc echo C`）。所有工具脚本已按此编写，改动时保持同一写法。
   - **ConvertTo-Json 长字符串坑**：PS 5.1 的 `ConvertTo-Json` 会把约 25KB 以上的长字符串错误序列化成嵌套对象（2026-08-14 曾致文件膨胀到 3.8MB、线上报 `Value is too long (maximum is 1048576 characters)`，`save!` 已回滚）。字段值一律走 raw 文件传输（`scp` 到 `/tmp/sync-field-value.txt`），JSON 只携带 name/target_id/type_id 小元数据（`sync-field.ps1` 已按此实现）。
3. **重烘焙**：直接改 `ThemeField#value` 后 `ensure_baked!` 不会重烘焙；必须先 `f.value_baked = nil` 再 `save!` 与 `ensure_baked!`（`sync-field.rb`/`restore-theme1.rb` 已内置）。
4. **缓存传播**：样式表有分布式缓存 + 匿名页面缓存，新 digest 异步传播；同步后立即验证的 URL 形如 `/stylesheets/mobile_theme_1_<digest>.css`。
   - **DistributedCache 纯内存坑**：样式 digest 缓存在各 worker 进程内存，runner 进程的 hash 为空、`clear_regex` 与显式 MessageBus delete 均无法可靠触达 web worker（实测无效）；`sync-field.ps1` 在 runner 执行后统一 `sv restart unicorn` 重置内存缓存并轮询等待新 digest（期间短暂 502），同时清 `ANON_CACHE_*` 匿名页面缓存。
   - **DistributedCache 纯内存坑**：样式 digest 缓存在各进程内存（MessageBus 同步），runner 进程 hash 为空时 `clear_regex` 清不到任何键；`sync-field.rb` 改为显式发 delete 消息，并清 `ANON_CACHE_*` 匿名页面缓存。若新 digest 迟迟不生效，`sv restart unicorn` 可强制收敛（短暂 502）。
5. **同步后验证**：按当前 MUI 票的验证脚本检查移动端（390/375/1440）`docOverflow=0`、`bodyOverflow=0` 与关键控件可见性；再让用户看前后截图。
6. 预览主题流程已废弃：今后本地 CSS 直接同步默认主题 ID 1，因此「同步前先备份」是强制步骤。

## 4. 历史记录

- 2026-08-14：清理预览主题 2/3/4/6；建立备份/同步/回滚工具；首次全量备份 `work/theme-backup/archives/theme-1-20260814-100109.json`。
- MUI-05/MUI-06：mobile.scss 已含 `// MUI-05:`、`// MUI-06:` 两个区块并线上验证通过。
- 2026-08-14：修复 `sync-field.ps1` 的 ConvertTo-Json 长字符串序列化 bug（MUI-07 同步时暴露），改为 raw 文件传输；线上 mobile 字段核对无损（23789 字节）。
