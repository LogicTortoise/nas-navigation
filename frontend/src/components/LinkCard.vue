<template>
  <!-- 方形卡片 -->
  <a
    v-if="cardType === 'square'"
    :href="link.url"
    target="_blank"
    class="glass-card group relative flex flex-col items-center justify-center p-6 rounded-2xl h-32 md:h-40 cursor-pointer"
  >
    <!-- 三点菜单按钮 -->
    <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity" ref="menuRef">
      <button
        @click.prevent="toggleMenu"
        class="text-slate-400 hover:text-slate-600 p-1.5 rounded-lg hover:bg-white/80 transition-colors"
      >
        <i class="fa-solid fa-ellipsis-vertical text-sm"></i>
      </button>

      <!-- 下拉菜单 -->
      <Transition
        enter-active-class="transition duration-150 ease-out"
        enter-from-class="opacity-0 scale-95"
        enter-to-class="opacity-100 scale-100"
        leave-active-class="transition duration-100 ease-in"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-95"
      >
        <div
          v-if="showMenu"
          class="absolute right-0 top-full mt-1 w-32 bg-white rounded-xl shadow-lg border border-slate-200 overflow-hidden z-50"
        >
          <button
            @click.prevent="handleEdit"
            class="w-full px-4 py-2.5 text-left text-sm text-slate-700 hover:bg-indigo-50 hover:text-indigo-600 transition-colors flex items-center gap-2"
          >
            <i class="fa-solid fa-edit text-xs"></i>
            编辑
          </button>
          <button
            v-if="link.restart_script"
            @click.prevent="handleRestart"
            class="w-full px-4 py-2.5 text-left text-sm text-slate-700 hover:bg-green-50 hover:text-green-600 transition-colors flex items-center gap-2"
          >
            <i class="fa-solid fa-rotate text-xs"></i>
            重启
          </button>
          <button
            @click.prevent="handleDelete"
            class="w-full px-4 py-2.5 text-left text-sm text-red-600 hover:bg-red-50 transition-colors flex items-center gap-2"
          >
            <i class="fa-solid fa-trash text-xs"></i>
            删除
          </button>
        </div>
      </Transition>
    </div>

    <!-- 图标 -->
    <div class="w-12 h-12 md:w-14 md:h-14 rounded-2xl bg-white shadow-md flex items-center justify-center mb-3 group-hover:scale-110 transition-transform duration-300">
      <img
        :src="getFaviconUrl()"
        :alt="link.name"
        class="w-8 h-8 object-contain"
        @error="handleIconError"
        v-if="!iconError"
      />
      <span v-else class="text-xl font-bold" :style="{ color: getIconColor() }">
        {{ link.name.charAt(0).toUpperCase() }}
      </span>
    </div>

    <!-- 名称 -->
    <span class="text-sm font-medium text-slate-700 group-hover:text-indigo-600 transition-colors truncate max-w-full">{{ link.name }}</span>
  </a>

  <!-- 横向卡片 -->
  <a
    v-else
    :href="link.url"
    target="_blank"
    class="glass-card flex items-center p-4 rounded-xl hover:bg-white transition-colors group"
  >
    <!-- 图标 -->
    <div :class="['w-10 h-10 rounded-lg flex items-center justify-center shrink-0', getIconBgClass()]">
      <img
        :src="getFaviconUrl()"
        :alt="link.name"
        class="w-6 h-6 object-contain"
        @error="handleIconError"
        v-if="!iconError"
      />
      <span v-else class="text-sm font-bold" :style="{ color: getIconColor() }">
        {{ link.name.charAt(0).toUpperCase() }}
      </span>
    </div>

    <!-- 文字信息 -->
    <div class="ml-3 overflow-hidden flex-1">
      <h3 class="text-sm font-semibold text-slate-800 truncate group-hover:text-indigo-600">{{ link.name }}</h3>
      <p class="text-xs text-slate-500 truncate">{{ link.description || link.url }}</p>
    </div>

    <!-- 三点菜单按钮 -->
    <div class="ml-2 opacity-0 group-hover:opacity-100 transition-opacity relative" ref="menuRefHorizontal">
      <button
        @click.prevent="toggleMenu"
        class="text-slate-400 hover:text-slate-600 p-1.5 rounded-lg hover:bg-white/80 transition-colors"
      >
        <i class="fa-solid fa-ellipsis-vertical text-sm"></i>
      </button>

      <!-- 下拉菜单 -->
      <Transition
        enter-active-class="transition duration-150 ease-out"
        enter-from-class="opacity-0 scale-95"
        enter-to-class="opacity-100 scale-100"
        leave-active-class="transition duration-100 ease-in"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-95"
      >
        <div
          v-if="showMenu"
          class="absolute right-0 top-full mt-1 w-32 bg-white rounded-xl shadow-lg border border-slate-200 overflow-hidden z-50"
        >
          <button
            @click.prevent="handleEdit"
            class="w-full px-4 py-2.5 text-left text-sm text-slate-700 hover:bg-indigo-50 hover:text-indigo-600 transition-colors flex items-center gap-2"
          >
            <i class="fa-solid fa-edit text-xs"></i>
            编辑
          </button>
          <button
            v-if="link.restart_script"
            @click.prevent="handleRestart"
            class="w-full px-4 py-2.5 text-left text-sm text-slate-700 hover:bg-green-50 hover:text-green-600 transition-colors flex items-center gap-2"
          >
            <i class="fa-solid fa-rotate text-xs"></i>
            重启
          </button>
          <button
            @click.prevent="handleDelete"
            class="w-full px-4 py-2.5 text-left text-sm text-red-600 hover:bg-red-50 transition-colors flex items-center gap-2"
          >
            <i class="fa-solid fa-trash text-xs"></i>
            删除
          </button>
        </div>
      </Transition>
    </div>
  </a>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  link: {
    type: Object,
    required: true
  },
  cardType: {
    type: String,
    default: 'square' // 'square' | 'horizontal'
  }
})

const emit = defineEmits(['edit', 'delete', 'restart'])

const iconError = ref(false)
const showMenu = ref(false)
const menuRef = ref(null)
const menuRefHorizontal = ref(null)

function toggleMenu() {
  showMenu.value = !showMenu.value
}

function handleEdit() {
  showMenu.value = false
  emit('edit', props.link)
}

function handleDelete() {
  showMenu.value = false
  emit('delete', props.link)
}

function handleRestart() {
  showMenu.value = false
  emit('restart', props.link)
}

function handleClickOutside(event) {
  const target = event.target
  if (menuRef.value && !menuRef.value.contains(target) &&
      menuRefHorizontal.value && !menuRefHorizontal.value.contains(target)) {
    showMenu.value = false
  }
}

function handleIconError() {
  iconError.value = true
}

function getFaviconUrl() {
  if (props.link.icon && props.link.icon.startsWith('http')) {
    return props.link.icon
  }
  try {
    const url = new URL(props.link.url)
    return `https://www.google.com/s2/favicons?domain=${url.hostname}&sz=64`
  } catch {
    return ''
  }
}

function getIconBgClass() {
  const colors = [
    'bg-blue-100', 'bg-green-100', 'bg-purple-100',
    'bg-pink-100', 'bg-yellow-100', 'bg-red-100',
    'bg-indigo-100', 'bg-cyan-100', 'bg-orange-100'
  ]
  const index = props.link.name.charCodeAt(0) % colors.length
  return colors[index]
}

function getIconColor() {
  const colors = [
    '#3b82f6', '#22c55e', '#a855f7',
    '#ec4899', '#eab308', '#ef4444',
    '#6366f1', '#06b6d4', '#f97316'
  ]
  const index = props.link.name.charCodeAt(0) % colors.length
  return colors[index]
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
