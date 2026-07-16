<template>
  <header class="sticky top-0 z-30 h-16 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md border-b border-gray-200 dark:border-gray-800">
    <div class="flex items-center justify-between h-full px-4">
      <div class="flex items-center gap-3">
        <UButton
          color="gray"
          variant="ghost"
          size="sm"
          icon="i-heroicons-bars-3"
          class="lg:hidden"
          @click="appStore.toggleMobileSidebar()"
        />

        <div class="hidden sm:flex items-center gap-1.5 text-sm">
          <UIcon name="i-heroicons-map-pin" class="w-4 h-4 text-brand-600" />
          <span class="text-gray-700 dark:text-gray-300 font-medium">
            {{ school.currentSchool?.name || $t('common.loading') }}
          </span>
          <span class="text-gray-400">|</span>
          <span class="text-gray-500 dark:text-gray-400">
            {{ school.currentAcademicYear?.name || '---' }}
          </span>
        </div>
      </div>

      <div class="flex items-center gap-1.5">
        <UButton
          color="gray"
          variant="ghost"
          size="sm"
          icon="i-heroicons-magnifying-glass"
          @click="appStore.openSearch()"
        >
          <span class="hidden lg:inline ml-1.5 text-xs text-gray-500">Ctrl+K</span>
        </UButton>

        <UButton
          color="gray"
          variant="ghost"
          size="sm"
          icon="i-heroicons-bell"
          class="relative"
          @click="toggleNotifications"
        >
          <span
            v-if="appStore.unreadNotifications > 0"
            class="absolute -top-0.5 -right-0.5 w-4 h-4 text-[10px] font-bold rounded-full bg-red-500 text-white flex items-center justify-center"
          >
            {{ appStore.unreadNotifications > 9 ? '9+' : appStore.unreadNotifications }}
          </span>
        </UButton>

        <div class="relative">
          <UButton
            color="gray"
            variant="ghost"
            size="sm"
            class="pl-1.5 pr-2 gap-2"
            @click="showUserMenu = !showUserMenu"
          >
            <div class="w-7 h-7 rounded-full bg-brand-100 dark:bg-brand-900/40 flex items-center justify-center shrink-0">
              <span class="text-xs font-bold text-brand-700 dark:text-brand-400">
                {{ avatarInitials }}
              </span>
            </div>
          </UButton>

          <div
            v-if="showUserMenu"
            class="absolute right-0 top-full mt-2 w-56 rounded-xl bg-white dark:bg-gray-800 shadow-lg border border-gray-200 dark:border-gray-700 py-1 animate-fade-in overflow-hidden"
            @click="showUserMenu = false"
          >
            <div class="px-4 py-3 border-b border-gray-100 dark:border-gray-700">
              <p class="text-sm font-semibold text-gray-900 dark:text-white">{{ authStore.user?.fullName }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400 truncate">{{ authStore.user?.email }}</p>
            </div>

            <NuxtLink
              to="/profile"
              class="flex items-center gap-2 px-4 py-2.5 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-750"
            >
              <UIcon name="i-heroicons-user" class="w-4 h-4" />
              {{ $t('menu.profile') }}
            </NuxtLink>

            <NuxtLink
              to="/settings"
              class="flex items-center gap-2 px-4 py-2.5 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-750"
            >
              <UIcon name="i-heroicons-cog-6-tooth" class="w-4 h-4" />
              {{ $t('menu.settings') }}
            </NuxtLink>

            <div class="border-t border-gray-100 dark:border-gray-700 my-1" />

            <button
              class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20"
              @click="handleLogout"
            >
              <UIcon name="i-heroicons-arrow-right-on-rectangle" class="w-4 h-4" />
              {{ $t('auth.logout') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <AppNotifications
      :is-open="showNotifications"
      @close="showNotifications = false"
    />

    <AppSearch
      :is-open="appStore.isSearchOpen"
      @close="appStore.closeSearch()"
    />
  </header>
</template>

<script setup lang="ts">
const appStore = useAppStore()
const authStore = useAuthStore()
const auth = useAuth()
const school = useSchool()

const showUserMenu = ref(false)
const showNotifications = ref(false)

const avatarInitials = computed(() => {
  const name = authStore.user?.fullName || 'User'
  return name
    .split(' ')
    .map(n => n[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
})

const toggleNotifications = () => {
  showNotifications.value = !showNotifications.value
  showUserMenu.value = false
}

const handleLogout = async () => {
  showUserMenu.value = false
  await auth.logout()
}

const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (!target.closest('.relative')) {
    showUserMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
