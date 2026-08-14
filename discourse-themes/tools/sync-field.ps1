param(
  [Parameter(Mandatory = $true)]
  [string]$FieldFile,
  [string]$Name = 'scss',
  [int]$TargetId = 2
)
# 将本地主题字段文件同步到默认主题 ID 1（默认: mobile scss, target=2）
$ErrorActionPreference = 'Stop'
if (-not (Test-Path -LiteralPath $FieldFile)) { throw "field file not found: $FieldFile" }
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$server = 'liyongquan@122.51.233.225'
$tmp = Join-Path $env:TEMP 'sync-field.json'

$payload = [ordered]@{
  name = $Name
  target_id = $TargetId
  type_id = 1
  value = (Get-Content -LiteralPath $FieldFile -Raw)
}
[System.IO.File]::WriteAllText($tmp, ($payload | ConvertTo-Json -Depth 5), (New-Object System.Text.UTF8Encoding($false)))

scp $tmp "${server}:/tmp/sync-field.json"
scp (Join-Path $root 'sync-field.rb') "${server}:/tmp/sync-field.rb"
ssh $server 'docker cp /tmp/sync-field.json app:/tmp/sync-field.json; docker cp /tmp/sync-field.rb app:/tmp/sync-field.rb'
ssh $server 'docker exec -u discourse -w /var/www/discourse app bash -lc ''RAILS_ENV=production bundle exec rails runner /tmp/sync-field.rb'''
