<template>
  <!-- 模态弹窗背景 -->
  <div
    class="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    @click.self="$emit('close')"
  >
    <!-- 添加链接弹窗 -->
    <div class="glass-panel rounded-2xl shadow-2xl w-full max-w-md lg:max-w-lg p-6 lg:p-8 modal-enter">

      <!-- 弹窗头部 -->
      <div class="flex justify-between items-center mb-6 lg:mb-8">
        <div class="flex items-center">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-indigo-600 to-purple-500 flex items-center justify-center text-white shadow-lg mr-3">
            <i class="fa-solid fa-plus text-lg"></i>
          </div>
          <h3 class="text-xl lg:text-2xl font-bold text-slate-800">
            {{ link ? '编辑链接' : '添加网站链接' }}
          </h3>
        </div>
        <button
          @click="$emit('close')"
          class="text-slate-400 hover:text-slate-600 transition-colors p-1 rounded-lg hover:bg-white/50"
        >
          <i class="fa-solid fa-xmark text-xl"></i>
        </button>
      </div>

      <!-- 表单内容 -->
      <form @submit.prevent="handleSubmit" class="space-y-6">

        <!-- 网站URL -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-slate-700">
            网站 URL <span class="text-rose-500">*</span>
          </label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <i class="fa-solid fa-link text-slate-400"></i>
            </div>
            <input
              type="url"
              v-model="form.url"
              @blur="fetchFavicon"
              placeholder="https://example.com"
              class="w-full pl-12 pr-4 py-3 bg-white/80 border border-slate-300 rounded-xl form-input-focus transition-all text-slate-800 placeholder-slate-400"
              required
            >
          </div>
          <p class="text-xs text-slate-500">
            输入网址后将自动获取网站图标
          </p>
        </div>

        <!-- 网站名称 -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-slate-700">
            网站名称 <span class="text-rose-500">*</span>
          </label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <i class="fa-solid fa-tag text-slate-400"></i>
            </div>
            <input
              type="text"
              v-model="form.name"
              placeholder="例如：Google、GitHub"
              class="w-full pl-12 pr-4 py-3 bg-white/80 border border-slate-300 rounded-xl form-input-focus transition-all text-slate-800 placeholder-slate-400"
              required
            >
          </div>
        </div>

        <!-- 网站图标预览 -->
        <div v-if="faviconUrl" class="space-y-2">
          <label class="block text-sm font-medium text-slate-700">网站图标</label>
          <div class="flex items-center space-x-4">
            <div class="w-12 h-12 rounded-lg bg-white shadow-md flex items-center justify-center">
              <img
                :src="faviconUrl"
                class="w-8 h-8 object-contain"
                @error="faviconError = true"
                v-if="!faviconError"
              />
              <i v-else class="fa-solid fa-globe text-slate-400 text-xl"></i>
            </div>
            <div class="flex-1">
              <p class="text-sm text-slate-600">{{ faviconError ? '无法获取图标' : '已自动获取网站图标' }}</p>
              <button
                type="button"
                @click="fetchFavicon"
                class="text-xs text-indigo-600 hover:text-indigo-700 transition-colors mt-1"
              >
                <i class="fa-solid fa-refresh mr-1"></i>重新获取
              </button>
            </div>
          </div>
        </div>

        <!-- 所属分类 -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-slate-700">
            所属分类 <span class="text-rose-500">*</span>
          </label>
          <CustomSelect
            v-model="form.category_id"
            :options="categoryOptions"
            placeholder="请选择分类"
            icon="fa-solid fa-folder"
          />
        </div>

        <!-- 描述（可选） -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-slate-700">
            描述 <span class="text-slate-400 text-xs">(可选)</span>
          </label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <i class="fa-solid fa-align-left text-slate-400"></i>
            </div>
            <input
              type="text"
              v-model="form.description"
              placeholder="简短描述..."
              class="w-full pl-12 pr-4 py-3 bg-white/80 border border-slate-300 rounded-xl form-input-focus transition-all text-slate-800 placeholder-slate-400"
            >
          </div>
        </div>

        <!-- 关联脚本（可选） -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-slate-700">
            关联脚本 <span class="text-slate-400 text-xs">(可选)</span>
          </label>
          <CustomSelect
            v-model="form.script_id"
            :options="scriptOptions"
            placeholder="选择脚本"
            icon="fa-solid fa-code"
            :allow-empty="true"
            empty-label="无"
          />
        </div>

        <!-- 表单底部操作按钮 -->
        <div class="pt-4 flex space-x-3 lg:space-x-4">
          <button
            type="button"
            @click="$emit('close')"
            class="flex-1 px-4 py-3 rounded-xl border border-slate-300 text-slate-700 hover:bg-slate-50 font-medium transition-all"
          >
            <i class="fa-solid fa-times mr-2"></i>取消
          </button>
          <button
            type="submit"
            class="flex-1 px-4 py-3 rounded-xl bg-indigo-600 text-white hover:bg-indigo-700 font-medium shadow-lg shadow-indigo-500/30 transition-all"
          >
            <i class="fa-solid fa-save mr-2"></i>保存链接
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import { useAppStore } from '../stores/app'
import CustomSelect from './CustomSelect.vue'

const props = defineProps({
  link: Object,
  categoryId: Number
})

const emit = defineEmits(['close', 'save'])
const store = useAppStore()

// 分类选项
const categoryOptions = computed(() => {
  return store.categories.map(cat => ({
    value: cat.id,
    label: cat.name,
    icon: getCategoryIcon(cat.icon),
    color: getCategoryColor(cat.icon)
  }))
})

// 脚本选项
const scriptOptions = computed(() => {
  return store.scripts.map(script => ({
    value: script.id,
    label: script.name,
    icon: 'fa-solid fa-terminal',
    color: '#10b981'
  }))
})

function getCategoryIcon(icon) {
  if (icon && icon.startsWith('fa-')) {
    return 'fa-solid ' + icon
  }
  const icons = {
    'star': 'fa-solid fa-star',
    'briefcase': 'fa-solid fa-briefcase',
    'palette': 'fa-solid fa-pen-nib',
    'folder': 'fa-solid fa-folder'
  }
  return icons[icon] || 'fa-solid fa-folder'
}

function getCategoryColor(icon) {
  const iconName = icon && icon.startsWith('fa-') ? icon.replace('fa-', '') : icon
  const colors = {
    'star': '#eab308',
    'briefcase': '#3b82f6',
    'palette': '#ec4899',
    'pen-nib': '#ec4899',
    'folder': '#f59e0b'
  }
  return colors[iconName] || '#64748b'
}

const form = reactive({
  name: '',
  url: '',
  category_id: null,
  icon: '',
  description: '',
  script_id: null,
  sort_order: 0
})

const faviconUrl = ref('')
const faviconError = ref(false)

function fetchFavicon() {
  if (!form.url) return

  try {
    const url = new URL(form.url)
    faviconUrl.value = `https://www.google.com/s2/favicons?domain=${url.hostname}&sz=64`
    faviconError.value = false
    // 也设置到 form.icon
    form.icon = faviconUrl.value
  } catch {
    faviconUrl.value = ''
    faviconError.value = true
  }
}

function handleSubmit() {
  emit('save', { ...form })
}

onMounted(() => {
  if (props.link) {
    Object.assign(form, {
      name: props.link.name,
      url: props.link.url,
      category_id: props.link.category_id,
      icon: props.link.icon || '',
      description: props.link.description || '',
      script_id: props.link.script_id,
      sort_order: props.link.sort_order || 0
    })
    if (props.link.url) {
      fetchFavicon()
    }
  } else if (props.categoryId) {
    form.category_id = props.categoryId
  } else if (store.categories.length > 0) {
    form.category_id = store.categories[0].id
  }
})
</script>
