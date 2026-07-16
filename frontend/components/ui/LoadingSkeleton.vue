<template>
  <div class="space-y-4">
    <div v-if="type === 'card'" class="card animate-pulse">
      <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/3 mb-4" />
      <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded w-full mb-2" />
      <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded w-2/3" />
    </div>

    <div v-else-if="type === 'table'" class="space-y-3">
      <div
        v-for="i in rows"
        :key="i"
        class="flex items-center gap-4"
      >
        <div
          v-for="j in columns"
          :key="j"
          class="h-4 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"
          :style="{ width: `${100 / columns - 4}%` }"
        />
      </div>
    </div>

    <div v-else-if="type === 'stat'" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div
        v-for="i in count"
        :key="i"
        class="card animate-pulse"
      >
        <div class="flex items-start justify-between">
          <div class="space-y-2 flex-1">
            <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded w-1/2" />
            <div class="h-6 bg-gray-200 dark:bg-gray-700 rounded w-2/3" />
          </div>
          <div class="w-10 h-10 rounded-lg bg-gray-200 dark:bg-gray-700" />
        </div>
      </div>
    </div>

    <div v-else-if="type === 'detail'" class="card animate-pulse space-y-4">
      <div class="flex items-center gap-3">
        <div class="w-12 h-12 rounded-full bg-gray-200 dark:bg-gray-700" />
        <div class="space-y-2 flex-1">
          <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/3" />
          <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded w-1/2" />
        </div>
      </div>
      <div v-for="i in 4" :key="i" class="h-3 bg-gray-200 dark:bg-gray-700 rounded w-full" />
    </div>

    <div v-else class="flex items-center justify-center py-8">
      <UIcon name="i-heroicons-arrow-path" class="w-6 h-6 text-gray-400 animate-spin" />
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    type?: 'card' | 'table' | 'stat' | 'detail' | 'spinner'
    rows?: number
    columns?: number
    count?: number
  }>(),
  {
    type: 'card',
    rows: 5,
    columns: 4,
    count: 4,
  }
)
</script>
