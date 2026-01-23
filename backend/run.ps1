# Backend 启动脚本
# 使用方式: .\run.ps1

$env:CGO_ENABLED = "0"
$env:GIN_MODE = "debug"

Write-Host "🚀 启动 Todo Backend 服务..." -ForegroundColor Green
Write-Host "📌 CGO_ENABLED=$env:CGO_ENABLED (使用纯 Go SQLite 驱动)" -ForegroundColor Cyan

go run main.go
