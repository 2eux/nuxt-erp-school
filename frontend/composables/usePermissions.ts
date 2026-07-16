export const usePermissions = () => {
  const authStore = useAuthStore()
  const { t } = useI18n()
  const toast = useToast()

  const userPermissions = computed(() => authStore.permissions)
  const isSuperAdmin = computed(() => authStore.user?.isSuperAdmin || false)
  const roleName = computed(() => authStore.user?.roleName || '')

  const can = (permission: string | string[]): boolean => {
    if (isSuperAdmin.value) return true

    const permissions = Array.isArray(permission) ? permission : [permission]
    return permissions.some(p => userPermissions.value.includes(p))
  }

  const canAll = (permissions: string[]): boolean => {
    if (isSuperAdmin.value) return true
    return permissions.every(p => userPermissions.value.includes(p))
  }

  const canAny = (permissions: string[]): boolean => {
    if (isSuperAdmin.value) return true
    return permissions.some(p => userPermissions.value.includes(p))
  }

  const assertCan = (permission: string, fallbackRoute = '/dashboard'): boolean => {
    if (!can(permission)) {
      toast.add({
        title: t('error.unauthorized'),
        description: t('error.no_permission'),
        color: 'error',
      })
      navigateTo(fallbackRoute)
      return false
    }
    return true
  }

  const hasModuleAccess = (module: string): boolean => {
    const modulePermissionMap: Record<string, string[]> = {
      dashboard: ['dashboard.view'],
      academic: ['academic.view', 'classes.view', 'subjects.view', 'schedules.view', 'attendance.view'],
      islamic: ['islamic.view', 'quran.view', 'mutabaah.view', 'halaqah.view'],
      students: ['students.view'],
      teachers: ['teachers.view', 'employees.view'],
      finance: ['finance.view', 'invoices.view', 'payments.view', 'journals.view', 'budget.view', 'payroll.view'],
      inventory: ['inventory.view', 'assets.view'],
      library: ['library.view'],
      medical: ['medical.view'],
      counseling: ['counseling.view'],
      admissions: ['admissions.view'],
      communication: ['announcements.view', 'messages.view'],
      documents: ['documents.view'],
      settings: ['settings.view'],
      ai: ['ai.view'],
    }

    const required = modulePermissionMap[module]
    if (!required) return false

    return canAny(required)
  }

  const filterMenuItemsByPermission = <T extends { permission?: string }>(items: T[]): T[] => {
    return items.filter(item => {
      if (!item.permission) return true
      return can(item.permission)
    })
  }

  return {
    userPermissions,
    isSuperAdmin,
    roleName,
    can,
    canAll,
    canAny,
    assertCan,
    hasModuleAccess,
    filterMenuItemsByPermission,
  }
}
