<template>
  <aside
    class="fixed top-0 left-0 z-50 h-full bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800 transition-all duration-300 flex flex-col"
    :class="[
      appStore.sidebar.isCollapsed ? 'w-20' : 'w-64',
      appStore.sidebar.isMobileOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0',
    ]"
  >
    <div class="flex items-center h-16 px-4 border-b border-gray-200 dark:border-gray-800 shrink-0">
      <div
        v-if="!appStore.sidebar.isCollapsed"
        class="flex items-center gap-3 flex-1 min-w-0"
      >
        <div class="w-9 h-9 rounded-lg bg-brand-600 flex items-center justify-center shrink-0">
          <UIcon name="i-heroicons-academic-cap" class="w-5 h-5 text-white" />
        </div>
        <span class="font-bold text-gray-900 dark:text-white truncate text-sm">
          {{ $t('app.short_name') }}
        </span>
      </div>
      <div
        v-else
        class="flex items-center justify-center w-full"
      >
        <div class="w-9 h-9 rounded-lg bg-brand-600 flex items-center justify-center">
          <UIcon name="i-heroicons-academic-cap" class="w-5 h-5 text-white" />
        </div>
      </div>
    </div>

    <nav class="flex-1 overflow-y-auto overflow-x-hidden py-2 px-2">
      <ul class="space-y-1">
        <li v-for="item in menuItems" :key="item.label">
          <template v-if="item.children && item.children.length > 0">
            <button
              class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors"
              :class="[
                appStore.sidebar.activeGroup === item.label
                  ? 'bg-brand-50 text-brand-700 dark:bg-brand-900/20 dark:text-brand-400'
                  : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800'
              ]"
              @click="toggleGroup(item.label)"
            >
              <UIcon :name="item.icon" class="w-5 h-5 shrink-0" />
              <span
                v-if="!appStore.sidebar.isCollapsed"
                class="flex-1 text-left"
              >{{ $t(`menu.${item.label.toLowerCase()}`) }}</span>
              <UIcon
                v-if="!appStore.sidebar.isCollapsed"
                name="i-heroicons-chevron-down"
                class="w-4 h-4 transition-transform"
                :class="appStore.sidebar.activeGroup === item.label ? 'rotate-180' : ''"
              />
            </button>

            <ul
              v-if="!appStore.sidebar.isCollapsed && appStore.sidebar.activeGroup === item.label"
              class="mt-1 ml-4 pl-6 border-l-2 border-gray-200 dark:border-gray-700 space-y-1"
            >
              <li v-for="child in item.children" :key="child.label">
                <NuxtLink
                  v-if="child.to && (!child.permission || hasPermission(child.permission))"
                  :to="child.to"
                  class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors text-gray-600 dark:text-gray-400 hover:text-brand-700 dark:hover:text-brand-400 hover:bg-brand-50 dark:hover:bg-brand-900/20"
                  :class="{ 'text-brand-700 bg-brand-50 dark:text-brand-400 dark:bg-brand-900/20': isActiveRoute(child.to) }"
                >
                  <UIcon :name="child.icon" class="w-4 h-4" />
                  <span>{{ $t(`menu.${child.label.toLowerCase().replace(/\s+/g, '_')}`) }}</span>
                  <span
                    v-if="child.badge"
                    class="ml-auto px-1.5 py-0.5 text-xs font-medium rounded-full bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400"
                  >
                    {{ child.badge }}
                  </span>
                </NuxtLink>
              </li>
            </ul>
          </template>

          <template v-else>
            <NuxtLink
              v-if="item.to && (!item.permission || hasPermission(item.permission))"
              :to="item.to"
              class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors"
              :class="[
                isActiveRoute(item.to)
                  ? 'bg-brand-50 text-brand-700 dark:bg-brand-900/20 dark:text-brand-400'
                  : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800'
              ]"
            >
              <UIcon :name="item.icon" class="w-5 h-5 shrink-0" />
              <span v-if="!appStore.sidebar.isCollapsed">{{ $t(`menu.${item.label.toLowerCase()}`) }}</span>
              <span
                v-if="item.badge"
                class="ml-auto px-1.5 py-0.5 text-xs font-medium rounded-full bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400"
              >
                {{ item.badge }}
              </span>
            </NuxtLink>
          </template>
        </li>
      </ul>
    </nav>

    <div class="border-t border-gray-200 dark:border-gray-800 p-3 shrink-0">
      <UButton
        color="gray"
        variant="ghost"
        :icon="appStore.sidebar.isCollapsed ? 'i-heroicons-chevron-double-right' : 'i-heroicons-chevron-double-left'"
        class="w-full justify-center"
        @click="appStore.toggleSidebar()"
      >
        <span v-if="!appStore.sidebar.isCollapsed" class="text-xs">{{ $t('common.collapse') }}</span>
      </UButton>
    </div>
  </aside>
</template>

<script setup lang="ts">
import type { MenuItem } from '~/types'

const appStore = useAppStore()
const authStore = useAuthStore()
const route = useRoute()

const menuItems: MenuItem[] = [
  {
    label: 'Dashboard',
    icon: 'i-heroicons-home',
    to: '/dashboard',
    exact: true,
  },
  {
    label: 'Academic',
    icon: 'i-heroicons-book-open',
    permission: 'academic.view',
    children: [
      { label: 'Classes', icon: 'i-heroicons-rectangle-group', to: '/academic/classes', permission: 'classes.view' },
      { label: 'Subjects', icon: 'i-heroicons-bookmark', to: '/academic/subjects', permission: 'subjects.view' },
      { label: 'Curriculum', icon: 'i-heroicons-clipboard-document-list', to: '/academic/curriculum', permission: 'curriculum.view' },
      { label: 'Schedules', icon: 'i-heroicons-clock', to: '/academic/schedules', permission: 'schedules.view' },
      { label: 'Attendance', icon: 'i-heroicons-check-badge', to: '/academic/attendance', permission: 'attendance.view' },
      { label: 'Gradebook', icon: 'i-heroicons-chart-bar', to: '/academic/gradebook', permission: 'gradebook.view' },
      { label: 'Exams', icon: 'i-heroicons-pencil-square', to: '/academic/exams', permission: 'exams.view' },
      { label: 'Report Cards', icon: 'i-heroicons-document-text', to: '/academic/report-cards', permission: 'report_cards.view' },
    ],
  },
  {
    label: 'Islamic',
    icon: 'i-heroicons-sparkles',
    permission: 'islamic.view',
    children: [
      { label: 'Quran & Tahfidz', icon: 'i-heroicons-book-open', to: '/islamic/quran', permission: 'quran.view' },
      { label: "Mutaba'ah", icon: 'i-heroicons-clipboard-document-check', to: '/islamic/mutabaah', permission: 'mutabaah.view' },
      { label: 'Prayer', icon: 'i-heroicons-sun', to: '/islamic/prayer', permission: 'prayer.view' },
      { label: 'Halaqah', icon: 'i-heroicons-user-group', to: '/islamic/halaqah', permission: 'halaqah.view' },
      { label: 'Islamic Events', icon: 'i-heroicons-calendar-days', to: '/islamic/events', permission: 'islamic_events.view' },
    ],
  },
  {
    label: 'Students',
    icon: 'i-heroicons-users',
    permission: 'students.view',
    children: [
      { label: 'All Students', icon: 'i-heroicons-list-bullet', to: '/students' },
      { label: 'Parents', icon: 'i-heroicons-heart', to: '/students/parents', permission: 'parents.view' },
      { label: 'Alumni', icon: 'i-heroicons-academic-cap', to: '/students/alumni', permission: 'alumni.view' },
    ],
  },
  {
    label: 'Teachers & HR',
    icon: 'i-heroicons-user-circle',
    permission: 'teachers.view',
    children: [
      { label: 'Teachers', icon: 'i-heroicons-user', to: '/teachers' },
      { label: 'Employees', icon: 'i-heroicons-briefcase', to: '/employees', permission: 'employees.view' },
      { label: 'Roles & Permissions', icon: 'i-heroicons-shield-check', to: '/roles', permission: 'roles.view' },
    ],
  },
  {
    label: 'Finance',
    icon: 'i-heroicons-currency-dollar',
    permission: 'finance.view',
    children: [
      { label: 'SPP', icon: 'i-heroicons-receipt-percent', to: '/finance/spp' },
      { label: 'Invoices', icon: 'i-heroicons-document-currency-dollar', to: '/finance/invoices' },
      { label: 'Payments', icon: 'i-heroicons-credit-card', to: '/finance/payments' },
      { label: 'Journal', icon: 'i-heroicons-book-open', to: '/finance/journal', permission: 'journals.view' },
      { label: 'Budget', icon: 'i-heroicons-scale', to: '/finance/budget', permission: 'budget.view' },
      { label: 'Payroll', icon: 'i-heroicons-banknotes', to: '/finance/payroll', permission: 'payroll.view' },
    ],
  },
  {
    label: 'Inventory',
    icon: 'i-heroicons-archive-box',
    permission: 'inventory.view',
    children: [
      { label: 'Items', icon: 'i-heroicons-tag', to: '/inventory/items' },
      { label: 'Assets', icon: 'i-heroicons-building-office', to: '/inventory/assets', permission: 'assets.view' },
    ],
  },
  {
    label: 'Library',
    icon: 'i-heroicons-building-library',
    to: '/library',
    permission: 'library.view',
  },
  {
    label: 'Medical',
    icon: 'i-heroicons-heart',
    to: '/medical',
    permission: 'medical.view',
  },
  {
    label: 'Counseling',
    icon: 'i-heroicons-chat-bubble-left-right',
    to: '/counseling',
    permission: 'counseling.view',
  },
  {
    label: 'Admissions',
    icon: 'i-heroicons-arrow-trending-up',
    to: '/admissions',
    permission: 'admissions.view',
  },
  {
    label: 'Communication',
    icon: 'i-heroicons-megaphone',
    permission: 'announcements.view',
    children: [
      { label: 'Announcements', icon: 'i-heroicons-speaker-wave', to: '/communication/announcements' },
      { label: 'Messages', icon: 'i-heroicons-envelope', to: '/communication/messages', permission: 'messages.view' },
      { label: 'Meetings', icon: 'i-heroicons-video-camera', to: '/communication/meetings', permission: 'meetings.view' },
    ],
  },
  {
    label: 'Documents & Letters',
    icon: 'i-heroicons-document-duplicate',
    to: '/documents',
    permission: 'documents.view',
  },
  {
    label: 'AI Assistant',
    icon: 'i-heroicons-cpu-chip',
    to: '/ai-assistant',
    permission: 'ai.view',
  },
  {
    label: 'Settings',
    icon: 'i-heroicons-cog-6-tooth',
    to: '/settings',
    permission: 'settings.view',
  },
]

const hasPermission = (perm: string): boolean => {
  if (authStore.user?.isSuperAdmin) return true
  return authStore.permissions.includes(perm)
}

const toggleGroup = (label: string): void => {
  appStore.setActiveGroup(appStore.sidebar.activeGroup === label ? null : label)
}

const isActiveRoute = (path: string): boolean => {
  return route.path.startsWith(path)
}

onMounted(() => {
  const currentGroup = menuItems.find(item =>
    item.children?.some(child => child.to && route.path.startsWith(child.to))
  )
  if (currentGroup) {
    appStore.setActiveGroup(currentGroup.label)
  }
})
</script>
