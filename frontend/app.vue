<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
  <Toaster
    :position="'top-right'"
    :rich-colors="true"
    :close-button="true"
    :duration="4000"
    :toast-options="{
      classes: {
        toast: 'group toast group-[.toast]:bg-white dark:group-[.toast]:bg-gray-800',
        title: 'text-sm font-semibold',
        description: 'text-xs text-gray-500 dark:text-gray-400',
      },
    }"
  />
</template>

<script setup lang="ts">
import { Toaster } from 'vue-sonner'

const colorMode = useColorMode()

onMounted(() => {
  colorMode.preference = colorMode.preference || 'system'
})

useHead({
  htmlAttrs: {
    class: computed(() => colorMode.value === 'dark' ? 'dark' : ''),
  },
})
</script>
