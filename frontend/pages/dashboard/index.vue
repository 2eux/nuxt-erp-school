<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ greeting }}, {{ authStore.user?.fullName || $t('common.user') }}!
        </h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          {{ $t('dashboard.subtitle') }}
        </p>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard
        :label="$t('dashboard.total_students')"
        :value="stats.totalStudents || 0"
        icon="i-heroicons-users"
        color="blue"
        :trend="stats.studentGrowth"
        :trend-label="$t('dashboard.vs_last_month')"
        :loading="loading"
      />
      <StatCard
        :label="$t('dashboard.total_teachers')"
        :value="stats.totalTeachers || 0"
        icon="i-heroicons-user-circle"
        color="emerald"
        :loading="loading"
      />
      <StatCard
        :label="$t('dashboard.total_revenue')"
        :value="formattedRevenue"
        icon="i-heroicons-currency-dollar"
        color="amber"
        :trend="stats.revenueGrowth"
        :trend-label="$t('dashboard.vs_last_month')"
        :loading="loading"
      />
      <StatCard
        :label="$t('dashboard.attendance_rate')"
        :value="stats.attendanceRate ? `${stats.attendanceRate}%` : '0%'"
        icon="i-heroicons-check-badge"
        color="purple"
        :loading="loading"
      />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 card">
        <div class="card-header">
          <h3 class="card-title">{{ $t('dashboard.enrollment_trend') }}</h3>
          <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-ellipsis-horizontal" />
        </div>
        <div class="h-72">
          <template v-if="loading">
            <LoadingSkeleton type="detail" />
          </template>
          <template v-else>
            <ApexChart
              type="area"
              height="100%"
              :options="enrollmentChart.options"
              :series="enrollmentChart.series"
            />
          </template>
        </div>
      </div>

      <div class="card">
        <div class="card-header">
          <h3 class="card-title">{{ $t('dashboard.prayer_attendance') }}</h3>
        </div>
        <div class="h-72">
          <template v-if="loading">
            <LoadingSkeleton type="detail" />
          </template>
          <template v-else>
            <ApexChart
              type="donut"
              height="100%"
              :options="prayerChart.options"
              :series="prayerChart.series"
            />
          </template>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="card">
        <div class="card-header">
          <h3 class="card-title">{{ $t('dashboard.recent_activities') }}</h3>
          <NuxtLink to="/activities" class="text-sm text-brand-600 hover:text-brand-500 font-medium">
            {{ $t('common.view_all') }}
          </NuxtLink>
        </div>

        <div v-if="loading" class="space-y-3">
          <div v-for="i in 4" :key="i" class="flex items-center gap-3">
            <div class="w-2 h-2 rounded-full bg-gray-200 dark:bg-gray-700 animate-pulse" />
            <div class="flex-1 space-y-1.5">
              <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded animate-pulse w-3/4" />
              <div class="h-2 bg-gray-200 dark:bg-gray-700 rounded animate-pulse w-1/2" />
            </div>
          </div>
        </div>

        <div v-else-if="recentActivities.length === 0" class="py-8">
          <EmptyState
            :title="$t('dashboard.no_recent_activities')"
            :description="$t('dashboard.no_recent_activities_description')"
            icon="i-heroicons-clock"
          />
        </div>

        <ul v-else class="space-y-3">
          <li v-for="activity in recentActivities" :key="activity.id" class="flex items-center gap-3">
            <div class="w-2 h-2 rounded-full bg-brand-500 mt-1 shrink-0" />
            <div class="flex-1 min-w-0">
              <p class="text-sm text-gray-900 dark:text-white">{{ activity.title }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400 truncate">{{ activity.description }}</p>
              <p class="text-xs text-gray-400 mt-0.5">{{ formatTime(activity.timestamp) }}</p>
            </div>
          </li>
        </ul>
      </div>

      <div class="card">
        <div class="card-header">
          <h3 class="card-title">{{ $t('dashboard.upcoming_events') }}</h3>
          <NuxtLink to="/communication/events" class="text-sm text-brand-600 hover:text-brand-500 font-medium">
            {{ $t('common.view_all') }}
          </NuxtLink>
        </div>

        <div v-if="loading" class="space-y-3">
          <div v-for="i in 3" :key="i" class="flex gap-3">
            <div class="w-12 h-12 rounded-lg bg-gray-200 dark:bg-gray-700 animate-pulse" />
            <div class="flex-1 space-y-1.5">
              <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded animate-pulse w-3/4" />
              <div class="h-2 bg-gray-200 dark:bg-gray-700 rounded animate-pulse w-1/2" />
            </div>
          </div>
        </div>

        <div v-else-if="upcomingEvents.length === 0" class="py-8">
          <EmptyState
            :title="$t('dashboard.no_upcoming_events')"
            icon="i-heroicons-calendar"
          />
        </div>

        <ul v-else class="space-y-3">
          <li v-for="event in upcomingEvents" :key="event.id" class="flex items-start gap-3">
            <div class="w-12 h-12 rounded-lg bg-brand-50 dark:bg-brand-900/20 flex flex-col items-center justify-center shrink-0">
              <span class="text-xs font-bold text-brand-600 dark:text-brand-400">
                {{ formatEventDate(event.date).day }}
              </span>
              <span class="text-[10px] text-brand-500 dark:text-brand-400 uppercase">
                {{ formatEventDate(event.date).month }}
              </span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900 dark:text-white">{{ event.title }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ event.description }}</p>
            </div>
          </li>
        </ul>
      </div>

      <div class="card">
        <div class="card-header">
          <h3 class="card-title">{{ $t('dashboard.tahfidz_progress') }}</h3>
        </div>

        <div v-if="loading" class="space-y-4">
          <div v-for="i in 4" :key="i" class="space-y-1.5">
            <div class="flex justify-between">
              <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded animate-pulse w-1/3" />
              <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded animate-pulse w-12" />
            </div>
            <div class="h-2 bg-gray-200 dark:bg-gray-700 rounded-full animate-pulse" />
          </div>
        </div>

        <div v-else class="space-y-4">
          <div v-for="item in tahfidzData" :key="item.name" class="space-y-1">
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-700 dark:text-gray-300">{{ item.name }}</span>
              <span class="text-sm font-semibold text-brand-600 dark:text-brand-400">
                {{ item.progress }}%
              </span>
            </div>
            <div class="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
              <div
                class="h-full bg-brand-500 rounded-full transition-all duration-500"
                :style="{ width: `${item.progress}%` }"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="card-header">
        <h3 class="card-title">{{ $t('dashboard.announcements') }}</h3>
        <NuxtLink to="/communication/announcements" class="text-sm text-brand-600 hover:text-brand-500 font-medium">
          {{ $t('common.view_all') }}
        </NuxtLink>
      </div>

      <div v-if="loading" class="space-y-3">
        <div v-for="i in 3" :key="i" class="flex items-start gap-3 p-3">
          <div class="w-8 h-8 rounded-full bg-gray-200 dark:bg-gray-700 animate-pulse shrink-0" />
          <div class="flex-1 space-y-1.5">
            <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded animate-pulse w-1/2" />
            <div class="h-2 bg-gray-200 dark:bg-gray-700 rounded animate-pulse w-full" />
          </div>
        </div>
      </div>

      <div v-else-if="announcements.length === 0" class="py-8">
        <EmptyState
          :title="$t('dashboard.no_announcements')"
          icon="i-heroicons-speaker-wave"
        />
      </div>

      <div v-else class="divide-y divide-gray-100 dark:divide-gray-700">
        <div
          v-for="announcement in announcements"
          :key="announcement.id"
          class="flex items-start gap-3 px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors rounded-lg cursor-pointer"
        >
          <div
            class="w-8 h-8 rounded-full flex items-center justify-center shrink-0 mt-0.5"
            :class="{
              'bg-red-100 dark:bg-red-900/30': announcement.type === 'urgent',
              'bg-blue-100 dark:bg-blue-900/30': announcement.type === 'academic',
              'bg-emerald-100 dark:bg-emerald-900/30': announcement.type === 'finance',
              'bg-amber-100 dark:bg-amber-900/30': announcement.type === 'event',
              'bg-gray-100 dark:bg-gray-700': announcement.type === 'general',
            }"
          >
            <UIcon
              :name="announcementTypeIcon[announcement.type] || 'i-heroicons-megaphone'"
              class="w-4 h-4"
              :class="{
                'text-red-600': announcement.type === 'urgent',
                'text-blue-600': announcement.type === 'academic',
                'text-emerald-600': announcement.type === 'finance',
                'text-amber-600': announcement.type === 'event',
                'text-gray-500': announcement.type === 'general',
              }"
            />
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-gray-900 dark:text-white">
              {{ announcement.title }}
            </p>
            <p class="text-xs text-gray-500 dark:text-gray-400 line-clamp-2 mt-0.5">
              {{ announcement.content }}
            </p>
            <div class="flex items-center gap-2 mt-1">
              <span class="text-xs text-gray-400">{{ formatTime(announcement.publishDate) }}</span>
              <StatusBadge :status="announcement.type" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DashboardStats, RecentActivity, Announcement } from '~/types'

definePageMeta({
  middleware: ['auth'],
})

const authStore = useAuthStore()
const api = useApi()
const { $dayjs } = useNuxtApp()
const { t } = useI18n()

const loading = ref(true)
const stats = ref<DashboardStats>({
  totalStudents: 0,
  totalTeachers: 0,
  totalEmployees: 0,
  totalClasses: 0,
  totalRevenue: 0,
  attendanceRate: 0,
  studentGrowth: 0,
  revenueGrowth: 0,
})
const recentActivities = ref<RecentActivity[]>([])
const upcomingEvents = ref<{ id: string; title: string; description: string; date: string; type: string }[]>([])
const announcements = ref<Announcement[]>([])

const announcementTypeIcon: Record<string, string> = {
  urgent: 'i-heroicons-exclamation-triangle',
  academic: 'i-heroicons-academic-cap',
  finance: 'i-heroicons-currency-dollar',
  event: 'i-heroicons-calendar-days',
  general: 'i-heroicons-megaphone',
}

const tahfidzData = computed(() => [
  { name: t('quran.juz_30'), progress: 85 },
  { name: t('quran.juz_29'), progress: 62 },
  { name: t('quran.juz_28'), progress: 40 },
  { name: t('quran.juz_1'), progress: 25 },
])

const greeting = computed(() => {
  const hour = $dayjs().hour()
  if (hour < 12) return t('dashboard.good_morning')
  if (hour < 15) return t('dashboard.good_afternoon')
  if (hour < 18) return t('dashboard.good_evening')
  return t('dashboard.good_night')
})

const formattedRevenue = computed(() => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    notation: 'compact',
    maximumFractionDigits: 1,
  }).format(stats.value.totalRevenue)
})

const enrollmentChart = computed(() => ({
  series: [
    {
      name: t('dashboard.new_students'),
      data: [45, 52, 38, 60, 55, 48, 58, 65, 50, 42, 55, 48],
    },
    {
      name: t('dashboard.enrolled'),
      data: [320, 332, 301, 334, 340, 350, 355, 360, 365, 358, 362, 368],
    },
  ],
  options: {
    chart: {
      type: 'area' as const,
      toolbar: { show: false },
      fontFamily: 'Inter, sans-serif',
      background: 'transparent',
    },
    dataLabels: { enabled: false },
    stroke: { curve: 'smooth' as const, width: 2 },
    xaxis: {
      categories: ['Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec', 'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
      labels: { style: { colors: '#9ca3af', fontSize: '12px' } },
    },
    yaxis: {
      labels: { style: { colors: '#9ca3af', fontSize: '12px' } },
    },
    colors: ['#059669', '#6366f1'],
    fill: {
      type: 'gradient',
      gradient: {
        shadeIntensity: 1,
        opacityFrom: 0.15,
        opacityTo: 0.05,
      },
    },
    grid: { borderColor: '#e5e7eb', strokeDashArray: 4 },
    legend: { show: false },
    tooltip: { theme: 'dark' },
  },
}))

const prayerChart = computed(() => ({
  series: [78, 12, 7, 3],
  options: {
    chart: {
      type: 'donut' as const,
      toolbar: { show: false },
      background: 'transparent',
    },
    labels: [t('attendance.ontime'), t('attendance.late'), t('attendance.absent'), t('attendance.excused')],
    colors: ['#059669', '#f59e0b', '#ef4444', '#6366f1'],
    legend: { show: true, position: 'bottom' as const, fontSize: '12px' },
    dataLabels: { enabled: false },
    plotOptions: {
      pie: {
        donut: {
          size: '65%',
          labels: {
            show: true,
            total: {
              show: true,
              label: t('attendance.present'),
              fontSize: '14px',
              color: '#059669',
            },
          },
        },
      },
    },
  },
}))

const formatTime = (date: string): string => {
  return $dayjs(date).fromNow()
}

const formatEventDate = (date: string): { day: string; month: string } => {
  return {
    day: $dayjs(date).format('DD'),
    month: $dayjs(date).format('MMM'),
  }
}

const fetchDashboardData = async () => {
  loading.value = true
  try {
    const [statsRes, activitiesRes, eventsRes, announcementsRes] = await Promise.all([
      api.get<DashboardStats>('/dashboard/stats').catch(() => stats.value),
      api.get<RecentActivity[]>('/dashboard/activities').catch(() => []),
      api.get<{ id: string; title: string; description: string; date: string; type: string }[]>('/dashboard/events').catch(() => []),
      api.get<Announcement[]>('/announcements', { limit: 3 }).catch(() => []),
    ])

    stats.value = statsRes as DashboardStats
    recentActivities.value = activitiesRes as RecentActivity[]
    upcomingEvents.value = eventsRes as { id: string; title: string; description: string; date: string; type: string }[]
    announcements.value = announcementsRes as Announcement[]
  } catch {
    // use defaults
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDashboardData()
})
</script>
