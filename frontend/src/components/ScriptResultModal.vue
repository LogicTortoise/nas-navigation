<template>
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" @click.self="$emit('close')">
    <div class="bg-white rounded-2xl w-full max-w-lg shadow-xl">
      <!-- 头部 -->
      <div class="flex items-center justify-between p-6 border-b">
        <div class="flex items-center gap-3">
          <div :class="[
            'w-10 h-10 rounded-full flex items-center justify-center',
            result.success ? 'bg-green-100' : 'bg-red-100'
          ]">
            <svg v-if="result.success" class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
            <svg v-else class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-semibold text-gray-900">
              {{ result.success ? '执行成功' : '执行失败' }}
            </h2>
            <p class="text-sm text-gray-500">耗时 {{ result.duration_ms }}ms</p>
          </div>
        </div>
        <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>

      <!-- 输出内容 -->
      <div class="p-6">
        <div v-if="result.output" class="mb-4">
          <h3 class="text-sm font-medium text-gray-700 mb-2">输出</h3>
          <pre class="bg-gray-900 text-green-400 p-4 rounded-lg text-sm overflow-auto max-h-64 font-mono">{{ result.output }}</pre>
        </div>

        <div v-if="result.error" class="mb-4">
          <h3 class="text-sm font-medium text-gray-700 mb-2">错误信息</h3>
          <pre class="bg-red-50 text-red-600 p-4 rounded-lg text-sm overflow-auto max-h-64 font-mono">{{ result.error }}</pre>
        </div>

        <div class="text-sm text-gray-500">
          退出码: {{ result.exit_code }}
        </div>
      </div>

      <!-- 按钮 -->
      <div class="flex justify-end p-6 border-t">
        <button
          @click="$emit('close')"
          class="px-6 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600 transition btn-press"
        >
          关闭
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  result: {
    type: Object,
    required: true
  }
})

defineEmits(['close'])
</script>
