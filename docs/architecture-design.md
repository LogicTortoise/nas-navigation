# NAS Navigation - 家用NAS导航网站架构设计

## 1. 项目概述

一个轻量级的家用NAS导航网站，提供快捷访问常用服务和一键执行脚本的功能。

### 1.1 核心功能
- **链接管理**: 添加、编辑、删除导航链接
- **分类管理**: 支持自定义分类，链接归类展示
- **脚本执行**: 支持配置启动脚本，一键重启服务
- **响应式设计**: 适配桌面端和移动端
- **搜索功能**: 快速搜索链接和服务

### 1.2 技术选型

| 组件 | 技术选择 | 理由 |
|------|----------|------|
| 后端 | Go + Gin | 轻量、高性能、单二进制部署 |
| 前端 | Vue 3 + Vite | 现代化、轻量、开发体验好 |
| 数据库 | SQLite | 轻量、无需额外服务、适合小型应用 |
| UI框架 | Tailwind CSS | 响应式、实用优先、定制性强 |

## 2. 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                        用户界面                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   Vue 3 SPA                          │   │
│  │  - 导航主页 (链接卡片、分类展示)                       │   │
│  │  - 管理页面 (链接/分类/脚本管理)                       │   │
│  │  - 响应式布局 (桌面/移动端适配)                        │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Go Backend (Gin)                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Link API   │  │ Category API │  │  Script API  │      │
│  │  - CRUD操作  │  │  - CRUD操作  │  │  - 执行脚本  │      │
│  │  - 搜索     │  │  - 排序     │  │  - 状态查询  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                              │                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                    Service Layer                     │   │
│  │  - 业务逻辑处理                                       │   │
│  │  - 脚本安全执行 (沙箱、超时控制)                       │   │
│  └─────────────────────────────────────────────────────┘   │
│                              │                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                    Data Layer                        │   │
│  │  - GORM ORM                                          │   │
│  │  - SQLite Driver                                     │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      SQLite Database                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    links     │  │  categories  │  │   scripts    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## 3. 数据模型

### 3.1 ER图

```
┌───────────────────┐       ┌───────────────────┐
│     Category      │       │       Link        │
├───────────────────┤       ├───────────────────┤
│ id (PK)           │       │ id (PK)           │
│ name              │◄──────│ category_id (FK)  │
│ icon              │   1:N │ name              │
│ sort_order        │       │ url               │
│ created_at        │       │ icon              │
│ updated_at        │       │ description       │
└───────────────────┘       │ sort_order        │
                            │ script_id (FK)    │──────┐
                            │ created_at        │      │
                            │ updated_at        │      │
                            └───────────────────┘      │
                                                       │
                            ┌───────────────────┐      │
                            │      Script       │      │
                            ├───────────────────┤      │
                            │ id (PK)           │◄─────┘
                            │ name              │   N:1
                            │ command           │
                            │ description       │
                            │ timeout           │
                            │ created_at        │
                            │ updated_at        │
                            └───────────────────┘
```

### 3.2 表结构

```sql
-- 分类表
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(100),
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 链接表
CREATE TABLE links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    url VARCHAR(500) NOT NULL,
    icon VARCHAR(500),
    description VARCHAR(500),
    sort_order INTEGER DEFAULT 0,
    script_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (script_id) REFERENCES scripts(id)
);

-- 脚本表
CREATE TABLE scripts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    command TEXT NOT NULL,
    description VARCHAR(500),
    timeout INTEGER DEFAULT 30,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 用户设置表
CREATE TABLE settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## 4. API设计

### 4.1 RESTful API

#### 分类管理
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/categories | 获取所有分类 |
| POST | /api/categories | 创建分类 |
| PUT | /api/categories/:id | 更新分类 |
| DELETE | /api/categories/:id | 删除分类 |

#### 链接管理
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/links | 获取所有链接 |
| GET | /api/links/search | 搜索链接 |
| POST | /api/links | 创建链接 |
| PUT | /api/links/:id | 更新链接 |
| DELETE | /api/links/:id | 删除链接 |

#### 脚本管理
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/scripts | 获取所有脚本 |
| POST | /api/scripts | 创建脚本 |
| PUT | /api/scripts/:id | 更新脚本 |
| DELETE | /api/scripts/:id | 删除脚本 |
| POST | /api/scripts/:id/execute | 执行脚本 |

#### 设置管理
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/settings | 获取设置 |
| PUT | /api/settings | 更新设置 |

## 5. 项目结构

```
nas-navigation/
├── backend/
│   ├── main.go                 # 入口文件
│   ├── config/
│   │   └── config.go           # 配置管理
│   ├── database/
│   │   ├── database.go         # 数据库初始化
│   │   └── migrations.go       # 数据迁移
│   ├── models/
│   │   ├── category.go         # 分类模型
│   │   ├── link.go             # 链接模型
│   │   ├── script.go           # 脚本模型
│   │   └── setting.go          # 设置模型
│   ├── handlers/
│   │   ├── category.go         # 分类处理器
│   │   ├── link.go             # 链接处理器
│   │   ├── script.go           # 脚本处理器
│   │   └── setting.go          # 设置处理器
│   ├── services/
│   │   ├── category.go         # 分类服务
│   │   ├── link.go             # 链接服务
│   │   ├── script.go           # 脚本服务 (含执行逻辑)
│   │   └── setting.go          # 设置服务
│   └── router/
│       └── router.go           # 路由配置
├── frontend/
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   ├── tailwind.config.js
│   └── src/
│       ├── main.js
│       ├── App.vue
│       ├── api/
│       │   └── index.js        # API封装
│       ├── components/
│       │   ├── Header.vue      # 头部组件
│       │   ├── SearchBar.vue   # 搜索栏
│       │   ├── CategorySection.vue  # 分类区块
│       │   ├── LinkCard.vue    # 链接卡片
│       │   ├── AddLinkModal.vue     # 添加链接弹窗
│       │   └── ScriptButton.vue     # 脚本执行按钮
│       ├── views/
│       │   ├── Home.vue        # 主页
│       │   └── Admin.vue       # 管理页
│       └── stores/
│           └── app.js          # 状态管理
├── data/
│   └── navigation.db           # SQLite数据库文件
├── scripts/                    # 用户脚本目录
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## 6. 安全考虑

### 6.1 脚本执行安全
- **超时控制**: 每个脚本设置最大执行时间
- **命令白名单**: 可配置允许执行的命令前缀
- **执行日志**: 记录所有脚本执行历史
- **权限隔离**: 脚本以低权限用户执行

### 6.2 输入验证
- URL格式验证
- 脚本命令过滤危险字符
- SQL注入防护 (ORM自带)
- XSS防护

## 7. 部署方案

### 7.1 Docker部署 (推荐)
```yaml
version: '3'
services:
  nas-navigation:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
      - ./scripts:/app/scripts
    restart: unless-stopped
```

### 7.2 二进制部署
```bash
# 构建
go build -o nas-navigation ./backend

# 运行
./nas-navigation --port 8080 --data ./data
```

## 8. 开发计划

### Phase 1: 基础功能
1. 后端框架搭建
2. 数据库模型实现
3. CRUD API实现
4. 前端基础页面

### Phase 2: 核心功能
1. 链接管理完善
2. 分类管理完善
3. 搜索功能
4. 响应式布局

### Phase 3: 高级功能
1. 脚本执行功能
2. 图标上传/选择
3. 拖拽排序
4. 设置页面

### Phase 4: 优化完善
1. 性能优化
2. 安全加固
3. Docker打包
4. 文档完善
