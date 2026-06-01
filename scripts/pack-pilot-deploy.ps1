# Pack project for VPS upload (excludes heavy dirs).
$ErrorActionPreference = "Stop"
$root = (Resolve-Path "$PSScriptRoot\..").Path
$out = Join-Path $root "pilot-deploy.tgz"
Push-Location $root
try {
  if (Test-Path pilot-deploy.tgz) { Remove-Item pilot-deploy.tgz -Force }
  & tar -czf pilot-deploy.tgz `
    --exclude=node_modules `
    --exclude=frontend/node_modules `
    --exclude=.git `
    --exclude=pilot-deploy.tgz `
    --exclude=reference `
    docker-compose.yml docker-compose.server.yml docker-compose.prod.yml docker-compose.cn-mirror.yml `
    .env.production.example .env.smoke.example `
    frontend mobile-web services scripts docs nginx infra migrations `
    DEPLOY.md README.md
  $mb = [math]::Round((Get-Item $out).Length / 1MB, 2)
  Write-Host "Created $out ($mb MB)"
} finally {
  Pop-Location
}
