<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-start justify-center pt-[15vh]"
  >
    <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="close" />

    <div class="relative w-full max-w-xl bg-white dark:bg-gray-800 rounded-xl shadow-2xl animate-slide-up overflow-hidden">
      <div class="flex items-center gap-3 px-4 py-3 border-b border-gray-200 dark:border-gray-700">
        <UIcon name="i-heroicons-magnifying-glass" class="w-5 h-5 text-gray-400" />
        <input
          ref="inputRef"
          v-model="query"
          type="text"
          :placeholder="$t('search.placeholder')"
          class="flex-1 bg-transparent border-none outline-none text-sm text-gray-900 dark:text-white placeholder:text-gray-400"
          @keydown.escape="close"
          @keydown.enter="handleSearch"
        />
        <UKbd size="sm">ESC</UKbd>
      </div>

      <div v-if="query.length > 0" class="max-h-80 overflow-y-auto">
        <div v-if="isLoading" class="flex items-center justify-center py-12">
          <UIcon name="i-heroicons-arrow-path" class="w-5 h-5 animate-spin text-gray-400" />
        </div>

        <template v-else-if="results.length > 0">
          <div
            v-for="group in groupedResults"
            :key="group.category"
            class="border-b border-gray-100 dark:border-gray-750 last:border-0"
          >
            <div class="px-4 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase bg-gray-50 dark:bg-gray-750">
              {{ group.category }}
            </div>
            <button
              v-for="result in group.items"
              :key="result.id"
              class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-left hover:bg-gray-50 dark:hover:bg-gray-750 transition-colors"
              @click="handleSelect(result)"
            >
              <UIcon :name="result.icon || 'i-heroicons-document'" class="w-4 h-4 text-gray-400 shrink-0" />
              <div class="flex-1 min-w-0">
                <p class="text-gray-900 dark:text-white truncate">{{ result.title }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400 truncate">{{ result.subtitle }}</p>
              </div>
            </button>
          </div>
        </template>

        <div v-else class="flex flex-col items-center justify-center py-8">
          <UIcon name="i-heroicons-magnifying-glass" class="w-8 h-8 text-gray-300 mb-3" />
          <p class="text-sm text-gray-500">{{ $t('search.no_results') }}</p>
        </div>
      </div>

      <div v-else class="py-6 px-4">
        <p class="text-sm text-gray-500 text-center mb-4">{{ $t('search.type_to_search') }}</p>
        <div class="grid grid-cols-2 gap-2">
          <button
            v-for="quick in quickActions"
            :key="quick.label"
            class="flex items-center gap-2 px-3 py-2 rounded-lg text-sm text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            @click="query = quick.query"
          >
            <UIcon :name="quick.icon" class="w-4 h-4" />
            <span>{{ $t(quick.label) }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const query = ref('')
const isLoading = ref(false)
const results = ref<SearchResult[]>([])
const inputRef = ref<HTMLInputElement>()
const router = useRouter()
const { t } = useI18n()

interface SearchResult {
  id: string
  title: string
  subtitle: string
  category: string
  icon: string
  to: string
}

const quickActions = [
  { label: 'search.students', icon: 'i-heroicons-users', query: t('search.students') },
  { label: 'search.teachers', icon: 'i-heroicons-user-circle', query: t('search.teachers') },
  { label: 'search.invoices', icon: 'i-heroicons-currency-dollar', query: t('search.invoices') },
  { label: 'search.settings', icon: 'i-heroicons-cog-6-tooth', query: t('search.settings') },
]

const groupedResults = computed(() => {
  const groups: Record<string, SearchResult[]> = {}
  for (const result of results.value) {
    if (!groups[result.category]) {
      groups[result.category] = []
    }
    groups[result.category].push(result)
  }
  return Object.entries(groups).map(([category, items]) => ({ category, items }))
})

const close = () => {
  emit('close')
  query.value = ''
}

const handleSearch = async () => {
  if (!query.value.trim()) return
  isLoading.value = true
  try {
    const api = useApi()
    const response = await api.get<SearchResult[]>('/search', { q: query.value })
    results.value = response
  } catch {
    results.value = []
  } finally {
    isLoading.value = false
  }
}

const handleSelect = (result: SearchResult) => {
  close()
  router.push(result.to)
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'k' && (e.metaKey || e.ctrlKey)) {
    e.preventDefault()
    if (props.isOpen) {
      close()
    } else {
      query.value = ''
      nextTick(() => inputRef.value?.focus())
    }
  }
}

watch(() => props.isOpen, (open) => {
  if (open) {
    nextTick(() => {
      inputRef.value?.focus()
    })
  } else {
    results.value = []
    query.value = ''
  }
})

let debounceTimer: ReturnType<typeof setTimeout>
watch(query, (val) => {
  clearTimeout(debounceTimer)
  if (val.length >= 2) {
    debounceTimer = setTimeout(() => {
      handleSearch()
    }, 300)
  } else {
    results.value = []
  }
})

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})
onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>
