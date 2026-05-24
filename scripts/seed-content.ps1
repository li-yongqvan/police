$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location (Join-Path $Root "..")

if (Test-Path ".env") {
    Get-Content ".env" | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
            Set-Item -Path "env:$($matches[1].Trim())" -Value $matches[2].Trim()
        }
    }
}

$DbUser = if ($env:POSTGRES_USER) { $env:POSTGRES_USER } else { "ai_forum" }
$DbName = if ($env:POSTGRES_DB) { $env:POSTGRES_DB } else { "ai_forum" }

for ($i = 0; $i -lt 90; $i++) {
    $ok = docker compose exec -T postgres psql -U $DbUser -d $DbName -tAc "SELECT to_regclass('schema_forum.posts') IS NOT NULL;" 2>$null
    if ($ok -match "t") { break }
    Start-Sleep -Seconds 2
}

Get-Content "scripts/seed/001_content.sql" -Raw | docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U $DbUser -d $DbName
Write-Host "Seed complete."
