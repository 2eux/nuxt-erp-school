import type { Notification } from '~/types'

export const useNotification = () => {
  const appStore = useAppStore()
  const api = useApi()
  const authStore = useAuthStore()
  const { t } = useI18n()
  const toast = useToast()

  const unreadCount = computed(() => appStore.unreadNotifications)
  const notifications = computed(() => appStore.notifications)
  const isLoading = ref(false)

  let pollingInterval: ReturnType<typeof setInterval> | null = null

  const fetchNotifications = async (page = 1, limit = 20): Promise<void> => {
    isLoading.value = true
    try {
      const response = await api.paginate<Notification>('/notifications', { page, limit })
      appStore.setNotifications(response.data)
      appStore.setUnreadCount(response.data.filter(n => !n.isRead).length)
    } finally {
      isLoading.value = false
    }
  }

  const markAsRead = async (id: string): Promise<void> => {
    try {
      await api.patch(`/notifications/${id}/read`)
      appStore.markNotificationRead(id)
    } catch {
      // handled by useApi
    }
  }

  const markAllAsRead = async (): Promise<void> => {
    try {
      await api.post('/notifications/read-all')
      appStore.markAllNotificationsRead()
    } catch {
      // handled by useApi
    }
  }

  const deleteNotification = async (id: string): Promise<void> => {
    try {
      await api.delete(`/notifications/${id}`)
      appStore.removeNotification(id)
    } catch {
      // handled by useApi
    }
  }

  const startPolling = (intervalMs = 30000): void => {
    if (pollingInterval) return

    pollingInterval = setInterval(async () => {
      if (authStore.isAuthenticated) {
        await fetchNotifications(1, 5)
      }
    }, intervalMs)
  }

  const stopPolling = (): void => {
    if (pollingInterval) {
      clearInterval(pollingInterval)
      pollingInterval = null
    }
  }

  const showToast = (
    title: string,
    description?: string,
    type: 'success' | 'error' | 'warning' | 'info' = 'info'
  ): void => {
    toast.add({ title, description, color: type === 'info' ? 'info' : type })
  }

  return {
    unreadCount,
    notifications,
    isLoading,
    fetchNotifications,
    markAsRead,
    markAllAsRead,
    deleteNotification,
    startPolling,
    stopPolling,
    showToast,
  }
}
