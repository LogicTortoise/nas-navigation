<template>
  <!-- 背景遮罩 -->
  <div class="fixed inset-0 bg-gradient-to-br from-indigo-100/30 to-purple-100/30 pointer-events-none z-0"></div>

  <!-- 主应用容器 -->
  <div class="relative z-10 w-full h-screen max-w-[1920px] mx-auto flex flex-col md:flex-row overflow-hidden md:rounded-3xl md:m-4 md:h-[calc(100vh-2rem)] shadow-2xl glass-panel">

    <!-- 左侧侧边栏 (分类导航) -->
    <aside class="hidden md:flex flex-col w-20 lg:w-64 h-full sidebar-border transition-all duration-300">
      <!-- Logo 区域 -->
      <div class="h-[80px] flex items-center justify-center lg:justify-start px-0 lg:px-8 border-b border-white/20">
        <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-indigo-600 to-purple-500 flex items-center justify-center text-white shadow-lg">
          <i class="fa-solid fa-compass text-xl"></i>
        </div>
        <span class="ml-3 text-xl font-bold tracking-tight text-slate-800 hidden lg:block">酷站导航</span>
      </div>

      <!-- 分类列表 -->
      <nav class="flex-1 overflow-y-auto py-6 px-3 space-y-2 hide-scrollbar">
        <button
          v-for="cat in store.categories"
          :key="cat.id"
          @click="currentCategory = cat.id"
          :class="[
            'w-full flex items-center justify-center lg:justify-start px-3 py-3 rounded-xl transition-all group',
            currentCategory === cat.id ? 'nav-item-active' : 'nav-item-inactive text-slate-600'
          ]"
        >
          <i :class="['text-lg w-6 text-center', getCategoryIconClass(cat.icon)]"></i>
          <span class="ml-3 font-medium hidden lg:block">{{ cat.name }}</span>
        </button>

        <div class="my-4 border-t border-white/20 w-full"></div>

        <button
          @click="openAddCategoryModal"
          class="w-full flex items-center justify-center lg:justify-start px-3 py-3 rounded-xl text-slate-500 hover:bg-slate-200/50 transition-all"
        >
          <i class="fa-solid fa-plus text-lg w-6 text-center"></i>
          <span class="ml-3 font-medium text-sm hidden lg:block">添加分类</span>
        </button>
      </nav>

      <!-- 底部用户信息和快捷操作 -->
      <div class="p-4 border-t border-white/20 space-y-3">
        <!-- 一键重启按钮 -->
        <button
          @click="handleRestartSystem"
          :disabled="isRestarting"
          class="w-full flex items-center justify-center lg:justify-start px-3 py-2.5 rounded-xl text-orange-600 hover:bg-orange-50 transition-all group"
        >
          <i :class="['fa-solid fa-rotate text-lg w-6 text-center', isRestarting ? 'animate-spin' : '']"></i>
          <span class="ml-3 font-medium text-sm hidden lg:block">{{ isRestarting ? '重启中...' : '一键重启' }}</span>
        </button>

        <!-- 用户信息 -->
        <button class="flex items-center justify-center lg:justify-start p-2 rounded-xl hover:bg-white/40 transition-colors cursor-pointer w-full">
          <img src="https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?ixlib=rb-4.0.3&auto=format&fit=crop&w=100&q=80"
               alt="用户头像"
               class="w-10 h-10 rounded-full border-2 border-white shadow-sm object-cover">
          <div class="ml-3 hidden lg:block overflow-hidden">
            <p class="text-sm font-bold text-slate-800 truncate">{{ store.userName }}</p>
            <p class="text-xs text-slate-500 truncate">点击进入设置</p>
          </div>
        </button>
      </div>
    </aside>

    <!-- 右侧主内容区 -->
    <main class="flex-1 flex flex-col h-full overflow-hidden relative">

      <!-- 顶部导航栏 -->
      <header class="h-[80px] flex items-center justify-between px-6 md:px-10 py-4 shrink-0 z-20">
        <!-- 移动端汉堡菜单 -->
        <button class="md:hidden p-2 text-slate-600 hover:text-indigo-600">
          <i class="fa-solid fa-bars text-2xl"></i>
        </button>

        <!-- 移动端 Logo -->
        <span class="md:hidden text-lg font-bold text-slate-800">酷站导航</span>

        <!-- 搜索框 -->
        <div class="hidden md:flex flex-1 max-w-2xl mx-auto relative group search-focus rounded-full transition-all duration-300 bg-white/60 backdrop-blur-md border border-white/50">
          <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
            <i class="fa-solid fa-search text-slate-400 group-focus-within:text-indigo-600 transition-colors"></i>
          </div>
          <input
            type="text"
            v-model="searchQuery"
            @input="handleSearch"
            class="block w-full pl-12 pr-4 py-3 bg-transparent border-none rounded-full text-slate-800 placeholder-slate-400 focus:ring-0 sm:text-sm outline-none"
            placeholder="搜索书签、Google 或输入网址...">
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center">
            <span class="text-xs text-slate-400 border border-slate-300 rounded px-1.5 py-0.5">⌘K</span>
          </div>
        </div>

        <!-- 右侧操作区 -->
        <div class="flex items-center space-x-4 ml-4">
          <!-- 时间显示 -->
          <div class="hidden xl:flex flex-col items-end mr-4 text-right">
            <span class="text-lg font-bold text-slate-700 leading-none">{{ currentTime }}</span>
            <span class="text-xs text-slate-500">{{ currentDate }}</span>
          </div>

          <!-- 添加按钮 -->
          <button
            @click="openAddLinkModal"
            class="flex items-center justify-center px-4 py-2.5 bg-slate-800 hover:bg-slate-700 text-white rounded-full shadow-lg hover:shadow-xl transition-all transform hover:-translate-y-0.5"
          >
            <i class="fa-solid fa-plus mr-2"></i>
            <span class="text-sm font-medium">添加</span>
          </button>
        </div>
      </header>

      <!-- 内容滚动区 -->
      <div class="flex-1 overflow-y-auto p-6 md:p-10 hide-scrollbar">

        <!-- 欢迎语 -->
        <div class="mb-8">
          <h1 class="text-2xl md:text-3xl font-bold text-slate-800 mb-2">
            {{ store.timeGreeting }}，{{ store.userName }}
            <span v-if="store.timeGreeting === '早安'">☀️</span>
            <span v-else-if="store.timeGreeting === '晚上好'">🌙</span>
          </h1>
          <p class="text-slate-500">{{ store.greeting }}</p>
        </div>

        <!-- 搜索结果 -->
        <section v-if="searchQuery && searchResults.length > 0" class="mb-10">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold text-slate-700">搜索结果</h2>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            <LinkCard
              v-for="link in searchResults"
              :key="link.id"
              :link="link"
              :card-type="'horizontal'"
              @edit="editLink"
              @delete="handleDeleteLink"
              @execute="handleExecuteScript"
            />
          </div>
        </section>

        <!-- 常用访问区域 -->
        <section v-else-if="!currentCategory || currentCategory === 'all'" class="mb-10">
          <!-- 所有分类 -->
          <div class="mb-10" v-for="cat in store.categoriesWithLinks" :key="cat.id">
            <div class="flex items-center justify-between mb-4">
              <h2 class="text-lg font-semibold text-slate-700 flex items-center">
                <i :class="['mr-2', getCategoryIconClass(cat.icon), getIconColor(cat.icon)]"></i>
                {{ cat.name }}
              </h2>
              <button
                @click="editCategory(cat)"
                class="text-xs text-slate-400 hover:text-indigo-600 transition-colors"
              >编辑</button>
            </div>

            <!-- 统一使用横向卡片样式 -->
            <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              <LinkCard
                v-for="link in cat.links"
                :key="link.id"
                :link="link"
                :card-type="'horizontal'"
                @edit="editLink"
                @delete="handleDeleteLink"
                @execute="handleExecuteScript"
              />
              <AddLinkButton :horizontal="true" @click="addLinkToCategory(cat.id)" />
            </div>
          </div>
        </section>

        <!-- 特定分类内容 -->
        <section v-else class="mb-10">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold text-slate-700">{{ currentCategoryData?.name }}</h2>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
            <LinkCard
              v-for="link in currentCategoryLinks"
              :key="link.id"
              :link="link"
              :card-type="'horizontal'"
              @edit="editLink"
              @delete="handleDeleteLink"
              @execute="handleExecuteScript"
            />
            <AddLinkButton :horizontal="true" @click="addLinkToCategory(currentCategory)" />
          </div>

          <!-- 无数据提示 -->
          <div v-if="currentCategoryLinks.length === 0" class="text-center py-16">
            <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-slate-100 flex items-center justify-center">
              <i class="fa-solid fa-link text-2xl text-slate-400"></i>
            </div>
            <h3 class="text-lg font-medium text-slate-700 mb-2">暂无网站链接</h3>
            <p class="text-slate-500 mb-6">开始添加您的第一个网站链接吧</p>
            <button
              @click="addLinkToCategory(currentCategory)"
              class="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
            >
              <i class="fa-solid fa-plus mr-2"></i>添加链接
            </button>
          </div>
        </section>
      </div>
    </main>
  </div>

  <!-- 添加链接弹窗 -->
  <AddLinkModal
    v-if="showAddModal"
    :link="currentEditLink"
    :category-id="currentCategoryId"
    @close="closeAddModal"
    @save="saveLink"
  />

  <!-- 分类管理弹窗 -->
  <CategoryModal
    v-if="showCategoryModal"
    :category="currentEditCategory"
    @close="closeCategoryModal"
    @save="saveCategory"
    @delete="handleDeleteCategory"
  />

  <!-- 脚本执行结果弹窗 -->
  <ScriptResultModal
    v-if="scriptResult"
    :result="scriptResult"
    @close="scriptResult = null"
  />

  <!-- 加载状态 -->
  <div v-if="store.loading" class="fixed inset-0 bg-white/80 backdrop-blur-sm flex items-center justify-center z-50">
    <div class="animate-spin rounded-full h-12 w-12 border-4 border-indigo-600 border-t-transparent"></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAppStore } from './stores/app'
import { useDebounceFn } from '@vueuse/core'
import LinkCard from './components/LinkCard.vue'
import AddLinkButton from './components/AddLinkButton.vue'
import AddLinkModal from './components/AddLinkModal.vue'
import CategoryModal from './components/CategoryModal.vue'
import ScriptResultModal from './components/ScriptResultModal.vue'

const store = useAppStore()

// 状态
const showAddModal = ref(false)
const showCategoryModal = ref(false)
const currentEditLink = ref(null)
const currentEditCategory = ref(null)
const currentCategoryId = ref(null)
const currentCategory = ref(null)
const scriptResult = ref(null)
const searchQuery = ref('')
const searchResults = ref([])
const isRestarting = ref(false)

// 时间显示
const currentTime = ref('')
const currentDate = ref('')

function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', hour12: false })
  currentDate.value = now.toLocaleDateString('zh-CN', { month: 'long', day: 'numeric', weekday: 'short' })
}

// 计算属性
const currentCategoryData = computed(() => {
  return store.categories.find(c => c.id === currentCategory.value)
})

const currentCategoryLinks = computed(() => {
  return store.links.filter(l => l.category_id === currentCategory.value)
})

// 搜索
const debouncedSearch = useDebounceFn(async (query) => {
  if (!query.trim()) {
    searchResults.value = []
    return
  }
  searchResults.value = await store.searchLinks(query)
}, 300)

function handleSearch() {
  debouncedSearch(searchQuery.value)
}

// 分类图标
function getCategoryIconClass(icon) {
  // 如果已经是 fa- 开头的图标名，直接返回
  if (icon && icon.startsWith('fa-')) {
    return 'fa-solid ' + icon
  }
  // 兼容旧的图标名称
  const icons = {
    'star': 'fa-solid fa-star',
    'briefcase': 'fa-solid fa-briefcase',
    'palette': 'fa-solid fa-pen-nib',
    'code': 'fa-solid fa-code',
    'music': 'fa-solid fa-music',
    'video': 'fa-solid fa-video',
    'book': 'fa-solid fa-graduation-cap',
    'game': 'fa-solid fa-gamepad',
    'cloud': 'fa-solid fa-cloud',
    'home': 'fa-solid fa-home',
    'folder': 'fa-solid fa-folder',
    'layer-group': 'fa-solid fa-layer-group'
  }
  return icons[icon] || 'fa-solid fa-folder'
}

function getIconColor(icon) {
  // 处理新的 fa- 开头图标名
  const iconName = icon && icon.startsWith('fa-') ? icon.replace('fa-', '') : icon
  const colors = {
    'star': 'text-yellow-400',
    'briefcase': 'text-blue-500',
    'palette': 'text-pink-500',
    'pen-nib': 'text-pink-500',
    'code': 'text-green-500',
    'game': 'text-purple-500',
    'gamepad': 'text-purple-500',
    'folder': 'text-yellow-500',
    'graduation-cap': 'text-orange-500',
    'heart': 'text-red-500',
    'music': 'text-indigo-500',
    'video': 'text-red-400'
  }
  return colors[iconName] || 'text-slate-500'
}

// 链接操作
function openAddLinkModal() {
  currentEditLink.value = null
  currentCategoryId.value = store.categories.length > 0 ? store.categories[0].id : null
  showAddModal.value = true
}

function addLinkToCategory(categoryId) {
  currentEditLink.value = null
  currentCategoryId.value = categoryId
  showAddModal.value = true
}

function editLink(link) {
  currentEditLink.value = link
  currentCategoryId.value = link.category_id
  showAddModal.value = true
}

function closeAddModal() {
  showAddModal.value = false
  currentEditLink.value = null
  currentCategoryId.value = null
}

async function saveLink(data) {
  try {
    if (currentEditLink.value) {
      await store.updateLink(currentEditLink.value.id, data)
    } else {
      await store.addLink(data)
    }
    closeAddModal()
  } catch (error) {
    alert('保存失败: ' + error.message)
  }
}

async function handleDeleteLink(link) {
  if (confirm(`确定删除"${link.name}"吗?`)) {
    try {
      await store.deleteLink(link.id)
    } catch (error) {
      alert('删除失败: ' + error.message)
    }
  }
}

// 分类操作
function openAddCategoryModal() {
  currentEditCategory.value = null
  showCategoryModal.value = true
}

function editCategory(category) {
  currentEditCategory.value = category
  showCategoryModal.value = true
}

function closeCategoryModal() {
  showCategoryModal.value = false
  currentEditCategory.value = null
}

async function saveCategory(data) {
  try {
    if (currentEditCategory.value) {
      await store.updateCategory(currentEditCategory.value.id, data)
    } else {
      await store.addCategory(data)
    }
    closeCategoryModal()
  } catch (error) {
    alert('保存失败: ' + error.message)
  }
}

async function handleDeleteCategory(category) {
  if (confirm(`确定删除分类"${category.name}"及其下所有链接吗?`)) {
    try {
      await store.deleteCategory(category.id)
      closeCategoryModal()
    } catch (error) {
      alert('删除失败: ' + error.message)
    }
  }
}

// 脚本执行
async function handleExecuteScript(link) {
  if (!link.script_id) return
  try {
    const result = await store.executeScript(link.script_id)
    scriptResult.value = result
  } catch (error) {
    scriptResult.value = {
      success: false,
      error: error.message,
      output: ''
    }
  }
}

// 一键重启所有服务
async function handleRestartSystem() {
  if (!confirm('确定要重启所有配置了重启脚本的服务吗？')) {
    return
  }

  isRestarting.value = true
  try {
    const result = await store.restartAllServices()

    // 构建输出信息
    let output = result.message + '\n\n'
    if (result.results && result.results.length > 0) {
      output += '执行详情:\n'
      result.results.forEach(r => {
        const status = r.success ? '✓' : '✗'
        output += `${status} ${r.service_name}: ${r.success ? (r.output || '成功') : r.error}\n`
      })
    }

    scriptResult.value = {
      success: result.failed_count === 0,
      output: output,
      error: result.failed_count > 0 ? `${result.failed_count} 个服务重启失败` : ''
    }
  } catch (error) {
    scriptResult.value = {
      success: false,
      error: error.response?.data?.error || error.message,
      output: ''
    }
  } finally {
    isRestarting.value = false
  }
}

onMounted(() => {
  store.loadData()
  updateTime()
  setInterval(updateTime, 60000)
})
</script>
