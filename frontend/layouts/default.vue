<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950">
    <AppSidebar />
    <div
      class="transition-all duration-300"
      :class="[
        appStore.sidebar.isCollapsed ? 'lg:ml-20' : 'lg:ml-64',
        'ml-0',
      ]"
    >
      <AppHeader />
      <AppBreadcrumb />

      <main class="p-4 sm:p-6 lg:p-8">
        <div class="max-w-[1600px] mx-auto">
          <slot />
        </div>
      </main>
    </div>

    <div
      v-if="appStore.sidebar.isMobileOpen"
      class="fixed inset-0 z-40 bg-black/50 lg:hidden"
      @click="appStore.setMobileSidebar(false)"
    />
  </div>
</template>

<script setup lang="ts">
const appStore = useAppStore()

const handleResize = () => {
  appStore.setMobile(window.innerWidth < 1024)
}

onMounted(() => {
  handleResize()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>
