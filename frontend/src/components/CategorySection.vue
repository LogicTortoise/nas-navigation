<template>
  <div class="mb-8">
    <!-- 分类标题 -->
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-2">
        <span class="text-xl" v-html="getCategoryIcon(category.icon)"></span>
        <h2 class="text-lg font-semibold text-gray-700">{{ category.name }}</h2>
      </div>
      <button
        @click="$emit('edit-category', category)"
        class="text-sm text-gray-400 hover:text-primary-500 transition"
      >
        编辑
      </button>
    </div>

    <!-- 横向布局（工作效率类） - 当分类名包含"效率"或sort_order为2时 -->
    <div v-if="isHorizontalLayout" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <LinkCard
        v-for="link in category.links"
        :key="link.id"
        :link="link"
        :horizontal="true"
        @edit="$emit('edit-link', link)"
        @delete="$emit('delete-link', link)"
        @execute="$emit('execute-script', link)"
      />

      <!-- 添加链接按钮 -->
      <button
        @click="$emit('add-link', category.id)"
        class="add-card flex items-center justify-center gap-2 p-4 rounded-2xl transition card-hover min-h-[72px]"
      >
        <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
        <span class="text-sm text-gray-400">添加链接</span>
      </button>
    </div>

    <!-- 方形卡片布局（常用访问、设计灵感类） -->
    <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
      <LinkCard
        v-for="link in category.links"
        :key="link.id"
        :link="link"
        :horizontal="false"
        @edit="$emit('edit-link', link)"
        @delete="$emit('delete-link', link)"
        @execute="$emit('execute-script', link)"
      />

      <!-- 添加链接按钮 -->
      <button
        @click="$emit('add-link', category.id)"
        class="add-card flex flex-col items-center justify-center p-5 rounded-2xl transition card-hover min-h-[120px]"
      >
        <div class="w-14 h-14 rounded-2xl bg-gray-100/50 flex items-center justify-center text-gray-400 mb-3">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
        </div>
        <span class="text-sm text-gray-400">添加链接</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import LinkCard from './LinkCard.vue'

const props = defineProps({
  category: {
    type: Object,
    required: true
  }
})

defineEmits(['edit-category', 'add-link', 'edit-link', 'delete-link', 'execute-script'])

// 判断是否使用横向布局
const isHorizontalLayout = computed(() => {
  return props.category.name.includes('效率') ||
         props.category.sort_order === 2 ||
         props.category.icon === 'briefcase'
})

function getCategoryIcon(icon) {
  const icons = {
    'star': '⭐',
    'briefcase': '💼',
    'palette': '🎨',
    'code': '💻',
    'music': '🎵',
    'video': '🎬',
    'book': '📚',
    'game': '🎮',
    'cloud': '☁️',
    'home': '🏠'
  }
  return icons[icon] || '📁'
}
</script>
