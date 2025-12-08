#!/bin/bash
# 停止所有导航服务

echo "正在停止服务..."

# 停止后端服务
pkill -f "nas-navigation" 2>/dev/null && echo "✓ 后端服务已停止" || echo "- 后端服务未运行"

# 停止前端开发服务
pkill -f "vite" 2>/dev/null && echo "✓ 前端服务已停止" || echo "- 前端服务未运行"

# 停止测试服务器
pkill -f "uvicorn main:app.*8877" 2>/dev/null && echo "✓ 测试服务器已停止" || echo "- 测试服务器未运行"

echo ""
echo "所有服务已停止"
