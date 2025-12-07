<template>
  <!-- 横向添加按钮 -->
  <button
    v-if="horizontal"
    class="glass-card flex items-center p-4 rounded-xl hover:bg-white transition-colors group border-2 border-dashed border-slate-200 hover:border-indigo-400"
    @click="$emit('click')"
  >
    <div class="w-10 h-10 rounded-lg bg-slate-100 text-slate-400 flex items-center justify-center shrink-0 group-hover:bg-indigo-100 group-hover:text-indigo-500 transition-colors">
      <i class="fa-solid fa-plus"></i>
    </div>
    <div class="ml-3 overflow-hidden flex-1">
      <span class="text-sm font-medium text-slate-400 group-hover:text-indigo-500 transition-colors">添加链接</span>
    </div>
  </button>

  <!-- 方形添加按钮 -->
  <button
    v-else
    class="glass-card group relative flex flex-col items-center justify-center p-6 rounded-2xl h-32 md:h-40 cursor-pointer border-2 border-dashed border-slate-200 hover:border-indigo-400 transition-all"
    @click="$emit('click')"
  >
    <div class="w-12 h-12 md:w-14 md:h-14 rounded-2xl bg-slate-100 flex items-center justify-center mb-3 group-hover:bg-indigo-100 transition-colors">
      <i class="fa-solid fa-plus text-slate-400 group-hover:text-indigo-500 transition-colors text-xl"></i>
    </div>
    <span class="text-sm font-medium text-slate-400 group-hover:text-indigo-500 transition-colors">添加链接</span>
  </button>
</template>

<script setup>
defineProps({
  horizontal: {
    type: Boolean,
    default: false
  }
})

defineEmits(['click'])
</script>
