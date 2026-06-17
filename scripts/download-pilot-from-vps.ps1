# Download ai-forum-migration.tar.gz from US VPS 107.172.138.10 (root).
# Run in PowerShell; enter root password when prompted (unless SSH key is set).
$ErrorActionPreference = "Stop"
$Host_ = "107.172.138.10"
$User = "root"
$File = "ai-forum-migration.tar.gz"
$Dest = (Resolve-Path "$PSScriptRoot\..").Path
$Local = Join-Path $Dest $File

$candidates = @(
  "/root/$File",
  "/tmp/$File",
  "/opt/ai-forum/$File",
  "/opt/$File"
)

Write-Host "Trying ${User}@${Host_} ..."
foreach ($remote in $candidates) {
  Write-Host "  -> $remote"
  try {
    scp -o StrictHostKeyChecking=accept-new "${User}@${Host_}:${remote}" $Local 2>$null
    if (Test-Path $Local) {
      $mb = [math]::Round((Get-Item $Local).Length / 1MB, 2)
      Write-Host "OK: $Local ($mb MB)"
      exit 0
    }
  } catch {
    # try next path
  }
}

Write-Host ""
Write-Host "Auto paths failed. Find file on server:"
Write-Host "  ssh root@107.172.138.10 `"find / -name '$File' 2>/dev/null`""
Write-Host ""
$custom = Read-Host "Remote full path (e.g. /root/ai-forum-migration.tar.gz)"
if (-not $custom.Trim()) { exit 1 }
scp -o StrictHostKeyChecking=accept-new "${User}@${Host_}:$($custom.Trim())" $Local
$mb = [math]::Round((Get-Item $Local).Length / 1MB, 2)
Write-Host "OK: $Local ($mb MB)"
