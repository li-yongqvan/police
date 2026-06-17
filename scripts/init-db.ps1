# Initialize PostgreSQL schemas. Run from repo root in PowerShell.
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
$ComposeFile = if ($env:COMPOSE_FILE) { $env:COMPOSE_FILE } else { "infra/docker-compose.yml" }

Write-Host "==> Starting Postgres + Redis ($ComposeFile)"
docker compose -f $ComposeFile up -d

Write-Host "==> Waiting for Postgres"
for ($i = 0; $i -lt 60; $i++) {
    docker compose -f $ComposeFile exec -T postgres pg_isready -U $DbUser -d $DbName 2>$null
    if ($LASTEXITCODE -eq 0) { break }
    Start-Sleep -Seconds 1
}

Write-Host "==> Applying schema init"
Get-Content "migrations/init/000_init_schemas.up.sql" -Raw | docker compose -f $ComposeFile exec -T postgres psql -v ON_ERROR_STOP=1 -U $DbUser -d $DbName

Write-Host ""
Write-Host "Done. Next: docker compose up -d --build; then scripts/seed-content.ps1"
