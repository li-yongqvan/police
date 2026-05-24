$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$envFile = Join-Path $root ".env"

if (-not (Test-Path $envFile)) {
    $localExample = Join-Path $root ".env.local.example"
    if (Test-Path $localExample) {
        Copy-Item $localExample $envFile
        Write-Host "Created .env from .env.local.example"
    }
}

function Start-GoService($dir, $schema) {
    $cmd = @"
Set-Location '$dir'
if (Test-Path '$envFile') {
  Get-Content '$envFile' | ForEach-Object {
    if (`$_ -match '^\s*([^#][^=]+)=(.*)$') {
      Set-Item -Path env:(`$matches[1].Trim()) -Value `$matches[2].Trim()
    }
  }
}
`$env:DB_SCHEMA = '$schema'
go run ./cmd/main.go
"@
    Start-Process powershell -ArgumentList '-NoLogo', '-NoProfile', '-Command', $cmd
}

Write-Host "1) Start infrastructure:"
Write-Host "   docker compose -f infra/docker-compose.yml up -d"
Write-Host "2) First-time DB schemas (if needed):"
Write-Host "   psql ... -f migrations/init/000_init_schemas.up.sql"
Write-Host ""

Start-GoService "$root\services\user" "schema_auth"
Start-GoService "$root\services\forum" "schema_forum"
Start-GoService "$root\services\admin" "schema_admin"
Start-Sleep -Seconds 2
Start-Process powershell -ArgumentList '-NoLogo','-NoProfile','-Command',"Set-Location '$root\frontend'; npm run dev -- --host 127.0.0.1 --port 4173"
