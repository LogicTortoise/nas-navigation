import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { categoryApi, linkApi, scriptApi, settingApi, systemApi } from '../api'

export const useAppStore = defineStore('app', () => {
  // 状态
  const categories = ref([])
  const links = ref([])
  const scripts = ref([])
  const settings = ref({})
  const loading = ref(false)
  const searchQuery = ref('')
  const searchResults = ref([])
  const editMode = ref(false)
  const showAddModal = ref(false)
  const showScriptModal = ref(false)
  const showCategoryModal = ref(false)
  const currentCategory = ref(null)
  const currentLink = ref(null)

  // 计算属性
  const userName = computed(() => settings.value.user_name || '用户')
  const greeting = computed(() => settings.value.greeting || '准备好开始高效的一天了吗?')

  const categoriesWithLinks = computed(() => {
    return categories.value.map(cat => ({
      ...cat,
      links: links.value.filter(link => link.category_id === cat.id)
    }))
  })

  // 获取当前时间问候语
  const timeGreeting = computed(() => {
    const hour = new Date().getHours()
    if (hour < 6) return '夜深了'
    if (hour < 9) return '早安'
    if (hour < 12) return '上午好'
    if (hour < 14) return '中午好'
    if (hour < 18) return '下午好'
    if (hour < 22) return '晚上好'
    return '夜深了'
  })

  // 方法
  async function loadData() {
    loading.value = true
    try {
      const [cats, lks, scrs, sets] = await Promise.all([
        categoryApi.getAll(),
        linkApi.getAll(),
        scriptApi.getAll(),
        settingApi.getAll()
      ])
      categories.value = cats
      links.value = lks
      scripts.value = scrs
      settings.value = sets
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      loading.value = false
    }
  }

  async function searchLinks(query) {
    if (!query.trim()) {
      searchResults.value = []
      return
    }
    try {
      searchResults.value = await linkApi.search(query)
    } catch (error) {
      console.error('Search failed:', error)
      searchResults.value = []
    }
  }

  // 分类操作
  async function addCategory(data) {
    const newCat = await categoryApi.create(data)
    categories.value.push(newCat)
    return newCat
  }

  async function updateCategory(id, data) {
    const updated = await categoryApi.update(id, data)
    const index = categories.value.findIndex(c => c.id === id)
    if (index !== -1) {
      categories.value[index] = updated
    }
    return updated
  }

  async function deleteCategory(id) {
    await categoryApi.delete(id)
    categories.value = categories.value.filter(c => c.id !== id)
    links.value = links.value.filter(l => l.category_id !== id)
  }

  // 链接操作
  async function addLink(data) {
    const newLink = await linkApi.create(data)
    links.value.push(newLink)
    return newLink
  }

  async function updateLink(id, data) {
    const updated = await linkApi.update(id, data)
    const index = links.value.findIndex(l => l.id === id)
    if (index !== -1) {
      links.value[index] = updated
    }
    return updated
  }

  async function deleteLink(id) {
    await linkApi.delete(id)
    links.value = links.value.filter(l => l.id !== id)
  }

  // 脚本操作
  async function addScript(data) {
    const newScript = await scriptApi.create(data)
    scripts.value.push(newScript)
    return newScript
  }

  async function updateScript(id, data) {
    const updated = await scriptApi.update(id, data)
    const index = scripts.value.findIndex(s => s.id === id)
    if (index !== -1) {
      scripts.value[index] = updated
    }
    return updated
  }

  async function deleteScript(id) {
    await scriptApi.delete(id)
    scripts.value = scripts.value.filter(s => s.id !== id)
  }

  async function executeScript(id) {
    return await scriptApi.execute(id)
  }

  // 设置操作
  async function updateSettings(data) {
    await settingApi.update(data)
    settings.value = { ...settings.value, ...data }
  }

  // 系统操作
  async function getSystemInfo() {
    return await systemApi.getInfo()
  }

  async function restartSystem(target = '') {
    return await systemApi.restart(target, true)
  }

  async function restartAllServices() {
    return await systemApi.restartAll()
  }

  async function restartLink(id) {
    return await systemApi.restartLink(id)
  }

  async function executeQuickAction(action) {
    return await systemApi.quickAction(action)
  }

  return {
    // 状态
    categories,
    links,
    scripts,
    settings,
    loading,
    searchQuery,
    searchResults,
    editMode,
    showAddModal,
    showScriptModal,
    showCategoryModal,
    currentCategory,
    currentLink,
    // 计算属性
    userName,
    greeting,
    categoriesWithLinks,
    timeGreeting,
    // 方法
    loadData,
    searchLinks,
    addCategory,
    updateCategory,
    deleteCategory,
    addLink,
    updateLink,
    deleteLink,
    addScript,
    updateScript,
    deleteScript,
    executeScript,
    updateSettings,
    getSystemInfo,
    restartSystem,
    restartAllServices,
    restartLink,
    executeQuickAction
  }
})
