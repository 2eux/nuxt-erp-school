import { defineStore } from 'pinia'
import type { Notification, SidebarState } from '~/types'

interface AppSettings {
  primaryColor: string
  accentColor: string
  borderRadius: 'none' | 'sm' | 'md' | 'lg'
  language: string
  density: 'compact' | 'comfortable' | 'spacious'
}

interface AppState {
  sidebar: SidebarState
  notifications: Notification[]
  unreadNotifications: number
  settings: AppSettings
  isSearchOpen: boolean
  isMobile: boolean
  loading: boolean
  breadcrumbs: { label: string; to?: string }[]
}

export const useAppStore = defineStore('app', {
  state: (): AppState => ({
    sidebar: {
      isCollapsed: false,
      isMobileOpen: false,
      activeGroup: null,
    },
    notifications: [],
    unreadNotifications: 0,
    settings: {
      primaryColor: 'emerald',
      accentColor: 'gold',
      borderRadius: 'md',
      language: 'id',
      density: 'comfortable',
    },
    isSearchOpen: false,
    isMobile: false,
    loading: false,
    breadcrumbs: [],
  }),

  getters: {
    sidebarCollapsed: (state) => state.sidebar.isCollapsed,
    sidebarMobileOpen: (state) => state.sidebar.isMobileOpen,
    currentBreadcrumbs: (state) => state.breadcrumbs,
  },

  actions: {
    toggleSidebar() {
      this.sidebar.isCollapsed = !this.sidebar.isCollapsed
      if (import.meta.client) {
        localStorage.setItem('sidebar_collapsed', String(this.sidebar.isCollapsed))
      }
    },

    setSidebarCollapsed(collapsed: boolean) {
      this.sidebar.isCollapsed = collapsed
    },

    toggleMobileSidebar() {
      this.sidebar.isMobileOpen = !this.sidebar.isMobileOpen
    },

    setMobileSidebar(open: boolean) {
      this.sidebar.isMobileOpen = open
    },

    setActiveGroup(group: string | null) {
      this.sidebar.activeGroup = group
    },

    setNotifications(notifications: Notification[]) {
      this.notifications = notifications
    },

    setUnreadCount(count: number) {
      this.unreadNotifications = count
    },

    markNotificationRead(id: string) {
      const notif = this.notifications.find(n => n.id === id)
      if (notif && !notif.isRead) {
        notif.isRead = true
        this.unreadNotifications = Math.max(0, this.unreadNotifications - 1)
      }
    },

    markAllNotificationsRead() {
      this.notifications.forEach(n => { n.isRead = true })
      this.unreadNotifications = 0
    },

    removeNotification(id: string) {
      const idx = this.notifications.findIndex(n => n.id === id)
      if (idx !== -1) {
        if (!this.notifications[idx].isRead) {
          this.unreadNotifications = Math.max(0, this.unreadNotifications - 1)
        }
        this.notifications.splice(idx, 1)
      }
    },

    updateSettings(settings: Partial<AppSettings>) {
      this.settings = { ...this.settings, ...settings }
    },

    toggleSearch() {
      this.isSearchOpen = !this.isSearchOpen
    },

    openSearch() {
      this.isSearchOpen = true
    },

    closeSearch() {
      this.isSearchOpen = false
    },

    setMobile(mobile: boolean) {
      this.isMobile = mobile
    },

    setLoading(loading: boolean) {
      this.loading = loading
    },

    setBreadcrumbs(breadcrumbs: { label: string; to?: string }[]) {
      this.breadcrumbs = breadcrumbs
    },

    clearBreadcrumbs() {
      this.breadcrumbs = []
    },

    initFromStorage() {
      if (import.meta.client) {
        const collapsed = localStorage.getItem('sidebar_collapsed')
        if (collapsed === 'true') {
          this.sidebar.isCollapsed = true
        }

        const primaryColor = localStorage.getItem('primary_color')
        if (primaryColor) {
          this.settings.primaryColor = primaryColor
        }

        const accentColor = localStorage.getItem('accent_color')
        if (accentColor) {
          this.settings.accentColor = accentColor
        }
      }
    },
  },
})
