<template>
  <div class="card-hover">
    <div class="flex items-start justify-between">
      <div class="flex-1 min-w-0">
        <p class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
          {{ label }}
        </p>
        <p class="text-2xl font-bold text-gray-900 dark:text-white mt-1 tabular-nums">
          <template v-if="loading">
            <span class="inline-block w-16 h-7 bg-gray-200 dark:bg-gray-700 rounded animate-pulse" />
          </template>
          <template v-else>
            {{ formattedValue }}
          </template>
        </p>
      </div>
      <div
        class="w-12 h-12 rounded-xl flex items-center justify-center shrink-0"
        :class="colorClasses.bg"
      >
        <UIcon :name="icon" class="w-6 h-6" :class="colorClasses.icon" />
      </div>
    </div>

    <div v-if="trend !== undefined && !loading" class="flex items-center gap-1 mt-3">
      <UIcon
        :name="trend >= 0 ? 'i-heroicons-arrow-trending-up' : 'i-heroicons-arrow-trending-down'"
        class="w-4 h-4"
        :class="trend >= 0 ? 'text-emerald-600' : 'text-red-600'"
      />
      <span
        class="text-sm font-medium"
        :class="trend >= 0 ? 'text-emerald-600' : 'text-red-600'"
      >
        {{ Math.abs(trend) }}%
      </span>
      <span class="text-sm text-gray-500 dark:text-gray-400 ml-1">{{ trendLabel }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    label: string
    value: string | number
    icon: string
    color?: 'emerald' | 'blue' | 'amber' | 'red' | 'purple' | 'indigo'
    trend?: number
    trendLabel?: string
    loading?: boolean
    prefix?: string
    suffix?: string
  }>(),
  {
    color: 'emerald',
    loading: false,
    prefix: '',
    suffix: '',
  }
)

const colorClasses = computed(() => {
  const map = {
    emerald: { bg: 'bg-emerald-100 dark:bg-emerald-900/30', icon: 'text-emerald-600 dark:text-emerald-400' },
    blue: { bg: 'bg-blue-100 dark:bg-blue-900/30', icon: 'text-blue-600 dark:text-blue-400' },
    amber: { bg: 'bg-amber-100 dark:bg-amber-900/30', icon: 'text-amber-600 dark:text-amber-400' },
    red: { bg: 'bg-red-100 dark:bg-red-900/30', icon: 'text-red-600 dark:text-red-400' },
    purple: { bg: 'bg-purple-100 dark:bg-purple-900/30', icon: 'text-purple-600 dark:text-purple-400' },
    indigo: { bg: 'bg-indigo-100 dark:bg-indigo-900/30', icon: 'text-indigo-600 dark:text-indigo-400' },
  }
  return map[props.color] || map.emerald
})

const formattedValue = computed(() => {
  if (typeof props.value === 'number') {
    return `${props.prefix}${props.value.toLocaleString()}${props.suffix}`
  }
  return `${props.prefix}${props.value}${props.suffix}`
})
</script>
