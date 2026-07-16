<template>
  <div
    v-if="isOpen"
    class="absolute right-4 top-full mt-2 w-96 bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 animate-fade-in overflow-hidden"
  >
    <div class="flex items-center justify-between px-4 py-3 border-b border-gray-100 dark:border-gray-700">
      <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
        {{ $t('notifications.title') }}
      </h3>
      <UButton
        v-if="notifications.length > 0"
        color="brand"
        variant="link"
        size="xs"
        @click="markAllAsRead"
      >
        {{ $t('notifications.mark_all_read') }}
      </UButton>
    </div>

    <div class="max-h-80 overflow-y-auto">
      <template v-if="notifications.length === 0">
        <div class="flex flex-col items-center justify-center py-8 px-4">
          <UIcon name="i-heroicons-bell-slash" class="w-8 h-8 text-gray-300 dark:text-gray-600 mb-3" />
          <p class="text-sm text-gray-500 dark:text-gray-400">
            {{ $t('notifications.empty') }}
          </p>
        </div>
      </template>

      <template v-else>
        <button
          v-for="notification in notifications"
          :key="notification.id"
          class="w-full flex items-start gap-3 px-4 py-3 text-left transition-colors hover:bg-gray-50 dark:hover:bg-gray-750 border-b border-gray-50 dark:border-gray-750 last:border-0"
          :class="{ 'bg-brand-50/50 dark:bg-brand-900/10': !notification.isRead }"
          @click="handleNotificationClick(notification)"
        >
          <div
            class="w-2 h-2 mt-1.5 rounded-full shrink-0"
            :class="{
              'bg-brand-500': notification.type === 'info',
              'bg-emerald-500': notification.type === 'success',
              'bg-amber-500': notification.type === 'warning',
              'bg-red-500': notification.type === 'error',
            }"
          />
          <div class="flex-1 min-w-0">
            <p
              class="text-sm text-gray-900 dark:text-white truncate"
              :class="{ 'font-semibold': !notification.isRead }"
            >
              {{ notification.title }}
            </p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">
              {{ notification.message }}
            </p>
            <p class="text-xs text-gray-400 dark:text-gray-500 mt-1">
              {{ formatTime(notification.createdAt) }}
            </p>
          </div>
          <UButton
            color="gray"
            variant="ghost"
            size="2xs"
            icon="i-heroicons-x-mark"
            class="shrink-0 -mt-0.5"
            @click.stop="handleDelete(notification.id)"
          />
        </button>
      </template>
    </div>

    <div class="border-t border-gray-100 dark:border-gray-700 p-2">
      <UButton
        color="gray"
        variant="ghost"
        size="sm"
        block
        to="/notifications"
        @click="$emit('close')"
      >
        {{ $t('notifications.view_all') }}
      </UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Notification } from '~/types'

const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const notif = useNotification()

const notifications = computed(() => notif.notifications.value || [])

const formatTime = (date: string): string => {
  const { $dayjs } = useNuxtApp()
  return $dayjs(date).fromNow()
}

const markAllAsRead = async () => {
  await notif.markAllAsRead()
}

const handleNotificationClick = async (notification: Notification) => {
  if (!notification.isRead) {
    await notif.markAsRead(notification.id)
  }
  emit('close')

  if (notification.referenceType && notification.referenceId) {
    // navigate based on referenceType
    const routes: Record<string, string> = {
      announcement: '/communication/announcements',
      message: '/communication/messages',
      invoice: '/finance/invoices',
      student: '/students',
    }
    const base = routes[notification.referenceType]
    if (base) {
      await navigateTo(`${base}/${notification.referenceId}`)
    }
  }
}

const handleDelete = async (id: string) => {
  await notif.deleteNotification(id)
}

onMounted(async () => {
  if (props.isOpen) {
    await notif.fetchNotifications(1, 10)
  }
})

watch(() => props.isOpen, async (open) => {
  if (open) {
    await notif.fetchNotifications(1, 10)
  }
})
</script>
