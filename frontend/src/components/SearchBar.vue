<template>
  <div class="relative">
    <div class="flex items-center glass-card rounded-2xl px-4 py-2.5 focus-within:ring-2 focus-within:ring-blue-400 transition">
      <svg class="w-5 h-5 text-gray-400 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
      </svg>
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜索书签、Google 或输入网址..."
        class="flex-1 bg-transparent outline-none text-gray-700 placeholder-gray-400 text-sm"
        @input="handleSearch"
      />
      <div v-if="searchQuery" class="flex items-center gap-2">
        <button @click="clearSearch" class="text-gray-400 hover:text-gray-600">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>
      <div class="hidden sm:flex items-center text-gray-400 text-xs ml-3 flex-shrink-0">
        <span class="px-1.5 py-0.5 bg-gray-200/60 rounded text-[10px] font-medium">⌘K</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAppStore } from '../stores/app'
import { useDebounceFn } from '@vueuse/core'

const store = useAppStore()
const searchQuery = ref('')

const debouncedSearch = useDebounceFn((query) => {
  store.searchLinks(query)
}, 300)

function handleSearch() {
  store.searchQuery = searchQuery.value
  debouncedSearch(searchQuery.value)
}

function clearSearch() {
  searchQuery.value = ''
  store.searchQuery = ''
  store.searchResults = []
}

// 监听键盘快捷键
if (typeof window !== 'undefined') {
  window.addEventListener('keydown', (e) => {
    if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
      e.preventDefault()
      document.querySelector('input[type="text"]')?.focus()
    }
  })
}
</script>
