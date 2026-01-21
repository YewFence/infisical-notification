#!/bin/bash
# Backend 启动脚本
# 使用方式: ./run.sh

export CGO_ENABLED=0
export GIN_MODE=debug

echo "🚀 启动 Todo Backend 服务..."
echo "📌 CGO_ENABLED=$CGO_ENABLED (使用纯 Go SQLite 驱动)"

go run main.go
