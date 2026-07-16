<template>
  <nav v-if="breadcrumbs.length > 0" class="px-4 sm:px-6 lg:px-8 py-2 bg-white/50 dark:bg-gray-900/50 border-b border-gray-100 dark:border-gray-800">
    <ol class="flex items-center gap-1.5 text-sm">
      <li>
        <NuxtLink
          to="/dashboard"
          class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
        >
          <UIcon name="i-heroicons-home" class="w-4 h-4" />
        </NuxtLink>
      </li>
      <template v-for="(crumb, index) in breadcrumbs" :key="index">
        <li class="text-gray-300 dark:text-gray-600">
          <UIcon name="i-heroicons-chevron-right" class="w-3.5 h-3.5" />
        </li>
        <li>
          <NuxtLink
            v-if="crumb.to && index < breadcrumbs.length - 1"
            :to="crumb.to"
            class="text-gray-500 dark:text-gray-400 hover:text-brand-600 dark:hover:text-brand-400 transition-colors"
          >
            {{ crumb.label }}
          </NuxtLink>
          <span
            v-else
            class="text-gray-700 dark:text-gray-300 font-medium"
          >
            {{ crumb.label }}
          </span>
        </li>
      </template>
    </ol>
  </nav>
</template>

<script setup lang="ts">
const route = useRoute()
const { t } = useI18n()

const breadcrumbs = computed(() => {
  const crumbs: { label: string; to?: string }[] = []
  const path = route.path

  const segments = path.split('/').filter(Boolean)

  let currentPath = ''

  const routeNameMap: Record<string, string> = {
    dashboard: t('menu.dashboard'),
    academic: t('menu.academic'),
    classes: t('menu.classes'),
    subjects: t('menu.subjects'),
    curriculum: t('menu.curriculum'),
    schedules: t('menu.schedules'),
    attendance: t('menu.attendance'),
    gradebook: t('menu.gradebook'),
    exams: t('menu.exams'),
    'report-cards': t('menu.report_cards'),
    islamic: t('menu.islamic'),
    quran: t('menu.quran_tahfidz'),
    mutabaah: t('menu.mutabaah'),
    prayer: t('menu.prayer'),
    halaqah: t('menu.halaqah'),
    events: t('menu.islamic_events'),
    students: t('menu.students'),
    parents: t('menu.parents'),
    alumni: t('menu.alumni'),
    teachers: t('menu.teachers'),
    employees: t('menu.employees'),
    roles: t('menu.roles_permissions'),
    finance: t('menu.finance'),
    spp: t('menu.spp'),
    invoices: t('menu.invoices'),
    payments: t('menu.payments'),
    journal: t('menu.journal'),
    budget: t('menu.budget'),
    payroll: t('menu.payroll'),
    inventory: t('menu.inventory'),
    items: t('menu.items'),
    assets: t('menu.assets'),
    library: t('menu.library'),
    medical: t('menu.medical'),
    counseling: t('menu.counseling'),
    admissions: t('menu.admissions'),
    communication: t('menu.communication'),
    announcements: t('menu.announcements'),
    messages: t('menu.messages'),
    meetings: t('menu.meetings'),
    documents: t('menu.documents_letters'),
    'ai-assistant': t('menu.ai_assistant'),
    settings: t('menu.settings'),
    profile: t('menu.profile'),
    notifications: t('notifications.title'),
  }

  for (const segment of segments) {
    currentPath += `/${segment}`
    const label = routeNameMap[segment] || segment.charAt(0).toUpperCase() + segment.slice(1)
    crumbs.push({ label, to: currentPath })
  }

  return crumbs
})
</script>
