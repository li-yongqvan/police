param(
  [Parameter(Mandatory = $true)]
  [string]$ArchivePath
)
# 从 work/theme-backup/archives 里的 JSON 回滚默认主题 ID 1
$ErrorActionPreference = 'Stop'
if (-not (Test-Path -LiteralPath $ArchivePath)) { throw "archive not found: $ArchivePath" }
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$server = 'liyongquan@122.51.233.225'

scp $ArchivePath "${server}:/tmp/restore-theme1.json"
scp (Join-Path $root 'restore-theme1.rb') "${server}:/tmp/restore-theme1.rb"
ssh $server 'docker cp /tmp/restore-theme1.json app:/tmp/restore-theme1.json; docker cp /tmp/restore-theme1.rb app:/tmp/restore-theme1.rb'
ssh $server 'docker exec -u discourse -w /var/www/discourse app bash -lc ''RAILS_ENV=production bundle exec rails runner /tmp/restore-theme1.rb'''
