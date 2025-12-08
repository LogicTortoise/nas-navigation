#!/bin/bash
# 重启 FastAPI 测试服务器
pkill -f "uvicorn main:app.*8877" 2>/dev/null
sleep 1
cd /Users/bytedance/Documents/Personal/github/navgation/test_server
# 使用 arm64 架构运行 Python
nohup arch -arm64 /usr/local/bin/python3 -m uvicorn main:app --host 0.0.0.0 --port 8877 > /tmp/test_server.log 2>&1 &
echo "Test server restarted on port 8877"
