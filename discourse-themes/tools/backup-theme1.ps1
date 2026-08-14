# 备份默认主题 ID 1：导出字段/设置/配色方案到本地 work/theme-backup/archives/
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$workRoot = Split-Path -Parent (Split-Path -Parent $root)
$archives = Join-Path $workRoot 'work\theme-backup\archives'
New-Item -ItemType Directory -Force -Path $archives | Out-Null

$server = 'liyongquan@122.51.233.225'
scp (Join-Path $root 'backup-theme1.rb') "${server}:/tmp/backup-theme1.rb"
ssh $server 'docker cp /tmp/backup-theme1.rb app:/tmp/backup-theme1.rb'
ssh $server 'docker exec -u discourse -w /var/www/discourse app bash -lc ''RAILS_ENV=production bundle exec rails runner /tmp/backup-theme1.rb'''
ssh $server 'docker cp app:/tmp/theme-1-backup.json /tmp/theme-1-backup.json'

$ts = Get-Date -Format 'yyyyMMdd-HHmmss'
$dest = Join-Path $archives "theme-1-$ts.json"
scp "${server}:/tmp/theme-1-backup.json" $dest

$j = Get-Content -LiteralPath $dest -Raw | ConvertFrom-Json
Write-Output "backup saved: $dest"
Write-Output "  fields=$($j.fields.Count) color_schemes=$($j.color_schemes.Count) settings=$($j.settings.Count) exported_at=$($j.exported_at)"
