export default defineNuxtRouteMiddleware(async (to) => {
  const authStore = useAuthStore()

  if (!authStore.isInitialized) {
    await authStore.initializeFromStorage()
  }

  const publicPaths = ['/auth/login', '/auth/forgot-password', '/auth/reset-password']

  if (!authStore.isAuthenticated && !publicPaths.includes(to.path)) {
    return navigateTo({
      path: '/auth/login',
      query: { redirect: to.fullPath },
    })
  }

  if (authStore.isAuthenticated && to.path === '/auth/login') {
    return navigateTo('/dashboard')
  }

  if (to.meta.permission) {
    const permissions = Array.isArray(to.meta.permission)
      ? to.meta.permission
      : [to.meta.permission]

    const hasAccess = authStore.isSuperAdmin || permissions.some(
      (p: string) => authStore.permissions.includes(p)
    )

    if (!hasAccess) {
      return navigateTo('/dashboard')
    }
  }

  return undefined
})
