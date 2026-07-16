<template>
  <div class="space-y-1 font-mono text-sm">
    <template v-if="isObject(data)">
      <div v-for="(entries, path) in groupByPrefix(data, 3)" :key="path" class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <button
          v-if="path"
          class="w-full flex items-center gap-2 px-3 py-2 bg-gray-50 dark:bg-gray-800/50 text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
          @click="toggleGroup(path)"
        >
          <UIcon
            :name="expandedGroups.has(path) ? 'i-heroicons-chevron-down' : 'i-heroicons-chevron-right'"
            class="w-4 h-4"
          />
          <span>{{ path }}</span>
          <span class="text-xs text-gray-400 ml-auto">{{ Object.keys(entries).length }} items</span>
        </button>
        <div v-if="!path || expandedGroups.has(path)" class="divide-y divide-gray-100 dark:divide-gray-750">
          <div
            v-for="(value, key) in entries"
            :key="key"
            class="flex items-start px-3 py-2 gap-4"
          >
            <span
              class="text-xs font-semibold shrink-0 min-w-[120px] truncate"
              :class="value === null ? 'text-gray-400' : 'text-brand-600 dark:text-brand-400'"
            >
              {{ path ? `${path}.${key}` : key }}
            </span>
            <span class="flex-1 break-all">
              <JsonValue :value="value" />
            </span>
          </div>
        </div>
      </div>
    </template>
    <template v-else>
      <JsonValue :value="data" />
    </template>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  data: unknown
  collapsible?: boolean
}>()

const expandedGroups = ref(new Set<string>())

const isObject = (val: unknown): val is Record<string, unknown> => {
  return val !== null && typeof val === 'object' && !Array.isArray(val)
}

const groupByPrefix = (obj: Record<string, unknown>, groupDepth: number): Record<string, Record<string, unknown>> => {
  const groups: Record<string, Record<string, unknown>> = { '': {} }

  for (const [key, value] of Object.entries(obj)) {
    const parts = key.split('.')
    const prefix = parts.length > groupDepth ? parts.slice(0, groupDepth).join('.') : ''

    if (!groups[prefix]) {
      groups[prefix] = {}
    }
    groups[prefix][key] = value
  }

  return groups
}

const toggleGroup = (path: string) => {
  if (expandedGroups.value.has(path)) {
    expandedGroups.value.delete(path)
  } else {
    expandedGroups.value.add(path)
  }
}

watch(() => props.data, () => {
  expandedGroups.value.clear()
})
</script>

<script lang="ts">
import { h, defineComponent } from 'vue'

const JsonValue = defineComponent({
  name: 'JsonValue',
  props: {
    value: { required: true },
  },
  setup(props) {
    return () => {
      const v = props.value
      if (v === null) return h('span', { class: 'text-gray-400 italic' }, 'null')
      if (v === undefined) return h('span', { class: 'text-gray-400 italic' }, 'undefined')
      if (typeof v === 'boolean') return h('span', { class: 'text-amber-600 dark:text-amber-400' }, String(v))
      if (typeof v === 'number') return h('span', { class: 'text-blue-600 dark:text-blue-400' }, String(v))
      if (typeof v === 'string') {
        if (v.startsWith('http') || v.startsWith('data:image')) {
          return h('a', { href: v, class: 'text-blue-500 underline break-all', target: '_blank' }, v)
        }
        if (v.length > 200) {
          return h('span', { class: 'text-emerald-600 dark:text-emerald-400' }, `"${v.slice(0, 200)}..."`)
        }
        return h('span', { class: 'text-emerald-600 dark:text-emerald-400' }, `"${v}"`)
      }
      if (Array.isArray(v)) {
        return h('span', { class: 'text-purple-600 dark:text-purple-400' }, `Array(${v.length})`)
      }
      if (typeof v === 'object') {
        const keys = Object.keys(v as Record<string, unknown>)
        return h('span', { class: 'text-gray-500' }, `{${keys.length} keys}`)
      }
      return h('span', {}, String(v))
    }
  },
})
</script>
