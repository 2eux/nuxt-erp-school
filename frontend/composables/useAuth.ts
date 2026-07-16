import type { AuthUser, LoginRequest, LoginResponse } from '~/types'

export const useAuth = () => {
  const authStore = useAuthStore()
  const api = useApi()
  const router = useRouter()
  const { t } = useI18n()
  const toast = useToast()

  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const currentUser = computed(() => authStore.user)
  const permissions = computed(() => authStore.permissions)

  const login = async (credentials: LoginRequest): Promise<LoginResponse> => {
    try {
      const response = await api.post<LoginResponse>('/auth/login', credentials)
      authStore.setAuth(response.user, response.tokens)

      if (credentials.rememberMe) {
        if (import.meta.client) {
          localStorage.setItem('remembered_email', credentials.email)
        }
      }

      toast.add({
        title: t('auth.login_success'),
        description: t('auth.welcome_back', { name: response.user.fullName }),
        color: 'success',
      })

      return response
    } catch (error) {
      toast.add({
        title: t('auth.login_failed'),
        description: t('auth.invalid_credentials'),
        color: 'error',
      })
      throw error
    }
  }

  const logout = async () => {
    try {
      await api.post('/auth/logout')
    } catch {
      // ignore logout API errors
    } finally {
      authStore.clearAuth()
      toast.add({
        title: t('auth.logout_success'),
        color: 'info',
      })
      await router.push('/auth/login')
    }
  }

  const refresh = async (): Promise<boolean> => {
    const refreshToken = authStore.refreshToken
    if (!refreshToken) return false

    try {
      const response = await api.post<LoginResponse>('/auth/refresh', { refreshToken })
      authStore.setAuth(response.user, response.tokens)
      return true
    } catch {
      authStore.clearAuth()
      return false
    }
  }

  const forgotPassword = async (email: string): Promise<void> => {
    await api.post('/auth/forgot-password', { email })
    toast.add({
      title: t('auth.reset_link_sent'),
      description: t('auth.check_email'),
      color: 'success',
    })
  }

  const resetPassword = async (token: string, password: string, passwordConfirmation: string): Promise<void> => {
    await api.post('/auth/reset-password', { token, password, passwordConfirmation })
    toast.add({
      title: t('auth.password_reset_success'),
      color: 'success',
    })
  }

  const changePassword = async (currentPassword: string, newPassword: string): Promise<void> => {
    await api.put('/auth/change-password', { currentPassword, newPassword })
    toast.add({
      title: t('auth.password_changed'),
      color: 'success',
    })
  }

  const updateProfile = async (data: Partial<AuthUser>): Promise<void> => {
    await api.put('/auth/profile', data)
    authStore.updateUser(data)
    toast.add({
      title: t('auth.profile_updated'),
      color: 'success',
    })
  }

  const hasPermission = (permission: string | string[]): boolean => {
    if (authStore.user?.isSuperAdmin) return true
    const perms = Array.isArray(permission) ? permission : [permission]
    return perms.some(p => authStore.permissions.includes(p))
  }

  const hasRole = (roleName: string | string[]): boolean => {
    const roles = Array.isArray(roleName) ? roleName : [roleName]
    return roles.some(r => authStore.user?.roleName === r)
  }

  return {
    isAuthenticated,
    currentUser,
    permissions,
    login,
    logout,
    refresh,
    forgotPassword,
    resetPassword,
    changePassword,
    updateProfile,
    hasPermission,
    hasRole,
  }
}
