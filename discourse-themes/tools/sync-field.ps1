param(
  [Parameter(Mandatory = $true)]
  [string]$FieldFile,
  [string]$Name = 'scss',
  [int]$TargetId = 2
)
# Sync one local theme field file to default theme id 1 (default: mobile scss, target=2).
# The field value travels as a RAW file (scp -> /tmp/sync-field-value.txt).
# JSON carries only tiny metadata (name/target_id/type_id), because PowerShell 5.1
# ConvertTo-Json mangles long strings into nested objects.
# After the Rails runner applies the value, unicorn is restarted so every web
# worker drops its in-memory stylesheet digest cache; the script then waits
# until the site serves the new mobile stylesheet digest.
$ErrorActionPreference = 'Stop'
if (-not (Test-Path -LiteralPath $FieldFile)) { throw "field file not found: $FieldFile" }
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$server = 'liyongquan@122.51.233.225'
$meta = Join-Path $env:TEMP 'sync-field-meta.json'

$payload = [ordered]@{
  name = $Name
  target_id = $TargetId
  type_id = 1
}
$enc = New-Object System.Text.UTF8Encoding($false)
[System.IO.File]::WriteAllText($meta, ($payload | ConvertTo-Json -Compress), $enc)

scp $meta "${server}:/tmp/sync-field-meta.json"
if ($LASTEXITCODE -ne 0) { throw 'scp meta failed' }
scp $FieldFile "${server}:/tmp/sync-field-value.txt"
if ($LASTEXITCODE -ne 0) { throw 'scp value failed' }
scp (Join-Path $root 'sync-field.rb') "${server}:/tmp/sync-field.rb"
if ($LASTEXITCODE -ne 0) { throw 'scp runner failed' }
ssh $server 'docker cp /tmp/sync-field-meta.json app:/tmp/sync-field-meta.json; docker cp /tmp/sync-field-value.txt app:/tmp/sync-field-value.txt; docker cp /tmp/sync-field.rb app:/tmp/sync-field.rb'
if ($LASTEXITCODE -ne 0) { throw 'docker cp failed' }
ssh $server 'docker exec -u discourse -w /var/www/discourse app bash -lc ''RAILS_ENV=production bundle exec rails runner /tmp/sync-field.rb'''
if ($LASTEXITCODE -ne 0) { throw 'rails runner failed' }

ssh $server 'docker exec app sv restart unicorn'
if ($LASTEXITCODE -ne 0) { throw 'unicorn restart failed' }

$probeUrl = 'http://122.51.233.225:8080/latest?theme-sync-probe=' + (Get-Date -Format 'yyyyMMddHHmmss')
$ready = $false
for ($i = 0; $i -lt 48; $i++) {
  Start-Sleep -Seconds 5
  try {
    $r = Invoke-WebRequest -Uri $probeUrl -UseBasicParsing -TimeoutSec 15
    $m = [regex]::Match($r.Content, '/stylesheets/mobile_theme_1_[a-f0-9]+\.css')
    if ($m.Success) {
      Write-Output ("served mobile css: " + $m.Value)
      $ready = $true
      break
    }
  } catch { }
}
if (-not $ready) { Write-Output 'WARN: site not ready or stylesheet link not found after unicorn restart' }