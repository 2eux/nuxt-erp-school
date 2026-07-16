<template>
  <div class="w-full">
    <div class="flex flex-wrap items-end gap-3">
      <div v-if="searchable" class="flex-1 min-w-[200px]">
        <label class="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
          {{ $t('common.search') }}
        </label>
        <UInput
          v-model="filters.search"
          :placeholder="$t('search.placeholder')"
          icon="i-heroicons-magnifying-glass"
          size="sm"
          color="gray"
          @input="debouncedApply"
        />
      </div>

      <div v-for="field in filterFields" :key="field.key" class="min-w-[150px]">
        <label class="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
          {{ $t(field.label) || field.label }}
        </label>

        <UInput
          v-if="field.type === 'text'"
          v-model="filters[field.key]"
          :placeholder="field.placeholder"
          size="sm"
          color="gray"
          @input="debouncedApply"
        />

        <USelect
          v-else-if="field.type === 'select'"
          v-model="filters[field.key]"
          :options="field.options || []"
          :placeholder="field.placeholder || $t('common.select')"
          size="sm"
          color="gray"
          @change="applyFilters"
        />

        <div
          v-else-if="field.type === 'date_range'"
          class="flex items-center gap-2"
        >
          <UInput
            v-model="filters[`${field.key}_start`]"
            type="date"
            size="sm"
            color="gray"
            @change="applyFilters"
          />
          <span class="text-gray-400 text-xs">-</span>
          <UInput
            v-model="filters[`${field.key}_end`]"
            type="date"
            size="sm"
            color="gray"
            @change="applyFilters"
          />
        </div>

        <UInput
          v-else-if="field.type === 'number'"
          v-model.number="filters[field.key]"
          type="number"
          size="sm"
          color="gray"
          @input="debouncedApply"
        />
      </div>

      <div class="flex items-center gap-2 ml-auto">
        <UButton
          v-if="hasActiveFilters"
          color="gray"
          variant="ghost"
          size="sm"
          icon="i-heroicons-x-mark"
          @click="clearFilters"
        >
          {{ $t('common.clear') }}
        </UButton>

        <div v-if="showPresets" class="flex items-center gap-1">
          <UButton
            color="gray"
            variant="ghost"
            size="xs"
            icon="i-heroicons-bookmark"
            @click="showSavePreset = true"
          >
            {{ $t('filter.save') }}
          </UButton>

          <USelect
            v-if="presets.length > 0"
            v-model="selectedPreset"
            :options="presetOptions"
            :placeholder="$t('filter.load_preset')"
            size="xs"
            color="gray"
            class="w-32"
            @change="loadPreset"
          />
        </div>
      </div>
    </div>

    <div
      v-if="showActiveFilters && activeFilterLabels.length > 0"
      class="flex flex-wrap gap-1.5 mt-2"
    >
      <span
        v-for="label in activeFilterLabels"
        :key="label.key"
        class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-brand-50 dark:bg-brand-900/20 text-xs text-brand-700 dark:text-brand-400"
      >
        {{ label.label }}
        <button @click="removeFilter(label.key)">
          <UIcon name="i-heroicons-x-mark" class="w-3 h-3" />
        </button>
      </span>
    </div>

    <UModal v-model:open="showSavePreset" title="Save Filter Preset">
      <template #body>
        <div class="space-y-3">
          <UInput
            v-model="newPresetName"
            label="Preset Name"
            placeholder="Enter preset name..."
          />
          <div class="flex justify-end gap-2">
            <UButton color="gray" variant="ghost" @click="showSavePreset = false">
              {{ $t('common.cancel') }}
            </UButton>
            <UButton color="primary" @click="savePreset">
              {{ $t('common.save') }}
            </UButton>
          </div>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
interface FilterField {
  key: string
  label: string
  type: 'text' | 'select' | 'date_range' | 'number'
  placeholder?: string
  options?: { label: string; value: string | number }[]
}

interface FilterPreset {
  name: string
  filters: Record<string, unknown>
}

const props = withDefaults(
  defineProps<{
    filterFields?: FilterField[]
    searchable?: boolean
    showPresets?: boolean
    showActiveFilters?: boolean
  }>(),
  {
    filterFields: () => [],
    searchable: true,
    showPresets: false,
    showActiveFilters: true,
  }
)

const emit = defineEmits<{
  'apply': [filters: Record<string, unknown>]
}>()

const filters = reactive<Record<string, unknown>>({})
const presets = ref<FilterPreset[]>([])
const selectedPreset = ref('')
const showSavePreset = ref(false)
const newPresetName = ref('')

let debounceTimer: ReturnType<typeof setTimeout>

const hasActiveFilters = computed(() => {
  return Object.values(filters).some(v => v !== undefined && v !== null && v !== '')
})

const activeFilterLabels = computed(() => {
  const labels: { key: string; label: string }[] = []
  for (const [key, value] of Object.entries(filters)) {
    if (value !== undefined && value !== null && value !== '') {
      const field = props.filterFields.find(f => f.key === key || key.startsWith(f.key))
      if (field) {
        labels.push({ key, label: `${field.label}: ${value}` })
      }
    }
  }
  return labels
})

const presetOptions = computed(() =>
  presets.value.map(p => ({ label: p.name, value: p.name }))
)

const applyFilters = () => {
  emit('apply', { ...filters })
}

const debouncedApply = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(applyFilters, 400)
}

const clearFilters = () => {
  for (const key of Object.keys(filters)) {
    delete filters[key]
  }
  applyFilters()
}

const removeFilter = (key: string) => {
  delete filters[key]
  applyFilters()
}

const savePreset = () => {
  if (!newPresetName.value.trim()) return
  presets.value.push({
    name: newPresetName.value.trim(),
    filters: { ...filters },
  })
  if (import.meta.client) {
    localStorage.setItem('filter_presets', JSON.stringify(presets.value))
  }
  newPresetName.value = ''
  showSavePreset.value = false
}

const loadPreset = (name: string) => {
  const preset = presets.value.find(p => p.name === name)
  if (preset) {
    for (const key of Object.keys(filters)) {
      delete filters[key]
    }
    Object.assign(filters, preset.filters)
    applyFilters()
  }
}

onMounted(() => {
  if (import.meta.client && props.showPresets) {
    const saved = localStorage.getItem('filter_presets')
    if (saved) {
      try {
        presets.value = JSON.parse(saved)
      } catch {
        // invalid JSON
      }
    }
  }
})
</script>
