<template>
  <div />
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth'],
})

const authStore = useAuthStore()

onBeforeMount(async () => {
  if (!authStore.isInitialized) {
    await authStore.initializeFromStorage()
  }

  if (authStore.isAuthenticated) {
    await navigateTo('/dashboard')
  } else {
    await navigateTo('/auth/login')
  }
})
</script>
