<template>
  <div class="flex flex-col items-center justify-center py-12 px-4">
    <div class="w-20 h-20 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center mb-4">
      <UIcon :name="icon || 'i-heroicons-inbox'" class="w-10 h-10 text-gray-400 dark:text-gray-500" />
    </div>
    <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-1">
      {{ title }}
    </h3>
    <p class="text-sm text-gray-500 dark:text-gray-400 text-center max-w-sm mb-4">
      {{ description }}
    </p>
    <UButton
      v-if="actionLabel"
      color="primary"
      size="sm"
      @click="$emit('action')"
    >
      {{ actionLabel }}
    </UButton>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  description?: string
  icon?: string
  actionLabel?: string
}>()

defineEmits<{
  action: []
}>()
</script>
