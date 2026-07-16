import { defineStore } from 'pinia'
import type { AuthUser, AuthTokens } from '~/types'

interface AuthState {
  user: AuthUser | null
  accessToken: string | null
  refreshToken: string | null
  permissions: string[]
  isAuthenticated: boolean
  isInitialized: boolean
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    accessToken: null,
    refreshToken: null,
    permissions: [],
    isAuthenticated: false,
    isInitialized: false,
  }),

  getters: {
    userId: (state) => state.user?.id || null,
    schoolId: (state) => state.user?.schoolId || null,
    roleName: (state) => state.user?.roleName || null,
    fullName: (state) => state.user?.fullName || '',
    avatar: (state) => state.user?.avatar || null,
  },

  actions: {
    setAuth(user: AuthUser, tokens: AuthTokens) {
      this.user = user
      this.accessToken = tokens.accessToken
      this.refreshToken = tokens.refreshToken
      this.permissions = user.permissions || []
      this.isAuthenticated = true
      this.isInitialized = true

      if (import.meta.client) {
        localStorage.setItem('access_token', tokens.accessToken)
        localStorage.setItem('refresh_token', tokens.refreshToken)
      }
    },

    setTokens(tokens: AuthTokens) {
      this.accessToken = tokens.accessToken
      this.refreshToken = tokens.refreshToken

      if (import.meta.client) {
        localStorage.setItem('access_token', tokens.accessToken)
        localStorage.setItem('refresh_token', tokens.refreshToken)
      }
    },

    updateUser(data: Partial<AuthUser>) {
      if (this.user) {
        this.user = { ...this.user, ...data }
      }
    },

    setPermissions(permissions: string[]) {
      this.permissions = permissions
    },

    async initializeFromStorage(): Promise<boolean> {
      if (!import.meta.client) return false

      const accessToken = localStorage.getItem('access_token')
      const refreshToken = localStorage.getItem('refresh_token')

      if (!accessToken || !refreshToken) {
        this.isInitialized = true
        return false
      }

      this.accessToken = accessToken
      this.refreshToken = refreshToken

      try {
        const api = useApi()
        const response = await api.get<{ user: AuthUser; tokens: AuthTokens }>('/auth/me')
        this.setAuth(response.user, response.tokens)
        return true
      } catch {
        this.clearAuth()
        return false
      }
    },

    clearAuth() {
      this.user = null
      this.accessToken = null
      this.refreshToken = null
      this.permissions = []
      this.isAuthenticated = false
      this.isInitialized = true

      if (import.meta.client) {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        localStorage.removeItem('remembered_email')
      }
    },

    markInitialized() {
      this.isInitialized = true
    },
  },
})
