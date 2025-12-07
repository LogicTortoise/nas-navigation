#!/bin/bash

# NAS Navigation 开发环境启动脚本

set -e

echo "Starting NAS Navigation development environment..."

# 切换到项目目录
cd "$(dirname "$0")"

# 启动后端
echo "Building and starting backend..."
cd backend
CGO_ENABLED=1 go build -o nas-navigation . 2>/dev/null || {
    echo "Installing dependencies..."
    GOPROXY=https://goproxy.cn,direct go mod tidy
    CGO_ENABLED=1 go build -o nas-navigation .
}
./nas-navigation &
BACKEND_PID=$!
echo "Backend started (PID: $BACKEND_PID)"
cd ..

# 启动前端
echo "Starting frontend..."
cd frontend
if [ ! -d "node_modules" ]; then
    echo "Installing frontend dependencies..."
    npm install
fi
npm run dev &
FRONTEND_PID=$!
echo "Frontend started (PID: $FRONTEND_PID)"
cd ..

echo ""
echo "========================================"
echo "NAS Navigation is running!"
echo "Frontend: http://localhost:3000"
echo "Backend:  http://localhost:8080"
echo "========================================"
echo ""
echo "Press Ctrl+C to stop all services"

# 捕获退出信号
trap "kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit" INT TERM

# 等待
wait
