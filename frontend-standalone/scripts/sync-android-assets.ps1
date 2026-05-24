param(
    [string]$ProjectRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
)

$dist = Join-Path $ProjectRoot "dist"
$assets = Join-Path $ProjectRoot "android\app\src\main\assets\dist"

if (-not (Test-Path $dist)) {
    Write-Host "Building frontend..."
    Push-Location $ProjectRoot
    npm run build
    Pop-Location
}

if (-not (Test-Path (Join-Path $dist "index.html"))) {
    throw "dist/index.html not found. Run npm run build first."
}

if (Test-Path $assets) {
    Remove-Item $assets -Recurse -Force
}

New-Item -ItemType Directory -Force -Path $assets | Out-Null
Copy-Item -Path (Join-Path $dist "*") -Destination $assets -Recurse -Force
Write-Host "Synced dist -> android/app/src/main/assets/dist"
