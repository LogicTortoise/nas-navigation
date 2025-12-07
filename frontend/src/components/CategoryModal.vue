<template>
  <!-- 模态弹窗遮罩 -->
  <div
    class="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    @click.self="$emit('close')"
  >
    <!-- 模态弹窗内容 -->
    <div class="bg-white rounded-2xl shadow-2xl overflow-hidden w-full max-w-md modal-enter">

      <!-- 弹窗头部 -->
      <div class="flex items-center justify-between p-6 border-b border-slate-200">
        <div class="flex items-center">
          <div class="w-8 h-8 rounded-lg bg-indigo-100 flex items-center justify-center mr-3">
            <i class="fa-solid fa-folder-plus text-indigo-600"></i>
          </div>
          <h2 class="text-xl font-bold text-slate-800">
            {{ category ? '编辑分类' : '添加分类' }}
          </h2>
        </div>
        <button
          @click="$emit('close')"
          class="text-slate-400 hover:text-slate-600 transition-colors p-1 rounded-lg hover:bg-slate-100"
        >
          <i class="fa-solid fa-xmark text-xl"></i>
        </button>
      </div>

      <!-- 弹窗内容区 -->
      <div class="p-6">
        <form @submit.prevent="handleSubmit" class="space-y-6">

          <!-- 分类名称 -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-slate-700">
              分类名称 <span class="text-rose-500">*</span>
            </label>
            <input
              type="text"
              v-model="form.name"
              class="w-full px-4 py-3 rounded-lg border border-slate-300 form-input-focus bg-white/80 backdrop-blur-sm"
              placeholder="请输入分类名称"
              required
              maxlength="20"
            >
            <p class="text-xs text-slate-400">最多20个字符</p>
          </div>

          <!-- 选择图标 -->
          <div class="space-y-3">
            <label class="block text-sm font-medium text-slate-700">选择图标</label>

            <!-- 图标选择器 -->
            <div class="icon-selector p-3 bg-slate-50 rounded-lg border border-slate-200">
              <button
                v-for="iconClass in iconList"
                :key="iconClass"
                type="button"
                @click="form.icon = iconClass"
                :class="[
                  'icon-option',
                  form.icon === iconClass ? 'selected' : ''
                ]"
              >
                <i :class="['fa-solid', iconClass, 'text-lg']"></i>
              </button>
            </div>

            <!-- 当前选中的图标预览 -->
            <div class="flex items-center space-x-3 p-3 bg-indigo-50 rounded-lg border border-indigo-200">
              <div class="w-10 h-10 rounded-lg bg-indigo-600 flex items-center justify-center text-white">
                <i :class="['fa-solid', form.icon, 'text-lg']"></i>
              </div>
              <div>
                <p class="text-sm font-medium text-slate-800">当前图标</p>
                <p class="text-xs text-slate-500">点击选择其他图标</p>
              </div>
            </div>
          </div>

        </form>
      </div>

      <!-- 弹窗底部 -->
      <div class="flex items-center justify-between p-6 border-t border-slate-200 bg-slate-50/50">
        <button
          v-if="category"
          type="button"
          @click="$emit('delete', category)"
          class="px-4 py-2.5 rounded-lg text-red-600 hover:bg-red-50 font-medium transition-all"
        >
          <i class="fa-solid fa-trash mr-2"></i>删除分类
        </button>
        <div v-else></div>

        <div class="flex items-center space-x-3">
          <button
            type="button"
            @click="$emit('close')"
            class="px-6 py-2.5 rounded-lg btn-secondary font-medium"
          >
            取消
          </button>
          <button
            type="submit"
            @click="handleSubmit"
            class="px-6 py-2.5 rounded-lg btn-primary font-medium shadow-lg shadow-indigo-500/20"
          >
            <i class="fa-solid fa-check mr-2"></i>
            保存分类
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, onMounted } from 'vue'

const props = defineProps({
  category: Object
})

const emit = defineEmits(['close', 'save', 'delete'])

// 预设图标列表 - 使用 Font Awesome 图标
const iconList = [
  'fa-folder', 'fa-briefcase', 'fa-graduation-cap', 'fa-pen-nib', 'fa-gamepad',
  'fa-heart', 'fa-star', 'fa-book', 'fa-music', 'fa-video',
  'fa-image', 'fa-code', 'fa-globe', 'fa-shopping-cart', 'fa-home',
  'fa-users', 'fa-cog', 'fa-chart-bar', 'fa-calendar', 'fa-envelope',
  'fa-file-alt', 'fa-search', 'fa-download', 'fa-upload', 'fa-share',
  'fa-link', 'fa-bookmark', 'fa-tag', 'fa-folder-open', 'fa-archive'
]

const form = reactive({
  name: '',
  icon: 'fa-folder',
  sort_order: 0
})

function handleSubmit() {
  emit('save', { ...form })
}

onMounted(() => {
  if (props.category) {
    Object.assign(form, {
      name: props.category.name,
      icon: props.category.icon || 'fa-folder',
      sort_order: props.category.sort_order || 0
    })
  }
})
</script>
