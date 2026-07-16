<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-gray-50 dark:bg-gray-950 p-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <div class="mx-auto w-16 h-16 rounded-2xl bg-brand-600 flex items-center justify-center mb-4">
          <UIcon name="i-heroicons-academic-cap" class="w-8 h-8 text-white" />
        </div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ $t('app.name') }}
        </h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          {{ $t('app.tagline') }}
        </p>
      </div>

      <div class="card">
        <slot />
      </div>

      <div class="text-center mt-6">
        <UButton
          color="gray"
          variant="ghost"
          size="xs"
          :icon="colorMode.preference === 'dark' ? 'i-heroicons-sun' : 'i-heroicons-moon'"
          @click="toggleMode"
        >
          {{ colorMode.preference === 'dark' ? $t('theme.light_mode') : $t('theme.dark_mode') }}
        </UButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const colorMode = useColorMode()
const { t } = useI18n()

const toggleMode = () => {
  colorMode.preference = colorMode.preference === 'dark' ? 'light' : 'dark'
}
</script>
