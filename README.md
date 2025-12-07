# NAS Navigation - 家用NAS导航网站

一个轻量级的家用NAS导航网站，提供快捷访问常用服务和一键执行脚本的功能。

## 功能特性

- **链接管理**: 添加、编辑、删除导航链接
- **分类管理**: 支持自定义分类，链接归类展示
- **脚本执行**: 支持配置启动脚本，一键重启服务
- **响应式设计**: 适配桌面端和移动端
- **搜索功能**: 快速搜索链接和服务

## 技术栈

- **后端**: Go + Gin + GORM
- **前端**: Vue 3 + Vite + Tailwind CSS
- **数据库**: SQLite

## 快速开始

### 开发环境

```bash
# 一键启动开发环境
./start-dev.sh

# 或分别启动
# 后端
cd backend
CGO_ENABLED=1 go build -o nas-navigation .
./nas-navigation

# 前端
cd frontend
npm install
npm run dev
```

访问 http://localhost:3000 查看应用。

### 生产部署

```bash
# 构建前端
cd frontend
npm run build

# 构建后端
cd ../backend
CGO_ENABLED=1 go build -o nas-navigation .

# 运行
./nas-navigation
```

访问 http://localhost:8080 查看应用。

## 项目结构

```
navgation/
├── backend/                # Go后端
│   ├── config/            # 配置管理
│   ├── database/          # 数据库
│   ├── handlers/          # HTTP处理器
│   ├── models/            # 数据模型
│   ├── router/            # 路由配置
│   ├── services/          # 业务逻辑
│   └── main.go            # 入口文件
├── frontend/              # Vue前端
│   ├── src/
│   │   ├── api/           # API封装
│   │   ├── components/    # Vue组件
│   │   ├── stores/        # 状态管理
│   │   └── App.vue        # 主组件
│   └── package.json
├── data/                  # 数据目录
│   └── navigation.db      # SQLite数据库
├── docs/                  # 文档
└── scripts/               # 用户脚本目录
```

## API文档

### 分类管理

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/categories | 获取所有分类 |
| POST | /api/categories | 创建分类 |
| PUT | /api/categories/:id | 更新分类 |
| DELETE | /api/categories/:id | 删除分类 |

### 链接管理

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/links | 获取所有链接 |
| GET | /api/links/search?q=xxx | 搜索链接 |
| POST | /api/links | 创建链接 |
| PUT | /api/links/:id | 更新链接 |
| DELETE | /api/links/:id | 删除链接 |

### 脚本管理

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/scripts | 获取所有脚本 |
| POST | /api/scripts | 创建脚本 |
| PUT | /api/scripts/:id | 更新脚本 |
| DELETE | /api/scripts/:id | 删除脚本 |
| POST | /api/scripts/:id/execute | 执行脚本 |

## 环境变量

| 变量 | 默认值 | 描述 |
|------|--------|------|
| PORT | 8080 | 服务端口 |
| DB_PATH | ../data/navigation.db | 数据库路径 |
| SCRIPT_DIR | ../scripts | 脚本目录 |
| STATIC_DIR | ../frontend/dist | 静态文件目录 |
| DEBUG | true | 调试模式 |

## License

MIT
