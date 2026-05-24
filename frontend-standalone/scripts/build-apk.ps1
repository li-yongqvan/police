param(
    [ValidateSet("debug", "release")]
    [string]$Variant = "debug"
)

$root = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$android = Join-Path $root "android"

& (Join-Path $PSScriptRoot "sync-android-assets.ps1") -ProjectRoot $root

Push-Location $android
try {
    if ($Variant -eq "release") {
        .\gradlew.bat assembleRelease
        $apk = Get-ChildItem -Path "app\build\outputs\apk\release" -Filter "*.apk" -Recurse | Select-Object -First 1
    } else {
        .\gradlew.bat assembleDebug
        $apk = Get-ChildItem -Path "app\build\outputs\apk\debug" -Filter "*.apk" -Recurse | Select-Object -First 1
    }

    if ($apk) {
        Write-Host ""
        Write-Host "APK built:" -ForegroundColor Green
        Write-Host $apk.FullName
    } else {
        throw "APK not found after build."
    }
} finally {
    Pop-Location
}
