# 测试指南

## 测试服务器

项目包含一个 FastAPI 测试服务器，用于验证链接管理和重启功能。

### 启动测试服务器

```bash
cd test_server
arch -arm64 /usr/local/bin/python3 -m uvicorn main:app --host 0.0.0.0 --port 8877
```

或使用重启脚本：

```bash
bash test_server/restart.sh
```

### 测试服务器端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/` | GET | 返回服务器信息和运行时间 |
| `/health` | GET | 健康检查端点 |

### 验证服务器运行

```bash
curl http://localhost:8877/
```

响应示例：
```json
{
  "message": "Hello from Test Server!",
  "start_time": "2025-12-08T00:00:12.657250",
  "uptime_seconds": 13.079728
}
```

## 功能测试

### 1. 添加带重启脚本的链接

1. 点击页面顶部的 `+ 添加` 按钮
2. 填写表单：
   - URL: `http://localhost:8877`
   - 名称: `Test Server`
   - 描述: `FastAPI测试服务`
   - 重启脚本: `/path/to/navgation/test_server/restart.sh`
3. 点击保存

### 2. 测试三点菜单

1. 将鼠标悬停在链接卡片上
2. 点击右上角的三点图标 (`⋮`)
3. 应显示下拉菜单：
   - **编辑** - 编辑链接信息
   - **重启** - 执行重启脚本（仅当配置了重启脚本时显示）
   - **删除** - 删除链接

### 3. 测试重启功能

1. 记录当前服务器 uptime: `curl http://localhost:8877/ | jq .uptime_seconds`
2. 点击三点菜单中的"重启"
3. 应显示"执行成功"对话框
4. 再次检查 uptime，应该重置为较小的值

### 4. 验证无重启脚本的链接

对于没有配置 `restart_script` 的链接（如 GitHub），三点菜单中不应显示"重启"选项。

## E2E 测试（使用 Playwright）

```javascript
// 导航到页面
await page.goto('http://localhost:3002');

// 点击添加按钮
await page.getByRole('button', { name: '+ 添加' }).click();

// 填写表单
await page.getByRole('textbox', { name: 'https://example.com' }).fill('http://localhost:8877');
await page.getByRole('textbox', { name: '例如：Google、GitHub' }).fill('Test Server');

// 保存
await page.getByRole('button', { name: '保存链接' }).click();

// 点击三点菜单
await page.getByRole('button').filter({ hasText: /^$/ }).first().click();

// 点击重启
await page.getByRole('button', { name: '重启' }).click();
```
