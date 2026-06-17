$root = Split-Path -Parent $MyInvocation.MyCommand.Path

Start-Process powershell -ArgumentList '-NoLogo','-NoProfile','-Command',"Set-Location '$root\services\user'; go run ./cmd/server" -WindowStyle Hidden
Start-Process powershell -ArgumentList '-NoLogo','-NoProfile','-Command',"Set-Location '$root\services\forum'; go run ./cmd/server" -WindowStyle Hidden
Start-Process powershell -ArgumentList '-NoLogo','-NoProfile','-Command',"Set-Location '$root\services\admin'; go run ./cmd/server" -WindowStyle Hidden
Start-Sleep -Seconds 2
Start-Process powershell -ArgumentList '-NoLogo','-NoProfile','-Command',"Set-Location '$root\frontend'; npm run dev -- --host 127.0.0.1 --port 4173" -WindowStyle Hidden
