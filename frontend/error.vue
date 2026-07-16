<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 p-4">
    <div class="max-w-lg w-full text-center">
      <div class="mb-8">
        <div
          class="mx-auto w-24 h-24 rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center mb-6"
        >
          <UIcon name="i-heroicons-exclamation-triangle" class="w-12 h-12 text-red-500" />
        </div>
        <h1 class="text-4xl font-bold text-gray-900 dark:text-white mb-2">
          {{ error.statusCode }}
        </h1>
        <h2 class="text-xl font-semibold text-gray-700 dark:text-gray-300 mb-4">
          {{ error.statusMessage || $t('error.title') }}
        </h2>
        <p class="text-gray-500 dark:text-gray-400 mb-8">
          {{ error.message || $t('error.description') }}
        </p>
      </div>
      <div class="flex items-center justify-center gap-4">
        <UButton
          v-if="error.statusCode !== 404"
          color="primary"
          size="lg"
          @click="handleRetry"
        >
          {{ $t('common.retry') }}
        </UButton>
        <UButton
          color="gray"
          variant="outline"
          size="lg"
          @click="handleClearError"
        >
          {{ $t('common.go_home') }}
        </UButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { NuxtError } from '#app'

const props = defineProps<{
  error: NuxtError
}>()

const { clearError } = useNuxtApp()
const { t } = useI18n()

const handleRetry = () => {
  clearError({ redirect: window.location.pathname })
}

const handleClearError = async () => {
  await clearError({ redirect: '/' })
}
</script>
