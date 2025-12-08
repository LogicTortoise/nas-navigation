# NAS Navigation - 家用NAS导航网站

一个现代化的家用NAS导航网站，提供快捷访问常用服务和一键重启服务的功能。

## 功能特性

- **链接管理**: 添加、编辑、删除导航链接
- **分类管理**: 支持自定义分类，链接归类展示
- **服务重启**: 为链接配置重启脚本，一键重启服务
- **三点菜单**: 现代化的下拉菜单设计，编辑/重启/删除操作
- **响应式设计**: 适配桌面端和移动端
- **搜索功能**: 快速搜索链接和服务
- **毛玻璃 UI**: 现代化的视觉设计

## 技术栈

### 后端
- **Go 1.21+** - 主要编程语言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **SQLite** - 数据库

### 前端
- **Vue 3** - 前端框架（Composition API + `<script setup>`）
- **Vite** - 构建工具
- **Pinia** - 状态管理
- **Tailwind CSS** - 样式框架
- **Font Awesome** - 图标库

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- npm 或 pnpm

### 开发环境

```bash
# 启动后端（在项目根目录）
CGO_ENABLED=1 go build -o nas-navigation .
./nas-navigation

# 启动前端（新终端）
cd frontend
npm install
npm run dev
```

- 后端服务: http://localhost:8080
- 前端开发服务: http://localhost:5173 (或下一个可用端口)

### 生产部署

```bash
# 构建前端
cd frontend
npm run build

# 构建后端
cd ..
CGO_ENABLED=1 go build -o nas-navigation .

# 运行（静态文件自动由后端提供）
./nas-navigation
```

访问 http://localhost:8080 查看应用。

## 项目结构

```
navgation/
├── backend/                 # Go 后端
│   ├── database/           # 数据库连接
│   ├── handlers/           # HTTP 处理器
│   │   ├── category.go     # 分类处理
│   │   ├── link.go         # 链接处理
│   │   ├── script.go       # 脚本处理
│   │   └── system.go       # 系统操作（重启等）
│   ├── models/             # 数据模型
│   │   ├── category.go
│   │   ├── link.go         # 包含 restart_script 字段
│   │   └── script.go
│   ├── router/             # 路由配置
│   └── services/           # 业务逻辑
├── frontend/               # Vue 前端
│   ├── src/
│   │   ├── api/           # API 调用封装
│   │   │   └── index.js
│   │   ├── components/    # Vue 组件
│   │   │   ├── LinkCard.vue       # 链接卡片（三点菜单）
│   │   │   ├── AddLinkModal.vue   # 添加/编辑链接
│   │   │   └── CustomSelect.vue   # 自定义下拉选择
│   │   ├── stores/        # Pinia 状态管理
│   │   │   └── app.js
│   │   └── App.vue        # 根组件
│   └── vite.config.js     # Vite 配置（含代理）
├── test_server/            # 测试服务器
│   ├── main.py            # FastAPI 测试服务
│   └── restart.sh         # 重启脚本示例
├── main.go                 # 后端入口
├── data.db                 # SQLite 数据库
├── TEST.md                 # 测试文档
└── README.md               # 项目文档
```

## API 接口

### 分类 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/categories` | 获取所有分类 |
| GET | `/api/categories/:id` | 获取单个分类 |
| POST | `/api/categories` | 创建分类 |
| PUT | `/api/categories/:id` | 更新分类 |
| DELETE | `/api/categories/:id` | 删除分类 |

### 链接 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/links` | 获取所有链接 |
| GET | `/api/links/search?q=xxx` | 搜索链接 |
| GET | `/api/links/:id` | 获取单个链接 |
| POST | `/api/links` | 创建链接 |
| PUT | `/api/links/:id` | 更新链接 |
| DELETE | `/api/links/:id` | 删除链接 |

### 系统 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/system/info` | 获取系统信息 |
| POST | `/api/system/restart` | 重启系统/服务 |
| POST | `/api/system/restart-link/:id` | 重启指定链接的服务 |
| POST | `/api/system/restart-all` | 重启所有服务 |

## 核心功能实现

### 1. 链接卡片三点菜单

鼠标悬停时显示三点按钮，点击展开下拉菜单：

```vue
<!-- frontend/src/components/LinkCard.vue -->
<template>
  <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100">
    <button @click.prevent="toggleMenu">
      <i class="fa-solid fa-ellipsis-vertical"></i>
    </button>
    <Transition>
      <div v-if="showMenu" class="dropdown-menu">
        <button @click="handleEdit">编辑</button>
        <button v-if="link.restart_script" @click="handleRestart">重启</button>
        <button @click="handleDelete">删除</button>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const showMenu = ref(false)
const menuRef = ref(null)

// 点击外部关闭菜单
function handleClickOutside(event) {
  if (menuRef.value && !menuRef.value.contains(event.target)) {
    showMenu.value = false
  }
}

onMounted(() => document.addEventListener('click', handleClickOutside))
onUnmounted(() => document.removeEventListener('click', handleClickOutside))
</script>
```

### 2. 服务重启机制

链接可以配置 `restart_script` 字段，存储重启脚本路径：

**数据模型 (backend/models/link.go)**
```go
type Link struct {
    ID            uint   `json:"id" gorm:"primaryKey"`
    CategoryID    uint   `json:"category_id"`
    Name          string `json:"name"`
    URL           string `json:"url"`
    RestartScript string `json:"restart_script"` // 重启脚本路径或命令
    // ...
}
```

**重启处理 (backend/handlers/system.go)**
```go
func RestartLink(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    link, err := linkService.GetByID(uint(id))

    if link.RestartScript == "" {
        c.JSON(400, gin.H{"error": "未配置重启脚本"})
        return
    }

    start := time.Now()
    cmd := exec.Command("bash", "-c", link.RestartScript)
    output, err := cmd.CombinedOutput()

    c.JSON(200, gin.H{
        "success":  err == nil,
        "output":   string(output),
        "duration": time.Since(start).Milliseconds(),
    })
}
```

**前端调用 (frontend/src/stores/app.js)**
```javascript
async function restartLink(id) {
  return await systemApi.restartLink(id)
}
```

### 3. 自定义下拉选择组件

替代原生 select，支持更好的样式和动画：

```vue
<!-- frontend/src/components/CustomSelect.vue -->
<template>
  <div class="relative">
    <button @click="toggleOpen" class="select-trigger">
      {{ selectedLabel }}
      <i class="fa-solid fa-chevron-down"></i>
    </button>
    <Transition>
      <div v-if="isOpen" class="absolute w-full bg-white rounded-xl shadow-lg">
        <div v-for="option in options"
             @click="selectOption(option)"
             class="px-4 py-2 hover:bg-indigo-50 cursor-pointer">
          {{ option.label }}
        </div>
      </div>
    </Transition>
  </div>
</template>
```

## 数据库结构

### categories 表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键 |
| name | TEXT | 分类名称 |
| icon | TEXT | 图标标识 |
| sort_order | INTEGER | 排序顺序 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### links 表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键 |
| category_id | INTEGER | 所属分类 ID |
| name | TEXT | 链接名称 |
| url | TEXT | 链接地址 |
| icon | TEXT | 图标 URL |
| description | TEXT | 描述信息 |
| restart_script | TEXT | 重启脚本路径或命令 |
| sort_order | INTEGER | 排序顺序 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

## 环境变量

| 变量 | 默认值 | 描述 |
|------|--------|------|
| PORT | 8080 | 服务端口 |
| DB_PATH | ./data.db | 数据库路径 |
| DEBUG | true | 调试模式 |

## 测试

详细测试说明请参考 [TEST.md](./TEST.md)。

### 快速验证

```bash
# 启动测试服务器
cd test_server
bash restart.sh

# 验证服务运行
curl http://localhost:8877/

# 返回示例
{
  "message": "Hello from Test Server!",
  "start_time": "2025-12-08T00:00:12.657250",
  "uptime_seconds": 13.079728
}
```

### 测试重启功能

1. 添加链接时配置 `restart_script` 为测试脚本路径
2. 点击链接卡片的三点菜单
3. 点击"重启"按钮
4. 观察服务器 uptime 是否重置

## 开发说明

### 添加新 API

1. `backend/handlers/` - 创建处理函数
2. `backend/router/router.go` - 注册路由
3. `frontend/src/api/index.js` - 添加 API 调用
4. `frontend/src/stores/app.js` - 添加状态方法

### 添加新组件

1. `frontend/src/components/` - 创建 `.vue` 文件
2. 使用 Composition API (`<script setup>`)
3. 使用 Tailwind CSS 进行样式设计
4. 使用 Vue Transition 添加动画效果

## License

MIT
