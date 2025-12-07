<template>
  <div class="relative" ref="selectRef">
    <!-- 选择器触发按钮 -->
    <button
      type="button"
      @click="isOpen = !isOpen"
      class="w-full pl-12 pr-10 py-3 bg-white/80 border border-slate-300 rounded-xl form-input-focus transition-all text-slate-800 text-left flex items-center justify-between"
      :class="{ 'border-indigo-500 ring-2 ring-indigo-500/20': isOpen }"
    >
      <span :class="selectedLabel ? 'text-slate-800' : 'text-slate-400'">
        {{ selectedLabel || placeholder }}
      </span>
      <i
        class="fa-solid fa-chevron-down text-slate-400 transition-transform duration-200"
        :class="{ 'rotate-180': isOpen }"
      ></i>
    </button>

    <!-- 左侧图标 -->
    <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
      <i :class="['text-slate-400', icon]"></i>
    </div>

    <!-- 下拉选项 -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div
        v-if="isOpen"
        class="absolute z-50 mt-2 w-full bg-white/95 backdrop-blur-xl border border-slate-200 rounded-xl shadow-xl overflow-hidden"
      >
        <div class="max-h-60 overflow-y-auto py-1">
          <!-- 空选项 -->
          <button
            v-if="allowEmpty"
            type="button"
            @click="selectOption(null, emptyLabel)"
            class="w-full px-4 py-2.5 text-left hover:bg-indigo-50 transition-colors flex items-center"
            :class="{ 'bg-indigo-50 text-indigo-600': modelValue === null }"
          >
            <i class="fa-solid fa-minus-circle text-slate-400 mr-3 w-5 text-center"></i>
            <span>{{ emptyLabel }}</span>
          </button>

          <!-- 选项列表 -->
          <button
            v-for="option in options"
            :key="option.value"
            type="button"
            @click="selectOption(option.value, option.label)"
            class="w-full px-4 py-2.5 text-left hover:bg-indigo-50 transition-colors flex items-center"
            :class="{ 'bg-indigo-50 text-indigo-600': modelValue === option.value }"
          >
            <i
              v-if="option.icon"
              :class="[option.icon, 'mr-3 w-5 text-center']"
              :style="{ color: option.color || '#64748b' }"
            ></i>
            <i
              v-else
              class="fa-solid fa-check mr-3 w-5 text-center"
              :class="modelValue === option.value ? 'text-indigo-600' : 'text-transparent'"
            ></i>
            <span>{{ option.label }}</span>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  modelValue: [String, Number, null],
  options: {
    type: Array,
    required: true
    // [{ value: 1, label: 'Option 1', icon: 'fa-solid fa-star', color: '#eab308' }]
  },
  placeholder: {
    type: String,
    default: '请选择'
  },
  icon: {
    type: String,
    default: 'fa-solid fa-folder'
  },
  allowEmpty: {
    type: Boolean,
    default: false
  },
  emptyLabel: {
    type: String,
    default: '无'
  }
})

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false)
const selectRef = ref(null)

const selectedLabel = computed(() => {
  if (props.modelValue === null || props.modelValue === undefined) {
    return props.allowEmpty ? null : null
  }
  const option = props.options.find(o => o.value === props.modelValue)
  return option ? option.label : null
})

function selectOption(value, label) {
  emit('update:modelValue', value)
  isOpen.value = false
}

// 点击外部关闭
function handleClickOutside(event) {
  if (selectRef.value && !selectRef.value.contains(event.target)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
